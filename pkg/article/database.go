package article

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/danhilltech/100kb.golang/pkg/utils"
)

func (engine *Engine) initDB(db *sql.DB) error {
	var err error
	engine.dbInsertPreparedArticle, err = db.Prepare("INSERT INTO articles(url, feedUrl, publishedAt) VALUES(?, ?, ?) ON CONFLICT(url) DO NOTHING;")
	if err != nil {
		return err
	}

	engine.dbUpdatePreparedArticle, err = db.Prepare("UPDATE articles SET lastFetchAt = ?, bodyRaw = ?, title = ?, description = ?, body = ?, wordCount = ?, firstPersonRatio = ?, sentenceEmbedding = ?, extractedKeywords = ? WHERE url = ?;")
	if err != nil {
		return err
	}

	engine.db = db
	return nil
}

func (engine *Engine) Insert(article *Article, feedUrl string, txchan *sql.Tx) error {
	_, err := txchan.Stmt(engine.dbInsertPreparedArticle).Exec(article.Url, feedUrl, article.PublishedAt)
	return err
}

func (engine *Engine) Update(article *Article, txchan *sql.Tx) error {
	var articleBodyRaw []byte
	var articleBody []byte
	var extractedKeywords []byte
	var sentenceEmbedding []byte
	var err error

	if len(article.BodyRaw) > 0 {
		articleBodyRaw, err = json.Marshal(article.BodyRaw)
		if err != nil {
			return err
		}
	}

	if len(article.Body) > 0 {
		articleBody, err = json.Marshal(article.Body)
		if err != nil {
			return err
		}
	}

	if len(article.SentenceEmbedding) > 0 {
		sentenceEmbedding, err = json.Marshal(article.SentenceEmbedding)
		if err != nil {
			return err
		}
	}
	if len(article.ExtractedKeywords) > 0 {
		extractedKeywords, err = json.Marshal(article.ExtractedKeywords)
		if err != nil {
			return err
		}
	}

	_, err = txchan.Stmt(engine.dbUpdatePreparedArticle).Exec(
		utils.NullInt64(article.LastFetchAt),
		utils.NullString(string(articleBodyRaw)),
		utils.NullString(article.Title),
		utils.NullString(article.Description),
		utils.NullString(string(articleBody)),
		utils.NullInt64(article.WordCount),
		utils.NullFloat64(article.FirstPersonRatio),
		utils.NullString(string(sentenceEmbedding)),
		utils.NullString(string(extractedKeywords)),
		article.Url,
	)
	return err
}

func (engine *Engine) getArticlesToIndex(txchan *sql.Tx) ([]*Article, error) {
	fmt.Printf("Getting articles to index...\t")
	defer fmt.Printf("✨\n")
	res, err := txchan.Query("SELECT url FROM articles WHERE lastFetchAt IS NULL;")
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var urls []*Article

	for res.Next() {
		var url string
		err = res.Scan(&url)
		if err != nil {
			return nil, err
		}

		urls = append(urls, &Article{Url: url})
	}
	if err := res.Err(); err != nil {
		return nil, err
	}
	return urls, nil
}

const ARTICLE_SELECT = `url, feedUrl, publishedAt, lastFetchAt, title, description, bodyRaw, body, sentenceEmbedding, extractedKeywords`

func articleRowScan(res *sql.Rows) (*Article, error) {
	var url string
	var feedUrl string
	var publishedAt int64
	var lastFetchAt sql.NullInt64
	var title sql.NullString
	var description sql.NullString
	var bodyRawJSON []byte
	var bodyJSON []byte
	var sentenceEmbeddingJSON []byte
	var extractedKeywordsJSON []byte

	err := res.Scan(
		&url,
		&feedUrl,
		&publishedAt,
		&lastFetchAt,
		&title,
		&description,
		&bodyRawJSON,
		&bodyJSON,
		&sentenceEmbeddingJSON,
		&extractedKeywordsJSON,
	)
	if err != nil {
		return nil, err
	}

	var bodyRaw []string
	if bodyRawJSON != nil {
		err = json.Unmarshal(bodyRawJSON, &bodyRaw)
		if err != nil {
			return nil, err
		}
	}

	var body []string
	if bodyJSON != nil {
		err = json.Unmarshal(bodyJSON, &body)
		if err != nil {
			return nil, err
		}
	}

	var sentenceEmbeding []float32
	if sentenceEmbeddingJSON != nil {
		err = json.Unmarshal(sentenceEmbeddingJSON, &sentenceEmbeding)
		if err != nil {
			return nil, err
		}
	}

	var extractedKeywords []*Keyword
	if extractedKeywordsJSON != nil {
		err = json.Unmarshal(extractedKeywordsJSON, &extractedKeywords)
		if err != nil {
			return nil, err
		}
	}

	article := &Article{
		Url:               url,
		FeedUrl:           feedUrl,
		PublishedAt:       publishedAt,
		LastFetchAt:       lastFetchAt.Int64,
		Title:             title.String,
		Description:       description.String,
		BodyRaw:           bodyRaw,
		Body:              body,
		SentenceEmbedding: sentenceEmbeding,
		ExtractedKeywords: extractedKeywords,
	}

	return article, nil
}

func (engine *Engine) getArticlesToMetaData(txchan *sql.Tx) ([]*Article, error) {
	res, err := txchan.Query(fmt.Sprintf("SELECT %s FROM articles WHERE body IS NULL AND bodyRaw IS NOT NULL AND lastFetchAt > 0;", ARTICLE_SELECT))
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var urls []*Article

	for res.Next() {

		article, err := articleRowScan(res)
		if err != nil {
			return nil, err
		}
		urls = append(urls, article)
	}
	if err := res.Err(); err != nil {
		return nil, err
	}
	return urls, nil
}

func (engine *Engine) getArticlesByFeed(txchan *sql.Tx, feed string, excludeUrl string) ([]*Article, error) {
	res, err := txchan.Query(fmt.Sprintf("SELECT %s FROM articles WHERE feedUrl = ? AND url != ?", ARTICLE_SELECT), feed, excludeUrl)

	if err != nil {
		return nil, err
	}
	defer res.Close()
	if err := res.Err(); err != nil {
		return nil, err
	}

	var urls []*Article

	for res.Next() {
		article, err := articleRowScan(res)
		if err != nil {
			return nil, err
		}
		urls = append(urls, article)
	}
	if err := res.Err(); err != nil {
		return nil, err
	}
	return urls, nil
}
