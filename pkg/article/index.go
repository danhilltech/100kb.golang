package article

import (
	"fmt"
	"io"
	"time"
	"unicode/utf8"
)

// Crawls
func (engine *Engine) articleIndex(article *Article) error {

	article.LastFetchAt = time.Now().Unix()
	article.Stage = STAGE_FAILED

	// crawl it
	resp, err := engine.http.GetWithSafety(article.Url)
	if err != nil {
		return err
	}
	if resp != nil && resp.Response != nil {
		defer resp.Response.Body.Close()

		byts, err := io.ReadAll(resp.Response.Body)
		if err != nil {
			return err
		}
		if len(byts) > 500000 { // Don't bother parsing anything over 500kb uncompressed
			fmt.Printf("Skipping %s as body too large at %d bytes\n", article.Url, len(byts))
			return nil
		}
		if !utf8.Valid(byts) {
			fmt.Printf("Skipping %s as body not valid utf8\n", article.Url)
			return nil
		}
		article.HTMLLength = int64(len(byts))

		article.Stage = STAGE_INDEXED

		// io.Copy(io.Discard, resp.Response.Body)
	}

	return nil

}
