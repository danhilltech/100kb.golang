package http

import (
	"database/sql"
	"io"
	"runtime"
	"testing"

	"github.com/danhilltech/100kb.golang/pkg/db"
)

func TestCachedRequest(t *testing.T) {
	engine, err := NewClient("/workspaces/100kb.golang/.cache")
	if err != nil {
		t.Error(err)
	}

	u := "https://gist.githubusercontent.com/danhilltech/2d085066e390f4bb74bf485a2b04b203/raw/0e1bcc0e64ddef9485c821a60f45f52fd84bf0ff/Dockerfile"

	resp, err := engine.Get(u)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	if len(data) != 3482 {
		t.Fail()
	}
}

func TestCachedRequestWithTxn(t *testing.T) {

	_ = db.DB_INIT_SCRIPT

	sqliteDatabase, err := sql.Open("sqlite3", ":memory:?mode=memory&_journal_mode=WAL&_sync=FULL") // Open the created SQLite File
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer sqliteDatabase.Close()

	// db.SetMaxOpenConns(1)

	_, err = sqliteDatabase.Exec(db.DB_INIT_SCRIPT)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	engine, err := NewClient("/workspaces/100kb.golang/.cache")
	if err != nil {
		t.Error(err)
		t.Error(err)
	}

	u := "https://gist.githubusercontent.com/danhilltech/2d085066e390f4bb74bf485a2b04b203/raw/0e1bcc0e64ddef9485c821a60f45f52fd84bf0ff/Dockerfile"

	txn, err := sqliteDatabase.Begin()
	if err != nil {
		t.Error(err)
	}
	defer txn.Rollback()

	resp, err := engine.GetWithTxn(u, txn)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer resp.Body.Close()

	err = txn.Commit()
	if err != nil {
		t.Error(err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	if len(data) != 3482 {

		t.Fail()
	}

	var url string
	var domain string
	var lastAttemptAt sql.NullInt64
	var status sql.NullString
	var contentType sql.NullString

	err = sqliteDatabase.QueryRow("SELECT url, domain, lastAttemptAt, status, contentType FROM url_requests WHERE url = ?", u).Scan(&url, &domain, &lastAttemptAt, &status, &contentType)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if domain != "gist.githubusercontent.com" {
		t.Error(domain)
		t.FailNow()
	}
}

func TestCachedHeadRequest(t *testing.T) {
	engine, err := NewClient("/workspaces/100kb.golang/.cache")
	if err != nil {
		t.Error(err)
	}

	u := "https://gist.githubusercontent.com/danhilltech/2d085066e390f4bb74bf485a2b04b203/raw/0e1bcc0e64ddef9485c821a60f45f52fd84bf0ff/Dockerfile"

	resp, err := engine.Head(u)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	t.Log(resp)

}

func BenchmarkRequest(b *testing.B) {
	engine, err := NewClient("/workspaces/100kb.golang/.cache")
	if err != nil {
		b.Error(err)
	}

	u := "https://gist.githubusercontent.com/danhilltech/2d085066e390f4bb74bf485a2b04b203/raw/0e1bcc0e64ddef9485c821a60f45f52fd84bf0ff/Dockerfile"

	b.ResetTimer()

	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	for i := 0; i < b.N; i++ {
		resp, err := engine.Get(u)
		if err != nil {
			b.Error(err)
		}

		resp.Body.Close()
	}
	runtime.ReadMemStats(&m2)
	b.Log("total:", m2.TotalAlloc-m1.TotalAlloc)
	b.Log("mallocs:", m2.Mallocs-m1.Mallocs)

}
