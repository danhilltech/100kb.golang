package article

import (
	"database/sql"
	"fmt"

	"github.com/danhilltech/100kb.golang/pkg/serialize"
	"github.com/danhilltech/100kb.golang/pkg/utils"
	"google.golang.org/protobuf/proto"
)

func (engine *Engine) initDB(db *sql.DB) error {
	var err error
	engine.dbInsertPreparedArticle, err = db.Prepare("INSERT INTO articles(url, feedUrl, domain, publishedAt) VALUES(?, ?, ?, ?) ON CONFLICT(url) DO NOTHING;")
	if err != nil {
		return err
	}

	engine.dbUpdatePreparedArticle, err = db.Prepare(`
	UPDATE articles SET 
	lastFetchAt = ?, 
	lastMetaAt = ?, 
	bodyRaw = ?, 
	title = ?, 
	description = ?, 
	body = ?, 
	sentenceEmbedding = ?, 
	extractedKeywords = ?, 
	lastContentExtractAt = ?, 
	badCount = ?, 
	badElementCount = ?,
	linkCount = ?,
	badLinkCount = ?,
	classifications = ?,
	htmlLength = ?,
	stage = ?
	WHERE url = ?;`)
	if err != nil {
		return err
	}

	engine.db = db
	return nil
}

func (engine *Engine) Insert(txn *sql.Tx, article *Article, feedUrl string, domain string) error {

	_, err := txn.Stmt(engine.dbInsertPreparedArticle).Exec(article.Url, feedUrl, domain, article.PublishedAt)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	return nil
}

func (engine *Engine) Update(txn *sql.Tx, article *Article) error {

	var articleBodyRaw []byte
	var articleBody []byte
	var extractedKeywords []byte
	var sentenceEmbedding []byte
	var classifications []byte
	var err error

	if article.BodyRaw != nil {
		articleBodyRaw, err = proto.Marshal(article.BodyRaw)
		if err != nil {
			return err
		}
	}

	if article.Body != nil {
		articleBody, err = proto.Marshal(article.Body)
		if err != nil {
			return err
		}
	}

	if article.SentenceEmbedding != nil {
		sentenceEmbedding, err = proto.Marshal(article.SentenceEmbedding)
		if err != nil {
			return err
		}
	}
	if article.ExtractedKeywords != nil {
		extractedKeywords, err = proto.Marshal(article.ExtractedKeywords)
		if err != nil {
			return err
		}
	}
	if article.Classifications != nil {
		classifications, err = proto.Marshal(article.Classifications)
		if err != nil {
			return err
		}
	}

	_, err = txn.Stmt(engine.dbUpdatePreparedArticle).Exec(
		utils.NullInt64(article.LastFetchAt),
		utils.NullInt64(article.LastMetaAt),
		utils.NullString(string(articleBodyRaw)),
		utils.NullString(article.Title),
		utils.NullString(article.Description),
		utils.NullString(string(articleBody)),
		utils.NullString(string(sentenceEmbedding)),
		utils.NullString(string(extractedKeywords)),
		utils.NullInt64(article.LastContentExtractAt),
		utils.NullInt64(article.BadCount),
		utils.NullInt64(article.BadElementCount),
		utils.NullInt64(article.LinkCount),
		utils.NullInt64(article.BadLinkCount),
		utils.NullString(string(classifications)),
		utils.NullInt64(article.HTMLLength),
		utils.NullInt64(article.Stage),
		article.Url,
	)
	return err
}

const ARTICLE_SELECT = `url, 
feedUrl, 
domain, 
publishedAt, 
lastFetchAt, 
lastMetaAt, 
title, 
description, 
bodyRaw, 
body, 
sentenceEmbedding, 
extractedKeywords,
lastContentExtractAt, 
badCount, 
badElementCount,
linkCount,
badLinkCount,
classifications, 
htmlLength,
stage`

func articleRowScan(res *sql.Rows) (*Article, error) {
	var url string
	var feedUrl string
	var domain string
	var publishedAt int64
	var lastFetchAt sql.NullInt64
	var lastMetaAt sql.NullInt64
	var lastContentExtractAt sql.NullInt64
	var title sql.NullString
	var description sql.NullString
	var bodyRawJSON []byte
	var bodyJSON []byte
	var sentenceEmbeddingJSON []byte
	var extractedKeywordsJSON []byte
	var classificationsJSON []byte

	var badCount, badElementCount, linkCount, badLinkCount sql.NullInt64

	var htmlLength sql.NullInt64

	var stage sql.NullInt64

	err := res.Scan(
		&url,
		&feedUrl,
		&domain,
		&publishedAt,
		&lastFetchAt,
		&lastMetaAt,
		&title,
		&description,
		&bodyRawJSON,
		&bodyJSON,
		&sentenceEmbeddingJSON,
		&extractedKeywordsJSON,
		&lastContentExtractAt,
		&badCount,
		&badElementCount,
		&linkCount,
		&badLinkCount,
		&classificationsJSON,
		&htmlLength,
		&stage,
	)
	if err != nil {
		return nil, err
	}

	var bodyRaw serialize.Content
	if bodyRawJSON != nil {
		err = proto.Unmarshal(bodyRawJSON, &bodyRaw)
		if err != nil {
			return nil, err
		}
	}

	var body serialize.Content
	if bodyJSON != nil {
		err = proto.Unmarshal(bodyJSON, &body)
		if err != nil {
			return nil, err
		}
	}

	var sentenceEmbeding serialize.Embeddings
	if sentenceEmbeddingJSON != nil {
		err = proto.Unmarshal(sentenceEmbeddingJSON, &sentenceEmbeding)
		if err != nil {
			return nil, err
		}
	}

	var extractedKeywords serialize.Keywords
	if extractedKeywordsJSON != nil {
		err = proto.Unmarshal(extractedKeywordsJSON, &extractedKeywords)
		if err != nil {
			return nil, err
		}
	}

	var classifications serialize.Keywords
	if classificationsJSON != nil {
		err = proto.Unmarshal(classificationsJSON, &classifications)
		if err != nil {
			return nil, err
		}
	}

	article := &Article{
		Url:                  url,
		FeedUrl:              feedUrl,
		PublishedAt:          publishedAt,
		LastFetchAt:          lastFetchAt.Int64,
		LastMetaAt:           lastMetaAt.Int64,
		LastContentExtractAt: lastContentExtractAt.Int64,
		Title:                title.String,
		Description:          description.String,
		BodyRaw:              &bodyRaw,
		Body:                 &body,
		SentenceEmbedding:    &sentenceEmbeding,
		ExtractedKeywords:    &extractedKeywords,
		Domain:               domain,
		BadCount:             badCount.Int64,
		BadElementCount:      badElementCount.Int64,
		LinkCount:            linkCount.Int64,
		BadLinkCount:         badLinkCount.Int64,
		Classifications:      &classifications,
		HTMLLength:           htmlLength.Int64,
		Stage:                stage.Int64,
	}

	return article, nil
}

func (engine *Engine) getArticlesToIndex() ([]*Article, error) {
	engine.log.Printf("Getting articles to index...\t")
	defer engine.log.Printf("✨\n")
	res, err := engine.db.Query(fmt.Sprintf("SELECT %s FROM articles WHERE lastFetchAt IS NULL OR lastFetchAt < %d ORDER BY lastFetchAt ASC LIMIT %d;", ARTICLE_SELECT, REFRESH_AGO_SECONDS, REFRESH_LIMIT))
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

func (engine *Engine) getArticlesToContentExtract() ([]*Article, error) {
	res, err := engine.db.Query(fmt.Sprintf("SELECT %s FROM articles WHERE stage=%d AND lastContentExtractAt IS NULL OR lastContentExtractAt < %d ORDER BY lastContentExtractAt ASC LIMIT %d;", ARTICLE_SELECT, STAGE_INDEXED, CONTENT_EXTRACT_AGO_SECONDS, CONTENT_EXTRACT_LIMIT))
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

func (engine *Engine) getArticlesToMetaDataAdvanved() ([]*Article, error) {
	res, err := engine.db.Query(fmt.Sprintf("SELECT %s FROM articles WHERE lastContentExtractAt > 0 AND lastMetaAt IS NULL AND stage = %d;", ARTICLE_SELECT, STAGE_VALID_CONTENT))
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

func (engine *Engine) getArticlesByFeed(txn *sql.Tx, feed string, excludeUrl string) ([]*Article, error) {
	res, err := txn.Query(fmt.Sprintf("SELECT %s FROM articles WHERE feedUrl = ? AND bodyRaw IS NOT NULL AND url != ? ORDER BY publishedAt DESC LIMIT 10", ARTICLE_SELECT), feed, excludeUrl)
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

func (engine *Engine) GetAllValid() ([]*Article, error) {
	res, err := engine.db.Query(fmt.Sprintf("SELECT %s FROM articles WHERE stage = %d", ARTICLE_SELECT, STAGE_COMPLETE))

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

func (engine *Engine) FindByURL(url string) (*Article, error) {
	res, err := engine.db.Query(fmt.Sprintf("SELECT %s FROM articles WHERE url = ? LIMIT 1", ARTICLE_SELECT), url)

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

	if len(urls) != 1 {
		return nil, nil
	}

	return urls[0], nil
}

func (engine *Engine) FindByFeedURL(feed string) ([]*Article, error) {
	res, err := engine.db.Query(fmt.Sprintf("SELECT %s FROM articles WHERE feedUrl = ? AND bodyRaw IS NOT NULL", ARTICLE_SELECT), feed)

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
