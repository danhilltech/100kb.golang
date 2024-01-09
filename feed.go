package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/mmcdole/gofeed"
)

type Feed struct {
	Url         string
	LastFetchAt int64
	Title       string
	Description string
	Language    string

	Articles []Article
}

func feedRefreshWorker(jobs <-chan *Feed, results chan<- *Feed) {
	for id := range jobs {
		err := feedRefresh(id)
		if err != nil {
			fmt.Println(err)
		}
		results <- id
	}
}

// Crawls
func feedRefresh(feed *Feed) error {
	// First check the URL isn't banned

	// crawl it
	resp, err := http.Get(feed.Url)
	if err != nil {
		return nil
	}

	defer resp.Body.Close()

	fp := gofeed.NewParser()
	rss, err := fp.Parse(resp.Body)
	if err != nil {
		return err
	}

	feed.Description = rss.Description
	feed.Language = rss.Language
	feed.Title = rss.Title

	feed.Articles = make([]Article, len(rss.Items))

	feed.LastFetchAt = time.Now().Unix()

	for i, item := range rss.Items {
		feed.Articles[i] = Article{
			Url:         item.Link,
			PublishedAt: item.PublishedParsed.Unix(),
		}
	}

	return nil

}

func feedRefreshes(hnurl []*Feed) error {

	workers := 10
	jobs := make(chan *Feed, len(hnurl))
	results := make(chan *Feed, len(hnurl))

	for w := 1; w <= workers; w++ {
		go feedRefreshWorker(jobs, results)
	}

	for j := 1; j <= len(hnurl); j++ {
		jobs <- hnurl[j-1]
	}
	close(jobs)

	items := make([]*Feed, len(hnurl))

	for a := 1; a <= len(hnurl); a++ {
		b := <-results
		items[a-1] = b
	}

	return nil

}
