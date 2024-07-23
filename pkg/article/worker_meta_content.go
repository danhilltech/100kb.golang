package article

import (
	"fmt"
	"runtime"
	"time"
)

func (engine *Engine) RunArticleMeta(chunkSize int) error {

	articles, err := engine.getArticlesToContentExtract()
	if err != nil {
		return err
	}

	fmt.Printf("Generating %d article metas\n", len(articles))

	done := 0
	lastA := 0
	t := time.Now().UnixMilli()

	printSize := 200

	txn, _ := engine.db.Begin()

	jobs := make(chan *Article, len(articles))
	results := make(chan *Article, len(articles))

	workers := runtime.NumCPU() - 2
	// workers := 1

	for w := 1; w <= workers; w++ {
		go engine.articeMetaWorker(jobs, results)
	}

	for j := 1; j <= len(articles); j++ {
		jobs <- articles[j-1]
	}
	close(jobs)

	for a := 0; a < len(articles); a++ {
		article := <-results

		err = engine.Update(txn, article)
		if err != nil {
			fmt.Println(article.Url, err)
		}

		if a > 0 && a%chunkSize == 0 {
			err := txn.Commit()
			if err != nil {
				return err
			}
			txn, _ = engine.db.Begin()
		}
		if a > 0 && a%printSize == 0 {
			diff := time.Now().UnixMilli() - t
			qps := (float64(done-lastA) / float64(diff)) * 1000
			lastA = done
			t = time.Now().UnixMilli()
			fmt.Printf("\tdone %d/%d at %0.2f/s\n", done, len(articles), qps)

		}
		done++
	}

	err = txn.Commit()
	if err != nil {
		return err
	}
	fmt.Printf("\tdone %d/%d\n", done, len(articles))

	return nil

}

func (engine *Engine) articeMetaWorker(jobs <-chan *Article, results chan<- *Article) {
	for id := range jobs {
		err := engine.articleExtractContent(id)
		if err != nil {
			fmt.Println(err)
		}
		results <- id
	}
}
