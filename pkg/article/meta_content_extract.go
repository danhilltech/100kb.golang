package article

import (
	"database/sql"
	"time"

	"github.com/danhilltech/100kb.golang/pkg/serialize"
	"golang.org/x/net/html"
)

func (engine *Engine) articleExtractContent(tx *sql.Tx, article *Article) error {
	// Check we have enough data
	article.LastContentExtractAt = time.Now().Unix()

	htmlStream, err := engine.cache.ReadStream(article.Url)
	if err != nil {
		return err
	}
	defer htmlStream.Close()

	htmlDoc, err := html.Parse(htmlStream)

	if err != nil {
		return err
	}

	_, _, badCount, err := engine.parser.IdentifyBadElements(htmlDoc, article.Url)
	if err != nil {
		return err
	}
	article.BadCount = int64(badCount)

	err = engine.parser.IdentifyGoodElements(htmlDoc, article.Url)
	if err != nil {
		return err
	}

	body, title, description, err := engine.parser.HtmlToText(htmlDoc)
	if err != nil {
		return err
	}

	article.BodyRaw = &serialize.Content{Content: body}
	article.Title = title
	article.Description = description

	return nil

}
