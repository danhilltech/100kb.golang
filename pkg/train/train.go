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
	"log"
	"net/url"
	"strings"
	"text/tabwriter"

	"github.com/danhilltech/100kb.golang/pkg/article"
	"github.com/danhilltech/100kb.golang/pkg/crawler"
	"github.com/danhilltech/100kb.golang/pkg/db"
	"github.com/danhilltech/100kb.golang/pkg/domain"
	"github.com/smira/go-statsd"

	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/ensemble"
	"github.com/sjwhitworth/golearn/evaluation"
	"github.com/sjwhitworth/golearn/filters"
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

func TrainSVM(cacheDir string) error {

	candidates := strings.Split(candidateList, "\n")

	entries := []Entry{}

	database, err := sql.Open("sqlite3", "file::memory:?cache=shared&_journal_mode=WAL&_sync=FULL") // Open the created SQLite File
	if err != nil {

		return err
	}

	statsdClient := statsd.NewClient("172.17.0.1:8125", statsd.MetricPrefix("100kb."))

	// db.SetMaxOpenConns(1)

	_, err = database.Exec(db.DB_INIT_SCRIPT)
	if err != nil {

		return err
	}

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
		})
		if err != nil {

			return err
		}
		entries = append(entries, Entry{url: g, score: 0, domain: u.Hostname()})
	}

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

		d.Tabulate(w)

	}
	w.Flush()
	fmt.Println(buf.String())

	// Training

	attrs := make([]base.Attribute, 10)

	attrs[0] = base.NewCategoricalAttribute()
	attrs[1] = base.NewCategoricalAttribute()
	attrs[2] = base.NewCategoricalAttribute()
	attrs[3] = base.NewCategoricalAttribute()
	attrs[4] = base.NewCategoricalAttribute()
	attrs[5] = base.NewCategoricalAttribute()
	attrs[6] = base.NewCategoricalAttribute()
	attrs[7] = base.NewCategoricalAttribute()
	attrs[8] = base.NewCategoricalAttribute()

	attrs[9] = base.NewFloatAttribute("fpr")

	instances := base.NewDenseInstances()

	// Add the attributes
	newSpecs := make([]base.AttributeSpec, len(attrs))
	for i, a := range attrs {
		newSpecs[i] = instances.AddAttribute(a)
	}

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

		if len(domain.Articles) < 3 {
			continue
		}

		labels := readJSON("labels.json")

		if labels[train.domain] > 0 {
			train.score = labels[train.domain]

			goodEntries = append(goodEntries, train)

		}
		continue

		fmt.Print("\033[H\033[2J")
		fmt.Print(Green)
		fmt.Println("################################################")
		fmt.Print(Gray)
		for _, b := range domain.Articles[0].Body.Content {
			fmt.Printf("%s\n\n", b.Text)
		}

		fmt.Print(Green)
		fmt.Println("-----")
		fmt.Print(Gray)
		for _, b := range domain.Articles[1].Body.Content {
			fmt.Printf("%s\n\n", b.Text)
		}

		fmt.Print(Green)
		fmt.Println("-----")
		fmt.Print(Gray)
		for _, b := range domain.Articles[2].Body.Content {
			fmt.Printf("%s\n\n", b.Text)
		}

		fmt.Print(Yellow)
		fmt.Println("Good=2 Bad=1")

		var score int64
		_, err := fmt.Scan(&score)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print(Reset)

		switch score {
		case 1:
			fmt.Println("Score BAD")
			train.score = 1
		case 2:
			fmt.Println("Score GOOD")
			train.score = 2
		}
		labels[train.domain] = train.score

		writeJSON("labels.json", labels)

	}

	instances.Extend(len(goodEntries))

	// 1 title begins with a number
	// 2 number of paragraphs with more than 40 words
	// 3 average sentence length
	// 4 number of code tags
	// 5 bad keyword density ("how to", "github")
	// 6 identify self help

	for i, train := range goodEntries {

		var domain *domain.Domain

		for _, d := range allDomains {
			if train.domain == d.Domain {
				domain = d
			}
		}

		if train.score == 2 {
			instances.Set(newSpecs[0], i, newSpecs[0].GetAttribute().GetSysValFromString("good"))
		}
		if train.score == 1 {
			instances.Set(newSpecs[0], i, newSpecs[0].GetAttribute().GetSysValFromString("bad"))
		}

		if domain.URLHumanName {
			instances.Set(newSpecs[1], i, newSpecs[1].GetAttribute().GetSysValFromString("True"))
		} else {
			instances.Set(newSpecs[1], i, newSpecs[1].GetAttribute().GetSysValFromString("False"))
		}
		if domain.URLNews {
			instances.Set(newSpecs[2], i, newSpecs[2].GetAttribute().GetSysValFromString("True"))
		} else {
			instances.Set(newSpecs[2], i, newSpecs[2].GetAttribute().GetSysValFromString("False"))
		}
		if domain.URLBlog {
			instances.Set(newSpecs[3], i, newSpecs[3].GetAttribute().GetSysValFromString("True"))
		} else {
			instances.Set(newSpecs[3], i, newSpecs[3].GetAttribute().GetSysValFromString("False"))
		}
		if domain.PageAbout {
			instances.Set(newSpecs[4], i, newSpecs[4].GetAttribute().GetSysValFromString("True"))
		} else {
			instances.Set(newSpecs[4], i, newSpecs[4].GetAttribute().GetSysValFromString("False"))
		}
		if domain.PageBlogRoll {
			instances.Set(newSpecs[5], i, newSpecs[5].GetAttribute().GetSysValFromString("True"))
		} else {
			instances.Set(newSpecs[5], i, newSpecs[5].GetAttribute().GetSysValFromString("False"))
		}
		if domain.PageNow {
			instances.Set(newSpecs[6], i, newSpecs[6].GetAttribute().GetSysValFromString("True"))
		} else {
			instances.Set(newSpecs[6], i, newSpecs[6].GetAttribute().GetSysValFromString("False"))
		}
		if domain.PageWriting {
			instances.Set(newSpecs[7], i, newSpecs[7].GetAttribute().GetSysValFromString("True"))
		} else {
			instances.Set(newSpecs[7], i, newSpecs[7].GetAttribute().GetSysValFromString("False"))
		}
		instances.Set(newSpecs[8], i, newSpecs[8].GetAttribute().GetSysValFromString(domain.DomainTLD))

		instances.Set(newSpecs[9], i, base.PackFloatToBytes(domain.GetFPR()))
		// instances.Set(newSpecs[2], i, base.PackFloatToBytes(float64(article.BadCount)/float64(article.HTMLLength)))
		// instances.Set(newSpecs[3], i, base.PackFloatToBytes(float64(article.WordCount)))

	}

	instances.AddClassAttribute(attrs[0])

	fmt.Println("Running Chi Merge...")
	filt := filters.NewChiMergeFilter(instances, 0.90)
	for _, a := range base.NonClassFloatAttributes(instances) {
		filt.AddAttribute(a)
	}
	fmt.Println("Training chi merge...")
	filt.Train()
	fmt.Println("Filtering with chi merge...")
	instf := base.NewLazilyFilteredInstances(instances, filt)

	trainData, testData := base.InstancesTrainTestSplit(instf, 0.7)

	fmt.Println(trainData)

	fmt.Println("Building model...")
	// cls := knn.NewKnnClassifier("euclidean", "linear", 2)
	cls := ensemble.NewRandomForest(70, 5)

	// Create a 60-40 training-test split

	err = cls.Fit(trainData)
	if err != nil {
		return err
	}

	fmt.Println("Predicting...")
	//Calculates the Euclidean distance and returns the most popular label
	predictions, err := cls.Predict(testData)
	if err != nil {
		panic(err)
	}

	// Prints precision/recall metrics
	confusionMat, err := evaluation.GetConfusionMatrix(testData, predictions)
	if err != nil {
		panic(fmt.Sprintf("Unable to get confusion matrix: %s", err.Error()))
	}
	fmt.Println(evaluation.GetSummary(confusionMat))

	return nil

}
