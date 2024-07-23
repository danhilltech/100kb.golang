package article

import (
	"fmt"
	"sync"
	"time"

	"github.com/danhilltech/100kb.golang/pkg/serialize"
	"github.com/pemistahl/lingua-go"
	"golang.org/x/net/html"
)

var mapLock sync.Mutex

func (engine *Engine) articleExtractContent(article *Article) error {
	// Check we have enough data
	article.LastContentExtractAt = time.Now().Unix()

	// check status ok

	htmlStream, err := engine.http.Get(article.Url)
	if err != nil || htmlStream.StatusCode > 400 {
		return fmt.Errorf("could not get article %w", err)
	}
	defer htmlStream.Body.Close()

	htmlDoc, err := html.Parse(htmlStream.Body)
	if err != nil {
		return fmt.Errorf("could not parse %w", err)
	}

	analysis, err := engine.parser.IdentifyElements(htmlDoc, article.Url)
	if err != nil {
		return fmt.Errorf("could not identify bad element %w", err)
	}
	article.BadCount = int64(len(analysis.BadUrls))
	article.BadElementCount = int64(len(analysis.BadElements))
	article.LinkCount = int64(len(analysis.Links))
	article.BadLinkCount = int64(len(analysis.BadLinkTitles))

	body, title, description, err := engine.parser.HtmlToText(htmlDoc)
	if err != nil {
		return fmt.Errorf("could not extract text %w", err)
	}

	if len(body) == 0 {
		article.Stage = STAGE_FAILED
		return fmt.Errorf("no body found %s", article.Url)
	}

	var considerText string

	for _, b := range body {
		considerText = fmt.Sprintf("%s %s", considerText, b.Text)
		if len(considerText) > 512 {
			break
		}
	}

	// Check its in English

	mapLock.Lock()
	defer mapLock.Unlock()
	if engine.langDomainCacheNonEng[article.Domain] > 2 {
		article.Stage = STAGE_FAILED
		return fmt.Errorf("not in english (cached) %s", article.Url)
	} else if engine.langDomainCacheEng[article.Domain] > 2 {
	} else {

		res, exists := engine.langId.DetectLanguageOf(considerText)

		if !exists || res != lingua.English {
			article.Stage = STAGE_FAILED
			engine.langDomainCacheNonEng[article.Domain]++
			return fmt.Errorf("not in english %s", article.Url)
		} else {
			engine.langDomainCacheEng[article.Domain]++
		}
	}

	article.BodyRaw = &serialize.Content{Content: body}
	article.Title = title
	article.Description = description

	article.Stage = STAGE_VALID_CONTENT

	return nil

}
