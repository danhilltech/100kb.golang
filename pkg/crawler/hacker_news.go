package crawler

import (
	"database/sql"
	"log"
	"net/http"
	"time"
)

const HN_BASE = "https://hacker-news.firebaseio.com/v0"

type Engine struct {
	client *http.Client
	db     *sql.DB
	log    *log.Logger

	dbInsertPreparedToCrawl *sql.Stmt
}

type HNItemType string

type ToCrawl struct {
	URL string

	HNID   int `json:"id"`
	Domain string
	By     string
	Type   HNItemType
	Time   int
	Score  int

	Text string
}

func NewEngine(log *log.Logger, db *sql.DB) (*Engine, error) {
	tr := &http.Transport{MaxIdleConnsPerHost: 1024, TLSHandshakeTimeout: 0 * time.Second}
	hnClient := &http.Client{Transport: tr}

	engine := Engine{
		client: hnClient,
		log:    log,
	}

	err := engine.initDB(db)
	if err != nil {
		return nil, err
	}

	return &engine, nil

}
