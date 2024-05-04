package article

import (
	"fmt"
	"net/url"
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

func (engine *Engine) BuildArticleSingle(urlRaw string) (*Article, error) {

	urlP, err := url.Parse(urlRaw)
	if err != nil {
		return nil, fmt.Errorf("could not parse url: %w", err)
	}

	article := &Article{Url: urlRaw, Domain: urlP.Hostname()}

	err = engine.articleIndex(article)
	if err != nil {
		return nil, fmt.Errorf("could not articleIndex: %w", err)
	}

	err = engine.articleExtractContent(article)
	if err != nil {
		return nil, fmt.Errorf("could not articleExtractContent: %w", err)
	}

	err = engine.articleMetaAdvanced(article)
	if err != nil {
		return nil, fmt.Errorf("could not articleMetaAdvanced: %w", err)
	}

	return article, nil
}
