package main

import "fmt"

func RunFeedRefresh() error {
	feeds, err := getFeedsToRefresh()
	if err != nil {
		return err
	}

	fmt.Printf("Getting %d feeds to refresh...\n", len(feeds))

	chunkSize := 100

	chunkIds := chunkFeedsToRefresh(feeds, chunkSize)

	for _, chunk := range chunkIds {
		err = doFeedRefreshChunk(chunk)
		if err != nil {
			return err
		}
	}

	return nil
}

func doFeedRefreshChunk(chunk []*Feed) error {

	fmt.Println("Starting chunk of feed refresh")

	insertTxn, err := db.Begin()
	if err != nil {
		return err
	}
	defer insertTxn.Rollback()

	err = feedRefreshes(chunk)
	if err != nil {
		return err
	}

	for _, feed := range chunk {
		// save it
		err = updateFeed(feed, insertTxn)
		if err != nil {
			return err
		}

		for _, article := range feed.Articles {
			err = addNewArticle(feed, &article, insertTxn)
			if err != nil {
				return err
			}
		}

		// save articles
	}

	err = insertTxn.Commit()
	if err != nil {
		return err
	}

	return nil
}
