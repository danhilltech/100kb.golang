package crawler

import (
	"database/sql"
	"net/url"
	"time"
)

func (engine *Engine) initDB(db *sql.DB) error {
	var err error

	engine.dbInsertPreparedCandidate, err = db.Prepare("INSERT INTO candidate_urls(url, domain, addedAt) VALUES(?, ?, ?) ON CONFLICT(url) DO NOTHING;")
	if err != nil {
		return err
	}

	engine.dbUpdatePreparedCandidate, err = db.Prepare("UPDATE candidate_urls SET lastCrawlAt = ? WHERE url = ?;")
	if err != nil {
		return err
	}
	return nil
}

func (engine *Engine) AddCandidateURL(u string, domain string, txchan *sql.Tx) error {
	_, err := txchan.Stmt(engine.dbInsertPreparedCandidate).Exec(u, domain, time.Now().Unix())
	return err
}

func (engine *Engine) TrackCandidateCrawl(item *UrlToCrawl, txchan *sql.Tx) error {
	_, err := txchan.Stmt(engine.dbUpdatePreparedCandidate).Exec(time.Now().Unix(), item.Url)
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
	missingFromHN, err := txchan.Query("SELECT DISTINCT hacker_news.url url FROM hacker_news LEFT JOIN candidate_urls c ON c.url = hacker_news.url WHERE c.url IS NULL AND hacker_news.url IS NOT NULL;")
	if err != nil {
		return nil, err
	}
	defer missingFromHN.Close()
	for missingFromHN.Next() {
		var urlScanned string
		err = missingFromHN.Scan(&urlScanned)
		if err != nil {
			return nil, err
		}
		u, err := url.Parse(urlScanned)
		if err != nil {
			continue
		}

		domain := u.Hostname()

		err = engine.AddCandidateURL(urlScanned, domain, txchan)
		if err != nil {
			return nil, err
		}

	}
	if err := missingFromHN.Err(); err != nil {
		return nil, err
	}

	res, err := txchan.Query(`SELECT url FROM (
		SELECT url, 
		ROW_NUMBER() OVER (PARTITION BY domain ORDER BY addedAt DESC) AS rn, 
		length(url) - LENGTH(REPLACE(url, "/", "")) as scnt
		FROM candidate_urls WHERE lastCrawlAt IS NULL AND domain IS NOT NULL
	) inn
	WHERE rn <= 2 
	AND scnt >= 4;`)

	if err != nil {
		return nil, err
	}
	defer res.Close()

	var urls []*UrlToCrawl

	for res.Next() {
		var url string
		err = res.Scan(&url)
		if err != nil {
			return nil, err
		}

		urls = append(urls, &UrlToCrawl{Url: url})
	}
	if err := res.Err(); err != nil {
		return nil, err
	}
	return urls, nil
}
