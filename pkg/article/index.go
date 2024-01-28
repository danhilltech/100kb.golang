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
	defer resp.Response.Body.Close()
	io.Copy(io.Discard, resp.Response.Body)

	return resp, nil

}
