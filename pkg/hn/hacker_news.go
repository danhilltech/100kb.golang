package hn

import (
	"database/sql"
	"net/http"
	"time"
)

const HN_BASE = "https://hacker-news.firebaseio.com/v0"

type Engine struct {
	client *http.Client
	db     *sql.DB

	dbInsertPreparedHN *sql.Stmt
}

type HNItemType string

type HNItem struct {
	ID    int
	URL   string
	By    string
	Type  HNItemType
	Time  int
	Score int
}

const (
	HNItemTypeStory = "story"
)

func NewEngine(db *sql.DB) (*Engine, error) {
	tr := &http.Transport{MaxIdleConnsPerHost: 1024, TLSHandshakeTimeout: 0 * time.Second}
	hnClient := &http.Client{Transport: tr}

	engine := Engine{
		client: hnClient,
	}

	err := engine.initDB(db)
	if err != nil {
		return nil, err
	}

	return &engine, nil

}
