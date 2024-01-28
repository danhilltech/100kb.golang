package feed

import (
	"fmt"
	"time"

	"github.com/danhilltech/100kb.golang/pkg/article"
	"github.com/danhilltech/100kb.golang/pkg/http"
	"github.com/mmcdole/gofeed"
)

type DomainWithHTTP struct {
	Domain   *Domain
	Response *http.URLRequest
}

func (engine *Engine) RunFeedRefresh(chunkSize int, workers int) error {
	txn, err := engine.db.Begin()
	if err != nil {
		return err
	}
	defer txn.Rollback()
	feeds, err := engine.getDomainsToRefresh(txn)
	if err != nil {
		return err
	}
	err = txn.Commit()
	if err != nil {
		return err
	}

	fmt.Printf("Checking %d feeds for new links\n", len(feeds))

	jobs := make(chan *Domain, len(feeds))
	results := make(chan *DomainWithHTTP, len(feeds))

	txn, err = engine.db.Begin()
	if err != nil {
		return err
	}
	defer txn.Rollback()

	for w := 1; w <= workers; w++ {
		go engine.feedRefreshWorker(jobs, results)
	}

	for j := 1; j <= len(feeds); j++ {
		jobs <- feeds[j-1]
	}
	close(jobs)

	t := time.Now().UnixMilli()

	for a := 1; a <= len(feeds); a++ {
		i := <-results
		if i.Response != nil {
			err = i.Response.Save(txn)
			if err != nil {
				return err
			}
		}

		err = engine.Update(i.Domain, txn)
		if err != nil {
			return err
		}

		for _, article := range i.Domain.Articles {

			err = engine.articleEngine.Insert(&article, i.Domain.FeedURL, i.Domain.Domain, txn)
			if err != nil {
				return err
			}

		}

		if a > 0 && a%chunkSize == 0 {
			diff := time.Now().UnixMilli() - t
			qps := (float64(chunkSize) / float64(diff)) * 1000
			t = time.Now().UnixMilli()
			fmt.Printf("\tdone %d/%d at %0.2f/s\n", a, len(feeds), qps)
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
	fmt.Printf("\tdone %d\n", len(feeds))

	err = txn.Commit()
	if err != nil {
		return err
	}

	return nil
}

// Crawls
func (engine *Engine) feedRefresh(feed *Domain) (*http.URLRequest, error) {
	feed.LastFetchAt = time.Now().Unix()
	// First check the URL isn't banned

	// crawl it
	resp, err := engine.httpCrawl.GetWithSafety(feed.FeedURL)
	if err != nil || resp == nil {
		return resp, nil
	}
	if resp.Response == nil {
		return resp, nil
	}

	defer resp.Response.Body.Close()

	fp := gofeed.NewParser()
	rss, err := fp.Parse(resp.Response.Body)
	if err != nil {
		return resp, err
	}

	feed.Language = rss.Language

	feed.Articles = make([]article.Article, len(rss.Items))

	for i, item := range rss.Items {
		feed.Articles[i] = article.Article{
			Url:    item.Link,
			Domain: feed.Domain,
		}
		if item.PublishedParsed != nil {
			feed.Articles[i].PublishedAt = item.PublishedParsed.Unix()
		}
	}

	return resp, nil

}

func (engine *Engine) feedRefreshWorker(jobs <-chan *Domain, results chan<- *DomainWithHTTP) {
	for id := range jobs {
		resp, err := engine.feedRefresh(id)
		if err != nil {
			fmt.Println(id.Domain, err)
		}
		results <- &DomainWithHTTP{Domain: id, Response: resp}
	}
}
