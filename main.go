package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/danhilltech/100kb.golang/pkg/article"
	"github.com/danhilltech/100kb.golang/pkg/crawler"
	"github.com/danhilltech/100kb.golang/pkg/db"
	"github.com/danhilltech/100kb.golang/pkg/feed"
	"github.com/danhilltech/100kb.golang/pkg/hn"
)

func main() {
	fmt.Println("Running\t\t\t🔥🔥🔥")

	db, err := db.InitDB("test2")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer db.StopDB()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		db.StopDB()
		fmt.Println("Interupt\t\t🔥🔥🔥")

		os.Exit(1)
	}()

	crawlEngine, err := crawler.NewEngine(db.DB)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	articleEngine, err := article.NewEngine(db.DB)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer articleEngine.Close()

	hnEngine, err := hn.NewEngine(db.DB)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	feedEngine, err := feed.NewEngine(db.DB, crawlEngine, articleEngine)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Engines loaded\t\t🚂🚂🚂")

	// Now run tasks

	httpWorkers := 5

	// 1. Get latest hackernews content
	err = hnEngine.RunRefresh(100, 10_00, httpWorkers)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// // 2. Check HN stories for any new feeds
	err = feedEngine.RunNewFeedSearch(200, httpWorkers)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// // 3. Get latest articles from our feeds
	err = feedEngine.RunFeedRefresh(200, httpWorkers)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 4. Crawl any new articles for content
	err = articleEngine.RunArticleIndex(200, httpWorkers)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 5. Generate metadata for articles
	err = articleEngine.RunArticleMeta(10, 1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	db.Tidy()

}
