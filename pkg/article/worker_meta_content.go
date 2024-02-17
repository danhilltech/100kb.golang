package article

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func (engine *Engine) RunArticleMeta(chunkSize int, workers int) error {
	txn, err := engine.db.Begin()
	if err != nil {
		return err
	}
	defer txn.Rollback()

	articles, err := engine.getArticlesToContentExtract(txn)
	if err != nil {
		return err
	}
	err = txn.Commit()
	if err != nil {
		return err
	}

	fmt.Printf("Generating %d article metas\n", len(articles))

	chunks := Chunk(articles, runtime.NumCPU()-2)

	var wg sync.WaitGroup

	for i := 1; i <= len(chunks); i++ {
		go func() {
			defer wg.Done()
			engine.runArticleMetaBatch(chunks[i], chunkSize)
		}()
	}

	wg.Wait()
	return nil

}

func (engine *Engine) runArticleMetaBatch(articles []*Article, chunkSize int) error {
	txn, err := engine.db.Begin()
	if err != nil {
		return err
	}
	defer txn.Rollback()

	a := 0
	t := time.Now().UnixMilli()

	for _, article := range articles {
		a++
		err := engine.articleExtractContent(txn, article)
		if err != nil {
			fmt.Println(err)
			err = engine.Update(article, txn)
			if err != nil {
				fmt.Println(article.Url, err)
			}

			continue
		}

		err = engine.Update(article, txn)
		if err != nil {
			fmt.Println(article.Url, err)
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
	fmt.Printf("\tdone %d/%d\n", a, len(articles))

	err = txn.Commit()
	if err != nil {
		return err
	}

	return nil
}
