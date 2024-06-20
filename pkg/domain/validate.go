package domain

import (
	"context"
	"fmt"
	"strings"
	"time"

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

	printSize := 100

	a := 0
	t := time.Now().UnixMilli()
	txn, _ := engine.db.Begin()
	for _, domain := range domains {

		err = engine.validateDomain(domain)
		if err != nil {
			fmt.Println(domain.Domain, err)
			continue
		}

		err = engine.Update(txn, domain)
		if err != nil {
			fmt.Println(domain.Domain, err)
			continue
		}

		if a > 0 && a%printSize == 0 {
			err := txn.Commit()
			if err != nil {
				return err
			}
			txn, _ = engine.db.Begin()
			diff := time.Now().UnixMilli() - t
			qps := (float64(printSize) / float64(diff)) * 1000
			t = time.Now().UnixMilli()
			fmt.Printf("\tdone %d/%d at %0.2f/s\n", a, len(domains), qps)

		}
		a++
	}
	err = txn.Commit()
	if err != nil {
		return err
	}
	fmt.Printf("\tdone %d/%d\n\n", a, len(domains))

	return nil
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

	// Now load it in a chrome window
	url, err := engine.getLatestArticleURL(domain)
	if err != nil {
		return err
	}

	body, err := engine.chrome.GetDomFromChrone(url)
	if err != nil {
		return err
	}

	if strings.Contains(body, "google_ads_") {
		domain.DomainGoogleAds = true
	}

	return nil

}

func startChrome() (*ChromeRunner, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.UserDataDir(".chrome"),
		chromedp.Flag("headless", "new"),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-first-run", true),
		chromedp.Flag("--blink-settings", "imagesEnabled=false"),
	)

	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)

	var contextOpts []chromedp.ContextOption

	// opts = append(opts, chromedp.WithDebugf(log.Printf))

	// opts = append(opts, chromedp.WithBrowserOption())

	// contextOpts = append(contextOpts, )

	ctx, cancel := chromedp.NewContext(allocCtx, contextOpts...)

	runner := ChromeRunner{
		AllocContext:       allocCtx,
		Context:            ctx,
		CancelAllocContext: allocCancel,
		CancelContext:      cancel,
	}

	return &runner, nil
}

func (chrome *ChromeRunner) Shutdown() error {

	chrome.CancelContext()
	chrome.CancelAllocContext()
	return nil
}

func (chrome *ChromeRunner) GetDomFromChrone(url string) (string, error) {
	ctx, cancel := chromedp.NewContext(chrome.Context)
	defer cancel()

	// create a timeout
	ctx, cancelTimeout := context.WithTimeout(ctx, time.Second*30)
	defer cancelTimeout()

	var body string

	if err := chromedp.Run(ctx,
		chromedp.Tasks{
			chromedp.Emulate(device.IPhone13),
			navigateAndWaitFor(url, "InteractiveTime"),
			chromedp.OuterHTML("html", &body),
		},
	); err != nil {
		return "", err
	}

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
	cctx, cancel := context.WithCancel(ctx)
	chromedp.ListenTarget(cctx, func(ev interface{}) {
		switch e := ev.(type) {
		case *page.EventLifecycleEvent:
			if e.Name == eventName {
				cancel()
				close(ch)
			}
		}
	})
	select {
	case <-ch:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}

}
