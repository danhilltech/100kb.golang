package main

import (
	"context"
	"flag"
	"fmt"
	"math"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	_ "net/http/pprof"

	"github.com/danhilltech/100kb.golang/pkg/article"
	"github.com/danhilltech/100kb.golang/pkg/crawler"
	"github.com/danhilltech/100kb.golang/pkg/db"
	"github.com/danhilltech/100kb.golang/pkg/domain"
	"github.com/danhilltech/100kb.golang/pkg/output"
	"github.com/danhilltech/100kb.golang/pkg/scorer"
	"github.com/danhilltech/100kb.golang/pkg/train"
	"github.com/smira/go-statsd"
)

var (
	MODE_INDEX  = "index"
	MODE_SEARCH = "search"
	MODE_META   = "meta"
	MODE_TRAIN  = "train"
	MODE_OUTPUT = "output"
)

func main() {
	fmt.Println("Running\t\t\t🔥🔥🔥")

	httpChunkSize := flag.Int("http-chunk-size", 100, "number of http chunks")
	hnFetchSize := flag.Int("hn-fetch-size", 100_000, "number of hn links to get")
	metaChunkSize := flag.Int("meta-chunk-size", 50, "number of meta chunks")
	mode := flag.String("mode", "index", "which process to run")
	cacheDir := flag.String("cache-dir", ".cache", "where to cache html")
	trainDir := flag.String("train-dir", "train", "where to cache html")
	utilization := flag.Float64("util", 1.0, "pcnt of cores to use")
	articleLoadML := flag.Bool("cuda", false, "use CUDA")

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

	statsdClient := statsd.NewClient("192.168.1.3:8125", statsd.MetricPrefix("100bk."))

	c := make(chan os.Signal, 1)
	signal.Notify(c,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	ctx, abort := context.WithCancel(context.Background())
	go func() {
		<-c
		fmt.Println("Interupt\t\t🔥🔥🔥")
		abort()

	}()

	err := runCoreLoop(
		ctx,
		*mode,
		*cacheDir,
		*trainDir,
		statsdClient,
		*articleLoadML,
		*httpChunkSize,
		httpWorkers,
		*hnFetchSize,
		*metaChunkSize,
	)

	fmt.Println("Done\t\t🔥🔥🔥")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func runCoreLoop(
	ctx context.Context,
	mode string,
	cacheDir string,
	trainDir string,
	statsdClient *statsd.Client,
	useML bool,
	httpChunkSize int,
	httpWorkers int,
	hnFetchSize int,
	metaChunkSize int,
) error {
	dbMode := "r"

	switch mode {
	case MODE_INDEX:
		dbMode = "rwc"
	case MODE_SEARCH:
		dbMode = "rwc"
	case MODE_META:
		dbMode = "rw"
	case MODE_OUTPUT:
		dbMode = "rw"
	case MODE_TRAIN:

		err := train.Train(ctx, cacheDir, trainDir)
		if err != nil {
			return err
		}
		return nil
	default:
		dbMode = "rw"
	}

	db, err := db.InitDB("/dbs/output", dbMode)
	if err != nil {
		return err
	}
	defer db.StopDB()

	dbVer, err := db.Version()
	if err != nil {
		return err
	}
	fmt.Printf("sqlite3 version: \t%s\n", dbVer)

	articleEngine, err := article.NewEngine(db.DB, statsdClient, cacheDir, useML)
	if err != nil {
		return err
	}
	defer articleEngine.Close()

	crawlEngine, err := crawler.NewEngine(db.DB)
	if err != nil {
		return err
	}

	feedEngine, err := domain.NewEngine(db.DB, articleEngine, statsdClient, cacheDir)
	if err != nil {
		return err
	}
	defer feedEngine.Close()

	fmt.Println("Engines loaded\t\t🚂🚂🚂")

	// Now run tasks
	switch mode {
	case MODE_INDEX:

		// // 3. Get latest articles from our feeds
		err = feedEngine.RunFeedRefresh(ctx, httpChunkSize, httpWorkers)
		if err != nil {
			return err
		}

		// 4. Crawl any new articles for content
		err = articleEngine.RunArticleIndex(ctx, httpChunkSize, httpWorkers)
		if err != nil {
			return err
		}

		db.Tidy()
	case MODE_SEARCH:
		// 1. Get latest hackernews content
		err = crawlEngine.RunHNRefresh(ctx, httpChunkSize*3, hnFetchSize, httpWorkers)
		if err != nil {
			return err
		}

		// // 2. Check HN stories for any new feeds
		err = feedEngine.RunNewFeedSearch(ctx, httpChunkSize, httpWorkers)
		if err != nil {
			return err
		}

		// // 2. Check HN stories for any new feeds
		err = feedEngine.RunKagiList(ctx)
		if err != nil {
			return err
		}

		db.Tidy()

	case MODE_META:
		// 5. Generate metadata for articles
		err = articleEngine.RunArticleMeta(ctx, metaChunkSize)
		if err != nil {
			return err
		}

		// 6. Second pass metas
		err = articleEngine.RunArticleMetaPassII(ctx)
		if err != nil {
			return err
		}

		// 7. Additional domain validation
		err = feedEngine.RunDomainValidate(ctx, metaChunkSize)
		if err != nil {
			return err
		}

		db.Tidy()

	case MODE_OUTPUT:

		articles, err := articleEngine.GetAllValid()
		if err != nil {
			return err
		}

		domains, err := feedEngine.GetAll()
		if err != nil {
			return err
		}

		model, err := scorer.LoadModel("models/model.json")
		if err != nil {
			return err
		}

		engine, err := output.NewRenderEnding("output", articles, domains, model, db.DB, articleEngine)
		if err != nil {
			return err
		}

		err = engine.Prepare()
		if err != nil {
			return err
		}

		err = engine.ArticleLists()
		if err != nil {
			return err
		}

		err = engine.StaticFiles()
		if err != nil {
			return err
		}

		go engine.RunHttp(ctx, "./output")

		<-ctx.Done()

	}

	return nil
}
