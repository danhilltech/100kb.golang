package article

import (
	"fmt"
	"strings"
	"time"
)

func (article *Article) GetTitle() string {
	return article.Title
}

func (article *Article) GetFPR() string {
	return fmt.Sprintf("%0.4f", article.FirstPersonRatio)
}

func (article *Article) GetScore() string {
	return fmt.Sprintf("%0.4f", article.Score())
}

func (article *Article) GetKeywords() string {
	b := strings.Builder{}

	for _, k := range article.ExtractedKeywords.Keywords {
		b.WriteString(k.Text)
		b.WriteString(fmt.Sprintf(" %0.3f", k.Score))
		b.WriteString(", ")
	}
	return b.String()
}

func (article *Article) GetPublishedAt() string {
	d := time.Unix(article.PublishedAt, 0)

	return d.Format(time.UnixDate)
}
