package domain

import (
	"context"
	_ "embed"
	"net/url"
	"strings"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

//go:embed data/names.txt
var namesList string

//go:embed data/popular-domains.txt
var popularDomainsList string

var PopularDomainList []string

func init() {
	PopularDomainList = strings.Split(popularDomainsList, "\n")
}

func (engine *Engine) identifyURL(baseUrl string) (bool, bool, bool, bool, error) {

	var urlHumanName, urlNews, urlBlog, popularDomain bool

	baseUrlP, err := url.Parse(baseUrl)
	if err != nil {
		return false, false, false, false, err
	}

	names := strings.Split(namesList, "\n")

	d := baseUrlP.Hostname()

	for _, name := range names {
		if len(name) > 0 && strings.Contains(d, strings.ToLower(name)) {
			urlHumanName = true
			break
		}
	}

	for _, domain := range PopularDomainList {
		if len(domain) > 0 && d == domain {
			popularDomain = true
			break
		}
	}

	urlNews = strings.Contains(strings.ToLower(baseUrlP.Hostname()), "news") || strings.Contains(strings.ToLower(baseUrlP.Hostname()), "daily") || strings.Contains(strings.ToLower(baseUrlP.Hostname()), "standard")

	urlBlog = strings.Contains(strings.ToLower(baseUrlP.Hostname()), "blog") || strings.Contains(strings.ToLower(baseUrlP.Path), "blog")

	return urlHumanName, urlNews, urlBlog, popularDomain, nil
}

func enableLifeCycleEvents() chromedp.ActionFunc {
	return func(ctx context.Context) error {
		err := page.Enable().Do(ctx)
		if err != nil {
			return err
		}

		return nil
	}
}

func navigateTo(url string) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		_, _, _, err := page.Navigate(url).Do(ctx)
		if err != nil {
			return err
		}

		err = page.SetLifecycleEventsEnabled(true).Do(ctx)
		if err != nil {
			return err
		}

		return nil
	}
}

func captureScreenshot(res *[]byte) chromedp.ActionFunc {
	return func(ctx context.Context) error {

		format := page.CaptureScreenshotFormatJpeg

		// capture screenshot
		var err error
		*res, err = page.CaptureScreenshot().
			WithFromSurface(true).
			WithFormat(format).
			WithQuality(int64(75)).
			Do(ctx)
		if err != nil {
			return err
		}
		return nil

	}
}

func acceptCookies() chromedp.EvaluateAction {

	return chromedp.EvaluateAsDevTools(`
		const buttons = Array.from(document.getElementsByTagName('button'));
		

		const button = buttons.filter((b) => b.textContent.includes('accept') | b.textContent.includes('agree'));
		for (const b of button) {
			b.click();
		};

		const iframes = Array.from(document.getElementsByTagName('iframe'));
		for (const iframe of iframes) {
			const buttons = Array.from(iframe.contentWindow.document.getElementsByTagName('button'));
			const button = buttons.filter((b) => b.textContent.includes('accept') || b.textContent.includes('agree'));
			for (const b of button) {
				b.click();
			};
		}

		`, nil)

}

func waitForEvent(eventName string, a *ChromeAnalysis) chromedp.ActionFunc {
	return func(ctx context.Context) error {

		return waitFor(ctx, eventName, a)
	}
}

func waitFor(ctx context.Context, eventName string, a *ChromeAnalysis) error {
	ch := make(chan struct{})
	cctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	chromedp.ListenTarget(cctx, func(ev interface{}) {
		switch e := ev.(type) {

		case *page.EventLifecycleEvent:
			{
				if e.Name == "InteractiveTime" {
					a.TTI = e.Timestamp.Time().UnixMilli() - a.BeganAt
				}
				if e.Name == eventName {
					cancel()
					close(ch)
				}
			}
		}
	})

	select {
	case <-ch:
		return nil
	case <-cctx.Done():
		return nil
	}

}
