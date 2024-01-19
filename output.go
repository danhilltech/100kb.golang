package main

import (
	"database/sql"
	"fmt"
	"math"
	"os"

	"github.com/danhilltech/100kb.golang/pkg/article"
)

var pageSize = 100

func CreateOutput(db *sql.DB) error {

	articleEngine, err := article.NewEngine(db)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer articleEngine.Close()

	txn, err := db.Begin()
	if err != nil {
		return err
	}
	defer txn.Rollback()

	articles, err := articleEngine.GetAllValid(txn)
	if err != nil {
		return err
	}
	err = txn.Commit()
	if err != nil {
		return err
	}

	articleCount := len(articles)

	numPages := int(math.Ceil(float64(articleCount) / float64(pageSize)))

	fmt.Printf("Articles:\t%d\n", articleCount)
	fmt.Printf("Page size:\t%d\n", pageSize)
	fmt.Printf("Pages:\t%d\n", numPages)

	return nil

}
