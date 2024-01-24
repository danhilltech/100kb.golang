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
	db        *sql.DB

	articles []*article.Article
}

type PageData struct {
	Title string
	Data  interface{}
}

var pageSize = 100

//go:embed views/*.html
var tmplFS embed.FS

func NewRenderEnding(outputDir string, articles []*article.Article, db *sql.DB) *RenderEngine {
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
		db:        db,
	}
}

func CreateOutput(db *sql.DB, cacheDir string) error {

	articleEngine, err := article.NewEngine(db, cacheDir)
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

	engine := NewRenderEnding("output", articles, db)

	err = engine.ArticleLists()
	if err != nil {
		return err
	}

	engine.runHttp()

	return nil

}

func (engine *RenderEngine) ArticleLists() error {

	articlesFiltered := []*article.Article{}

	for _, a := range engine.articles {
		if a.FirstPersonRatio > 0.02 {
			articlesFiltered = append(articlesFiltered, a)
		}
	}

	articleCount := len(articlesFiltered)

	numPages := int(math.Ceil(float64(articleCount) / float64(pageSize)))
	fmt.Printf("Articles:\t%d\n", articleCount)
	fmt.Printf("Page size:\t%d\n", pageSize)
	fmt.Printf("Pages:\t%d\n", numPages)

	sort.Slice(articlesFiltered, func(i, j int) bool {
		return articlesFiltered[i].Score() > articlesFiltered[j].Score()
	})

	for page := 0; page < numPages-1; page++ {
		start := page * pageSize
		end := (page + 1) * pageSize
		pageArticles := articlesFiltered[start:end]

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
