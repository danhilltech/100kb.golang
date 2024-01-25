package article

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/danhilltech/100kb.golang/pkg/ai"
	retryhttp "github.com/danhilltech/100kb.golang/pkg/http"
	"github.com/danhilltech/100kb.golang/pkg/parsing"
	"github.com/danhilltech/100kb.golang/pkg/serialize"
	"github.com/peterbourgon/diskv/v3"
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
	cache *diskv.Diskv
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

	d := diskv.New(diskv.Options{
		BasePath: cachePath,
		// CacheSizeMax:      1024 * 1024,
		CacheSizeMax:      10_737_418_240, // 1 GB
		AdvancedTransform: AdvancedTransformExample,
		InverseTransform:  InverseTransformExample,
	})

	engine.cache = d

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

func AdvancedTransformExample(key string) *diskv.PathKey {
	path := strings.Split(key, "/")
	last := len(path) - 1
	return &diskv.PathKey{
		Path:     path[:last],
		FileName: path[last] + ".txt",
	}
}

// If you provide an AdvancedTransform, you must also provide its
// inverse:

func InverseTransformExample(pathKey *diskv.PathKey) (key string) {
	txt := pathKey.FileName[len(pathKey.FileName)-4:]
	if txt != ".txt" {
		panic("Invalid file found in storage folder!")
	}
	return strings.Join(pathKey.Path, "/") + pathKey.FileName[:len(pathKey.FileName)-4]
}
