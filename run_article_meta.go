package main

import "fmt"

func RunArticleMeta() error {
	articles, err := getArticlesToMetaData()
	if err != nil {
		return err
	}

	fmt.Printf("Getting %d articles to meta...\n", len(articles))

	chunkSize := 100

	chunkIds := chunkArticlesToRefresh(articles, chunkSize)

	for _, chunk := range chunkIds {
		err = doArticleMeta(chunk)
		if err != nil {
			return err
		}
	}

	return nil
}

func doArticleMeta(chunk []*Article) error {

	fmt.Println("Starting chunk of article meta")

	insertTxn, err := db.Begin()
	if err != nil {
		return err
	}
	defer insertTxn.Rollback()

	err = articleMetas(insertTxn, chunk)
	if err != nil {
		return err
	}

	for _, article := range chunk {
		// save it
		err = updateArticle(article, insertTxn)
		if err != nil {
			return err
		}

	}

	err = insertTxn.Commit()
	if err != nil {
		return err
	}

	return nil
}
