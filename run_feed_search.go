package main

import "fmt"

func RunNewFeedSearch() error {
	urls, err := getURLsToCrawlFromHN()
	if err != nil {
		return err
	}

	fmt.Printf("Getting %d feed searches...\n", len(urls))

	chunkSize := 100

	chunkIds := chunkHNUrlToCrawl(urls, chunkSize)

	for _, chunk := range chunkIds {
		err = doFeedSearchChunk(chunk)
		if err != nil {
			return err
		}
	}

	return nil
}

func doFeedSearchChunk(chunk []*HNUrlToCrawl) error {

	fmt.Println("Starting chunk of feed search")

	insertTxn, err := db.Begin()
	if err != nil {
		return err
	}
	defer insertTxn.Rollback()

	err = crawlURLsForFeeds(chunk)
	if err != nil {
		fmt.Println(err)
	}
	for _, i := range chunk {
		if i.Feed != "" {
			addNewFeed(i.Feed, insertTxn)
			if err != nil {
				fmt.Println(err)
			}
		}

		err = saveCrawl(i, insertTxn)
		if err != nil {
			fmt.Println(err)
		}
	}

	err = insertTxn.Commit()
	if err != nil {
		return err
	}
	return nil
}
