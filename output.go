package main

import (
	"database/sql"
	"embed"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"text/template"

	"github.com/danhilltech/100kb.golang/pkg/article"
)

type RenderEngine struct {
	templates *template.Template
	outputDir string

	articles []*article.Article
}

type PageData struct {
	Title string
	Data  interface{}
}

var pageSize = 100

//go:embed views/*.html
var tmplFS embed.FS

func NewRenderEnding(outputDir string, articles []*article.Article) *RenderEngine {
	funcMap := template.FuncMap{
		// "inc": inc,
	}

	templates := template.Must(template.New("").Funcs(funcMap).ParseFS(tmplFS, "views/*.html"))

	for _, t := range templates.Templates() {
		fmt.Println(t.Name())
	}

	return &RenderEngine{
		templates: templates,
		articles:  articles,
		outputDir: outputDir,
	}
}

func CreateOutput(db *sql.DB) error {

	articleEngine, err := article.NewEngine(db)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer articleEngine.Close()

	txn, err := db.Begin()
	if err != nil {
		return err
	}
	defer txn.Rollback()

	articles, err := articleEngine.GetAllValid(txn)
	if err != nil {
		return err
	}
	err = txn.Commit()
	if err != nil {
		return err
	}

	engine := NewRenderEnding("output", articles)

	err = engine.ArticleLists()
	if err != nil {
		return err
	}

	return nil

}

func (engine *RenderEngine) ArticleLists() error {

	articleCount := len(engine.articles)

	numPages := int(math.Ceil(float64(articleCount) / float64(pageSize)))
	fmt.Printf("Articles:\t%d\n", articleCount)
	fmt.Printf("Page size:\t%d\n", pageSize)
	fmt.Printf("Pages:\t%d\n", numPages)

	sort.Slice(engine.articles, func(i, j int) bool {
		return engine.articles[i].PublishedAt > engine.articles[j].PublishedAt
	})

	for page := 0; page < numPages; page++ {
		start := page * pageSize
		end := (page + 1) * pageSize
		pageArticles := engine.articles[start:end]

		err := engine.articleListsPage(page, pageArticles)
		if err != nil {
			return err
		}

	}

	return nil
}

func (engine *RenderEngine) articleListsPage(page int, articles []*article.Article) error {

	err := os.MkdirAll(filepath.Join(engine.outputDir, "page"), os.ModePerm)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(engine.getFilePath(fmt.Sprintf("page/%d.html", page)), os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Println(articles[0].Title, articles[0].Description)

	pageData := PageData{
		Title: "test",
		Data:  articles,
	}

	err = engine.templates.ExecuteTemplate(f, "articleList.html", pageData)
	if err != nil {
		return err
	}

	return nil
}

func (engine *RenderEngine) getFilePath(file string) string {
	return filepath.Join(engine.outputDir, file)

}
