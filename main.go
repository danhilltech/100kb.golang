package main

import (
	"flag"
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

	httpChunkSize := flag.Int("http-chunk-size", 100, "number of http chunks")
	httpWorkers := flag.Int("http-workers", 20, "number of http workers")
	hnFetchSize := flag.Int("hn-fetch-size", 10_000, "number of hn links to get")
	metaChunkSize := flag.Int("meta-chunk-size", 50, "number of meta chunks")
	metaWorkers := flag.Int("meta-workers", 4, "number of meta workers")

	flag.Parse()

	fmt.Println("Config:")
	fmt.Printf("\thttpChunkSize:\t\t%d\n", *httpChunkSize)
	fmt.Printf("\tttpWorkers:\t\t%d\n", *httpWorkers)
	fmt.Printf("\thnFetchSize:\t\t%d\n", *hnFetchSize)
	fmt.Printf("\tmetaChunkSize:\t\t%d\n", *metaChunkSize)
	fmt.Printf("\tmetaWorkers:\t\t%d\n", *metaWorkers)

	db, err := db.InitDB("/dbs/output")
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

	// 1. Get latest hackernews content
	err = hnEngine.RunRefresh(*httpChunkSize, *hnFetchSize, *httpWorkers)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// // 2. Check HN stories for any new feeds
	err = feedEngine.RunNewFeedSearch(*httpChunkSize, *httpWorkers)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// // 3. Get latest articles from our feeds
	err = feedEngine.RunFeedRefresh(*httpChunkSize, *httpWorkers)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 4. Crawl any new articles for content
	err = articleEngine.RunArticleIndex(*httpChunkSize, *httpWorkers)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 5. Generate metadata for articles
	err = articleEngine.RunArticleMeta(*metaChunkSize, *metaWorkers)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	db.Tidy()

}
