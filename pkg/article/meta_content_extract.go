package article

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"time"
	"unicode/utf8"

	"github.com/danhilltech/100kb.golang/pkg/serialize"
	"golang.org/x/net/html"
	"golang.org/x/net/publicsuffix"
)

func (engine *Engine) articleExtractContent(tx *sql.Tx, article *Article) error {
	// Check we have enough data
	article.LastContentExtractAt = time.Now().Unix()

	htmlStream, err := engine.http.Get(article.Url)
	if err != nil {
		return fmt.Errorf("could not get article %w", err)
	}
	defer htmlStream.Body.Close()

	fullBody, err := io.ReadAll(htmlStream.Body)
	if err != nil {
		return err
	}
	if len(fullBody) > 500000 { // Don't bother parsing anything over 500kb uncompressed
		fmt.Printf("Skipping %s as body too large at %d bytes\n", article.Url, len(fullBody))
		return nil
	}
	if !utf8.Valid(fullBody) {
		fmt.Printf("Skipping %s as body not valid utf8\n", article.Url)
		return nil
	}

	r := bytes.NewReader(fullBody)

	htmlDoc, err := html.Parse(r)

	if err != nil {
		return fmt.Errorf("could not parse %w", err)
	}

	_, _, badCount, err := engine.parser.IdentifyBadElements(htmlDoc, article.Url)
	if err != nil {
		return fmt.Errorf("could not identify bad element %w", err)
	}
	article.BadCount = int64(badCount)

	err = engine.parser.IdentifyGoodElements(htmlDoc, article.Url)
	if err != nil {
		return fmt.Errorf("could not identify good elements %w", err)
	}

	body, title, description, err := engine.parser.HtmlToText(htmlDoc)
	if err != nil {
		return fmt.Errorf("could not extract text %w", err)
	}

	article.BodyRaw = &serialize.Content{Content: body}
	article.Title = title
	article.Description = description

	hasAbout, hasBlogRoll, hasWriting, err := engine.parser.IdentifyInternalPages(htmlDoc, article.Url)
	if err != nil {
		return fmt.Errorf("could not identify interal pages %w", err)
	}

	article.PageAbout = hasAbout
	article.PageBlogRoll = hasBlogRoll
	article.PageWriting = hasWriting

	tld, icann := publicsuffix.PublicSuffix(article.Domain)
	if icann {
		article.DomainTLD = tld
	}

	urlHumanName, urlNews, urlBlog, popularDomain, err := engine.parser.IdentifyURL(article.Url)
	if err != nil {
		return fmt.Errorf("could not identify url %w", err)
	}

	article.URLBlog = urlBlog
	article.URLHumanName = urlHumanName
	article.URLNews = urlNews
	article.DomainIsPopular = popularDomain

	return nil

}
