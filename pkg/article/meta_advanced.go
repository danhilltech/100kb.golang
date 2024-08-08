package article

import (
	"database/sql"
	"hash/fnv"
	"strings"
	"time"

	"github.com/danhilltech/100kb.golang/pkg/serialize"
)

var zeroShotLabels = []string{
	"technology",
	"life",
	"family",
	"science",
	"politics",
	"news",
	"religion",
	"god",
	"programming",
	"food",
	"crypto",
	"investing",
	"management",
	"nature",
}

func (engine *Engine) articleMetaAdvanced(txn *sql.Tx, article *Article) error {

	var t1, t2 int64

	t1 = time.Now().UnixMilli()

	// Check we have enough data
	article.LastMetaAt = time.Now().Unix()

	var feedArticles []*Article

	engine.cacheArticlesMutex.RLock()
	if engine.cacheArticles[article.Domain] != nil {
		feedArticles = engine.cacheArticles[article.Domain]
		engine.cacheArticlesMutex.RUnlock()
	} else {
		engine.cacheArticlesMutex.RUnlock()
		var err error
		feedArticles, err = engine.getArticlesByFeed(txn, article.FeedUrl, article.Url)
		if err != nil {
			return err

		}
		engine.cacheArticlesMutex.Lock()
		engine.cacheArticles[article.Domain] = feedArticles
		engine.cacheArticlesMutex.Unlock()
	}

	currentCanon := make(map[uint64]bool)

	for _, feed := range feedArticles {
		for _, para := range feed.BodyRaw.Content {
			keyHash := fnv.New64()
			keyHash.Write([]byte(strings.ToLower(para.Text)))
			currentCanon[keyHash.Sum64()] = true
		}
	}

	t2 = time.Now().UnixMilli() - t1
	engine.sd.Timing("articleMetaAdvanced.currentCanon", t2)
	t1 = time.Now().UnixMilli()

	// unique the content

	var uniqueContent []*serialize.FlatNode

	for _, line := range article.BodyRaw.Content {
		if len(line.Text) > 0 {

			keyHash := fnv.New64()
			keyHash.Write([]byte(strings.ToLower(line.Text)))

			found := currentCanon[keyHash.Sum64()]

			if !found {
				uniqueContent = append(uniqueContent, line)
			}
		}
	}

	article.Body = &serialize.Content{Content: uniqueContent}

	t2 = time.Now().UnixMilli() - t1
	engine.sd.Timing("articleMetaAdvanced.uniqueContent", t2)
	t1 = time.Now().UnixMilli()

	if len(article.Body.Content) == 0 {
		return nil
	}

	var summaryTexts []string
	hasFirstPara := false
	for _, c := range uniqueContent {
		if c.Type == "h1" || c.Type == "h2" || c.Type == "h3" {
			summaryTexts = append(summaryTexts, c.Text)
		}
		if c.Type == "p" && !hasFirstPara {
			summaryTexts = append(summaryTexts, c.Text)
		}
		if len(summaryTexts) >= 5 {
			break
		}
	}

	t2 = time.Now().UnixMilli() - t1
	engine.sd.Timing("articleMetaAdvanced.htmlTags", t2)

	if len(summaryTexts) > 0 &&
		engine.sentenceEmbeddingModel != nil &&
		engine.zeroShotModel != nil {
		t1 = time.Now().UnixMilli()
		// AI
		var startTime, diff int64
		startTime = time.Now().UnixMilli()
		vecs, err := engine.sentenceEmbeddingModel.Embeddings(summaryTexts)
		if err != nil {
			return err
		}
		diff = time.Now().UnixMilli() - startTime
		if diff > 500 {
			engine.log.Printf("SLOW sentence embedding %d %s\n", diff, article.Url)
		}

		if len(vecs) > 0 {
			article.SentenceEmbedding = &serialize.Embeddings{}
			for _, vec := range vecs {
				emnd := serialize.Embedding{Vectors: vec.Vectors}
				article.SentenceEmbedding.Embeddings = append(article.SentenceEmbedding.Embeddings, &emnd)
			}
		}

		startTime = time.Now().UnixMilli()
		ess, err := engine.keywordExtractionModel.Extract(summaryTexts)
		if err != nil {
			return err
		}
		diff = time.Now().UnixMilli() - startTime
		if diff > 500 {
			engine.log.Printf("SLOW keyword extraction %d %s\n", diff, article.Url)
		}

		kwds := map[string][]float32{}

		if len(ess) > 0 {
			for _, es := range ess {
				for _, k := range es.Keywords {
					if kwds[string(k.Text)] == nil {
						kwds[string(k.Text)] = []float32{}
					}
					kwds[string(k.Text)] = append(kwds[string(k.Text)], k.Score)

				}
			}

			for k, ss := range kwds {

				score := float32(0.0)
				for _, s := range ss {
					score += s
				}
				score = score / float32(len(ss))

				article.ExtractedKeywords.Keywords = append(article.ExtractedKeywords.Keywords, &serialize.Keyword{Text: k, Score: score})
			}

		}

		startTime = time.Now().UnixMilli()
		zcs, err := engine.zeroShotModel.Predict(summaryTexts, zeroShotLabels)
		if err != nil {
			return err
		}
		diff = time.Now().UnixMilli() - startTime
		if diff > 500 {
			engine.log.Printf("SLOW zero shot %d %s\n%+v\n", diff, article.Url, summaryTexts)
		}

		zeroshots := map[string][]float32{}

		if len(zcs) > 0 {
			article.Classifications = &serialize.Keywords{}

			for _, es := range zcs {
				for _, k := range es.Classifications {
					if zeroshots[string(k.Label)] == nil {
						zeroshots[string(k.Label)] = []float32{}
					}
					zeroshots[string(k.Label)] = append(zeroshots[string(k.Label)], k.Score)

				}
			}

			article.Classifications.Keywords = []*serialize.Keyword{}

			for k, ss := range zeroshots {

				score := float32(0.0)
				for _, s := range ss {
					score += s
				}
				score = score / float32(len(ss))

				article.Classifications.Keywords = append(article.Classifications.Keywords, &serialize.Keyword{Text: k, Score: score})
			}

		}

		t2 = time.Now().UnixMilli() - t1
		engine.sd.Timing("articleMetaAdvanced.ai", t2)
	}

	article.Stage = STAGE_COMPLETE

	return nil

}
