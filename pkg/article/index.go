package article

import (
	"fmt"
	"io"
	"time"
	"unicode/utf8"
)

var ErrTooManyFailingOnDomain = fmt.Errorf("too many fails on domain")

// Crawls
func (engine *Engine) articleIndex(article *Article) error {

	article.LastFetchAt = time.Now().Unix()
	article.Stage = STAGE_FAILED

	if engine.cacheFailingDomains[article.Domain] > 5 {
		return ErrTooManyFailingOnDomain
	}

	// crawl it
	resp, err := engine.http.Get(article.Url)
	if err != nil {
		engine.cacheFailingDomains[article.Domain]++
		return err
	}
	if resp != nil {
		defer resp.Body.Close()
		if resp.StatusCode > 400 {
			engine.cacheFailingDomains[article.Domain]++
			return nil
		}

		byts, err := io.ReadAll(resp.Body)
		if err != nil {
			engine.cacheFailingDomains[article.Domain]++
			return err
		}

		if !utf8.Valid(byts) {
			engine.cacheFailingDomains[article.Domain]++
			return nil
		}
		article.HTMLLength = int64(len(byts))

		article.Stage = STAGE_INDEXED

		// io.Copy(io.Discard, resp.Response.Body)
	}

	return nil

}
