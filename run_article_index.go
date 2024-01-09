package main

import "fmt"

func RunArticleRefresh() error {
	articles, err := getArticlesToIndex()
	if err != nil {
		return err
	}

	fmt.Printf("Getting %d articles to index...\n", len(articles))

	chunkSize := 100

	chunkIds := chunkArticlesToRefresh(articles, chunkSize)

	for _, chunk := range chunkIds {
		err = doFeedArticleIndex(chunk)
		if err != nil {
			return err
		}
	}

	return nil
}

func doFeedArticleIndex(chunk []*Article) error {

	fmt.Println("Starting chunk of article indexing")

	insertTxn, err := db.Begin()
	if err != nil {
		return err
	}
	defer insertTxn.Rollback()

	err = articleIndexes(chunk)
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
