package article

import (
	"database/sql"
	"net/http"

	"github.com/danhilltech/100kb.golang/pkg/ai"
	retryhttp "github.com/danhilltech/100kb.golang/pkg/http"
	"github.com/danhilltech/100kb.golang/pkg/serialize"
)

type Engine struct {
	dbInsertPreparedArticle *sql.Stmt
	dbUpdatePreparedArticle *sql.Stmt
	db                      *sql.DB

	sentenceEmbeddingModel *ai.SentenceEmbeddingModel
	keywordExtractionModel *ai.KeywordExtractionModel

	// aiMutex sync.Mutex

	http *http.Client
}

type Article struct {
	Url         string
	FeedUrl     string
	PublishedAt int64
	Html        []byte
	BodyRaw     *serialize.Content
	LastFetchAt int64
	LastMetaAt  int64
	Title       string
	Description string

	Body              *serialize.Content
	WordCount         int64
	H1Count           int64
	HNCount           int64
	PCount            int64
	FirstPersonRatio  float64
	SentenceEmbedding *serialize.Embeddings
	ExtractedKeywords *serialize.Keywords
}

func NewEngine(db *sql.DB) (*Engine, error) {
	engine := Engine{}

	engine.initDB(db)

	engine.http = retryhttp.NewRetryableClient()

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
