package article

import (
	"database/sql"

	"github.com/danhilltech/100kb.golang/pkg/ai"
	retryhttp "github.com/danhilltech/100kb.golang/pkg/http"
	"github.com/danhilltech/100kb.golang/pkg/parsing"
	"github.com/danhilltech/100kb.golang/pkg/serialize"
)

type Engine struct {
	dbInsertPreparedArticle *sql.Stmt
	dbUpdatePreparedArticle *sql.Stmt
	db                      *sql.DB

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

	Body             *serialize.Content
	WordCount        int64
	H1Count          int64
	HNCount          int64
	PCount           int64
	BadCount         int64
	FirstPersonRatio float64

	// NEW
	HTMLLength int64

	PageAbout    bool
	PageBlogRoll bool
	PageWriting  bool

	URLNews      bool
	URLBlog      bool
	URLHumanName bool

	DomainIsPopular bool
	DomainTLD       string

	// END

	SentenceEmbedding *serialize.Embeddings
	ExtractedKeywords *serialize.Keywords
	Classifications   *serialize.Keywords
}

func NewEngine(db *sql.DB, cachePath string, withModels bool) (*Engine, error) {
	engine := Engine{}
	var err error

	err = engine.initDB(db)
	if err != nil {
		return nil, err
	}

	engine.http, err = retryhttp.NewClient(cachePath)
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
