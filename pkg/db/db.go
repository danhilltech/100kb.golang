package db

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// var db *sql.DB

type Database struct {
	DB *sql.DB
}

var (
	ErrDatabaseNotOpen = errors.New("database not open")
)

const DB_INIT_SCRIPT = `
CREATE TABLE IF NOT EXISTS to_crawl (
	url TEXT,
	hn_id INTEGER,
	domain TEXT,
    author TEXT,
	type TEXT,
    addedAt INTEGER NOT NULL,
	postedAt INTEGER,
	score INTEGER,
	PRIMARY KEY (url, hn_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS to_crawl_hn_id ON to_crawl(hn_id);
CREATE INDEX IF NOT EXISTS to_crawl_domain ON to_crawl(domain);
CREATE INDEX IF NOT EXISTS to_crawl_score ON to_crawl(score);


CREATE TABLE IF NOT EXISTS domains (
	domain TEXT PRIMARY KEY,
	feedUrl TEXT,
	lastFetchAt INTEGER,
	lastValidateAt INTEGER,
	feedTitle TEXT,
	language TEXT,
	platform TEXT,

	latestPostAt INTEGER,
	domainIsPopular INTEGER,
	domainTLD TEXT,
	pageAbout INTEGER,
	pageBlogRoll INTEGER,
	pageWriting INTEGER,
	pageNow INTEGER,
	urlNews INTEGER,
	urlBlog INTEGER,
	urlHumanName INTEGER,

	chromeAnalysis TEXT
);

CREATE INDEX IF NOT EXISTS domains_feedUrl ON domains(feedUrl);
CREATE INDEX IF NOT EXISTS domains_lastFetchAt ON domains(lastFetchAt);

CREATE TABLE IF NOT EXISTS articles (
	url TEXT PRIMARY KEY,
	feedUrl TEXT,
	domain TEXT,
	publishedAt INTEGER,
	lastFetchAt INTEGER,
	lastMetaAt INTEGER,
	lastContentExtractAt INTEGER,
	title TEXT,
	description TEXT,
	bodyRaw BLOB,
	body BLOB,
	badCount INTEGER,
	badElementCount INTEGER,
	linkCount INTEGER,
	badLinkCount INTEGER,
	sentenceEmbedding BLOB,
	extractedKeywords BLOB,
	classifications BLOB,
	htmlLength INTEGER,
	stage INTEGER
);

CREATE INDEX IF NOT EXISTS articles_feedUrl ON articles(feedUrl);

`

/*
ALTER TABLE articles ADD COLUMN htmlLength INTEGER;
ALTER TABLE articles ADD COLUMN pageAbout INTEGER;
ALTER TABLE articles ADD COLUMN pageBlogRoll INTEGER;
ALTER TABLE articles ADD COLUMN pageWriting INTEGER;
ALTER TABLE articles ADD COLUMN urlNews INTEGER;
ALTER TABLE articles ADD COLUMN urlBlog INTEGER;
ALTER TABLE articles ADD COLUMN urlHumanName INTEGER;
ALTER TABLE articles ADD COLUMN domainIsPopular INTEGER;
ALTER TABLE articles ADD COLUMN domainTLD TEXT;
ALTER TABLE articles ADD COLUMN stage INTEGER;

*/

func InitDB(name string, mode string) (*Database, error) {
	fmt.Printf("creating database with mode %s\t", mode)
	defer fmt.Printf("❄️\n")

	// file, err := os.Create("/dbs/output.db") // Create SQLite file
	// if err != nil {
	// 	return nil, err
	// }
	// file.Close()

	sqliteDatabase, err := sql.Open("sqlite3", fmt.Sprintf("file:%s.db?mode=%s&_journal_mode=WAL&_sync=FULL", name, mode)) // Open the created SQLite File
	if err != nil {
		return nil, err
	}
	db := sqliteDatabase

	// db.SetMaxOpenConns(1)

	_, err = db.Exec(DB_INIT_SCRIPT)
	if err != nil {
		return nil, err
	}

	return &Database{DB: db}, nil
}

func (db *Database) StopDB() {
	if db != nil {
		db.DB.Close()
	}
}

func (db *Database) Version() (string, error) {
	var ver string
	err := db.DB.QueryRow(`select sqlite_version();`).Scan(&ver)
	return ver, err

}

func (db *Database) Tidy() {
	db.DB.Exec("VACUUM;")
}
