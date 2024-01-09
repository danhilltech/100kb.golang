package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

var dbInsertPreparedHN *sql.Stmt
var dbInsertPreparedCrawl *sql.Stmt
var dbInsertPreparedFeed *sql.Stmt
var dbUpdatePreparedFeed *sql.Stmt
var dbInsertPreparedArticle *sql.Stmt
var dbUpdatePreparedArticle *sql.Stmt

var (
	ErrDatabaseNotOpen = errors.New("database not open")
)

const DB_INIT_SCRIPT = `
CREATE TABLE IF NOT EXISTS hacker_news (
	id INTEGER PRIMARY KEY,
    url TEXT,
    author TEXT,
	type TEXT,
    addedAt INTEGER NOT NULL,
	postedAt INTEGER,
	score INTEGER
);

CREATE UNIQUE INDEX IF NOT EXISTS hacker_news_url ON hacker_news(url);

CREATE TABLE IF NOT EXISTS crawls (
	url TEXT PRIMARY KEY,
	hackerNewsId INTEGER,
	lastCrawlAt INTEGER
);

CREATE INDEX IF NOT EXISTS crawls_hackerNewsId ON crawls(hackerNewsId);

CREATE TABLE IF NOT EXISTS feeds (
	url TEXT PRIMARY KEY,
	lastFetchAt INTEGER,
	title TEXT,
	description TEXT,
	language TEXT
);

CREATE TABLE IF NOT EXISTS articles (
	url TEXT PRIMARY KEY,
	feedUrl TEXT,
	publishedAt INTEGER,
	lastFetchAt INTEGER,
	title TEXT,
	description TEXT,
	bodyRaw TEXT,
	body TEXT,
	wordCount INTEGER,
	firstPersonRatio REAL
);

CREATE INDEX IF NOT EXISTS articles_feedUrl ON articles(feedUrl);
`

func initDB() error {
	fmt.Println("creating database...")

	sqliteDatabase, err := sql.Open("sqlite3", "file:test.db?cache=shared&mode=rwc") // Open the created SQLite File
	if err != nil {
		return err
	}
	db = sqliteDatabase

	db.SetMaxOpenConns(1)

	_, err = db.Exec(DB_INIT_SCRIPT)
	if err != nil {
		return err
	}

	dbInsertPreparedHN, err = db.Prepare("INSERT INTO hacker_news(id, url, author, type, addedAt, postedAt, score) VALUES(?, ?, ?, ?, ?, ?, ?)  ON CONFLICT(url) DO NOTHING")
	if err != nil {
		return err
	}

	dbInsertPreparedCrawl, err = db.Prepare("INSERT INTO crawls(url, hackerNewsId, lastCrawlAt) VALUES(?, ?, ?)  ON CONFLICT(url) DO UPDATE SET lastCrawlAt = excluded.lastCrawlAt;")
	if err != nil {
		return err
	}

	dbInsertPreparedFeed, err = db.Prepare("INSERT INTO feeds(url) VALUES(?) ON CONFLICT(url) DO NOTHING;")
	if err != nil {
		return err
	}

	dbUpdatePreparedFeed, err = db.Prepare("UPDATE feeds SET lastFetchAt = ?, title = ?, description = ?, language = ? WHERE url = ?;")
	if err != nil {
		return err
	}

	dbInsertPreparedArticle, err = db.Prepare("INSERT INTO articles(url, feedUrl, publishedAt) VALUES(?, ?, ?) ON CONFLICT(url) DO NOTHING;")
	if err != nil {
		return err
	}

	dbUpdatePreparedArticle, err = db.Prepare("UPDATE articles SET lastFetchAt = ?, bodyRaw = ?, title = ?, description = ?, body = ?, wordCount = ?, firstPersonRatio = ? WHERE url = ?;")
	if err != nil {
		return err
	}

	return nil
}

func stopDB() {
	if dbInsertPreparedHN != nil {
		dbInsertPreparedHN.Close()
	}
	if db != nil {
		db.Close()
	}
}

func saveHNItem(item *HNItem, txchan *sql.Tx) error {
	_, err := txchan.Stmt(dbInsertPreparedHN).Exec(item.ID, nullString(item.URL), nullString(item.By), item.Type, time.Now().Unix(), item.Time, item.Score)
	return err
}

func saveCrawl(item *HNUrlToCrawl, txchan *sql.Tx) error {
	_, err := txchan.Stmt(dbInsertPreparedCrawl).Exec(item.Url, item.ID, time.Now().Unix())
	return err
}

func addNewFeed(url string, txchan *sql.Tx) error {
	_, err := txchan.Stmt(dbInsertPreparedFeed).Exec(url)
	return err
}

func updateFeed(feed *Feed, txchan *sql.Tx) error {
	_, err := txchan.Stmt(dbUpdatePreparedFeed).Exec(nullInt64(feed.LastFetchAt), nullString(feed.Title), nullString(feed.Description), nullString(feed.Language), feed.Url)
	return err
}

func addNewArticle(feed *Feed, article *Article, txchan *sql.Tx) error {
	_, err := txchan.Stmt(dbInsertPreparedArticle).Exec(article.Url, feed.Url, article.PublishedAt)
	return err
}

func updateArticle(article *Article, txchan *sql.Tx) error {
	var articleBody []byte
	var err error

	if len(article.BodyRaw) > 0 {

		articleBody, err = json.Marshal(article.BodyRaw)
		if err != nil {
			return err
		}

	}

	_, err = txchan.Stmt(dbUpdatePreparedArticle).Exec(
		nullInt64(article.LastFetchAt),
		nullString(string(articleBody)),
		nullString(article.Title),
		nullString(article.Description),
		nullString(article.Body),
		nullInt64(article.WordCount),
		nullFloat64(article.FirstPersonRatio),
		article.Url,
	)
	return err
}

func getExistingHNIDs() ([]int, error) {
	res, err := db.Query("SELECT id FROM hacker_news")
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

func getURLsToCrawlFromHN() ([]*HNUrlToCrawl, error) {
	res, err := db.Query("SELECT id, hacker_news.url FROM hacker_news LEFT JOIN crawls on crawls.hackerNewsId = hacker_news.id WHERE hacker_news.url IS NOT NULL AND crawls.lastCrawlAt IS NULL;")
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var urls []*HNUrlToCrawl

	for res.Next() {
		var id int
		var url string
		err = res.Scan(&id, &url)
		if err != nil {
			return nil, err
		}

		urls = append(urls, &HNUrlToCrawl{ID: id, Url: url})
	}
	if err := res.Err(); err != nil {
		return nil, err
	}
	return urls, nil
}

func getFeedsToRefresh() ([]*Feed, error) {
	res, err := db.Query("SELECT url FROM feeds WHERE lastFetchAt IS NULL;")
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var urls []*Feed

	for res.Next() {
		var url string
		err = res.Scan(&url)
		if err != nil {
			return nil, err
		}

		urls = append(urls, &Feed{Url: url})
	}
	if err := res.Err(); err != nil {
		return nil, err
	}
	return urls, nil
}

func getArticlesToIndex() ([]*Article, error) {
	res, err := db.Query("SELECT url FROM articles WHERE lastFetchAt IS NULL;")
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var urls []*Article

	for res.Next() {
		var url string
		err = res.Scan(&url)
		if err != nil {
			return nil, err
		}

		urls = append(urls, &Article{Url: url})
	}
	if err := res.Err(); err != nil {
		return nil, err
	}
	return urls, nil
}

func articleRowScan(res *sql.Rows) (*Article, error) {
	var url string
	var feedUrl string
	var publishedAt int64
	var lastFetchAt sql.NullInt64
	var title sql.NullString
	var description sql.NullString
	var bodyRawJSON []byte
	var body sql.NullString
	err := res.Scan(&url, &feedUrl, &publishedAt, &lastFetchAt, &title, &description, &bodyRawJSON, &body)
	if err != nil {
		return nil, err
	}

	var bodyRaw []string
	if bodyRawJSON != nil {
		err = json.Unmarshal(bodyRawJSON, &bodyRaw)
		if err != nil {
			return nil, err
		}
	}

	article := &Article{Url: url, FeedUrl: feedUrl, PublishedAt: publishedAt, LastFetchAt: lastFetchAt.Int64, Title: title.String, Description: description.String, BodyRaw: bodyRaw, Body: body.String}

	return article, nil
}

func getArticlesToMetaData() ([]*Article, error) {
	res, err := db.Query("SELECT url, feedUrl, publishedAt, lastFetchAt, title, description, bodyRaw, body FROM articles WHERE body IS NULL AND bodyRaw IS NOT NULL AND lastFetchAt > 0;")
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var urls []*Article

	for res.Next() {

		article, err := articleRowScan(res)
		if err != nil {
			return nil, err
		}
		urls = append(urls, article)
	}
	if err := res.Err(); err != nil {
		return nil, err
	}
	return urls, nil
}

func getArticlesByFeed(tx *sql.Tx, feed string, excludeUrl string) ([]*Article, error) {
	res, err := tx.Query("SELECT url, feedUrl, publishedAt, lastFetchAt, title, description, bodyRaw, body FROM articles WHERE feedUrl = ? AND url != ?", feed, excludeUrl)

	if err != nil {
		return nil, err
	}
	if err := res.Err(); err != nil {
		return nil, err
	}
	defer res.Close()

	var urls []*Article

	for res.Next() {
		article, err := articleRowScan(res)
		if err != nil {
			return nil, err
		}
		urls = append(urls, article)
	}
	if err := res.Err(); err != nil {
		return nil, err
	}
	return urls, nil
}

func nullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func nullInt64(s int64) sql.NullInt64 {
	if s == 0 {
		return sql.NullInt64{}
	}
	return sql.NullInt64{
		Int64: s,
		Valid: true,
	}
}

func nullFloat64(s float64) sql.NullFloat64 {
	if s == 0 {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{
		Float64: s,
		Valid:   true,
	}
}

func dbTidy() {
	db.Exec("VACUUM;")
}
