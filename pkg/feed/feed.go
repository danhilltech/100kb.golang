package feed

import (
	"database/sql"

	"github.com/danhilltech/100kb.golang/pkg/article"
	retryhttp "github.com/danhilltech/100kb.golang/pkg/http"
)

type Engine struct {
	dbInsertPreparedFeed *sql.Stmt
	dbUpdatePreparedFeed *sql.Stmt
	db                   *sql.DB
	articleEngine        *article.Engine
	httpCrawl            *retryhttp.Client
}

type Domain struct {
	Domain      string
	FeedURL     string
	LastFetchAt int64
	FeedTitle   string
	Language    string

	Articles []article.Article
}

func NewEngine(db *sql.DB, articleEngine *article.Engine, cacheDir string) (*Engine, error) {
	engine := Engine{articleEngine: articleEngine}

	err := engine.initDB(db)
	if err != nil {
		return nil, err
	}

	// tr := &http.Transport{MaxIdleConnsPerHost: 1024, TLSHandshakeTimeout: 0 * time.Second}
	// hnClient := &http.Client{Transport: tr}

	// engine.http = hnClient

	engine.httpCrawl, err = retryhttp.NewClient(cacheDir)
	if err != nil {
		return nil, err
	}

	return &engine, nil
}
