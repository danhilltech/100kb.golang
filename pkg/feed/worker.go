package feed

import (
	"fmt"
	"math/rand"

	"github.com/danhilltech/100kb.golang/pkg/crawler"
)

func (engine *Engine) RunNewFeedSearch(chunkSize int, workers int) error {
	txn, err := engine.db.Begin()
	if err != nil {
		return err
	}
	defer txn.Rollback()

	urls, err := engine.crawlEngine.GetURLsToCrawl(txn)
	if err != nil {
		return err
	}
	err = txn.Commit()
	if err != nil {
		return err
	}

	rand.Shuffle(len(urls), func(i, j int) { urls[i], urls[j] = urls[j], urls[i] })

	chunkIds := crawler.Chunk(urls, chunkSize)

	fmt.Printf("Checking %d HN urls for feeds in %d chunks\n", len(urls), len(chunkIds))

	for _, chunk := range chunkIds {
		err = engine.doFeedSearchChunk(chunk, workers)
		if err != nil {
			return err
		}
	}

	return nil
}

func (engine *Engine) doFeedSearchChunk(chunk []*crawler.UrlToCrawl, workers int) error {
	fmt.Printf("Chunk...\t\t")
	defer fmt.Printf("✅\n")

	insertTxn, err := engine.db.Begin()
	if err != nil {
		return err
	}
	defer insertTxn.Rollback()

	err = engine.crawlEngine.CrawlURLsForFeeds(chunk, workers)
	if err != nil {
		fmt.Println(err)
	}
	for _, i := range chunk {
		if i.Feed != "" {
			err = engine.Insert(i.Feed, insertTxn)
			if err != nil {
				fmt.Println(err)
			}
		}

		err = engine.crawlEngine.AddURL(i, insertTxn)
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

func (engine *Engine) RunFeedRefresh(chunkSize int, workers int) error {
	txn, err := engine.db.Begin()
	if err != nil {
		return err
	}
	defer txn.Rollback()
	feeds, err := engine.getFeedsToRefresh(txn)
	if err != nil {
		return err
	}
	err = txn.Commit()
	if err != nil {
		return err
	}

	chunkIds := Chunk(feeds, chunkSize)

	fmt.Printf("Checking %d feeds for new links in %d chunks\n", len(feeds), len(chunkIds))

	for _, chunk := range chunkIds {
		err = engine.doFeedRefreshChunk(chunk, workers)
		if err != nil {
			return err
		}
	}

	return nil
}

func (engine *Engine) doFeedRefreshChunk(chunk []*Feed, workers int) error {
	fmt.Printf("Chunk...\t\t")
	defer fmt.Printf("✅\n")
	insertTxn, err := engine.db.Begin()
	if err != nil {
		return err
	}
	defer insertTxn.Rollback()

	err = engine.refresh(chunk, workers)
	if err != nil {
		return err
	}

	for _, feed := range chunk {
		// save it
		err = engine.Update(feed, insertTxn)
		if err != nil {
			return err
		}

		for _, article := range feed.Articles {
			err = engine.articleEngine.Insert(&article, feed.Url, insertTxn)
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
