package feed

import (
	"database/sql"

	"github.com/danhilltech/100kb.golang/pkg/utils"
)

func (engine *Engine) initDB(db *sql.DB) error {
	var err error
	engine.dbInsertPreparedFeed, err = db.Prepare("INSERT INTO feeds(url) VALUES(?) ON CONFLICT(url) DO NOTHING;")
	if err != nil {
		return err
	}

	engine.dbUpdatePreparedFeed, err = db.Prepare("UPDATE feeds SET lastFetchAt = ?, title = ?, description = ?, language = ? WHERE url = ?;")
	if err != nil {
		return err
	}

	engine.db = db

	return nil
}

func (engine *Engine) Insert(url string, txchan *sql.Tx) error {
	_, err := txchan.Stmt(engine.dbInsertPreparedFeed).Exec(url)
	return err
}

func (engine *Engine) Update(feed *Feed, txchan *sql.Tx) error {
	_, err := txchan.Stmt(engine.dbUpdatePreparedFeed).Exec(
		utils.NullInt64(feed.LastFetchAt),
		utils.NullString(feed.Title),
		utils.NullString(feed.Description),
		utils.NullString(feed.Language),
		feed.Url,
	)
	return err
}

func (engine *Engine) getFeedsToRefresh(txchan *sql.Tx) ([]*Feed, error) {
	res, err := txchan.Query("SELECT url FROM feeds WHERE lastFetchAt IS NULL;")
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
