package crawler

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/danhilltech/100kb.golang/pkg/utils"
)

func (engine *Engine) initDB(db *sql.DB) error {
	var err error
	engine.dbInsertPreparedToCrawl, err = db.Prepare("INSERT INTO to_crawl(url, hn_id, author, type, addedAt, postedAt, score, domain) VALUES(?, ?, ?, ?, ?, ?, ?, ?) ON CONFLICT(hn_id) DO UPDATE SET score=excluded.score")
	if err != nil {

		return err
	}

	engine.db = db

	return nil
}

func (engine *Engine) InsertToCrawl(txn *sql.Tx, item *ToCrawl) error {
	_, err := txn.Stmt(engine.dbInsertPreparedToCrawl).Exec(item.URL, utils.NullInt64(int64(item.HNID)), utils.NullString(item.By), item.Type, time.Now().Unix(), item.Time, item.Score, item.Domain)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	return nil
}

func (engine *Engine) getExistingIDs() ([]int, error) {
	fmt.Printf("Getting existing IDs\t")
	defer fmt.Printf("❄️\n")
	res, err := engine.db.Query("SELECT hn_id FROM to_crawl WHERE hn_id > 0 AND postedAt < ?", time.Now().Unix()-60*6048)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var ids []int

	for res.Next() {
		var id sql.NullInt64
		err = res.Scan(&id)
		if err != nil {
			return nil, err
		}
		if id.Valid {
			ids = append(ids, int(id.Int64))
		}
	}
	if err := res.Err(); err != nil {
		return nil, err
	}
	return ids, nil
}
