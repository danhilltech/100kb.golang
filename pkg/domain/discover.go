package domain

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
)

func (engine *Engine) RunNewFeedSearch(ctx context.Context, chunkSize int, workers int) error {

	urls, err := engine.getURLsToCrawl()
	if err != nil {
		return err
	}

	rand.Shuffle(len(urls), func(i, j int) { urls[i], urls[j] = urls[j], urls[i] })

	engine.log.Printf("Checking %d urls for feeds \n", len(urls))

	jobs := make(chan string, len(urls))
	results := make(chan string, len(urls))

	for w := 1; w <= workers; w++ {
		go engine.crawlURLForFeedWorker(jobs, results)
	}

	for j := 1; j <= len(urls); j++ {
		jobs <- urls[j-1]
	}
	close(jobs)

	txn, _ := engine.db.Begin()
	defer txn.Rollback()

	t := time.Now().UnixMilli()
	for a := 1; a <= len(urls); a++ {
		select {
		case <-ctx.Done():
			txn.Commit()
			return ctx.Err()
		case feed := <-results:

			if feed != "" {
				u, err := url.Parse(feed)

				if err == nil {
					err = engine.Insert(txn, u.Hostname(), feed)
					if err != nil {
						engine.log.Println(err)
					}
				}
			}

			if a > 0 && a%chunkSize == 0 {
				err := txn.Commit()
				if err != nil {
					return err
				}
				txn, _ = engine.db.Begin()
				diff := time.Now().UnixMilli() - t
				qps := (float64(chunkSize) / float64(diff)) * 1000
				t = time.Now().UnixMilli()
				engine.log.Printf("\tdone %d/%d at %0.2f/s\n", a, len(urls), qps)

			}

		}
	}

	txn.Commit()
	engine.log.Printf("\tdone %d\n", len(urls))

	return nil
}

func (engine *Engine) crawlURLForFeedWorker(jobs <-chan string, results chan<- string) {
	for id := range jobs {
		feed, err := engine.extractFeed(id)
		if err != nil {
			engine.log.Println(id, err)
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
					if attr.Key == "title" && strings.Contains(strings.ToLower(attr.Val), "comments") {
						isRSS = false
					}

				}
				if isRSS {
					return url
				}

			}
		}
	}
}

func (engine *Engine) RunKagiList(ctx context.Context) error {

	engine.log.Println("Getting Kagi list...")

	resp, err := http.Get("https://raw.githubusercontent.com/kagisearch/smallweb/main/smallweb.txt")
	// handle the error if there is one
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	d := 0

	txn, _ := engine.db.Begin()
	defer txn.Rollback()

	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			txn.Commit()
			return ctx.Err()
		default:

			feed := scanner.Text()

			u, err := url.Parse(feed)

			if err == nil {
				err = engine.Insert(txn, u.Hostname(), feed)
				if err != nil {
					engine.log.Println(err)
				}
			}
			d++
		}
	}

	engine.log.Printf("\tdone %d\n", d)
	txn.Commit()

	return nil
}

func (engine *Engine) RunBearBlog(ctx context.Context) error {

	engine.log.Println("Getting BearBlog list...")

	for i := 0; i < 10; i++ {
		err := engine.runBearInner(ctx, i)
		if err != nil {
			return err
		}
	}
	return nil

}

func (engine *Engine) runBearInner(ctx context.Context, page int) error {
	resp, err := http.Get(fmt.Sprintf("https://bearblog.dev/discover/?page=%d", page))
	// handle the error if there is one
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	d := 0

	txn, _ := engine.db.Begin()
	defer txn.Rollback()

	z, err := html.Parse(resp.Body)
	if err != nil {
		return err
	}

	sel, err := cascadia.Parse("span > a")
	if err != nil {
		return err
	}

	for _, a := range cascadia.QueryAll(z, sel) {
		select {
		case <-ctx.Done():
			txn.Commit()
			return ctx.Err()
		default:

			for _, at := range a.Attr {
				if at.Key == "href" {

					err = engine.Insert(txn, strings.ReplaceAll(at.Val, "//", ""), fmt.Sprintf("http:%s/feed", at.Val))
					if err != nil {
						engine.log.Println(err)
					}
				}
			}

			// u, err := url.Parse(feed)

			// if err == nil {
			// 	err = engine.Insert(txn, u.Hostname(), feed)
			// 	if err != nil {
			// 		engine.log.Println(err)
			// 	}
			// }
			d++
		}
	}

	engine.log.Printf("\tdone %d\n", d)
	txn.Commit()

	return nil
}
