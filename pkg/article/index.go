package article

import (
	"io"
	"time"
	"unicode/utf8"
)

// Crawls
func (engine *Engine) articleIndex(article *Article) error {

	article.LastFetchAt = time.Now().Unix()
	article.Stage = STAGE_FAILED

	// crawl it
	resp, err := engine.http.Get(article.Url)
	if err != nil || resp.StatusCode > 400 {
		return err
	}
	if resp != nil {
		defer resp.Body.Close()

		byts, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		if !utf8.Valid(byts) {
			return nil
		}
		article.HTMLLength = int64(len(byts))

		article.Stage = STAGE_INDEXED

		// io.Copy(io.Discard, resp.Response.Body)
	}

	return nil

}
