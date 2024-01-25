package article

import (
	"database/sql"
	"net/http"

	"github.com/danhilltech/100kb.golang/pkg/ai"
	retryhttp "github.com/danhilltech/100kb.golang/pkg/http"
	"github.com/danhilltech/100kb.golang/pkg/parsing"
	"github.com/danhilltech/100kb.golang/pkg/serialize"
	"github.com/danhilltech/100kb.golang/pkg/utils"
)

type Engine struct {
	dbInsertPreparedArticle *sql.Stmt
	dbUpdatePreparedArticle *sql.Stmt
	db                      *sql.DB

	sentenceEmbeddingModel *ai.SentenceEmbeddingModel
	keywordExtractionModel *ai.KeywordExtractionModel
	parser                 *parsing.Engine

	// aiMutex sync.Mutex

	http  *http.Client
	cache *utils.Cache
}

type Article struct {
	Url                  string
	FeedUrl              string
	PublishedAt          int64
	BodyRaw              *serialize.Content
	LastFetchAt          int64
	LastMetaAt           int64
	LastContentExtractAt int64
	Title                string
	Description          string

	Body              *serialize.Content
	WordCount         int64
	H1Count           int64
	HNCount           int64
	PCount            int64
	BadCount          int64
	FirstPersonRatio  float64
	SentenceEmbedding *serialize.Embeddings
	ExtractedKeywords *serialize.Keywords

	HumanClassification int64
}

func NewEngine(db *sql.DB, cachePath string) (*Engine, error) {
	engine := Engine{}
	var err error

	err = engine.initDB(db)
	if err != nil {
		return nil, err
	}

	engine.http = retryhttp.NewRetryableClient()

	engine.sentenceEmbeddingModel, err = ai.NewSentenceEmbeddingModel()
	if err != nil {
		return nil, err
	}

	engine.keywordExtractionModel, err = ai.NewKeywordExtractionModel()
	if err != nil {
		return nil, err
	}

	engine.parser, err = parsing.NewEngine()
	if err != nil {
		return nil, err
	}

	engine.cache = utils.NewDiskCache(cachePath)

	return &engine, nil
}

func (engine *Engine) Close() {
	if engine.sentenceEmbeddingModel != nil {
		engine.sentenceEmbeddingModel.Close()
	}
	if engine.keywordExtractionModel != nil {
		engine.keywordExtractionModel.Close()
	}

	if engine.parser != nil {
		engine.parser.Close()
	}
}
