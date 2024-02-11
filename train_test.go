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

func TestBuildSVM(t *testing.T) {

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

	fmt.Println("Starting engine")

	articleEngine, err := article.NewEngine(sqliteDatabase, "/app/.cache", true)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	scores := readCsvFile("/app/scoring/scored.csv")

	txn, _ := sqliteDatabase.Begin()
	defer txn.Rollback()

	svm, err := os.Create("svm.csv")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer svm.Close()

	fmt.Println("Starting articles")
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

			svm.WriteString(fmt.Sprintf("12:%d ", boolToInt(article.URLBlog)))
			svm.WriteString(fmt.Sprintf("13:%d ", boolToInt(article.URLHumanName)))
			svm.WriteString(fmt.Sprintf("14:%d ", boolToInt(article.URLNews)))

			svm.WriteString(fmt.Sprintf("15:%d ", boolToInt(article.PageAbout)))
			svm.WriteString(fmt.Sprintf("16:%d ", boolToInt(article.PageBlogRoll)))
			svm.WriteString(fmt.Sprintf("17:%d ", boolToInt(article.PageWriting)))

			svm.WriteString(fmt.Sprintf("18:%d ", boolToInt(article.DomainIsPopular)))

			svm.WriteString(fmt.Sprintf("19:%d ", tldToInt(article.DomainTLD, "com")))
			svm.WriteString(fmt.Sprintf("20:%d ", tldToInt(article.DomainTLD, "co.uk")))
			svm.WriteString(fmt.Sprintf("21:%d ", tldToInt(article.DomainTLD, "gov")))
			svm.WriteString(fmt.Sprintf("22:%d ", tldToInt(article.DomainTLD, "edu")))
			svm.WriteString(fmt.Sprintf("23:%d ", tldToInt(article.DomainTLD, "net")))
			svm.WriteString(fmt.Sprintf("24:%d ", tldToInt(article.DomainTLD, "org")))
			svm.WriteString(fmt.Sprintf("25:%d ", tldToInt(article.DomainTLD, "ai")))
			svm.WriteString(fmt.Sprintf("26:%d ", tldToInt(article.DomainTLD, "io")))
			svm.WriteString(fmt.Sprintf("27:%d ", tldToInt(article.DomainTLD, "blog")))

			svm.WriteString(fmt.Sprintf("28:%0.8f ", getClassificationScore(article, "technology")))
			svm.WriteString(fmt.Sprintf("29:%0.8f ", getClassificationScore(article, "life")))
			svm.WriteString(fmt.Sprintf("30:%0.8f ", getClassificationScore(article, "family")))
			svm.WriteString(fmt.Sprintf("31:%0.8f ", getClassificationScore(article, "science")))
			svm.WriteString(fmt.Sprintf("32:%0.8f ", getClassificationScore(article, "politics")))
			svm.WriteString(fmt.Sprintf("33:%0.8f ", getClassificationScore(article, "news")))
			svm.WriteString(fmt.Sprintf("34:%0.8f ", getClassificationScore(article, "religion")))
			svm.WriteString(fmt.Sprintf("35:%0.8f ", getClassificationScore(article, "god")))
			svm.WriteString(fmt.Sprintf("36:%0.8f ", getClassificationScore(article, "programming")))
			svm.WriteString(fmt.Sprintf("37:%0.8f ", getClassificationScore(article, "food")))
			svm.WriteString(fmt.Sprintf("38:%0.8f ", getClassificationScore(article, "crypto")))
			svm.WriteString(fmt.Sprintf("39:%0.8f ", getClassificationScore(article, "investing")))
			svm.WriteString(fmt.Sprintf("40:%0.8f ", getClassificationScore(article, "management")))
			svm.WriteString(fmt.Sprintf("41:%0.8f ", getClassificationScore(article, "nature")))

			svm.WriteString("\n")

			fmt.Println(article.Classifications.GetKeywords())
		}
	}

}

func getClassificationScore(article *article.Article, keyword string) float32 {
	kws := article.Classifications.GetKeywords()

	for _, k := range kws {
		if k.Text == keyword {
			return k.Score
		}
	}
	return 0.0
}

func boolToInt(b bool) int {
	if b {
		return 1
	} else {
		return -1
	}
}

func tldToInt(tld string, match string) int {
	if tld == match {
		return 1
	} else {
		return -1
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
