package domain

import (
	"database/sql"
	"fmt"
	"net/url"

	"github.com/danhilltech/100kb.golang/pkg/article"
	"github.com/danhilltech/100kb.golang/pkg/http"
	"github.com/danhilltech/100kb.golang/pkg/utils"
)

const DOMAIN_SELECT = `
domains.domain, 
domains.feedUrl, 
domains.lastFetchAt, 
domains.lastValidateAt, 
feedTitle, 
language,
urlNews,
urlBlog,
urlHumanName,
domainIsPopular,
domainTLD,
platform,
domainGoogleAds
`

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
	lastValidateAt = ?,
	feedTitle = ?, 
	language = ?,
	urlNews = ?,
	urlBlog = ?,
	urlHumanName = ?,
	domainIsPopular = ?,
	domainTLD = ?,
	platform = ?,
	domainGoogleAds = ?
	
	WHERE domain = ?;`)
	if err != nil {
		return err
	}

	engine.db = db

	return nil
}

func (engine *Engine) Insert(txn *sql.Tx, domain string, feedurl string) error {

	_, err := txn.Stmt(engine.dbInsertPreparedFeed).Exec(domain, feedurl)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	return nil

}

func (engine *Engine) Update(txn *sql.Tx, feed *Domain) error {
	_, err := txn.Stmt(engine.dbUpdatePreparedFeed).Exec(
		utils.NullInt64(feed.LastFetchAt),
		utils.NullInt64(feed.LastValidateAt),
		utils.NullString(feed.FeedTitle),
		utils.NullString(feed.Language),
		utils.NullBool(feed.URLNews),
		utils.NullBool(feed.URLBlog),
		utils.NullBool(feed.URLHumanName),
		utils.NullBool(feed.DomainIsPopular),
		utils.NullString(feed.DomainTLD),
		utils.NullString(feed.Platform),
		utils.NullBool(feed.DomainGoogleAds),
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
	var lastFetchAt, lastValidateAt sql.NullInt64
	var feedTitle, language sql.NullString

	var urlNews, urlBlog, urlHumanName, domainIsPopular, domainGoogleAds sql.NullInt64
	var domainTLD, platform sql.NullString

	err := res.Scan(
		&domain,
		&feedUrl,
		&lastFetchAt,
		&lastValidateAt,
		&feedTitle,
		&language,
		&urlNews,
		&urlBlog,
		&urlHumanName,
		&domainIsPopular,
		&domainTLD,
		&platform,
		&domainGoogleAds,
	)
	if err != nil {
		return nil, err
	}

	d := &Domain{
		Domain:          domain,
		FeedURL:         feedUrl,
		LastFetchAt:     lastFetchAt.Int64,
		LastValidateAt:  lastValidateAt.Int64,
		FeedTitle:       feedTitle.String,
		Language:        language.String,
		URLNews:         urlNews.Int64 > 0,
		URLBlog:         urlBlog.Int64 > 0,
		URLHumanName:    urlHumanName.Int64 > 0,
		DomainIsPopular: domainIsPopular.Int64 > 0,
		DomainTLD:       domainTLD.String,
		Platform:        platform.String,
		DomainGoogleAds: domainGoogleAds.Int64 > 0,
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

func (engine *Engine) getDomainsToValidate() ([]*Domain, error) {
	res, err := engine.db.Query(fmt.Sprintf("SELECT %s FROM domains WHERE lastValidateAt IS NULL;", DOMAIN_SELECT))
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

func (engine *Engine) getLatestArticleURL(d *Domain) (string, error) {
	res := engine.db.QueryRow(fmt.Sprintf("SELECT url FROM articles WHERE domain = ? AND stage = ? ORDER BY lastContentExtractAt DESC LIMIT 1;", d.Domain, article.STAGE_COMPLETE))

	var url string
	err := res.Scan(&url)
	if err != nil {
		return "", err
	}

	if err := res.Err(); err != nil {
		return "", err
	}
	return url, nil
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
		WHERE h.url IS NOT NULL AND d.domain IS NULL AND h.score >= 2
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
