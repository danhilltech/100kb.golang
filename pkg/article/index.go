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

	engine.cacheFailingDomainsRWMutex.RLock()
	if engine.cacheFailingDomains[article.Domain] > 5 {
		engine.cacheFailingDomainsRWMutex.RUnlock()
		return ErrTooManyFailingOnDomain
	}
	engine.cacheFailingDomainsRWMutex.RUnlock()

	// crawl it
	resp, err := engine.http.Get(article.Url)
	if err != nil {
		engine.cacheFailingDomainsRWMutex.Lock()
		engine.cacheFailingDomains[article.Domain]++
		engine.cacheFailingDomainsRWMutex.Unlock()
		return err
	}
	if resp != nil {
		defer resp.Body.Close()
		if resp.StatusCode > 400 {
			engine.cacheFailingDomainsRWMutex.Lock()
			engine.cacheFailingDomains[article.Domain]++
			engine.cacheFailingDomainsRWMutex.Unlock()
			return nil
		}

		byts, err := io.ReadAll(resp.Body)
		if err != nil {
			engine.cacheFailingDomainsRWMutex.Lock()
			engine.cacheFailingDomains[article.Domain]++
			engine.cacheFailingDomainsRWMutex.Unlock()
			return err
		}

		if !utf8.Valid(byts) {
			engine.cacheFailingDomainsRWMutex.Lock()
			engine.cacheFailingDomains[article.Domain]++
			engine.cacheFailingDomainsRWMutex.Unlock()
			return nil
		}
		article.HTMLLength = int64(len(byts))

		article.Stage = STAGE_INDEXED

		// io.Copy(io.Discard, resp.Response.Body)
	}

	return nil

}
