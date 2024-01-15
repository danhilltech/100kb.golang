package article

import (
	"database/sql"
	"fmt"
	"strings"
)

func (engine *Engine) articleMeta(tx *sql.Tx, article *Article) error {
	// Check we have enough data
	if article.BodyRaw == nil || len(article.BodyRaw) == 0 {
		return nil
	}

	feedArticles, err := engine.getArticlesByFeed(tx, article.FeedUrl, article.Url)
	if err != nil {
		return err
	}

	var currentCanon []string

	for _, feed := range feedArticles {
		currentCanon = append(currentCanon, feed.BodyRaw...)
	}

	// unique the content

	var uniqueContent []string

	for _, line := range article.BodyRaw {
		found := false
		for _, currLine := range currentCanon {
			if line == currLine {
				found = true
			}
		}
		if !found && len(line) > 0 {
			uniqueContent = append(uniqueContent, line)
		}
	}

	article.Body = uniqueContent

	if len(article.Body) == 0 {
		return nil
	}

	bodyConcat := strings.Join(uniqueContent, " ")

	// Word count
	article.WordCount = int64(len(strings.Split(bodyConcat, " ")))

	firstPersonCount := 0

	firstPersonCount += strings.Count(bodyConcat, "I ")
	firstPersonCount += strings.Count(bodyConcat, " my ")
	firstPersonCount += strings.Count(bodyConcat, " My ")
	firstPersonCount += strings.Count(bodyConcat, " me ")
	firstPersonCount += strings.Count(bodyConcat, " mine ")
	firstPersonCount += strings.Count(bodyConcat, " we ")
	firstPersonCount += strings.Count(bodyConcat, " We ")
	firstPersonCount += strings.Count(bodyConcat, " us ")
	firstPersonCount += strings.Count(bodyConcat, " our ")
	firstPersonCount += strings.Count(bodyConcat, " Our ")

	if article.WordCount > 0 && firstPersonCount > 0 {
		article.FirstPersonRatio = float64(firstPersonCount) / float64(article.WordCount)
	} else {
		article.FirstPersonRatio = 0
	}

	engine.aiMutex.Lock()
	defer engine.aiMutex.Unlock()

	// AI
	vec, err := engine.sentenceEmbeddingModel.Embeddings([]string{bodyConcat})
	if err != nil {
		return err
	}

	article.SentenceEmbedding = vec[0].Vectors

	es, err := engine.keywordExtractionModel.Extract([]string{bodyConcat})
	if err != nil {
		return err
	}

	fmt.Println(es[0].Keywords)

	article.ExtractedKeywords = es[0].Keywords

	return nil

}
