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
CREATE TABLE IF NOT EXISTS hacker_news (
	id INTEGER PRIMARY KEY,
    url TEXT,
	domain TEXT,
    author TEXT,
	type TEXT,
    addedAt INTEGER NOT NULL,
	postedAt INTEGER,
	score INTEGER
);

CREATE INDEX IF NOT EXISTS hacker_news_url ON hacker_news(url);

CREATE TABLE IF NOT EXISTS url_requests (
	url TEXT PRIMARY KEY,
	domain TEXT NOT NULL,
	lastAttemptAt INTEGER,
	status TEXT,
	contentType TEXT
);

CREATE TABLE IF NOT EXISTS domains (
	domain TEXT PRIMARY KEY,
	feedUrl TEXT,
	lastFetchAt INTEGER,
	feedTitle TEXT,
	language TEXT

);

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
	wordCount INTEGER,
	h1Count INTEGER,
	hnCount INTEGER,
	pCount INTEGER,
	badCount INTEGER,
	firstPersonRatio REAL,
	sentenceEmbedding BLOB,
	extractedKeywords BLOB,
	classifications BLOB
);

CREATE INDEX IF NOT EXISTS articles_feedUrl ON articles(feedUrl);
`

func InitDB(name string, mode string) (*Database, error) {
	fmt.Printf("creating database\t")
	defer fmt.Printf("❄️\n")

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
