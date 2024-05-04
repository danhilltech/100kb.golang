package domain

import (
	"database/sql"

	"github.com/danhilltech/100kb.golang/pkg/article"
	retryhttp "github.com/danhilltech/100kb.golang/pkg/http"
	"github.com/danhilltech/100kb.golang/pkg/parsing"
	"github.com/smira/go-statsd"
)

type Engine struct {
	dbInsertPreparedFeed *sql.Stmt
	dbUpdatePreparedFeed *sql.Stmt
	db                   *sql.DB
	articleEngine        *article.Engine
	httpCrawl            *retryhttp.Client
	parser               *parsing.Engine
}

type Domain struct {
	Domain      string
	FeedURL     string
	LastFetchAt int64
	FeedTitle   string
	Language    string

	PageAbout    bool
	PageBlogRoll bool
	PageWriting  bool
	PageNow      bool

	URLNews      bool
	URLBlog      bool
	URLHumanName bool

	DomainIsPopular bool
	DomainTLD       string

	Articles []*article.Article
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

	engine.httpCrawl, err = retryhttp.NewClient(cacheDir, db, sd)
	if err != nil {
		return nil, err
	}

	engine.parser, err = parsing.NewEngine()
	if err != nil {
		return nil, err
	}

	return &engine, nil
}
