package article

import (
	"database/sql"
	"fmt"
)

func Chunk(slice []*Article, chunkSize int) [][]*Article {
	var chunks [][]*Article
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize

		// necessary check to avoid slicing beyond
		// slice capacity
		if end > len(slice) {
			end = len(slice)
		}

		chunks = append(chunks, slice[i:end])
	}

	return chunks
}

func (engine *Engine) BuildArticleSingle(txn *sql.Tx, url string) (*Article, error) {
	article := &Article{Url: url}

	resp, err := engine.articleIndex(article)
	if err != nil {
		return nil, fmt.Errorf("could not articleIndex: %w", err)
	}

	if resp.Response != nil {

		defer resp.Response.Body.Close()
	}

	err = engine.articleExtractContent(txn, article)
	if err != nil {
		return nil, fmt.Errorf("could not articleExtractContent: %w", err)
	}

	err = engine.articleMetaAdvanced(txn, article)
	if err != nil {
		return nil, fmt.Errorf("could not articleMetaAdvanced: %w", err)
	}

	return article, nil
}
