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
	"github.com/sjwhitworth/golearn/knn"
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
	// candidates := []string{"https://herbertlui.net/"}

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
			continue

		}
		continue

		fmt.Println(domain.Articles[0].Url)

		// fmt.Print("\033[H\033[2J")
		// fmt.Print(Green)2
		// fmt.Println("################################################")
		// fmt.Print(Gray)
		// for _, b := range domain.Articles[0].Body.Content {
		// 	fmt.Printf("%s\n\n", b.Text)
		// }

		// fmt.Print(Green)
		// fmt.Println("-----")
		// fmt.Print(Gray)
		// for _, b := range domain.Articles[1].Body.Content {
		// 	fmt.Printf("%s\n\n", b.Text)
		// }

		// fmt.Print(Green)
		// fmt.Println("-----")
		// fmt.Print(Gray)
		// for _, b := range domain.Articles[2].Body.Content {
		// 	fmt.Printf("%s\n\n", b.Text)
		// }

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

	err = trainKNN(goodEntries, allDomains)
	if err != nil {
		return err
	}

	err = trainRF(goodEntries, allDomains)
	if err != nil {
		return err
	}

	return nil

}

func trainRF(goodEntries []Entry, allDomains []*domain.Domain) error {
	// Training

	attrCount := 12

	attrs := make([]base.Attribute, attrCount)

	n := 0
	attrs[n] = base.NewCategoricalAttribute()
	n++

	attrs[n] = base.NewCategoricalAttribute()
	attrs[n].SetName("fpr")
	n++
	attrs[n] = base.NewCategoricalAttribute()
	attrs[n].SetName("wordCount")
	n++
	attrs[n] = base.NewCategoricalAttribute()
	attrs[n].SetName("wordsPerByte")
	n++

	attrs[n] = base.NewCategoricalAttribute()
	attrs[n].SetName("wordsPerP")
	n++
	attrs[n] = base.NewCategoricalAttribute()
	attrs[n].SetName("goodPcnt")

	n++
	attrs[n] = base.NewCategoricalAttribute()
	attrs[n].SetName("urlHuman")

	n++
	attrs[n] = base.NewCategoricalAttribute()
	attrs[n].SetName("urlNews")

	n++
	attrs[n] = base.NewCategoricalAttribute()
	attrs[n].SetName("pageNow")

	n++
	attrs[n] = base.NewCategoricalAttribute()
	attrs[n].SetName("pageAbout")

	n++
	attrs[n] = base.NewCategoricalAttribute()
	attrs[n].SetName("pageBlogRoll")

	n++
	attrs[n] = base.NewCategoricalAttribute()
	attrs[n].SetName("pageWriting")

	instances := base.NewDenseInstances()

	// Add the attributes
	newSpecs := make([]base.AttributeSpec, len(attrs))
	for i, a := range attrs {
		newSpecs[i] = instances.AddAttribute(a)
	}

	instances.Extend(len(goodEntries))

	// 1 title begins with a number
	// 2 number of paragraphs with more than 40 words
	// 3 average sentence length
	// 4 number of code tags
	// 5 bad keyword density ("how to", "github")
	// 6 identify self help
	// 7 youtube/podcasts
	// https://webring.xxiivv.com/#vitbaisa
	// https://frankmeeuwsen.com/blogroll/
	// title uniqueness/levenstien
	GetCodeTagCount

	for i, train := range goodEntries {

		var domain *domain.Domain

		for _, d := range allDomains {
			if train.domain == d.Domain {
				domain = d
			}
		}

		n = 0

		if train.score == 2 {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("good"))
		}
		if train.score == 1 {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("bad"))
		}
		n++

		if domain.GetFPR() > 0.08 {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("5"))
		} else if domain.GetFPR() > 0.04 {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("4"))
		} else if domain.GetFPR() > 0.02 {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("3"))

		} else {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("0"))
		}
		n++

		if domain.GetWordCount() > 1200 {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("2"))
		} else if domain.GetWordCount() > 300 {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("1"))
		} else {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("0"))
		}
		n++

		if domain.GetWordsPerByte() > 0.05 {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("2"))
		} else if domain.GetWordsPerByte() > 0.01 {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("1"))
		} else {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("0"))
		}
		n++

		if domain.GetWordsPerParagraph() > 200 {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("2"))
		} else if domain.GetWordsPerParagraph() > 40 {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("1"))
		} else {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("0"))
		}
		n++

		if domain.GetGoodBadTagRatio() > 0.95 {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("1"))
		} else if domain.GetGoodBadTagRatio() > 0.8 {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("1"))
		} else {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("0"))
		}
		n++

		if domain.URLHumanName {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("1"))
		} else {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("0"))
		}
		n++

		if domain.URLNews {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("1"))
		} else {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("0"))
		}
		n++

		if domain.PageNow {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("1"))
		} else {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("0"))
		}
		n++

		if domain.PageAbout {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("1"))
		} else {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("0"))
		}
		n++

		if domain.PageBlogRoll {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("1"))
		} else {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("0"))
		}
		n++

		if domain.PageWriting {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("1"))
		} else {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("0"))
		}

	}

	instances.AddClassAttribute(attrs[0])

	// fmt.Println("Running Chi Merge...")
	// filt := filters.NewChiMergeFilter(instances, 0.99)
	// for _, a := range base.NonClassFloatAttributes(instances) {
	// 	filt.AddAttribute(a)
	// }
	// fmt.Println("Training chi merge...")
	// filt.Train()
	// fmt.Println("Filtering with chi merge...")
	// instf := base.NewLazilyFilteredInstances(instances, filt)

	trainData, testData := base.InstancesTrainTestSplit(instances, 0.6)

	fmt.Println(trainData)

	fmt.Println("Building model...")

	cls := ensemble.NewRandomForest(80, 4)

	// Create a 60-40 training-test split

	err := cls.Fit(trainData)
	if err != nil {
		return err
	}

	fmt.Println("Predicting...")
	//Calculates the Euclidean distance and returns the most popular label
	predictions, err := cls.Predict(testData)
	if err != nil {
		panic(err)
	}

	// var predictions base.FixedDataGrid

	// Prints precision/recall metrics
	confusionMat, err := evaluation.GetConfusionMatrix(testData, predictions)
	if err != nil {
		panic(fmt.Sprintf("Unable to get confusion matrix: %s", err.Error()))
	}
	fmt.Println(evaluation.GetSummary(confusionMat))

	return nil
}

func trainKNN(goodEntries []Entry, allDomains []*domain.Domain) error {
	// Training

	attrCount := 6

	attrs := make([]base.Attribute, attrCount)

	n := 0
	attrs[n] = base.NewCategoricalAttribute()
	n++

	// attrs[n] = base.NewFloatAttribute("urlHumanName")
	// n++
	// attrs[n] = base.NewFloatAttribute("urlNews")
	// n++
	// attrs[n] = base.NewFloatAttribute("urlBlog")
	// n++
	// attrs[n] = base.NewFloatAttribute("pageAbout")
	// n++
	// attrs[n] = base.NewFloatAttribute("pageBlogRoll")
	// n++
	// attrs[n] = base.NewFloatAttribute("pageNow")
	// n++
	// attrs[n] = base.NewFloatAttribute("pageWriting")
	// n++
	attrs[n] = base.NewFloatAttribute("fpr")
	n++
	attrs[n] = base.NewFloatAttribute("wordCount")
	n++
	attrs[n] = base.NewFloatAttribute("wordsPerByte")
	n++
	// attrs[n] = base.NewFloatAttribute("goodCount")
	// n++
	// attrs[n] = base.NewFloatAttribute("badCount")
	// n++
	attrs[n] = base.NewFloatAttribute("wordsPerP")
	n++
	attrs[n] = base.NewFloatAttribute("goodPcnt")

	instances := base.NewDenseInstances()

	// Add the attributes
	newSpecs := make([]base.AttributeSpec, len(attrs))
	for i, a := range attrs {
		newSpecs[i] = instances.AddAttribute(a)
	}

	instances.Extend(len(goodEntries))

	// 1 title begins with a number
	// 2 number of paragraphs with more than 40 words
	// 3 average sentence length
	// 4 number of code tags
	// 5 bad keyword density ("how to", "github")
	// 6 identify self help
	// 7 youtube/podcasts
	// https://webring.xxiivv.com/#vitbaisa
	// https://frankmeeuwsen.com/blogroll/
	// title uniqueness/levenstien

	for i, train := range goodEntries {

		var domain *domain.Domain

		for _, d := range allDomains {
			if train.domain == d.Domain {
				domain = d
			}
		}

		n = 0

		if train.score == 2 {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("good"))
		}
		if train.score == 1 {
			instances.Set(newSpecs[n], i, newSpecs[n].GetAttribute().GetSysValFromString("bad"))
		}
		n++

		// if domain.URLHumanName {
		// 	instances.Set(newSpecs[n], i, base.PackFloatToBytes(1.0))
		// } else {
		// 	instances.Set(newSpecs[n], i, base.PackFloatToBytes(-1.0))
		// }
		// n++

		// if domain.URLNews {
		// 	instances.Set(newSpecs[n], i, base.PackFloatToBytes(1.0))
		// } else {
		// 	instances.Set(newSpecs[n], i, base.PackFloatToBytes(-1.0))
		// }
		// n++

		// if domain.URLBlog {
		// 	instances.Set(newSpecs[n], i, base.PackFloatToBytes(1.0))
		// } else {
		// 	instances.Set(newSpecs[n], i, base.PackFloatToBytes(-1.0))
		// }
		// n++

		// if domain.PageAbout {
		// 	instances.Set(newSpecs[n], i, base.PackFloatToBytes(1.0))
		// } else {
		// 	instances.Set(newSpecs[n], i, base.PackFloatToBytes(-1.0))
		// }
		// n++

		// if domain.PageBlogRoll {
		// 	instances.Set(newSpecs[n], i, base.PackFloatToBytes(1.0))
		// } else {
		// 	instances.Set(newSpecs[n], i, base.PackFloatToBytes(-1.0))
		// }
		// n++

		// if domain.PageNow {
		// 	instances.Set(newSpecs[n], i, base.PackFloatToBytes(1.0))
		// } else {
		// 	instances.Set(newSpecs[n], i, base.PackFloatToBytes(-1.0))
		// }
		// n++

		// if domain.PageWriting {
		// 	instances.Set(newSpecs[n], i, base.PackFloatToBytes(1.0))
		// } else {
		// 	instances.Set(newSpecs[n], i, base.PackFloatToBytes(-1.0))
		// }
		// n++

		instances.Set(newSpecs[n], i, base.PackFloatToBytes(domain.GetFPR()))
		n++
		instances.Set(newSpecs[n], i, base.PackFloatToBytes(float64(domain.GetWordCount())))
		n++
		instances.Set(newSpecs[n], i, base.PackFloatToBytes(domain.GetWordsPerByte()))
		n++

		// instances.Set(newSpecs[n], i, base.PackFloatToBytes(float64(domain.GetGoodTagCount())))
		// n++
		// instances.Set(newSpecs[n], i, base.PackFloatToBytes(float64(domain.GetBadTagCount())))
		// n++
		instances.Set(newSpecs[n], i, base.PackFloatToBytes(float64(domain.GetWordsPerParagraph())))
		n++
		instances.Set(newSpecs[n], i, base.PackFloatToBytes(float64(domain.GetGoodBadTagRatio())))

	}

	maxFloats := make([]float64, attrCount)
	minFloats := make([]float64, attrCount)
	for i := 1; i < attrCount; i++ {
		for row := 0; row < len(goodEntries); row++ {
			byteVal := instances.Get(newSpecs[i], row)

			fltVal := base.UnpackBytesToFloat(byteVal)

			if fltVal > maxFloats[i] {
				maxFloats[i] = fltVal
			}
			if fltVal < minFloats[i] {
				minFloats[i] = fltVal
			}

		}
	}

	for row := 0; row < len(goodEntries); row++ {
		for i := 1; i < attrCount; i++ {
			byteVal := instances.Get(newSpecs[i], row)

			fltVal := base.UnpackBytesToFloat(byteVal)

			rng := maxFloats[i] - minFloats[i]

			nrmVal := ((fltVal - minFloats[i]) / rng) * 100

			instances.Set(newSpecs[i], row, base.PackFloatToBytes(nrmVal))

		}

	}

	instances.AddClassAttribute(attrs[0])

	// fmt.Println("Running Chi Merge...")
	// filt := filters.NewChiMergeFilter(instances, 0.99)
	// for _, a := range base.NonClassFloatAttributes(instances) {
	// 	filt.AddAttribute(a)
	// }
	// fmt.Println("Training chi merge...")
	// filt.Train()
	// fmt.Println("Filtering with chi merge...")
	// instf := base.NewLazilyFilteredInstances(instances, filt)

	trainData, testData := base.InstancesTrainTestSplit(instances, 0.75)

	fmt.Println(trainData)

	fmt.Println("Building model...")
	cls := knn.NewKnnClassifier("cosine", "linear", 3)
	// cls := trees.NewID3DecisionTree(0.6)

	// cls := ensemble.NewRandomForest(40, 3)

	// Create a 60-40 training-test split

	err := cls.Fit(trainData)
	if err != nil {
		return err
	}

	fmt.Println("Predicting...")
	//Calculates the Euclidean distance and returns the most popular label
	predictions, err := cls.Predict(testData)
	if err != nil {
		panic(err)
	}

	// var predictions base.FixedDataGrid

	// Prints precision/recall metrics
	confusionMat, err := evaluation.GetConfusionMatrix(testData, predictions)
	if err != nil {
		panic(fmt.Sprintf("Unable to get confusion matrix: %s", err.Error()))
	}
	fmt.Println(evaluation.GetSummary(confusionMat))

	return nil
}
