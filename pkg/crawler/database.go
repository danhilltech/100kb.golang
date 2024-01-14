package crawler

import (
	"database/sql"
	"time"
)

func (engine *Engine) initDB(db *sql.DB) error {
	var err error

	engine.dbInsertPreparedCrawl, err = db.Prepare("INSERT INTO crawls(url, hackerNewsId, lastCrawlAt) VALUES(?, ?, ?)  ON CONFLICT(url) DO UPDATE SET lastCrawlAt = excluded.lastCrawlAt;")
	if err != nil {
		return err
	}
	return nil
}

func (engine *Engine) AddURL(item *UrlToCrawl, txchan *sql.Tx) error {
	_, err := txchan.Stmt(engine.dbInsertPreparedCrawl).Exec(item.Url, item.HackerNewsID, time.Now().Unix())
	return err
}

func (engine *Engine) GetExistingIDs(txchan *sql.Tx) ([]int, error) {
	res, err := txchan.Query("SELECT id FROM hacker_news")
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var ids []int

	for res.Next() {
		var id int
		err = res.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := res.Err(); err != nil {
		return nil, err
	}
	return ids, nil
}

func (engine *Engine) GetURLsToCrawl(txchan *sql.Tx) ([]*UrlToCrawl, error) {
	res, err := txchan.Query("SELECT id, hacker_news.url FROM hacker_news LEFT JOIN crawls on crawls.hackerNewsId = hacker_news.id WHERE hacker_news.url IS NOT NULL AND crawls.lastCrawlAt IS NULL;")
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var urls []*UrlToCrawl

	for res.Next() {
		var id int
		var url string
		err = res.Scan(&id, &url)
		if err != nil {
			return nil, err
		}

		urls = append(urls, &UrlToCrawl{HackerNewsID: id, Url: url})
	}
	if err := res.Err(); err != nil {
		return nil, err
	}
	return urls, nil
}
