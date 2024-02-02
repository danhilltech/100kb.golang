package main

import (
	"database/sql"
	"embed"
	"encoding/csv"
	"fmt"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"text/template"

	"github.com/danhilltech/100kb.golang/pkg/article"
)

type RenderEngine struct {
	templates map[string]*template.Template
	outputDir string
	db        *sql.DB

	articles []*article.Article
}

type PageData struct {
	Title string
	Data  interface{}
}

var pageSize = 100

//go:embed views/*.html views/layouts/*.html
var tmplFS embed.FS

//go:embed views/static/*
var staticFS embed.FS

func NewRenderEnding(outputDir string, articles []*article.Article, db *sql.DB) (*RenderEngine, error) {
	// funcMap := template.FuncMap{
	// 	// "inc": inc,
	// }

	templates := make(map[string]*template.Template)

	tmplFiles, err := fs.ReadDir(tmplFS, "views")
	if err != nil {
		return nil, err
	}

	for _, tmpl := range tmplFiles {
		if tmpl.IsDir() {
			continue
		}

		pt, err := template.ParseFS(tmplFS, "views/"+tmpl.Name(), "views/layouts/*html")
		if err != nil {
			return nil, err
		}

		templates[tmpl.Name()] = pt
	}
	return &RenderEngine{
		templates: templates,
		articles:  articles,
		outputDir: outputDir,
		db:        db,
	}, nil
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

	engine, err := NewRenderEnding("output", articles, db)
	if err != nil {
		return err
	}

	err = engine.WriteCSV()
	if err != nil {
		return err
	}

	err = engine.ArticleLists()
	if err != nil {
		return err
	}

	err = engine.StaticFiles()
	if err != nil {
		return err
	}

	engine.runHttp()

	return nil

}

func (engine *RenderEngine) WriteCSV() error {
	file, err := os.Create(engine.getFilePath("articles.csv"))
	if err != nil {
		return err
	}
	defer file.Close()

	csvwriter := csv.NewWriter(file)

	var data [][]string

	row := []string{
		"url",
		"domain",
		"feedUrl",
		"title",
		"wordCount",
		"pCount",
		"h1Count",
		"hnCount",
		"badCount",
		"fpr",
	}
	data = append(data, row)

	for _, a := range engine.articles {
		row := []string{
			a.Url,
			a.Domain,
			a.FeedUrl,
			a.Title,
			strconv.Itoa(int(a.WordCount)),
			strconv.Itoa(int(a.PCount)),
			strconv.Itoa(int(a.H1Count)),
			strconv.Itoa(int(a.HNCount)),
			strconv.Itoa(int(a.BadCount)),
			strconv.FormatFloat(a.FirstPersonRatio, 'f', 4, 64),
		}
		data = append(data, row)
	}
	return csvwriter.WriteAll(data)
}

func (engine *RenderEngine) StaticFiles() error {
	files, err := staticFS.ReadDir("views/static")
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Join(engine.outputDir, "static"), os.ModePerm)
	if err != nil {
		return err
	}

	for _, f := range files {
		bytes, err := staticFS.ReadFile("views/static/" + f.Name())
		if err != nil {
			return err
		}

		f, err := os.Create(engine.getFilePath(fmt.Sprintf("static/%s", f.Name())))
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = f.Write(bytes)
		if err != nil {
			return err
		}

	}

	return nil
}

func (engine *RenderEngine) ArticleLists() error {

	// articlesFiltered := []*article.Article{}

	// for _, a := range engine.articles {
	// 	if a.FirstPersonRatio > 0.02 && a.Score() > 0 {
	// 		articlesFiltered = append(articlesFiltered, a)
	// 	}
	// }

	articleCount := len(engine.articles)

	numPages := int(math.Ceil(float64(articleCount) / float64(pageSize)))
	fmt.Printf("Articles:\t%d\n", articleCount)
	fmt.Printf("Page size:\t%d\n", pageSize)
	fmt.Printf("Pages:\t%d\n", numPages)

	sort.Slice(engine.articles, func(i, j int) bool {
		return engine.articles[i].Score() > engine.articles[j].Score()
	})

	for page := 0; page < numPages-1; page++ {
		start := page * pageSize
		end := (page + 1) * pageSize
		pageArticles := engine.articles[start:end]

		err := engine.articleListsPage(page, pageArticles)
		if err != nil {
			return err
		}
	}

	for _, article := range engine.articles {
		err := engine.articlePage(article)
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

	f, err := os.Create(engine.getFilePath(fmt.Sprintf("page/%d.html", page)))
	if err != nil {
		return err
	}
	defer f.Close()

	pageData := PageData{
		Title: "test",
		Data:  articles,
	}

	err = engine.templates["articleList.html"].Execute(f, pageData)
	if err != nil {
		// return err
		fmt.Println("DAN", err)
	}

	return nil
}

func (engine *RenderEngine) articlePage(article *article.Article) error {

	err := os.MkdirAll(filepath.Join(engine.outputDir, "article"), os.ModePerm)
	if err != nil {
		return err
	}

	f, err := os.Create(engine.getFilePath(fmt.Sprintf("article/%s", article.GetSlug())))
	if err != nil {
		return err
	}
	defer f.Close()

	pageData := PageData{
		Title: article.Title,
		Data:  article,
	}

	err = engine.templates["articleInfo.html"].Execute(f, pageData)
	if err != nil {
		return err
	}

	return nil
}

func (engine *RenderEngine) getFilePath(file string) string {
	return filepath.Join(engine.outputDir, file)

}
