package domain

import (
	"database/sql"

	"github.com/danhilltech/100kb.golang/pkg/article"
	retryhttp "github.com/danhilltech/100kb.golang/pkg/http"
	"github.com/smira/go-statsd"
)

const REFRESH_AGO_SECONDS = 60 * 60 * 24 * 7
const REFRESH_LIMIT = 5000

const VALIDATE_AGO_SECONDS = 60 * 60 * 24 * 21
const VALIDATE_LIMIT = 5000

type Engine struct {
	dbInsertPreparedFeed *sql.Stmt
	dbUpdatePreparedFeed *sql.Stmt
	db                   *sql.DB
	articleEngine        *article.Engine
	httpCrawl            *retryhttp.Client

	chrome *ChromeRunner
}

type Domain struct {
	Domain         string
	FeedURL        string
	LastFetchAt    int64
	LastValidateAt int64
	FeedTitle      string
	Language       string

	URLNews      bool
	URLBlog      bool
	URLHumanName bool

	DomainIsPopular bool
	DomainTLD       string

	ChromeAnalysis *ChromeAnalysis

	Platform string

	Articles []*article.Article

	// Only used at runtime/output
	LiveScore            float64
	LiveLatestArticleURL string
}

func NewEngine(db *sql.DB, articleEngine *article.Engine, sd *statsd.Client, cacheDir string) (*Engine, error) {
	engine := Engine{articleEngine: articleEngine}

	err := engine.initDB(db)
	if err != nil {
		return nil, err
	}

	// tr := &http.Transport{MaxIdleConnsPerHost: 1024, TLSHandshakeTimeout: 0 * time.Second}
	// hnClient := &http.Client{Transport: tr}

	// engine.http = hnClient

	engine.httpCrawl, err = retryhttp.NewClient(cacheDir, sd)
	if err != nil {
		return nil, err
	}

	engine.chrome, err = startChrome(cacheDir)
	if err != nil {
		return nil, err
	}

	return &engine, nil
}

func (d *Engine) Close() error {
	return d.chrome.Shutdown()

}
