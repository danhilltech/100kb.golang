package train

/**
docker run \
 --name graphite\
 --restart=always\
 -p 8080:80\
 -p 2003-2004:2003-2004\
 -p 2023-2024:2023-2024\
 -p 8125:8125/udp\
 -p 8126:8126\
 graphiteapp/graphite-statsd
*/

import (
	"bytes"
	"database/sql"
	_ "embed"
	"fmt"
	"net/url"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/danhilltech/100kb.golang/pkg/article"
	"github.com/danhilltech/100kb.golang/pkg/crawler"
	"github.com/danhilltech/100kb.golang/pkg/db"
	"github.com/danhilltech/100kb.golang/pkg/domain"
	"github.com/danhilltech/100kb.golang/pkg/output"
	"github.com/smira/go-statsd"
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Magenta = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

//go:embed scoring/candidates.txt
var candidateList string

type Entry struct {
	url    string
	score  int
	domain string
}

/*
docker run \
 --name graphite\
 --restart=always\
 -p 8080:80\
 -p 2003-2004:2003-2004\
 -p 2023-2024:2023-2024\
 -p 8125:8125/udp\
 -p 8126:8126\
 graphiteapp/graphite-statsd
*/

func TrainSVM(cacheDir string) error {

	candidates := strings.Split(candidateList, "\n")
	// candidates := []string{"https://herbertlui.net/"}

	entries := []Entry{}

	existing := false

	_, err := os.Stat("train.db")

	if err == nil {
		existing = true
	}

	fmt.Printf("Existing: %t\n", existing)

	database, err := sql.Open("sqlite3", "file:train.db?cache=shared&_journal_mode=WAL&_sync=FULL") // Open the created SQLite File
	if err != nil {

		return err
	}

	if !existing {
		_, err = database.Exec(db.DB_INIT_SCRIPT)
		if err != nil {
			return err
		}
	}

	statsdClient := statsd.NewClient("172.17.0.1:8125", statsd.MetricPrefix("100kb."))

	// db.SetMaxOpenConns(1)

	crawlEngine, err := crawler.NewEngine(database)
	if err != nil {

		return err
	}

	articleEngine, err := article.NewEngine(database, statsdClient, cacheDir, false)
	if err != nil {

		return err
	}
	defer articleEngine.Close()

	domainEngine, err := domain.NewEngine(database, articleEngine, statsdClient, cacheDir)
	if err != nil {

		return err
	}

	for _, g := range candidates {
		u, err := url.Parse(g)
		if err != nil {
			return err
		}
		err = crawlEngine.InsertToCrawl(&crawler.ToCrawl{
			URL:    g,
			Domain: u.Hostname(),
			Score:  10,
		})
		if err != nil {

			return err
		}
		entries = append(entries, Entry{url: g, score: 0, domain: u.Hostname()})
	}

	if !existing {

		httpChunkSize := 100
		httpWorkers := 40
		metaChunkSize := 10

		// // 2. Check HN stories for any new feeds
		err = domainEngine.RunNewFeedSearch(httpChunkSize, httpWorkers)
		if err != nil {
			return err
		}

		// // 3. Get latest articles from our feeds
		err = domainEngine.RunFeedRefresh(httpChunkSize, httpWorkers)
		if err != nil {
			return err
		}

		// 4. Crawl any new articles for content
		err = articleEngine.RunArticleIndex(httpChunkSize, httpWorkers)
		if err != nil {
			return err
		}

		err = articleEngine.RunArticleMeta(metaChunkSize)
		if err != nil {
			return err
		}

		// 6. Second pass metas
		err = articleEngine.RunArticleMetaPassII(metaChunkSize)
		if err != nil {
			return err
		}
	}

	allDomains, err := domainEngine.GetAll()
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	w := tabwriter.NewWriter(buf, 0, 0, 1, ' ', 0)
	allDomains[0].TabulateHeader(w)
	for _, d := range allDomains {
		articles, err := articleEngine.FindByFeedURL(d.FeedURL)
		if err != nil {
			return err
		}
		d.Articles = append(d.Articles, articles...)

		if len(d.GetLatestArticlesToScore()) >= 3 {
			d.Tabulate(w)
		}

	}
	w.Flush()
	fmt.Println(buf.String())

	var goodEntries []Entry

	// label the data

	for _, train := range entries {

		var domain *domain.Domain

		for _, d := range allDomains {
			if train.domain == d.Domain {
				domain = d
			}
		}

		if domain == nil {
			// fmt.Printf("Can't find: %s\n", train.domain)
			continue
		}

		if len(domain.GetLatestArticlesToScore()) < 3 {
			continue
		}

		labels := readJSON("labels.json")

		if labels[train.domain] > 0 {
			train.score = labels[train.domain]

			goodEntries = append(goodEntries, train)

		}

	}

	err = trainKNN(goodEntries, allDomains)
	if err != nil {
		return err
	}

	// err = trainRF(goodEntries, allDomains)
	// if err != nil {
	// 	return err
	// }

	// err = trainSVM(goodEntries, allDomains)
	// if err != nil {
	// 	return err
	// }

	// err = trainBespoke(goodEntries, allDomains)
	// if err != nil {
	// 	return err
	// }

	// err = trainPerceptron(goodEntries, allDomains)
	// if err != nil {
	// 	return err
	// }

	mdl, err := trainLogistic(goodEntries, allDomains)
	if err != nil {
		return err
	}

	err = mdl.Save("model.json")
	if err != nil {
		return err
	}

	// Output
	articles, err := articleEngine.GetAllValid()
	if err != nil {
		return err
	}

	engine, err := output.NewRenderEnding("output-train", articles, allDomains, mdl, database, articleEngine)
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

	engine.RunHttp("./output-train")

	return nil

}
