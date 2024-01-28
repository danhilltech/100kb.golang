package feed

import (
	"database/sql"
	"net/url"

	"github.com/danhilltech/100kb.golang/pkg/http"
	"github.com/danhilltech/100kb.golang/pkg/utils"
)

func (engine *Engine) initDB(db *sql.DB) error {
	var err error
	engine.dbInsertPreparedFeed, err = db.Prepare("INSERT INTO domains(domain, feedUrl) VALUES(?, ?) ON CONFLICT(domain) DO NOTHING;")
	if err != nil {
		return err
	}

	engine.dbUpdatePreparedFeed, err = db.Prepare("UPDATE domains SET lastFetchAt = ?, feedTitle = ?, language = ? WHERE domain = ?;")
	if err != nil {
		return err
	}

	engine.db = db

	return nil
}

func (engine *Engine) Insert(domain string, feedurl string, txchan *sql.Tx) error {
	_, err := txchan.Stmt(engine.dbInsertPreparedFeed).Exec(domain, feedurl)
	return err
}

func (engine *Engine) Update(feed *Domain, txchan *sql.Tx) error {
	_, err := txchan.Stmt(engine.dbUpdatePreparedFeed).Exec(
		utils.NullInt64(feed.LastFetchAt),
		utils.NullString(feed.FeedTitle),
		utils.NullString(feed.Language),
		feed.Domain,
	)
	return err
}

func (engine *Engine) getDomainsToRefresh(txchan *sql.Tx) ([]*Domain, error) {
	res, err := txchan.Query("SELECT domain, feedUrl, lastFetchAt, feedTitle, language FROM domains WHERE lastFetchAt IS NULL;")
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var urls []*Domain

	for res.Next() {
		var feedUrl string
		var domain string
		var lastFetchAt sql.NullInt64
		var feedTitle, language sql.NullString
		err = res.Scan(&domain, &feedUrl, &lastFetchAt, &feedTitle, &language)
		if err != nil {
			return nil, err
		}

		urls = append(urls, &Domain{Domain: domain, FeedURL: feedUrl, LastFetchAt: lastFetchAt.Int64, FeedTitle: feedTitle.String, Language: language.String})
	}
	if err := res.Err(); err != nil {
		return nil, err
	}
	return urls, nil
}

func (engine *Engine) getURLsToCrawl(txchan *sql.Tx) ([]string, error) {
	missingFromHN, err := txchan.Query(`SELECT url
	FROM (
		SELECT h.url, 
		ROW_NUMBER() OVER (PARTITION BY h.domain ORDER BY addedAt DESC) AS rn, 
		length(h.url) - LENGTH(REPLACE(h.url, "/", "")) as scnt
		FROM hacker_news h
		LEFT JOIN domains d ON d.domain = h.domain
		LEFT JOIN url_requests u on u.domain = h.domain
		WHERE h.url IS NOT NULL AND d.domain IS NULL AND u.domain IS NULL
	) raw 
	
	WHERE rn <= 2 
	AND scnt >= 4;`)
	if err != nil {
		return nil, err
	}
	defer missingFromHN.Close()

	var urls []string

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

		// // TODO check domain isn't in the feeds already

		domain := u.Hostname()

		isBad := false
		//TODO improve this
		for _, bad := range http.BANNED_URLS {
			if bad == domain {
				isBad = true
			}
		}

		if !isBad {
			urls = append(urls, urlScanned)
		}

	}
	if err := missingFromHN.Err(); err != nil {
		return nil, err
	}

	return urls, nil
}
