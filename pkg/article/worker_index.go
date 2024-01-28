package article

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/danhilltech/100kb.golang/pkg/http"
)

type ArticleWithHttp struct {
	Article  *Article
	Response *http.URLRequest
}

func (engine *Engine) RunArticleIndex(chunkSize int, workers int) error {
	txn, err := engine.db.Begin()
	if err != nil {
		return err
	}
	defer txn.Rollback()
	articles, err := engine.getArticlesToIndex(txn)
	if err != nil {
		return err
	}
	err = txn.Commit()
	if err != nil {
		return err
	}

	rand.Shuffle(len(articles), func(i, j int) { articles[i], articles[j] = articles[j], articles[i] })

	fmt.Printf("Crawling %d new articles\n", len(articles))

	jobs := make(chan *Article, len(articles))
	results := make(chan *ArticleWithHttp, len(articles))

	for w := 1; w <= workers; w++ {
		go engine.articleIndexWorker(jobs, results)
	}

	for j := 1; j <= len(articles); j++ {
		jobs <- articles[j-1]
	}
	close(jobs)

	txn, err = engine.db.Begin()
	if err != nil {
		return err
	}

	t := time.Now().UnixMilli()
	for a := 0; a < len(articles); a++ {
		articleWithHttp := <-results

		if articleWithHttp.Response != nil {
			err = articleWithHttp.Response.Save(txn)
			if err != nil {
				return err
			}
		}

		// save it
		err = engine.Update(articleWithHttp.Article, txn)
		if err != nil {
			fmt.Println(articleWithHttp.Article.Url, err)
			continue
		}

		if a > 0 && a%chunkSize == 0 {
			diff := time.Now().UnixMilli() - t
			qps := (float64(chunkSize) / float64(diff)) * 1000
			t = time.Now().UnixMilli()
			fmt.Printf("\tdone %d/%d at %0.2f/s\n", a, len(articles), qps)
			err = txn.Commit()
			if err != nil {
				return err
			}
			txn, err = engine.db.Begin()
			if err != nil {
				return err
			}
		}
	}

	err = txn.Commit()
	if err != nil {
		return err
	}
	fmt.Printf("\tdone %d\n", len(articles))

	return nil
}

func (engine *Engine) articleIndexWorker(jobs <-chan *Article, results chan<- *ArticleWithHttp) {
	for id := range jobs {
		res, err := engine.articleIndex(id)
		if err != nil {
			fmt.Println(err)
		}
		results <- &ArticleWithHttp{Article: id, Response: res}
	}
}
