package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

type HNUrlToCrawl struct {
	ID   int
	Url  string
	Feed string
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
}

func crawlURLForFeedWorker(jobs <-chan *HNUrlToCrawl, results chan<- *HNUrlToCrawl) {
	for id := range jobs {
		err := crawlURLForFeed(id)
		if err != nil {
			fmt.Println(err)
		}
		results <- id
	}
}

// Crawls
func crawlURLForFeed(hnurl *HNUrlToCrawl) error {
	// First check the URL isn't banned

	parsedUrl, err := url.Parse(hnurl.Url)
	if err != nil {
		return nil
	}

	for _, ban := range BANNED_URLS {
		if ban == parsedUrl.Hostname() {
			return nil
		}
	}

	// crawl it
	resp, err := http.Get(hnurl.Url)
	if err != nil {
		return nil
	}

	defer resp.Body.Close()

	feed := extractFeedURL(resp.Body)

	// Check for malformed
	if strings.HasPrefix(feed, "//") {
		return nil
	}

	if feed != "" {
		feedUrl, err := url.Parse(feed)

		if err != nil {
			return nil
		}
		cleanFeed := parsedUrl.ResolveReference(feedUrl)

		hnurl.Feed = cleanFeed.String()
	}

	return nil

}

func crawlURLsForFeeds(hnurl []*HNUrlToCrawl) error {

	workers := 10
	jobs := make(chan *HNUrlToCrawl, len(hnurl))
	results := make(chan *HNUrlToCrawl, len(hnurl))

	for w := 1; w <= workers; w++ {
		go crawlURLForFeedWorker(jobs, results)
	}

	for j := 1; j <= len(hnurl); j++ {
		jobs <- hnurl[j-1]
	}
	close(jobs)

	items := make([]*HNUrlToCrawl, len(hnurl))

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
