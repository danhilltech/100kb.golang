package article

import (
	"database/sql"
	"strings"
	"time"

	"github.com/danhilltech/100kb.golang/pkg/serialize"
)

func (engine *Engine) articleMetaAdvanced(tx *sql.Tx, article *Article) error {
	// Check we have enough data
	article.LastMetaAt = time.Now().Unix()

	feedArticles, err := engine.getArticlesByFeed(tx, article.FeedUrl, article.Url)
	if err != nil {
		return err
	}

	var currentCanon []string

	for _, feed := range feedArticles {
		for _, para := range feed.BodyRaw.Content {
			currentCanon = append(currentCanon, strings.ToLower(para.Text))
		}
	}

	// unique the content

	var uniqueContent []*serialize.FlatNode

	for _, line := range article.BodyRaw.Content {
		if len(line.Text) > 0 {
			found := false
			for _, currLine := range currentCanon {
				if strings.ToLower(line.Text) == currLine {
					found = true
				}
			}
			if !found {
				uniqueContent = append(uniqueContent, line)
			}
		}
	}

	article.Body = &serialize.Content{Content: uniqueContent}

	if len(article.Body.Content) == 0 {
		return nil
	}

	bodyBuilder := strings.Builder{}

	for _, para := range uniqueContent {
		bodyBuilder.WriteString(para.Text)
		bodyBuilder.WriteRune('\n')
	}

	bodyConcat := bodyBuilder.String()

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

	var h1Count, hnCount, pCount int64

	for _, sec := range uniqueContent {
		switch sec.Type {
		case "h1":
			h1Count++
			continue
		case "h2", "h3":
			hnCount++
			continue
		case "p", "li":
			pCount++
			continue
		}
	}

	article.H1Count = h1Count
	article.HNCount = hnCount
	article.PCount = pCount

	var firstPara string
	for _, c := range uniqueContent {
		if c.Type == "p" && len(c.Text) >= 150 {
			firstPara = c.Text
			break
		}
	}
	if len(firstPara) >= 150 {
		// AI
		vec, err := engine.sentenceEmbeddingModel.Embeddings([]string{firstPara})
		if err != nil {
			return err
		}

		if len(vec) > 0 {
			emnd := serialize.Embedding{Vectors: vec[0].Vectors}
			article.SentenceEmbedding = &serialize.Embeddings{}
			article.SentenceEmbedding.Embeddings = append(article.SentenceEmbedding.Embeddings, &emnd)
		}

		es, err := engine.keywordExtractionModel.Extract([]string{firstPara})
		if err != nil {
			return err
		}

		// article.ExtractedKeywords = es[0].Keywords
		if len(es) > 0 {
			for _, k := range es[0].Keywords {
				article.ExtractedKeywords.Keywords = append(article.ExtractedKeywords.Keywords, &serialize.Keyword{Text: string(k.Text), Score: k.Score})
			}
		}
	}

	return nil

}
