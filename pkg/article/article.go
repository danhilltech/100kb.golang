package article

import (
	"database/sql"
	"sync"

	"github.com/danhilltech/100kb.golang/pkg/ai"
)

type Engine struct {
	dbInsertPreparedArticle *sql.Stmt
	dbUpdatePreparedArticle *sql.Stmt
	db                      *sql.DB

	sentenceEmbeddingModel *ai.SentenceEmbeddingModel
	keywordExtractionModel *ai.KeywordExtractionModel

	aiMutex sync.Mutex
}

type Keyword struct {
	Text  string
	Score float32
}

type Article struct {
	Url         string
	FeedUrl     string
	PublishedAt int64
	Html        []byte
	BodyRaw     []string
	LastFetchAt int64
	Title       string
	Description string

	Body              []string
	WordCount         int64
	FirstPersonRatio  float64
	SentenceEmbedding []float32
	ExtractedKeywords []Keyword
}

func NewEngine(db *sql.DB) (*Engine, error) {
	engine := Engine{}

	engine.initDB(db)

	var err error

	engine.sentenceEmbeddingModel, err = ai.NewSentenceEmbeddingModel()
	if err != nil {
		return nil, err
	}

	engine.keywordExtractionModel, err = ai.NewKeywordExtractionModel()
	if err != nil {
		return nil, err
	}

	return &engine, nil
}

func (engine *Engine) Close() {
	if engine.sentenceEmbeddingModel != nil {
		engine.sentenceEmbeddingModel.Close()
	}
	if engine.keywordExtractionModel != nil {
		engine.keywordExtractionModel.Close()
	}
}
