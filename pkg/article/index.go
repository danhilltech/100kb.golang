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
	resp, err := engine.http.Get(article.Url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	html, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	article.HTML = html

	return nil

}
