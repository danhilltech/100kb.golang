package article

import (
	"fmt"
	"sync"
	"time"
)

func (engine *Engine) RunArticleMeta(chunkSize int) error {

	articles, err := engine.getArticlesToContentExtract()
	if err != nil {
		return err
	}

	fmt.Printf("Generating %d article metas\n", len(articles))

	a := 0
	lastA := 0
	t := time.Now().UnixMilli()

	var wg sync.WaitGroup

	printSize := 200

	for _, article := range articles {

		wg.Add(1)

		go func(article *Article) {
			defer wg.Done()

			err := engine.articleExtractContent(article)
			if err != nil {
				// fmt.Println(article.Url, err)
				err = engine.Update(article)
				if err != nil {
					fmt.Println(article.Url, err)
				}

				return
			}

			err = engine.Update(article)
			if err != nil {
				fmt.Println(article.Url, err)
				return
			}
		}(article)

		if a > 0 && a%chunkSize == 0 {
			wg.Wait()
		}
		if a > 0 && a%printSize == 0 {
			diff := time.Now().UnixMilli() - t
			qps := (float64(a-lastA) / float64(diff)) * 1000
			lastA = a
			t = time.Now().UnixMilli()
			fmt.Printf("\tdone %d/%d at %0.2f/s\n", a, len(articles), qps)

		}
		a++
	}
	wg.Wait()
	fmt.Printf("\tdone %d/%d\n", a, len(articles))

	return nil

}
