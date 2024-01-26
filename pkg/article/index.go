package article

import (
	"io"
	"strings"
	"time"
)

// Crawls
func (engine *Engine) articleIndex(article *Article) error {

	article.LastFetchAt = time.Now().Unix()

	// Check its a good url
	if strings.HasSuffix(article.Url, ".mp4") {
		return nil
	}
	if strings.HasSuffix(article.Url, ".mp3") {
		return nil
	}
	if strings.HasSuffix(article.Url, ".pdf") {
		return nil
	}
	if !strings.HasPrefix(article.Url, "http") {
		return nil
	}

	head, err := engine.http.Head(article.Url)
	if err != nil {
		return nil
	}
	defer head.Body.Close()
	if !strings.Contains(head.Header.Get("Content-Type"), "text/html") {
		return nil
	}

	// crawl it
	res, err := engine.cache.Get(article.Url, engine.http)
	if err != nil {
		return err
	}
	defer res.Close()
	io.Copy(io.Discard, res)

	return nil

}
