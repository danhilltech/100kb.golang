package domain

import (
	"context"
	"database/sql"
	"fmt"
	"hash/fnv"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
)

type ChromeRunner struct {
	AllocContext       context.Context
	Context            context.Context
	CancelAllocContext context.CancelFunc
	CancelContext      context.CancelFunc

	cacheDir string
}

func (engine *Engine) RunDomainValidate(chunkSize int) error {

	domains, err := engine.getDomainsToValidate()
	if err != nil {
		return err
	}

	fmt.Printf("Validating %d domains\n", len(domains))

	for _, domain := range domains {
		// Now load it in a chrome window
		url, err := engine.getLatestArticleURL(domain)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		domain.LiveLatestArticleURL = url

	}

	printSize := 100

	jobs := make(chan *Domain, len(domains))
	results := make(chan *Domain, len(domains))

	workers := runtime.NumCPU() * 2

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
	for a := 0; a < len(domains); a++ {
		domain := <-results

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

	err = txn.Commit()
	if err != nil {
		return err
	}
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
		body, err := engine.chrome.GetDomFromChrone(domain.LiveLatestArticleURL)
		if err != nil {
			return err
		}

		if strings.Contains(body, "google_ads_") {
			domain.DomainGoogleAds = true
		}
	}

	return nil

}

func startChrome(cacheDir string) (*ChromeRunner, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.UserDataDir(".chrome"),
		chromedp.Flag("headless", "new"),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-first-run", true),
		chromedp.Flag("--blink-settings", "imagesEnabled=false"),
		chromedp.Flag("--disable-gl-drawing-for-tests", true),
		chromedp.Flag("--disable-gl-drawing-for-tests", true),
		chromedp.Flag("--hide-scrollbars", true),
		chromedp.Flag("--mute-audio", true),
		chromedp.Flag("--no-sandbox", true),
		chromedp.Flag("--disable-setuid-sandbox", true),
		chromedp.Flag("--disable-translate", true),
		chromedp.Flag("--disable-extensions", true),
	)

	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)

	var contextOpts []chromedp.ContextOption

	// opts = append(opts, chromedp.WithDebugf(log.Printf))

	// opts = append(opts, chromedp.WithBrowserOption())

	// contextOpts = append(contextOpts, )

	ctx, cancel := chromedp.NewContext(allocCtx, contextOpts...)

	err := chromedp.Run(ctx)
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

func (chrome *ChromeRunner) GetDomFromChrone(urlToGet string) (string, error) {

	keyHash := fnv.New64()

	u, err := url.Parse(urlToGet)
	if err != nil {
		return "", err
	}

	keyHash.Write([]byte(u.Hostname()))

	key := keyHash.Sum64()

	cacheFile := fmt.Sprintf("%s/dom/%d.txt", chrome.cacheDir, key)

	fmt.Println(u, cacheFile)

	existing, err := os.ReadFile(cacheFile)
	if err == nil && existing != nil {
		return string(existing), nil
	}

	ctx, cancel := chromedp.NewContext(chrome.Context)
	defer cancel()

	// create a timeout
	// ctx, cancelTimeout := context.WithTimeout(ctx, time.Second*10)
	// defer cancelTimeout()

	var body string

	if err := chromedp.Run(ctx,
		chromedp.Tasks{
			chromedp.Emulate(device.IPhone13),
			navigateAndWaitFor(urlToGet, "InteractiveTime"),
			chromedp.OuterHTML("html", &body),
		},
	); err != nil {
		return body, err
	}

	os.Mkdir(fmt.Sprintf("%s/dom", chrome.cacheDir), os.ModePerm)
	os.WriteFile(cacheFile, []byte(body), os.ModePerm)

	return body, nil
}

func navigateAndWaitFor(url string, eventName string) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		_, _, _, err := page.Navigate(url).Do(ctx)
		if err != nil {
			return err
		}

		return waitFor(ctx, eventName)
	}
}

func waitFor(ctx context.Context, eventName string) error {
	ch := make(chan struct{})
	cctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	chromedp.ListenTarget(cctx, func(ev interface{}) {
		switch e := ev.(type) {
		case *page.EventLifecycleEvent:
			{
				if e.Name == eventName {
					cancel()
					close(ch)
				}
			}
		case *network.EventRequestWillBeSent:
			{
				if strings.Contains(e.Request.URL, "google") {
					fmt.Println(e.Request.URL)
				}
			}
		}
	})

	select {
	case <-ch:
		return nil
	case <-cctx.Done():
		fmt.Println("done")

		return nil
	}

}
