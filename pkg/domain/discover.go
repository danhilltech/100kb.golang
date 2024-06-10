package domain

import (
	"fmt"
	"io"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

func (engine *Engine) RunNewFeedSearch(chunkSize int, workers int) error {

	urls, err := engine.getURLsToCrawl()
	if err != nil {
		return err
	}

	rand.Shuffle(len(urls), func(i, j int) { urls[i], urls[j] = urls[j], urls[i] })

	fmt.Printf("Checking %d urls for feeds \n", len(urls))

	jobs := make(chan string, len(urls))
	results := make(chan string, len(urls))

	for w := 1; w <= workers; w++ {
		go engine.crawlURLForFeedWorker(jobs, results)
	}

	for j := 1; j <= len(urls); j++ {
		jobs <- urls[j-1]
	}
	close(jobs)

	t := time.Now().UnixMilli()
	for a := 1; a <= len(urls); a++ {
		feed := <-results

		if feed != "" {
			u, err := url.Parse(feed)

			if err == nil {
				err = engine.Insert(u.Hostname(), feed)
				if err != nil {
					fmt.Println(err)
				}
			}
		}

		if a > 0 && a%chunkSize == 0 {
			diff := time.Now().UnixMilli() - t
			qps := (float64(chunkSize) / float64(diff)) * 1000
			t = time.Now().UnixMilli()
			fmt.Printf("\tdone %d/%d at %0.2f/s\n", a, len(urls), qps)

		}

	}
	fmt.Printf("\tdone %d\n", len(urls))

	return nil
}

func (engine *Engine) crawlURLForFeedWorker(jobs <-chan string, results chan<- string) {
	for id := range jobs {
		feed, err := engine.extractFeed(id)
		if err != nil {
			fmt.Println(err)
		}
		results <- feed

	}
}

// Crawls
func (engine *Engine) extractFeed(candidate string) (string, error) {
	// First check the URL isn't banned

	parsedUrl, err := url.Parse(candidate)
	if err != nil {
		return "", err
	}

	if parsedUrl == nil || parsedUrl.Hostname() == "" {
		return "", nil
	}

	// crawl it
	res, err := engine.httpCrawl.Get(candidate)
	if err != nil || res.StatusCode > 400 {
		return "", err
	}
	if res == nil {
		return "", err
	}
	defer res.Body.Close()

	feed := extractFeedURL(res.Body)

	// Check for malformed
	if strings.HasPrefix(feed, "//") {
		return "", nil
	}

	if feed != "" {
		feedUrl, err := url.Parse(feed)

		if err != nil {
			return "", err
		}
		cleanFeed := parsedUrl.ResolveReference(feedUrl)

		v1 := cleanFeed.String()

		h1, err := engine.httpCrawl.Head(v1)
		if err != nil {
			return "", err
		}

		if h1.StatusCode < 400 {
			return v1, nil
		}

		// possibles := []string{"/feed", "/rss", "/rss.xml", "/blog/feed", "/blog/rss", "/blog/rss.xml"}

		// for _, poss := range possibles {
		// 	u := url.URL{}
		// 	u.Path = poss
		// 	clean := parsedUrl.ResolveReference(&u)

		// 	v := clean.String()
		// 	h, err := engine.httpCrawl.Head(v)
		// 	if err != nil {
		// 		return "", err
		// 	}
		// 	if h.StatusCode < 400 &&
		// 		(strings.Contains(h.Header.Get("Content-Type"), "application/rss+xml") ||
		// 			strings.Contains(h.Header.Get("Content-Type"), "application/atom+xml") ||
		// 			strings.Contains(h.Header.Get("Content-Type"), "text/xml") ||
		// 			strings.Contains(h.Header.Get("Content-Type"), "application/xml")) {
		// 		return v1, nil
		// 	}
		// }

	}

	return "", nil
}

func extractFeedURL(resp io.Reader) string {
	z := html.NewTokenizer(resp)

	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return ""
		case html.StartTagToken, html.SelfClosingTagToken:
			t := z.Token()
			if t.Data == "link" {
				isRSS := false
				url := ""
				for _, attr := range t.Attr {
					if attr.Key == "type" && (attr.Val == "application/rss+xml" || attr.Val == "application/atom+xml") {
						isRSS = true
					}
					if attr.Key == "href" {
						url = attr.Val
					}

				}
				if isRSS {
					return url
				}

			}
		}
	}
}
