package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "net/http/pprof"

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
	debug := flag.Bool("debug", false, "run debugging tools")
	mode := flag.String("mode", "index", "which process to run")

	flag.Parse()

	fmt.Println("Config:")
	fmt.Printf("  httpChunkSize:\t%d\n", *httpChunkSize)
	fmt.Printf("  httpWorkers:\t\t%d\n", *httpWorkers)
	fmt.Printf("  hnFetchSize:\t\t%d\n", *hnFetchSize)
	fmt.Printf("  metaChunkSize:\t%d\n", *metaChunkSize)
	fmt.Printf("  metaWorkers:\t\t%d\n", *metaWorkers)
	fmt.Printf("Mode\t%s", *mode)

	if *debug {
		// go tool pprof -top http://localhost:6060/debug/pprof/heap
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}

	switch *mode {
	case "index":

		db, err := db.InitDB("/dbs/output", "rwc")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer db.StopDB()

		dbVer, err := db.Version()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("sqlite3 version: \t%s\n", dbVer)

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

		if *debug {
			<-c
		}

	case "output":
		db, err := db.InitDB("/dbs/output", "ro")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer db.StopDB()

		dbVer, err := db.Version()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("sqlite3 version: \t%s\n", dbVer)

		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			db.StopDB()
			fmt.Println("Interupt\t\t🔥🔥🔥")

			os.Exit(1)
		}()

		err = CreateOutput(db.DB)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if *debug {
			<-c
		}

	}

}
