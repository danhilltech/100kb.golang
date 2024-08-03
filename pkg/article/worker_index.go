package article

import (
	"context"
	"math/rand"
	"time"
)

func (engine *Engine) RunArticleIndex(ctx context.Context, chunkSize int, workers int) error {

	articles, err := engine.getArticlesToIndex()
	if err != nil {
		return err
	}

	rand.Shuffle(len(articles), func(i, j int) { articles[i], articles[j] = articles[j], articles[i] })

	engine.log.Printf("Crawling %d new articles\n", len(articles))

	jobs := make(chan *Article, len(articles))
	results := make(chan *Article, len(articles))

	for w := 1; w <= workers; w++ {
		go engine.articleIndexWorker(jobs, results)
	}

	for j := 1; j <= len(articles); j++ {
		jobs <- articles[j-1]
	}
	close(jobs)

	txn, _ := engine.db.Begin()
	defer txn.Rollback()

	t := time.Now().UnixMilli()
	for a := 0; a < len(articles); a++ {
		select {
		case <-ctx.Done():
			txn.Commit()
			return ctx.Err()
		case article := <-results:

			// save it
			err = engine.Update(txn, article)
			if err != nil {
				engine.log.Println(article.Url, err)
				continue
			}

			if a > 0 && a%chunkSize == 0 {
				err := txn.Commit()
				if err != nil {
					return err
				}
				txn, _ = engine.db.Begin()
				diff := time.Now().UnixMilli() - t
				qps := (float64(chunkSize) / float64(diff)) * 1000
				t = time.Now().UnixMilli()
				engine.log.Printf("\tdone %d/%d at %0.2f/s\n", a, len(articles), qps)

			}
		}
	}

	txn.Commit()
	engine.log.Printf("\tdone %d\n", len(articles))

	return nil
}

func (engine *Engine) articleIndexWorker(jobs <-chan *Article, results chan<- *Article) {
	for id := range jobs {
		err := engine.articleIndex(id)
		if err != nil {
			engine.log.Println(id.Url, err)
		}
		results <- id
	}
}
