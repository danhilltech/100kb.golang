package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	_ "net/http/pprof"

	"github.com/danhilltech/100kb.golang/pkg/article"
	"github.com/danhilltech/100kb.golang/pkg/crawler"
	"github.com/danhilltech/100kb.golang/pkg/db"
	"github.com/danhilltech/100kb.golang/pkg/domain"
	"github.com/danhilltech/100kb.golang/pkg/output"
	"github.com/danhilltech/100kb.golang/pkg/train"
	"github.com/smira/go-statsd"
)

func main() {
	fmt.Println("Running\t\t\t🔥🔥🔥")

	httpChunkSize := flag.Int("http-chunk-size", 100, "number of http chunks")
	hnFetchSize := flag.Int("hn-fetch-size", 100_000, "number of hn links to get")
	metaChunkSize := flag.Int("meta-chunk-size", 50, "number of meta chunks")
	debug := flag.Bool("debug", false, "run debugging tools")
	mode := flag.String("mode", "index", "which process to run")
	cacheDir := flag.String("cache-dir", ".cache", "where to cache html")
	utilization := flag.Float64("util", 1.0, "pcnt of cores to use")

	flag.Parse()

	cores := runtime.NumCPU()

	fmt.Println("Config:")
	fmt.Printf("  httpChunkSize:\t%d\n", *httpChunkSize)
	fmt.Printf("  hnFetchSize:\t\t%d\n", *hnFetchSize)
	fmt.Printf("  metaChunkSize:\t%d\n", *metaChunkSize)

	fmt.Printf("  utilization:\t\t%0.2f\n", *utilization)

	httpWorkers := int(math.Floor(float64(cores) * *utilization * 8.0))
	metaWorkers := int(math.Floor(float64(cores) * *utilization * 0.5))

	fmt.Printf("  cores:\t\t%d\n", cores)
	fmt.Printf("  httpWorkers:\t\t%d\n", httpWorkers)
	fmt.Printf("  metaWorkers:\t\t%d\n", metaWorkers)

	fmt.Printf("Mode\t%s\n", *mode)

	debugPrinterCtx, cancelDebugPrinter := context.WithCancel(context.Background())

	statsdClient := statsd.NewClient("192.168.1.3:8125", statsd.MetricPrefix("100bk."))

	if *debug {
		// go tool pprof -top http://localhost:6060/debug/pprof/heap
		fmt.Println("Starting debug pprof...")
		go func() {
			log.Println(http.ListenAndServe(":6060", nil))
		}()

		go debugPrinter(debugPrinterCtx)

	}

	dbMode := "r"
	articleLoadML := false

	switch *mode {
	case "index":
		dbMode = "rwc"
		articleLoadML = false
	case "meta":
		dbMode = "rw"
		articleLoadML = true
	case "output":
		dbMode = "rw"
		articleLoadML = false
	case "train":

		err := train.TrainSVM(*cacheDir)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		return
	default:
		dbMode = "rw"
		articleLoadML = false
	}

	db, err := db.InitDB("/dbs/output", dbMode)
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
		fmt.Println("Interupt\t\t🔥🔥🔥")
		db.StopDB()

		os.Exit(1)
	}()

	articleEngine, err := article.NewEngine(db.DB, statsdClient, *cacheDir, articleLoadML)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer articleEngine.Close()

	crawlEngine, err := crawler.NewEngine(db.DB)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	feedEngine, err := domain.NewEngine(db.DB, articleEngine, statsdClient, *cacheDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Engines loaded\t\t🚂🚂🚂")

	// Now run tasks
	switch *mode {
	case "index":
		// 1. Get latest hackernews content
		err = crawlEngine.RunHNRefresh(*httpChunkSize, *hnFetchSize, httpWorkers)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// // 2. Check HN stories for any new feeds
		err = feedEngine.RunNewFeedSearch(*httpChunkSize, httpWorkers)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// // 3. Get latest articles from our feeds
		err = feedEngine.RunFeedRefresh(*httpChunkSize, httpWorkers)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// 4. Crawl any new articles for content
		err = articleEngine.RunArticleIndex(*httpChunkSize, httpWorkers)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		db.Tidy()

	case "meta":
		// 5. Generate metadata for articles
		err = articleEngine.RunArticleMeta(*metaChunkSize)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// 6. Second pass metas
		err = articleEngine.RunArticleMetaPassII(*metaChunkSize)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		db.Tidy()

	case "output":

		articles, err := articleEngine.GetAllValid()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		domains, err := feedEngine.GetAll()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		engine, err := output.NewRenderEnding("output", articles, domains, db.DB, articleEngine)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = engine.ArticleLists()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = engine.StaticFiles()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		engine.RunHttp()

	}

	if *debug {
		cancelDebugPrinter()
		<-c
	}

}

func debugPrinter(ctx context.Context) {

	i := 0

	for {
		select {
		case <-ctx.Done():
			return
		default:
			if i == 6 {
				printMemUsage()
				i = 0
			}

			time.Sleep(500 * time.Millisecond)
			i++
		}
	}
}

func printMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", m.Alloc/1024/1024)
	fmt.Printf("\tTotalAlloc = %v MiB", m.TotalAlloc/1024/1024)
	fmt.Printf("\tSys = %v MiB", m.Sys/1024/1024)
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}
