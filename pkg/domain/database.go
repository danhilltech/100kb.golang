package domain

import (
	"database/sql"
	"fmt"
	"net/url"
	"sync"

	"github.com/danhilltech/100kb.golang/pkg/http"
	"github.com/danhilltech/100kb.golang/pkg/utils"
)

const DOMAIN_SELECT = `
domains.domain, 
domains.feedUrl, 
domains.lastFetchAt, 
feedTitle, 
language,
pageAbout,
pageBlogRoll,
pageWriting,
pageNow,
urlNews,
urlBlog,
urlHumanName,
domainIsPopular,
domainTLD,
platform
`

var mu sync.Mutex

func (engine *Engine) initDB(db *sql.DB) error {
	var err error
	engine.dbInsertPreparedFeed, err = db.Prepare("INSERT INTO domains(domain, feedUrl) VALUES(?, ?) ON CONFLICT(domain) DO NOTHING;")
	if err != nil {
		return err
	}

	engine.dbUpdatePreparedFeed, err = db.Prepare(`UPDATE 
	domains 
	SET 
	lastFetchAt = ?, 
	feedTitle = ?, 
	language = ?,
	pageAbout = ?,
	pageBlogRoll = ?,
	pageWriting = ?,
	pageNow = ?,
	urlNews = ?,
	urlBlog = ?,
	urlHumanName = ?,
	domainIsPopular = ?,
	domainTLD = ?,
	platform = ?
	
	WHERE domain = ?;`)
	if err != nil {
		return err
	}

	engine.db = db

	return nil
}

func (engine *Engine) Insert(domain string, feedurl string) error {
	mu.Lock()
	defer mu.Unlock()
	_, err := engine.dbInsertPreparedFeed.Exec(domain, feedurl)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	return nil

}

func (engine *Engine) Update(feed *Domain) error {
	mu.Lock()
	defer mu.Unlock()
	_, err := engine.dbUpdatePreparedFeed.Exec(
		utils.NullInt64(feed.LastFetchAt),
		utils.NullString(feed.FeedTitle),
		utils.NullString(feed.Language),
		utils.NullBool(feed.PageAbout),
		utils.NullBool(feed.PageBlogRoll),
		utils.NullBool(feed.PageWriting),
		utils.NullBool(feed.PageNow),
		utils.NullBool(feed.URLNews),
		utils.NullBool(feed.URLBlog),
		utils.NullBool(feed.URLHumanName),
		utils.NullBool(feed.DomainIsPopular),
		utils.NullString(feed.DomainTLD),
		utils.NullString(feed.Platform),
		feed.Domain,
	)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	return nil
}

func domainRowScan(res *sql.Rows) (*Domain, error) {
	var feedUrl string
	var domain string
	var lastFetchAt sql.NullInt64
	var feedTitle, language sql.NullString

	var pageAbout, pageBlogRoll, pageWriting, pageNow, urlNews, urlBlog, urlHumanName, domainIsPopular sql.NullInt64
	var domainTLD, platform sql.NullString

	err := res.Scan(
		&domain,
		&feedUrl,
		&lastFetchAt,
		&feedTitle,
		&language,
		&pageAbout,
		&pageBlogRoll,
		&pageWriting,
		&pageNow,
		&urlNews,
		&urlBlog,
		&urlHumanName,
		&domainIsPopular,
		&domainTLD,
		&platform,
	)
	if err != nil {
		return nil, err
	}

	d := &Domain{
		Domain:          domain,
		FeedURL:         feedUrl,
		LastFetchAt:     lastFetchAt.Int64,
		FeedTitle:       feedTitle.String,
		Language:        language.String,
		PageAbout:       pageAbout.Int64 > 0,
		PageBlogRoll:    pageBlogRoll.Int64 > 0,
		PageWriting:     pageWriting.Int64 > 0,
		PageNow:         pageNow.Int64 > 0,
		URLNews:         urlNews.Int64 > 0,
		URLBlog:         urlBlog.Int64 > 0,
		URLHumanName:    urlHumanName.Int64 > 0,
		DomainIsPopular: domainIsPopular.Int64 > 0,
		DomainTLD:       domainTLD.String,
		Platform:        platform.String,
	}
	return d, nil
}

func (engine *Engine) getDomainsToRefresh() ([]*Domain, error) {
	res, err := engine.db.Query(fmt.Sprintf("SELECT %s FROM domains WHERE lastFetchAt IS NULL;", DOMAIN_SELECT))
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var urls []*Domain

	for res.Next() {
		domain, err := domainRowScan(res)
		if err != nil {
			return nil, err
		}
		urls = append(urls, domain)

	}
	if err := res.Err(); err != nil {
		return nil, err
	}
	return urls, nil
}

func (engine *Engine) GetAll() ([]*Domain, error) {
	res, err := engine.db.Query(fmt.Sprintf("SELECT %s FROM domains;", DOMAIN_SELECT))
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var urls []*Domain

	for res.Next() {
		domain, err := domainRowScan(res)
		if err != nil {
			return nil, err
		}

		urls = append(urls, domain)

	}
	if err := res.Err(); err != nil {
		return nil, err
	}
	return urls, nil
}

func (engine *Engine) getURLsToCrawl() ([]string, error) {
	// length(h.url) - LENGTH(REPLACE(h.url, "/", "")) as scnt
	// AND scnt >= 4
	missingFromHN, err := engine.db.Query(`SELECT url
	FROM (
		SELECT h.url, 
		ROW_NUMBER() OVER (PARTITION BY h.domain ORDER BY addedAt DESC) AS rn 
		
		FROM to_crawl h
		LEFT JOIN domains d ON d.domain = h.domain
		LEFT JOIN url_requests u on u.domain = h.domain
		WHERE h.url IS NOT NULL AND d.domain IS NULL AND u.domain IS NULL AND h.score > 2
	) raw 
	
	WHERE rn <= 2;`)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	defer missingFromHN.Close()

	if err == sql.ErrNoRows {
		return nil, nil
	}

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
