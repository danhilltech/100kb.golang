package domain

import (
	"math"
	"sort"
	"time"

	"github.com/danhilltech/100kb.golang/pkg/article"
)

func (d *Domain) GetLatestArticlesToScore() []*article.Article {
	goodArticles := []*article.Article{}
	if d.Articles == nil || len(d.Articles) == 0 {
		return goodArticles
	}

	sort.Slice(d.Articles, func(i, j int) bool {
		return d.Articles[i].PublishedAt > d.Articles[j].PublishedAt
	})

	for _, a := range d.Articles {
		if a.PublishedAt > (time.Now().Unix()-60*60*24*365) && len(goodArticles) <= 5 {
			goodArticles = append(goodArticles, a)
		}
	}

	return goodArticles
}

func (d *Domain) GetFPR() float64 {
	fpr := 0.0
	arts := d.GetLatestArticlesToScore()
	if len(arts) == 0 {
		return 0.0
	}

	for _, a := range arts {
		fpr += a.FirstPersonRatio
	}

	val := fpr / float64(len(arts))

	if math.IsNaN(val) || math.IsInf(val, 0) {
		return 0
	}

	return val
}

func (d *Domain) GetWordCount() uint64 {
	var wordCount int64
	arts := d.GetLatestArticlesToScore()
	if len(arts) == 0 {
		return 0
	}

	for _, a := range arts {
		wordCount += a.WordCount
	}

	val := float64(wordCount) / float64(len(arts))

	if math.IsNaN(val) || math.IsInf(val, 0) {
		return 0
	}

	return uint64(val)
}

func (d *Domain) GetGoodTagCount() uint64 {
	var wordCount int64
	arts := d.GetLatestArticlesToScore()
	if len(arts) == 0 {
		return 0
	}

	for _, a := range arts {
		wordCount += a.PCount + a.H1Count + a.HNCount
	}

	val := float64(wordCount) / float64(len(arts))

	if math.IsNaN(val) || math.IsInf(val, 0) {
		return 0
	}

	return uint64(val)
}

func (d *Domain) GetBadTagCount() uint64 {
	var wordCount int64
	arts := d.GetLatestArticlesToScore()
	if len(arts) == 0 {
		return 0
	}

	for _, a := range arts {
		wordCount += a.BadCount
	}

	val := float64(wordCount) / float64(len(arts))

	if math.IsNaN(val) || math.IsInf(val, 0) {
		return 0
	}

	return uint64(val)
}

func (d *Domain) GetWordsPerByte() float64 {
	var wordCount float64
	arts := d.GetLatestArticlesToScore()
	if len(arts) == 0 {
		return 0
	}

	for _, a := range arts {
		wordCount += float64(a.WordCount) / float64(a.HTMLLength)
	}

	val := float64(wordCount) / float64(len(arts))

	if math.IsNaN(val) || math.IsInf(val, 0) {
		return 0.0
	}

	return val
}

func (d *Domain) GetWordsPerParagraph() float64 {
	var wordCount float64
	arts := d.GetLatestArticlesToScore()
	if len(arts) == 0 {
		return 0
	}

	for _, a := range arts {
		wordCount += (float64(a.WordCount) / float64(a.PCount))
	}

	val := float64(wordCount) / float64(len(arts))

	if math.IsNaN(val) || math.IsInf(val, 0) {
		return 0.0
	}

	return val
}

func (d *Domain) GetGoodBadTagRatio() float64 {

	val := float64(d.GetGoodTagCount()) / float64(d.GetBadTagCount()+d.GetGoodTagCount())

	if math.IsNaN(val) || math.IsInf(val, 0) {
		return 0.0
	}

	return val
}
