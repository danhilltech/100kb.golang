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

func (article *Article) GetScore() string {
	return fmt.Sprintf("%0.4f", article.Score())
}

func (article *Article) GetURL() string {
	return article.Url
}

func (article *Article) GetDomain() string {
	return article.Domain
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

	return d.Format(time.UnixDate)
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
