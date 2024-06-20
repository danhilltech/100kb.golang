package domain

import (
	"math"
	"sort"
	"strings"
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
		if a.PublishedAt > (time.Now().Unix()-60*60*24*365) && len(goodArticles) <= 10 {
			goodArticles = append(goodArticles, a)
		}
	}

	return goodArticles
}

func (d *Domain) GetLatestArticlesPerMonth() float64 {
	goodArticles := []*article.Article{}
	if d.Articles == nil || len(d.Articles) == 0 {
		return 0
	}

	sort.Slice(d.Articles, func(i, j int) bool {
		return d.Articles[i].PublishedAt > d.Articles[j].PublishedAt
	})

	for _, a := range d.Articles {
		if a.PublishedAt > (time.Now().Unix() - 60*60*24*180) {
			goodArticles = append(goodArticles, a)
		}
	}

	val := float64(len(goodArticles)) / 6.0

	if math.IsNaN(val) || math.IsInf(val, 0) {
		return 0
	}

	return val
}

func (d *Domain) GetFPR() float64 {
	arts := d.GetLatestArticlesToScore()
	if len(arts) == 0 {
		return 0.0
	}

	var words, fprs int

	for _, a := range arts {
		for _, t := range a.Body.GetContent() {
			words += len(strings.Fields(t.Text))

			fprs += strings.Count(t.Text, "I ")
			fprs += strings.Count(t.Text, "I'm ")
			fprs += strings.Count(t.Text, "I'll ")
			fprs += strings.Count(t.Text, "I've ")
			fprs += strings.Count(t.Text, " my ")
			fprs += strings.Count(t.Text, " My ")
			fprs += strings.Count(t.Text, " me ")
			fprs += strings.Count(t.Text, " mine ")
			if strings.HasPrefix(t.Text, "I ") {
				fprs++
			}

		}
	}

	val := float64(fprs) / float64(words)

	if math.IsNaN(val) || math.IsInf(val, 0) {
		return 0
	}

	return val
}

func (d *Domain) GetWordCount() float64 {
	arts := d.GetLatestArticlesToScore()
	if len(arts) == 0 {
		return 0.0
	}

	var words int

	for _, a := range arts {
		for _, t := range a.Body.GetContent() {
			words += len(strings.Fields(t.Text))

		}
	}

	val := float64(words) / float64(len(arts))

	if math.IsNaN(val) || math.IsInf(val, 0) {
		return 0
	}

	return val
}

func (d *Domain) GetPTagToGoodRatio() float64 {
	var count int64
	arts := d.GetLatestArticlesToScore()
	if len(arts) == 0 {
		return 0
	}

	goodCount := 0

	for _, a := range arts {
		for _, t := range a.Body.GetContent() {
			if t.Type == "p" || t.Type == "li" {
				count++
			}
			for _, t := range a.Body.GetContent() {
				if t.Type == "h1" || t.Type == "h2" || t.Type == "h3" || t.Type == "p" || t.Type == "li" {
					goodCount++
				}
			}
		}

	}

	val := float64(count) / float64(goodCount)

	if math.IsNaN(val) || math.IsInf(val, 0) {
		return 0
	}

	return val
}

func (d *Domain) GetHCount() uint64 {
	var count int64
	arts := d.GetLatestArticlesToScore()
	if len(arts) == 0 {
		return 0
	}

	for _, a := range arts {
		for _, t := range a.Body.GetContent() {
			if t.Type == "h1" || t.Type == "h2" || t.Type == "h3" {
				count++
			}
		}
	}

	val := float64(count) / float64(len(arts))

	if math.IsNaN(val) || math.IsInf(val, 0) {
		return 0
	}

	return uint64(val)
}

func (d *Domain) GetGoodTagCount() float64 {
	var count int64
	arts := d.GetLatestArticlesToScore()
	if len(arts) == 0 {
		return 0
	}

	for _, a := range arts {
		for _, t := range a.Body.GetContent() {
			if t.Type == "h1" || t.Type == "h2" || t.Type == "h3" || t.Type == "p" || t.Type == "li" {
				count++
			}
		}
	}

	val := float64(count) / float64(len(arts))

	if math.IsNaN(val) || math.IsInf(val, 0) {
		return 0
	}

	return val
}

func (d *Domain) GetCodeTagCount() float64 {
	var count int64
	arts := d.GetLatestArticlesToScore()
	if len(arts) == 0 {
		return 0
	}

	for _, a := range arts {
		for _, t := range a.Body.GetContent() {
			if t.Type == "code" {
				count++
			}
		}
	}

	val := float64(count) / float64(len(arts))

	if math.IsNaN(val) || math.IsInf(val, 0) {
		return 0
	}

	return val
}

func (d *Domain) GetBadTagCount() float64 {
	var wordCount int64
	arts := d.GetLatestArticlesToScore()
	if len(arts) == 0 {
		return 0
	}

	for _, a := range arts {
		wordCount += a.BadElementCount
	}

	val := float64(wordCount) / float64(len(arts))

	if math.IsNaN(val) || math.IsInf(val, 0) {
		return 0
	}

	return val
}

func (d *Domain) GetBadURLCount() float64 {
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

	return val
}

func (d *Domain) GetBadLinkCount() float64 {
	var wordCount int64
	arts := d.GetLatestArticlesToScore()
	if len(arts) == 0 {
		return 0
	}

	for _, a := range arts {
		wordCount += a.BadLinkCount
	}

	val := float64(wordCount) / float64(len(arts))

	if math.IsNaN(val) || math.IsInf(val, 0) {
		return 0
	}

	return val
}

func (d *Domain) GetLinkCount() float64 {
	var wordCount int64
	arts := d.GetLatestArticlesToScore()
	if len(arts) == 0 {
		return 0
	}

	for _, a := range arts {
		wordCount += a.LinkCount
	}

	val := float64(wordCount) / float64(len(arts))

	if math.IsNaN(val) || math.IsInf(val, 0) {
		return 0
	}

	return val
}

func (d *Domain) GetWordsPerByte() float64 {
	var bytSize float64
	arts := d.GetLatestArticlesToScore()
	if len(arts) == 0 {
		return 0
	}

	for _, a := range arts {
		bytSize += float64(a.HTMLLength)
	}

	wordCount := d.GetWordCount()

	val := float64(wordCount) / (bytSize / float64(len(arts)))

	if math.IsNaN(val) || math.IsInf(val, 0) {
		return 0.0
	}

	return val
}

func (d *Domain) GetGoodTagsPerByte() float64 {
	var count int64
	var htmlLen int64
	arts := d.GetLatestArticlesToScore()
	if len(arts) == 0 {
		return 0
	}

	for _, a := range arts {
		for _, t := range a.Body.GetContent() {
			if t.Type == "h1" || t.Type == "h2" || t.Type == "h3" || t.Type == "p" || t.Type == "li" {
				count++
			}
		}
		htmlLen = htmlLen + a.HTMLLength
	}

	val := (float64(count) / float64(htmlLen)) / float64(len(arts))

	if math.IsNaN(val) || math.IsInf(val, 0) {
		return 0
	}

	return val
}

func (d *Domain) GetWordsPerParagraph() float64 {
	var wordCount int
	var pCount int
	arts := d.GetLatestArticlesToScore()
	if len(arts) == 0 {
		return 0
	}

	for _, a := range arts {
		for _, t := range a.Body.GetContent() {
			if t.Type == "p" || t.Type == "li" {
				wordCount += len(strings.Fields(t.Text))
				pCount++
			}
		}
	}

	val := float64(wordCount) / float64(pCount)

	if math.IsNaN(val) || math.IsInf(val, 0) {
		return 0.0
	}

	return val
}

func (d *Domain) ContainsGoogleTagManager() bool {
	arts := d.GetLatestArticlesToScore()
	if len(arts) == 0 {
		return false
	}

	tagCount := int64(0)

	for _, a := range arts {
		tagCount += a.ContainsGoogleTagManager

	}

	val := float64(tagCount) / float64(len(arts))

	if math.IsNaN(val) || math.IsInf(val, 0) {
		return false
	}

	return val > 0.5
}

func (d *Domain) GetGoodBadTagRatio() float64 {

	val := float64(d.GetGoodTagCount()) / float64(d.GetBadTagCount()+d.GetGoodTagCount())

	if math.IsNaN(val) || math.IsInf(val, 0) {
		return 0.0
	}

	return val
}

func (d *Domain) GetBadLinkRatio() float64 {

	val := float64(d.GetBadLinkCount()) / float64(d.GetLinkCount())

	if math.IsNaN(val) || math.IsInf(val, 0) {
		return 0.0
	}

	return val
}

func (domain *Domain) GetFloatFeatureNames() []string {
	names := []string{}

	names = append(names, "fpr")
	names = append(names, "wordCount")
	names = append(names, "wordsPerByte")
	names = append(names, "wordsPerP")
	names = append(names, "goodPcnt")
	names = append(names, "goodTagsPerByte")
	names = append(names, "articlesPerMonth")
	// names = append(names, "goodTagCount")
	// names = append(names, "badTagCount")
	names = append(names, "badUrlCount")
	names = append(names, "pToGoodRatio")
	names = append(names, "badLinkRatio")
	// names = append(names, "codeCount")
	names = append(names, "urlHuman")
	names = append(names, "urlBlog")
	names = append(names, "urlNews")
	names = append(names, "popularDomain")
	names = append(names, "containsGoogleTagManager")
	names = append(names, "containsGoogleAds")

	return names
}

func safeLog(v float64) float64 {

	val := math.Log(v)

	if math.IsNaN(val) || math.IsInf(val, 0) {
		return 0
	}

	return val
}

func (domain *Domain) GetFloatFeatures() []float64 {
	features := []float64{}

	features = append(features, domain.GetFPR())
	features = append(features, safeLog(domain.GetWordCount()))
	features = append(features, domain.GetWordsPerByte())
	features = append(features, domain.GetWordsPerParagraph())
	features = append(features, domain.GetGoodBadTagRatio())
	features = append(features, domain.GetGoodTagsPerByte())
	features = append(features, domain.GetLatestArticlesPerMonth())
	// features = append(features, safeLog(domain.GetGoodTagCount()))
	// features = append(features, safeLog(domain.GetBadTagCount()))
	features = append(features, safeLog(domain.GetBadURLCount()))
	features = append(features, domain.GetPTagToGoodRatio())
	// features = append(features, safeLog(domain.GetCodeTagCount()))

	features = append(features, domain.GetBadLinkRatio())

	if domain.URLHumanName {
		features = append(features, 1.0)
	} else {
		features = append(features, 0)
	}

	if domain.URLBlog {
		features = append(features, 1.0)
	} else {
		features = append(features, 0)
	}

	if domain.URLNews {
		features = append(features, 1.0)
	} else {
		features = append(features, 0)
	}
	if domain.DomainIsPopular {
		features = append(features, 1.0)
	} else {
		features = append(features, 0)
	}
	if domain.ContainsGoogleTagManager() {
		features = append(features, 1.0)
	} else {
		features = append(features, 0)
	}

	if domain.DomainGoogleAds {
		features = append(features, 1.0)
	} else {
		features = append(features, 0)
	}

	return features
}
