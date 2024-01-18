package feed

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/danhilltech/100kb.golang/pkg/article"
	"github.com/danhilltech/100kb.golang/pkg/crawler"
	retryhttp "github.com/danhilltech/100kb.golang/pkg/http"
	"github.com/mmcdole/gofeed"
)

type Engine struct {
	dbInsertPreparedFeed *sql.Stmt
	dbUpdatePreparedFeed *sql.Stmt
	db                   *sql.DB
	crawlEngine          *crawler.Engine
	articleEngine        *article.Engine
	http                 *http.Client
}

type Feed struct {
	Url         string
	LastFetchAt int64
	Title       string
	Description string
	Language    string

	Articles []article.Article
}

func NewEngine(db *sql.DB, crawlEngine *crawler.Engine, articleEngine *article.Engine) (*Engine, error) {
	engine := Engine{crawlEngine: crawlEngine, articleEngine: articleEngine}

	engine.initDB(db)

	engine.http = retryhttp.NewRetryableClient()

	return &engine, nil
}

func (engine *Engine) feedRefreshWorker(jobs <-chan *Feed, results chan<- *Feed) {
	for id := range jobs {
		err := engine.feedRefresh(id)
		if err != nil {
			fmt.Println(id.Url, err)
		}
		results <- id
	}
}

// Crawls
func (engine *Engine) feedRefresh(feed *Feed) error {
	feed.LastFetchAt = time.Now().Unix()
	// First check the URL isn't banned

	// crawl it
	resp, err := engine.http.Get(feed.Url)
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

	feed.Articles = make([]article.Article, len(rss.Items))

	for i, item := range rss.Items {
		feed.Articles[i] = article.Article{
			Url: item.Link,
		}
		if item.PublishedParsed != nil {
			feed.Articles[i].PublishedAt = item.PublishedParsed.Unix()
		}
	}

	return nil

}

func (engine *Engine) refresh(hnurl []*Feed, workers int) error {
	jobs := make(chan *Feed, len(hnurl))
	results := make(chan *Feed, len(hnurl))

	for w := 1; w <= workers; w++ {
		go engine.feedRefreshWorker(jobs, results)
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
