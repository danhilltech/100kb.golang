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

	// chunkSizeNew := float64(len(articles)) / float64(runtime.NumCPU()-2)

	// chunks := Chunk(articles, int(chunkSizeNew))

	// fmt.Printf("Generating %d article metas\n", len(articles))

	//

	// for _, chunk := range chunks {
	// 	wg.Add(1)
	// 	go func(chunk []*Article) {
	// 		defer wg.Done()
	// 		err := engine.runArticleMetaBatch(chunk, chunkSize)
	// 		if err != nil {
	// 			fmt.Println("runArticleMetaBatch error %w", err)
	// 		}
	// 	}(chunk)
	// }

	// wg.Wait()

	txn, err = engine.db.Begin()
	if err != nil {
		return err
	}
	defer txn.Rollback()

	a := 0
	lastA := 0
	t := time.Now().UnixMilli()

	var wg sync.WaitGroup

	atOnce := runtime.NumCPU() - 2

	for _, article := range articles {

		wg.Add(1)

		go func(article *Article) {
			defer wg.Done()
			a++
			err := engine.articleExtractContent(txn, article)
			if err != nil {
				fmt.Println(err)
				err = engine.Update(article, txn)
				if err != nil {
					fmt.Println(article.Url, err)
				}

				return
			}

			err = engine.Update(article, txn)
			if err != nil {
				fmt.Println(article.Url, err)
				return
			}
		}(article)

		if a > 0 && a%atOnce == 0 {
			wg.Wait()
		}

		if a > 0 && a%chunkSize == 0 {
			wg.Wait()
			diff := time.Now().UnixMilli() - t
			qps := (float64(a-lastA) / float64(diff)) * 1000
			lastA = a
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

func (engine *Engine) runArticleMetaBatch(articles []*Article, chunkSize int) error {

	return nil
}
