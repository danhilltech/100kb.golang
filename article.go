package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Article struct {
	Url         string
	FeedUrl     string
	PublishedAt int64
	Html        []byte
	BodyRaw     []string
	LastFetchAt int64
	Title       string
	Description string
	Body        string

	WordCount        int64
	FirstPersonRatio float64
}

func articleIndexWorker(jobs <-chan *Article, results chan<- *Article) {
	for id := range jobs {
		err := articleIndex(id)
		if err != nil {
			fmt.Println(err)
		}
		results <- id
	}
}

// Crawls
func articleIndex(article *Article) error {
	// crawl it
	resp, err := http.Get(article.Url)
	if err != nil {
		return nil
	}

	defer resp.Body.Close()

	html, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	article.Html = html

	reader := bytes.NewReader(html)

	body, title, description, err := HtmlToText(reader)
	if err != nil {
		return nil
	}

	article.BodyRaw = body
	article.LastFetchAt = time.Now().Unix()
	article.Title = title
	article.Description = description

	return nil

}

func articleIndexes(articles []*Article) error {

	workers := 10
	jobs := make(chan *Article, len(articles))
	results := make(chan *Article, len(articles))

	for w := 1; w <= workers; w++ {
		go articleIndexWorker(jobs, results)
	}

	for j := 1; j <= len(articles); j++ {
		jobs <- articles[j-1]
	}
	close(jobs)

	items := make([]*Article, len(articles))

	for a := 1; a <= len(articles); a++ {
		b := <-results
		items[a-1] = b
	}

	return nil

}

func articleMetaWorker(tx *sql.Tx, jobs <-chan *Article, results chan<- *Article) {
	for id := range jobs {
		err := articleMeta(tx, id)
		if err != nil {
			fmt.Println(err)
		}
		results <- id
	}
}

// Crawls
func articleMeta(tx *sql.Tx, article *Article) error {
	feedArticles, err := getArticlesByFeed(tx, article.FeedUrl, article.Url)
	if err != nil {
		return err
	}

	var currentCanon []string

	for _, feed := range feedArticles {
		currentCanon = append(currentCanon, feed.BodyRaw...)
	}

	// unique the content

	var uniqueContent []string

	for _, line := range article.BodyRaw {
		found := false
		for _, currLine := range currentCanon {
			if line == currLine {
				found = true
			}
		}
		if !found {
			uniqueContent = append(uniqueContent, line)
		}
	}

	article.Body = strings.Join(uniqueContent, " ")

	// Word count
	article.WordCount = int64(len(strings.Split(article.Body, " ")))

	firstPersonCount := 0

	firstPersonCount += strings.Count(article.Body, "I ")
	firstPersonCount += strings.Count(article.Body, " my ")
	firstPersonCount += strings.Count(article.Body, " me ")
	firstPersonCount += strings.Count(article.Body, " mine ")
	firstPersonCount += strings.Count(article.Body, " we ")
	firstPersonCount += strings.Count(article.Body, " us ")
	firstPersonCount += strings.Count(article.Body, " our ")

	if article.WordCount > 0 && firstPersonCount > 0 {
		article.FirstPersonRatio = float64(firstPersonCount) / float64(article.WordCount)
	} else {
		article.FirstPersonRatio = 0
	}

	return nil

}

func articleMetas(tx *sql.Tx, articles []*Article) error {

	workers := 10
	jobs := make(chan *Article, len(articles))
	results := make(chan *Article, len(articles))

	for w := 1; w <= workers; w++ {
		go articleMetaWorker(tx, jobs, results)
	}

	for j := 1; j <= len(articles); j++ {
		jobs <- articles[j-1]
	}
	close(jobs)

	items := make([]*Article, len(articles))

	for a := 1; a <= len(articles); a++ {
		b := <-results
		items[a-1] = b
	}

	return nil

}
