package article

import (
	"database/sql"
	"fmt"
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

	chunkIds := Chunk(articles, chunkSize)

	fmt.Printf("Generating %d article metas %d chunks\n", len(articles), len(chunkIds))

	for _, chunk := range chunkIds {
		err = engine.doArticleMeta(chunk, workers)
		if err != nil {
			fmt.Println(err)
		}
	}

	return nil
}

func (engine *Engine) doArticleMeta(chunk []*Article, workers int) error {
	fmt.Printf("Chunk...\t\t")
	defer fmt.Printf("✅\n")

	insertTxn, err := engine.db.Begin()
	if err != nil {
		return err
	}
	defer insertTxn.Rollback()

	err = engine.articleMetas(insertTxn, chunk, workers)
	if err != nil {
		return err
	}

	for _, article := range chunk {
		// save it
		err = engine.Update(article, insertTxn)
		if err != nil {
			fmt.Println(article.Url, err)
		}
	}

	err = insertTxn.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (engine *Engine) articleMetaWorker(tx *sql.Tx, jobs <-chan *Article, results chan<- *Article) {
	for id := range jobs {
		err := engine.articleExtractContent(tx, id)
		if err != nil {
			fmt.Println(err)
		}
		results <- id
	}
}

// Crawls

func (engine *Engine) articleMetas(tx *sql.Tx, articles []*Article, workers int) error {

	jobs := make(chan *Article, len(articles))
	results := make(chan *Article, len(articles))

	for w := 1; w <= workers; w++ {
		go engine.articleMetaWorker(tx, jobs, results)
	}

	for j := 1; j <= len(articles); j++ {
		jobs <- articles[j-1]
	}
	close(jobs)

	items := make([]*Article, len(articles))

	for a := 1; a <= len(articles); a++ {
		b := <-results
		items[a-1] = b
	}

	return nil

}
