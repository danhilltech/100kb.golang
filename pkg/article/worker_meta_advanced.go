package article

import (
	"fmt"
	"time"
)

func (engine *Engine) RunArticleMetaPassII(chunkSize int) error {

	articles, err := engine.getArticlesToMetaDataAdvanved()
	if err != nil {
		return err
	}

	fmt.Printf("Generating %d article advanced metas\n", len(articles))

	printSize := 100

	a := 0
	t := time.Now().UnixMilli()
	txn, _ := engine.db.Begin()
	for _, article := range articles {

		err := engine.articleMetaAdvanced(txn, article)
		if err != nil {
			fmt.Println(err)
			continue
		}

		err = engine.Update(txn, article)
		if err != nil {
			fmt.Println(article.Url, err)
			continue
		}

		if a > 0 && a%printSize == 0 {
			err := txn.Commit()
			if err != nil {
				return err
			}
			txn, _ = engine.db.Begin()
			diff := time.Now().UnixMilli() - t
			qps := (float64(printSize) / float64(diff)) * 1000
			t = time.Now().UnixMilli()
			fmt.Printf("\tdone %d/%d at %0.2f/s\n", a, len(articles), qps)

		}
		a++
	}
	err = txn.Commit()
	if err != nil {
		return err
	}
	fmt.Printf("\tdone %d/%d\n\n", a, len(articles))

	return nil
}
