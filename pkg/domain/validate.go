package domain

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

type ChromeRunner struct {
	AllocContext       context.Context
	Context            context.Context
	CancelAllocContext context.CancelFunc
	CancelContext      context.CancelFunc

	cacheDir string
}

func (engine *Engine) RunDomainValidate(ctx context.Context, chunkSize int) error {

	domains, err := engine.getDomainsToValidate()
	if err != nil {
		return err
	}

	fmt.Printf("Validating %d domains\n", len(domains))

	latestUrls, err := engine.getLatestArticleURLs()
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	for _, domain := range domains {

		domain.LiveLatestArticleURL = latestUrls[domain.Domain]

	}

	fmt.Printf("Starting workers\n")

	printSize := 100

	jobs := make(chan *Domain, len(domains))
	results := make(chan *Domain, len(domains))

	workers := runtime.NumCPU() - 2

	for w := 1; w <= workers; w++ {
		go engine.validateDomainWorker(jobs, results)
	}

	for j := 1; j <= len(domains); j++ {
		jobs <- domains[j-1]
	}
	close(jobs)

	done := 0
	t := time.Now().UnixMilli()
	txn, _ := engine.db.Begin()
	defer txn.Rollback()

	for a := 0; a < len(domains); a++ {
		select {
		case <-ctx.Done():
			txn.Commit()
			return ctx.Err()
		case domain := <-results:

			err = engine.Update(txn, domain)
			if err != nil {
				fmt.Println(domain.Domain, err)
				continue
			}

			if a > 0 && a%chunkSize == 0 {

				err := txn.Commit()
				if err != nil {
					return err
				}
				txn, _ = engine.db.Begin()
			}

			if a > 0 && a%printSize == 0 {
				diff := time.Now().UnixMilli() - t
				qps := (float64(printSize) / float64(diff)) * 1000
				t = time.Now().UnixMilli()
				fmt.Printf("\tdone %d/%d at %0.2f/s\n", done, len(domains), qps)

			}
			done++
		}
	}

	txn.Commit()
	fmt.Printf("\tdone %d/%d\n\n", done, len(domains))

	return nil
}

func (engine *Engine) validateDomainWorker(jobs <-chan *Domain, results chan<- *Domain) {
	for id := range jobs {
		err := engine.validateDomain(id)
		if err != nil {
			fmt.Println(id.Domain, err)
		}
		results <- id
	}
}

func (engine *Engine) validateDomain(domain *Domain) error {
	domain.LastValidateAt = time.Now().Unix()
	// First check the URL isn't banned

	fullDomain := fmt.Sprintf("https://%s", domain.Domain)

	urlHumanName, urlNews, urlBlog, popularDomain, err := engine.identifyURL(fullDomain)
	if err != nil {
		return fmt.Errorf("could not identify url %w", err)
	}

	domain.URLBlog = urlBlog
	domain.URLHumanName = urlHumanName
	domain.URLNews = urlNews
	domain.DomainIsPopular = popularDomain

	if domain.LiveLatestArticleURL != "" {
		analysis, err := engine.chrome.GetChromeAnalysis(domain.LiveLatestArticleURL)
		if err != nil {
			return err
		}

		domain.ChromeAnalysis = analysis

	}

	return nil

}

func startChrome(cacheDir string) (*ChromeRunner, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.UserDataDir("/workspaces/100kb.golang/.chrome"),
		chromedp.Flag("headless", "new"),
		chromedp.Flag("no-first-run", true),
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.Flag("disable-gl-drawing-for-tests", true),
		chromedp.Flag("disable-web-security", true),
		chromedp.Flag("disable-site-isolation-trials", true),
		chromedp.Flag("disable-site-isolation-for-policy", true),
		chromedp.Flag("disable-features", "StrictOriginIsolation,IsolateOrigins"),
		chromedp.Flag("hide-scrollbars", true),
		chromedp.Flag("mute-audio", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-setuid-sandbox", true),
		chromedp.Flag("disable-translate", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disk-cache-size", "0"),
	)

	_, err := os.Stat("/chrome/chrome-linux64/chrome")
	if err == nil {
		opts = append(opts, chromedp.ExecPath("/chrome/chrome-linux64/chrome"))
	}

	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)

	var contextOpts []chromedp.ContextOption

	// contextOpts = append(contextOpts, chromedp.WithDebugf(log.Printf))

	ctx, cancel := chromedp.NewContext(allocCtx, contextOpts...)

	err = chromedp.Run(ctx)
	if err != nil {
		return nil, err
	}

	runner := ChromeRunner{
		AllocContext:       allocCtx,
		Context:            ctx,
		CancelAllocContext: allocCancel,
		CancelContext:      cancel,
		cacheDir:           cacheDir,
	}

	return &runner, nil
}

func (chrome *ChromeRunner) Shutdown() error {

	chrome.CancelContext()
	chrome.CancelAllocContext()
	return nil
}

func (chrome *ChromeRunner) GetChromeAnalysis(urlToGet string) (*ChromeAnalysis, error) {

	keyHash := fnv.New64()

	u, err := url.Parse(urlToGet)
	if err != nil {
		return nil, err
	}

	keyHash.Write([]byte(u.Hostname()))

	key := keyHash.Sum64()

	cacheFile := fmt.Sprintf("%s/dom/%d.json", chrome.cacheDir, key)

	existing, err := os.ReadFile(cacheFile)
	if err == nil && existing != nil {

		var existingParsed *ChromeAnalysis

		err = json.Unmarshal(existing, &existingParsed)
		if err != nil {
			return nil, err
		}
		fmt.Printf("chrome cache hit for %s\n", urlToGet)
		return existingParsed, nil

	}

	ctx, cancel := chromedp.NewContext(chrome.Context)
	defer cancel()

	start := time.Now().UnixMilli()

	analysis := ChromeAnalysis{
		BeganAt: start,
		TTI:     1000 * 60,
	}

	chromeRequests := map[network.RequestID]*ChromeRequest{}

	chromedp.ListenTarget(ctx, func(ev interface{}) {
		switch e := ev.(type) {
		case *network.EventLoadingFinished:
			{
				if chromeRequests[e.RequestID] != nil {
					chromeRequests[e.RequestID].Size = e.EncodedDataLength
				} else {
					fmt.Println("request missing", e.RequestID)
				}
			}
		case *network.EventRequestWillBeSent:
			{
				chromeRequests[e.RequestID] = &ChromeRequest{
					Type: e.Type.String(),
					URL:  e.Request.URL,
				}

			}
		}
	})

	cctx, ccancel := context.WithTimeout(ctx, 35*time.Second)
	defer ccancel()

	chromedp.Run(cctx,
		chromedp.Tasks{
			chromedp.EmulateViewport(1440, 900),
			enableLifeCycleEvents(),
			navigateTo(urlToGet),
			waitForEvent("InteractiveTime", &analysis),
			acceptCookies(),
			chromedp.Sleep(1 * time.Second),
			waitForEvent("networkIdle", &analysis),
			chromedp.Sleep(3 * time.Second),
			// captureScreenshot(&screenshot),
			// chromedp.OuterHTML("html", &body),
		},
	)

	analysis.Requests = []*ChromeRequest{}

	for _, req := range chromeRequests {
		if req != nil && strings.HasPrefix(req.URL, "http") {
			analysis.Requests = append(analysis.Requests, req)
		}
	}

	os.Mkdir(fmt.Sprintf("%s/dom", chrome.cacheDir), os.ModePerm)

	cacheWrite, err := json.Marshal(analysis)
	if err != nil {
		return nil, err
	}

	os.WriteFile(cacheFile, cacheWrite, os.ModePerm)

	end := time.Now().UnixMilli()

	fmt.Printf("chrome for %s took %dms\n", urlToGet, end-start)

	return &analysis, nil
}
