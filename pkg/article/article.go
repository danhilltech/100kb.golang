package article

import (
	"database/sql"

	"github.com/danhilltech/100kb.golang/pkg/ai"
	retryhttp "github.com/danhilltech/100kb.golang/pkg/http"
	"github.com/danhilltech/100kb.golang/pkg/parsing"
	"github.com/danhilltech/100kb.golang/pkg/serialize"
	"github.com/pemistahl/lingua-go"
	"github.com/smira/go-statsd"
)

const STAGE_FAILED = 0
const STAGE_INDEXED = 1
const STAGE_VALID_CONTENT = 2
const STAGE_COMPLETE = 10

type Engine struct {
	dbInsertPreparedArticle *sql.Stmt
	dbUpdatePreparedArticle *sql.Stmt
	db                      *sql.DB
	sd                      *statsd.Client
	langId                  lingua.LanguageDetector

	sentenceEmbeddingModel *ai.SentenceEmbeddingModel
	keywordExtractionModel *ai.KeywordExtractionModel
	zeroShotModel          *ai.ZeroShotModel
	parser                 *parsing.Engine

	http *retryhttp.Client
}

type Article struct {
	Url                  string
	FeedUrl              string
	Domain               string
	PublishedAt          int64
	BodyRaw              *serialize.Content
	LastFetchAt          int64
	LastMetaAt           int64
	LastContentExtractAt int64
	Title                string
	Description          string

	Body     *serialize.Content
	BadCount int64

	HTMLLength int64

	Stage int64

	SentenceEmbedding *serialize.Embeddings
	ExtractedKeywords *serialize.Keywords
	Classifications   *serialize.Keywords

	ContainsGoogleTagManager int64

	// Used in live/output
	DomainScore float64
}

func NewEngine(db *sql.DB, sd *statsd.Client, cachePath string, withModels bool) (*Engine, error) {
	engine := Engine{}
	var err error

	err = engine.initDB(db)
	if err != nil {
		return nil, err
	}

	engine.http, err = retryhttp.NewClient(cachePath, db, sd)
	if err != nil {
		return nil, err
	}

	if withModels {
		engine.sentenceEmbeddingModel, err = ai.NewSentenceEmbeddingModel()
		if err != nil {
			return nil, err
		}

		engine.keywordExtractionModel, err = ai.NewKeywordExtractionModel()
		if err != nil {
			return nil, err
		}

		engine.zeroShotModel, err = ai.NewZeroShotModel()
		if err != nil {
			return nil, err
		}
	}

	engine.parser, err = parsing.NewEngine()
	if err != nil {
		return nil, err
	}

	engine.sd = sd

	languages := []lingua.Language{
		lingua.English,
		lingua.French,
		lingua.German,
		lingua.Spanish,
	}

	engine.langId = lingua.NewLanguageDetectorBuilder().
		FromLanguages(languages...).
		WithLowAccuracyMode().
		Build()

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
