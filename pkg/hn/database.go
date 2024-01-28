package hn

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/danhilltech/100kb.golang/pkg/utils"
)

func (engine *Engine) initDB(db *sql.DB) error {
	var err error
	engine.dbInsertPreparedHN, err = db.Prepare("INSERT INTO hacker_news(id, url, author, type, addedAt, postedAt, score, domain) VALUES(?, ?, ?, ?, ?, ?, ?, ?)  ON CONFLICT(id) DO NOTHING")
	if err != nil {
		return err
	}

	engine.db = db

	return nil
}

func (engine *Engine) save(item *HNItem, txchan *sql.Tx) error {
	_, err := txchan.Stmt(engine.dbInsertPreparedHN).Exec(item.ID, utils.NullString(item.URL), utils.NullString(item.By), item.Type, time.Now().Unix(), item.Time, item.Score, item.Domain)
	return err
}

func (engine *Engine) getExistingIDs(txchan *sql.Tx) ([]int, error) {
	fmt.Printf("Getting existing IDs\t")
	defer fmt.Printf("❄️\n")
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
