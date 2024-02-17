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
	wordCount = ?, 
	h1Count = ?, 
	hnCount = ?, 
	pCount = ?, 
	firstPersonRatio = ?, 
	sentenceEmbedding = ?, 
	extractedKeywords = ?, 
	lastContentExtractAt = ?, 
	badCount = ?, 
	classifications = ?,
	htmlLength = ?,
	pageAbout = ?,
	pageBlogRoll = ?,
	pageWriting = ?,
	urlNews = ?,
	urlBlog = ?,
	urlHumanName = ?,
	domainIsPopular = ?,
	domainTLD = ?,
	stage = ?
	WHERE url = ?;`)
	if err != nil {
		return err
	}

	engine.db = db
	return nil
}

func (engine *Engine) Insert(article *Article, feedUrl string, domain string, txchan *sql.Tx) error {
	_, err := txchan.Stmt(engine.dbInsertPreparedArticle).Exec(article.Url, feedUrl, domain, article.PublishedAt)
	return err
}

func (engine *Engine) Update(article *Article, txchan *sql.Tx) error {
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

	_, err = txchan.Stmt(engine.dbUpdatePreparedArticle).Exec(
		utils.NullInt64(article.LastFetchAt),
		utils.NullInt64(article.LastMetaAt),
		utils.NullString(string(articleBodyRaw)),
		utils.NullString(article.Title),
		utils.NullString(article.Description),
		utils.NullString(string(articleBody)),
		utils.NullInt64(article.WordCount),
		utils.NullInt64(article.H1Count),
		utils.NullInt64(article.HNCount),
		utils.NullInt64(article.PCount),
		utils.NullFloat64(article.FirstPersonRatio),
		utils.NullString(string(sentenceEmbedding)),
		utils.NullString(string(extractedKeywords)),
		utils.NullInt64(article.LastContentExtractAt),
		utils.NullInt64(article.BadCount),
		utils.NullString(string(classifications)),
		utils.NullInt64(article.HTMLLength),
		utils.NullBool(article.PageAbout),
		utils.NullBool(article.PageBlogRoll),
		utils.NullBool(article.PageWriting),
		utils.NullBool(article.URLNews),
		utils.NullBool(article.URLBlog),
		utils.NullBool(article.URLHumanName),
		utils.NullBool(article.DomainIsPopular),
		utils.NullString(article.DomainTLD),
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
wordCount, 
h1Count, 
hnCount, 
pCount, 
firstPersonRatio, 
lastContentExtractAt, 
badCount, 
classifications, 
htmlLength,
pageAbout,
pageBlogRoll,
pageWriting,
urlNews,
urlBlog,
urlHumanName,
domainIsPopular,
domainTLD,
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

	var wordCount, h1Count, hnCount, pCount, badCount sql.NullInt64
	var firstPersonRatio sql.NullFloat64

	var htmlLength, pageAbout, pageBlogRoll, pageWriting, urlNews, urlBlog, urlHumanName, domainIsPopular sql.NullInt64
	var domainTLD sql.NullString

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
		&wordCount,
		&h1Count,
		&hnCount,
		&pCount,
		&firstPersonRatio,
		&lastContentExtractAt,
		&badCount,
		&classificationsJSON,
		&htmlLength,
		&pageAbout,
		&pageBlogRoll,
		&pageWriting,
		&urlNews,
		&urlBlog,
		&urlHumanName,
		&domainIsPopular,
		&domainTLD,
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
		WordCount:            wordCount.Int64,
		H1Count:              h1Count.Int64,
		HNCount:              hnCount.Int64,
		PCount:               pCount.Int64,
		FirstPersonRatio:     firstPersonRatio.Float64,
		Domain:               domain,
		BadCount:             badCount.Int64,
		Classifications:      &classifications,
		HTMLLength:           htmlLength.Int64,
		PageAbout:            pageAbout.Int64 > 0,
		PageBlogRoll:         pageBlogRoll.Int64 > 0,
		PageWriting:          pageWriting.Int64 > 0,
		URLNews:              urlNews.Int64 > 0,
		URLBlog:              urlBlog.Int64 > 0,
		URLHumanName:         urlHumanName.Int64 > 0,
		DomainIsPopular:      domainIsPopular.Int64 > 0,
		DomainTLD:            domainTLD.String,
		Stage:                stage.Int64,
	}

	return article, nil
}

func (engine *Engine) getArticlesToIndex(txchan *sql.Tx) ([]*Article, error) {
	fmt.Printf("Getting articles to index...\t")
	defer fmt.Printf("✨\n")
	res, err := txchan.Query(fmt.Sprintf("SELECT %s FROM articles WHERE lastFetchAt IS NULL;", ARTICLE_SELECT))
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

func (engine *Engine) getArticlesToContentExtract(txchan *sql.Tx) ([]*Article, error) {
	res, err := txchan.Query(fmt.Sprintf("SELECT %s FROM articles WHERE lastFetchAt > 0 AND lastContentExtractAt IS NULL;", ARTICLE_SELECT))
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

func (engine *Engine) getArticlesToMetaDataAdvanved(txchan *sql.Tx) ([]*Article, error) {
	res, err := txchan.Query(fmt.Sprintf("SELECT %s FROM articles WHERE lastContentExtractAt > 0 AND lastMetaAt IS NULL AND stage >= 2;", ARTICLE_SELECT))
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
	res, err := txchan.Query(fmt.Sprintf("SELECT %s FROM articles WHERE feedUrl = ? AND bodyRaw IS NOT NULL AND url != ?", ARTICLE_SELECT), feed, excludeUrl)

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

func (engine *Engine) GetAllValid(txchan *sql.Tx) ([]*Article, error) {
	res, err := txchan.Query(fmt.Sprintf("SELECT %s FROM articles WHERE lastMetaAt IS NOT NULL AND wordCount > 10", ARTICLE_SELECT))

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

func (engine *Engine) FindByURL(txchan *sql.Tx, url string) (*Article, error) {
	res, err := txchan.Query(fmt.Sprintf("SELECT %s FROM articles WHERE url = ? LIMIT 1", ARTICLE_SELECT), url)

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
