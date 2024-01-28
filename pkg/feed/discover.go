package feed

import (
	"fmt"
	"io"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"github.com/danhilltech/100kb.golang/pkg/http"
	"golang.org/x/net/html"
)

type DiscoverWithHTTP struct {
	Feed     string
	Response *http.URLRequest
}

func (engine *Engine) RunNewFeedSearch(chunkSize int, workers int) error {
	txn, err := engine.db.Begin()
	if err != nil {
		return err
	}
	defer txn.Rollback()

	urls, err := engine.getURLsToCrawl(txn)
	if err != nil {
		return err
	}
	err = txn.Commit()
	if err != nil {
		return err
	}

	rand.Shuffle(len(urls), func(i, j int) { urls[i], urls[j] = urls[j], urls[i] })

	fmt.Printf("Checking %d urls for feeds \n", len(urls))

	jobs := make(chan string, len(urls))
	results := make(chan *DiscoverWithHTTP, len(urls))

	txn, err = engine.db.Begin()
	if err != nil {
		return err
	}
	defer txn.Rollback()

	for w := 1; w <= workers; w++ {
		go engine.crawlURLForFeedWorker(jobs, results)
	}

	for j := 1; j <= len(urls); j++ {
		jobs <- urls[j-1]
	}
	close(jobs)

	t := time.Now().UnixMilli()
	for a := 1; a <= len(urls); a++ {
		i := <-results

		if i.Response != nil {
			err = i.Response.Save(txn)
			if err != nil {
				return err
			}
		}

		if i.Feed != "" {
			u, err := url.Parse(i.Feed)

			if err == nil {
				err = engine.Insert(u.Hostname(), i.Feed, txn)
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
			err = txn.Commit()
			if err != nil {
				return err
			}
			txn, err = engine.db.Begin()
			if err != nil {
				return err
			}
		}

	}
	fmt.Printf("\tdone %d\n", len(urls))

	err = txn.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (engine *Engine) crawlURLForFeedWorker(jobs <-chan string, results chan<- *DiscoverWithHTTP) {
	for id := range jobs {
		feed, ht, err := engine.extractFeed(id)
		if err != nil {
			fmt.Println(id, err)
		}
		results <- &DiscoverWithHTTP{Feed: feed, Response: ht}
	}
}

// Crawls
func (engine *Engine) extractFeed(candidate string) (string, *http.URLRequest, error) {
	// First check the URL isn't banned

	urlRequest := &http.URLRequest{
		Url:           candidate,
		Status:        "8xx",
		LastAttemptAt: time.Now().Unix(),
	}

	parsedUrl, err := url.Parse(candidate)
	if err != nil {
		return "", urlRequest, err
	}

	urlRequest.Domain = parsedUrl.Hostname()

	if parsedUrl == nil || parsedUrl.Hostname() == "" {
		return "", urlRequest, nil
	}

	// crawl it
	res, err := engine.httpCrawl.GetWithSafety(candidate)
	if err != nil {
		return "", urlRequest, err
	}
	if res == nil || res.Response == nil {
		return "", urlRequest, err
	}
	defer res.Response.Body.Close()

	feed := extractFeedURL(res.Response.Body)

	// Check for malformed
	if strings.HasPrefix(feed, "//") {
		return "", res, nil
	}

	if feed != "" {
		feedUrl, err := url.Parse(feed)

		if err != nil {
			return "", res, err
		}
		cleanFeed := parsedUrl.ResolveReference(feedUrl)

		v1 := cleanFeed.String()

		h1, err := engine.httpCrawl.Head(v1)
		if err != nil {
			return "", res, err
		}

		if h1.StatusCode < 400 && (strings.Contains(h1.Header.Get("Content-Type"), "application/rss+xml") || strings.Contains(h1.Header.Get("Content-Type"), "application/atom+xml")) {
			return v1, res, nil
		}

		possibles := []string{"/feed", "/rss", "/rss.xml", "/blog/feed", "/blog/rss", "/blog/rss.xml"}

		for _, poss := range possibles {
			u := url.URL{}
			u.Path = poss
			clean := parsedUrl.ResolveReference(&u)

			v := clean.String()
			h, err := engine.httpCrawl.Head(v)
			if err != nil {
				return "", res, err
			}
			if h.StatusCode < 400 && (strings.Contains(h1.Header.Get("Content-Type"), "application/rss+xml") || strings.Contains(h1.Header.Get("Content-Type"), "application/atom+xml")) {
				return v, res, nil
			}
		}

	}

	return "", res, nil
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
