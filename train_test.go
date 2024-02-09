package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/danhilltech/100kb.golang/pkg/article"
	"github.com/danhilltech/100kb.golang/pkg/db"
)

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

	articleEngine, err := article.NewEngine(sqliteDatabase, "/workspaces/100kb.golang/.cache", false)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	scores := readCsvFile("/workspaces/100kb.golang/scoring/scored.csv")

	txn, _ := sqliteDatabase.Begin()
	defer txn.Rollback()

	svm, err := os.Create("svm.csv")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer svm.Close()

	for _, row := range scores {
		url := row[0]
		score, err := strconv.ParseInt(row[1], 10, 64)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		article, err := articleEngine.BuildArticleSingle(txn, url)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		cleanScore := 0
		switch score {
		case 1:
			cleanScore = 1
		case 2:
			cleanScore = 1
		case 3:
			cleanScore = 1
		case 4:
			cleanScore = 1
		case 5:
			cleanScore = 2
		}

		if article.HTMLLength > 0 {

			svm.WriteString(fmt.Sprintf("%d ", cleanScore))
			svm.WriteString(fmt.Sprintf("1:%d ", article.BadCount))
			svm.WriteString(fmt.Sprintf("2:%d ", article.PCount))
			svm.WriteString(fmt.Sprintf("3:%d ", article.H1Count))
			svm.WriteString(fmt.Sprintf("4:%0.8f ", article.FirstPersonRatio))
			svm.WriteString(fmt.Sprintf("5:%0.8f ", float64(article.PCount)/float64(article.HTMLLength)))
			svm.WriteString(fmt.Sprintf("6:%0.8f ", float64(article.BadCount)/float64(article.HTMLLength)))

			svm.WriteString(fmt.Sprintf("7:%d ", article.HNCount))
			svm.WriteString(fmt.Sprintf("8:%0.8f ", float64(article.PCount)/float64(article.H1Count)))
			svm.WriteString(fmt.Sprintf("9:%d ", article.WordCount))
			svm.WriteString(fmt.Sprintf("10:%0.8f ", float64(article.PCount)/float64(article.WordCount)))
			svm.WriteString(fmt.Sprintf("11:%0.8f ", float64(article.WordCount)/float64(article.HTMLLength)))
			svm.WriteString("\n")
		}
	}

}

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}
