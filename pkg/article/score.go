package article

import (
	"math"
	"strings"
	"time"
)

func (article *Article) Score() float64 {
	score := 1.0
	if strings.Contains(article.Url, "forum") {
		score = 0
	}
	if strings.Contains(article.Url, "news") {
		score = 0
	}

	now := time.Now().Unix()

	diff := now - article.PublishedAt

	timeWeight := 0.2 / math.Log(float64(diff)/3600)

	score = score * timeWeight

	score = score * article.FirstPersonRatio

	score = score * math.Log(float64(article.WordCount))

	return score
}
