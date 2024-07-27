package domain

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

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
chromeAnalysis
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
	chromeAnalysis = ?
	
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

// move adblock filter to validate domain

func (engine *Engine) Update(txn *sql.Tx, feed *Domain) error {

	chromeStr, err := json.Marshal(feed.ChromeAnalysis)
	if err != nil {
		return err
	}

	_, err = txn.Stmt(engine.dbUpdatePreparedFeed).Exec(
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
		chromeStr,
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

	var urlNews, urlBlog, urlHumanName, domainIsPopular sql.NullInt64
	var domainTLD, platform sql.NullString
	var chromeAnalysisStr sql.NullString

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
		&chromeAnalysisStr,
	)
	if err != nil {
		return nil, err
	}

	var reqOut ChromeAnalysis
	if chromeAnalysisStr.Valid {
		err := json.Unmarshal([]byte(chromeAnalysisStr.String), &reqOut)
		if err != nil {
			return nil, err
		}

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
		ChromeAnalysis:  &reqOut,
	}
	return d, nil
}

func (engine *Engine) getDomainsToRefresh() ([]*Domain, error) {
	res, err := engine.db.Query(fmt.Sprintf("SELECT %s FROM domains WHERE lastFetchAt IS NULL OR lastFetchAt < %d ORDER BY lastFetchAt ASC LIMIT %d;", DOMAIN_SELECT, time.Now().Unix()-REFRESH_AGO_SECONDS, REFRESH_LIMIT))
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
	res, err := engine.db.Query(fmt.Sprintf("SELECT %s FROM domains WHERE lastValidateAt IS NULL OR lastValidateAt < %d ORDER BY lastValidateAt ASC LIMIT %d;", DOMAIN_SELECT, VALIDATE_AGO_SECONDS, VALIDATE_LIMIT))
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

func (engine *Engine) getLatestArticleURLs() (map[string]string, error) {

	// var url string
	// if err := engine.db.QueryRow("SELECT url FROM articles WHERE domain = ? AND stage = ? ORDER BY lastContentExtractAt DESC LIMIT 1;", d.Domain, article.STAGE_COMPLETE).Scan(&url); err != nil {
	// 	return "", err
	// }
	// return url, nil

	res, err := engine.db.Query(`SELECT domain, url FROM (
	SELECT
		domain,
		url, 
		ROW_NUMBER() OVER (PARTITION BY domain ORDER BY publishedAt DESC) AS rn 
		
		FROM articles
		WHERE stage = ?
	) as inn WHERE rn = 1;`, article.STAGE_COMPLETE)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	out := make(map[string]string)

	for res.Next() {
		var domain, url string
		err := res.Scan(
			&domain,
			&url,
		)
		if err != nil {
			return nil, err
		}

		out[domain] = url

	}
	if err := res.Err(); err != nil {
		return nil, err
	}
	return out, nil

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
		WHERE h.url IS NOT NULL AND d.domain IS NULL AND h.score >= 4
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

		domain := u.Hostname()

		isBad := false
		//TODO improve this
		for _, bad := range http.BANNED_DOMAINS {
			if strings.HasSuffix(domain, bad) {
				isBad = true
			}
		}

		if !isBad {
			for _, bad := range PopularDomainList {
				if len(bad) > 0 && strings.HasSuffix(domain, bad) {
					isBad = true
				}
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
