package domain

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/danhilltech/100kb.golang/pkg/article"
	"github.com/mmcdole/gofeed"
	"golang.org/x/net/publicsuffix"
)

func (engine *Engine) RunFeedRefresh(chunkSize int, workers int) error {
	feeds, err := engine.getDomainsToRefresh()
	if err != nil {
		return err
	}

	fmt.Printf("Checking %d feeds for new links\n", len(feeds))

	jobs := make(chan *Domain, len(feeds))
	results := make(chan *Domain, len(feeds))

	for w := 1; w <= workers; w++ {
		go engine.feedRefreshWorker(jobs, results)
	}

	for j := 1; j <= len(feeds); j++ {
		jobs <- feeds[j-1]
	}
	close(jobs)

	t := time.Now().UnixMilli()

	for a := 1; a <= len(feeds); a++ {
		domain := <-results

		err = engine.Update(domain)
		if err != nil {
			return err
		}

		for _, article := range domain.Articles {

			err = engine.articleEngine.Insert(article, domain.FeedURL, domain.Domain)
			if err != nil {
				return err
			}

		}

		if a > 0 && a%chunkSize == 0 {
			diff := time.Now().UnixMilli() - t
			qps := (float64(chunkSize) / float64(diff)) * 1000
			t = time.Now().UnixMilli()
			fmt.Printf("\tdone %d/%d at %0.2f/s\n", a, len(feeds), qps)

		}

	}
	fmt.Printf("\tdone %d\n", len(feeds))

	return nil
}

// Crawls
func (engine *Engine) feedRefresh(feed *Domain) error {
	feed.LastFetchAt = time.Now().Unix()
	// First check the URL isn't banned

	// crawl it
	resp, err := engine.httpCrawl.GetWithSafety(feed.FeedURL)
	if err != nil || resp == nil {
		return nil
	}
	if resp.Response == nil {
		return nil
	}

	defer resp.Response.Body.Close()

	fp := gofeed.NewParser()
	rss, err := fp.Parse(resp.Response.Body)
	if err != nil {
		return err
	}

	feed.Language = rss.Language

	feed.Articles = []*article.Article{}

	fullDomain := fmt.Sprintf("https://%s", feed.Domain)

	headResp, err := engine.httpCrawl.Head(fmt.Sprintf("%s/about", fullDomain))
	if err == nil && headResp.StatusCode < 400 {
		feed.PageAbout = true
	}
	if !feed.PageAbout {
		headResp, err = engine.httpCrawl.Head(fmt.Sprintf("%s/me", fullDomain))
		if err == nil && headResp.StatusCode < 400 {
			feed.PageAbout = true
		}
	}
	if !feed.PageAbout {
		headResp, err = engine.httpCrawl.Head(fmt.Sprintf("%s/about.html", fullDomain))
		if err == nil && headResp.StatusCode < 400 {
			feed.PageAbout = true
		}
	}
	if !feed.PageAbout {
		headResp, err = engine.httpCrawl.Head(fmt.Sprintf("%s/whoami", fullDomain))
		if err == nil && headResp.StatusCode < 400 {
			feed.PageAbout = true
		}
	}
	headResp, err = engine.httpCrawl.Head(fmt.Sprintf("%s/blogroll", fullDomain))
	if err == nil && headResp.StatusCode < 400 {
		feed.PageBlogRoll = true
	}
	headResp, err = engine.httpCrawl.Head(fmt.Sprintf("%s/archive", fullDomain))
	if err == nil && headResp.StatusCode < 400 {
		feed.PageWriting = true
	}

	tld, _ := publicsuffix.PublicSuffix(fullDomain)
	if tld != "" {
		feed.DomainTLD = tld
	}

	urlHumanName, urlNews, urlBlog, popularDomain, err := engine.parser.IdentifyURL(fullDomain)
	if err != nil {
		return fmt.Errorf("could not identify url %w", err)
	}

	feed.URLBlog = urlBlog
	feed.URLHumanName = urlHumanName
	feed.URLNews = urlNews
	feed.DomainIsPopular = popularDomain

	for _, item := range rss.Items {

		if item.Link == "" {
			continue
		}
		if !strings.HasPrefix(item.Link, "http") && !strings.HasPrefix(item.Link, "/") {
			continue
		}

		baseUrl, err := url.Parse(feed.FeedURL)

		if err != nil {
			return err
		}

		newUrl, err := url.Parse(item.Link)
		if err != nil {
			return err
		}

		baseUrl.Scheme = "https"

		cleanUrl := baseUrl.ResolveReference(newUrl)

		cleanUrlString := cleanUrl.String()

		cleanUrlString = strings.ReplaceAll(cleanUrlString, "http://", "https://")

		if !strings.HasPrefix(cleanUrlString, "https://") || strings.Contains(cleanUrlString, "https:///") {
			fmt.Println(cleanUrlString, feed.Domain, item.Link, feed.FeedURL)
			continue
		}

		art := &article.Article{
			Url:    cleanUrlString,
			Domain: feed.Domain,
		}
		if item.PublishedParsed != nil {
			art.PublishedAt = item.PublishedParsed.Unix()
		}

		feed.Articles = append(feed.Articles, art)
	}

	return nil

}

func (engine *Engine) feedRefreshWorker(jobs <-chan *Domain, results chan<- *Domain) {
	for id := range jobs {
		err := engine.feedRefresh(id)
		if err != nil {
			fmt.Println(id.Domain, err)
		}
		results <- id
	}
}
