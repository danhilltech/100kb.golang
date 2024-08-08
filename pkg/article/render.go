package article

import (
	"fmt"
	"hash/fnv"
	"html/template"
	"strings"
	"time"
)

func (article *Article) GetTitle() string {
	return article.Title
}

func (article *Article) GetURL() string {
	return article.Url
}

func (article *Article) GetDomain() string {
	return article.Domain
}

func (article *Article) GetDomainScore() float64 {
	return article.DomainScore
}

func (article *Article) GetDomainClassName() string {
	if article.DomainScore > 0.8 {
		return "score-excellent"
	} else if article.DomainScore > 0.5 {
		return "score-good"
	} else if article.DomainScore > 0.2 {
		return "score-bad"
	} else {
		return "score-poor"
	}
}

func (article *Article) GetSlug() string {
	keyHash := fnv.New64()

	keyHash.Write([]byte(article.Url))
	return fmt.Sprintf("article-%d.html", keyHash.Sum64())
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

	return d.Format("2006-01-02")
}

func (article *Article) GetPos() string {
	return fmt.Sprintf("%d", article.DayPosition)
}

func (article *Article) GetTags() []string {
	tags := []string{}

	if article.ExtractedKeywords == nil || len(article.ExtractedKeywords.Keywords) == 0 {
		return tags
	}

	for _, t := range article.ExtractedKeywords.Keywords {
		tags = append(tags, fmt.Sprintf("%s (%0.3f)", t.Text, t.Score))
	}
	return tags
}

func (article *Article) GetZeroShot() []string {
	tags := []string{}

	if article.Classifications == nil || len(article.Classifications.Keywords) == 0 {
		return tags
	}

	for _, t := range article.Classifications.Keywords {
		tags = append(tags, fmt.Sprintf("%s (%0.3f)", t.Text, t.Score))
	}
	return tags
}

func (article *Article) GetHTML() template.HTML {
	w := strings.Builder{}

	for _, c := range article.Body.Content {
		switch c.Type {
		case "p":
			w.WriteString("<p>")
			w.WriteString(c.Text)
			w.WriteString("</p>")
		case "h1":
			w.WriteString("<h1>")
			w.WriteString(c.Text)
			w.WriteString("</h1>")
		case "h2":
			w.WriteString("<h2>")
			w.WriteString(c.Text)
			w.WriteString("</h2>")
		case "h3":
			w.WriteString("<h3>")
			w.WriteString(c.Text)
			w.WriteString("</h3>")
		case "li":
			w.WriteString("<li>")
			w.WriteString(c.Text)
			w.WriteString("</li>")
		default:
			w.WriteString("<")
			w.WriteString(c.Type)
			w.WriteString(">")

			w.WriteString(c.Text)

			w.WriteString("</")
			w.WriteString(c.Type)
			w.WriteString(">")
		}
	}
	return template.HTML(w.String())
}
