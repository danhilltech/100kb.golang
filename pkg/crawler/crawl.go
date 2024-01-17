package crawler

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	retryhttp "github.com/danhilltech/100kb.golang/pkg/http"
	"golang.org/x/net/html"
)

type Engine struct {
	dbInsertPreparedCrawl *sql.Stmt
	http                  *http.Client
}

type UrlToCrawl struct {
	HackerNewsID int
	Url          string
	Feed         string
}

var BANNED_URLS = []string{
	"youtube.com",
	"www.youtube.com",
	"nytimes.com",
	"www.nytimes.com",
	"en.wikipedia.org",
	"github.com",
	"medium.com",
	"reddit.com",
	"old.reddit.com",
	"arstechnica.com",
	"x.com",
	"twitter.com",
	"theguardian.com",
	"www.theguardian.com",
	"www.theatlantic.com",
	"npr.org",
	"www.nature.com",
	"www.newyorker.com",
	"forbes.com",
	"www.forbes.com",
	"dailymail.co.uk",
	"www.dailymail.co.uk",
	"coindesk.com",
	"www.coindesk.com",
	"mailchi.mp",
	"techcrunch.com",
}

func NewEngine(db *sql.DB) (*Engine, error) {
	engine := Engine{}

	err := engine.initDB(db)
	if err != nil {
		return nil, err
	}

	engine.http = retryhttp.NewRetryableClient()

	return &engine, nil
}

func (engine *Engine) crawlURLForFeedWorker(jobs <-chan *UrlToCrawl, results chan<- *UrlToCrawl) {
	for id := range jobs {
		err := engine.crawlURLForFeed(id)
		if err != nil {
			fmt.Println(id.Url, err)
		}
		results <- id
	}
}

// Crawls
func (engine *Engine) crawlURLForFeed(hnurl *UrlToCrawl) error {
	// First check the URL isn't banned

	parsedUrl, err := url.Parse(hnurl.Url)
	if err != nil {
		return err
	}

	if parsedUrl == nil || parsedUrl.Hostname() == "" {
		return nil
	}

	for _, ban := range BANNED_URLS {
		if ban == parsedUrl.Hostname() {
			return nil
		}
	}

	// crawl it
	resp, err := engine.http.Get(hnurl.Url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	feed := extractFeedURL(resp.Body)

	// Make sure we flush it
	io.Copy(io.Discard, resp.Body)

	// Check for malformed
	if strings.HasPrefix(feed, "//") {
		return nil
	}

	if feed != "" {
		feedUrl, err := url.Parse(feed)

		if err != nil {
			return err
		}
		cleanFeed := parsedUrl.ResolveReference(feedUrl)

		v1 := cleanFeed.String()

		h1, err := engine.http.Head(v1)
		if err != nil {
			return err
		}

		if h1.StatusCode < 400 && strings.Contains(h1.Header.Get("Content-Type"), "application/rss+xml") {
			hnurl.Feed = v1
			return nil
		}

		possibles := []string{"/feed", "/rss", "/rss.xml", "/blog/feed", "/blog/rss", "/blog/rss.xml"}

		for _, poss := range possibles {
			u := url.URL{}
			u.Path = poss
			clean := parsedUrl.ResolveReference(&u)

			v := clean.String()
			h, err := engine.http.Head(v)
			if err != nil {
				return err
			}
			if h.StatusCode < 400 && strings.Contains(h.Header.Get("Content-Type"), "application/rss+xml") {
				hnurl.Feed = v
				return nil
			}
		}

	}

	return nil

}

func (engine *Engine) CrawlURLsForFeeds(hnurl []*UrlToCrawl, workers int) error {

	jobs := make(chan *UrlToCrawl, len(hnurl))
	results := make(chan *UrlToCrawl, len(hnurl))

	for w := 1; w <= workers; w++ {
		go engine.crawlURLForFeedWorker(jobs, results)
	}

	for j := 1; j <= len(hnurl); j++ {
		jobs <- hnurl[j-1]
	}
	close(jobs)

	items := make([]*UrlToCrawl, len(hnurl))

	for a := 1; a <= len(hnurl); a++ {
		b := <-results
		items[a-1] = b
	}

	return nil

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
					if attr.Key == "type" && attr.Val == "application/rss+xml" {
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
