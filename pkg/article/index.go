package article

import (
	"io"
	"time"

	"github.com/danhilltech/100kb.golang/pkg/http"
)

// Crawls
func (engine *Engine) articleIndex(article *Article) (*http.URLRequest, error) {

	article.LastFetchAt = time.Now().Unix()

	// crawl it
	resp, err := engine.http.GetWithSafety(article.Url)
	if err != nil {
		return resp, err
	}
	if resp != nil && resp.Response != nil {
		defer resp.Response.Body.Close()

		byts, err := io.ReadAll(resp.Response.Body)
		if err != nil {
			return nil, err
		}
		article.HTMLLength = int64(len(byts))

		// io.Copy(io.Discard, resp.Response.Body)
	}

	return resp, nil

}
