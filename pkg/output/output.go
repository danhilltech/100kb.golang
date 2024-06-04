package output

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"sort"
	"text/template"

	"github.com/danhilltech/100kb.golang/pkg/article"
	"github.com/danhilltech/100kb.golang/pkg/domain"
)

type RenderEngine struct {
	templates     map[string]*template.Template
	outputDir     string
	db            *sql.DB
	articleEngine *article.Engine

	articles []*article.Article
	domains  []*domain.Domain
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

func NewRenderEnding(outputDir string, articles []*article.Article, domains []*domain.Domain, db *sql.DB, articleEngine *article.Engine) (*RenderEngine, error) {

	templates := make(map[string]*template.Template)

	tmplFiles, err := fs.ReadDir(tmplFS, "views")
	if err != nil {
		return nil, err
	}

	for _, tmpl := range tmplFiles {
		if tmpl.IsDir() {
			continue
		}

		pt, err := template.New(tmpl.Name()).Funcs(template.FuncMap{
			"GetDomain": func(a article.Article) *domain.Domain {
				for _, d := range domains {
					if d.Domain == a.Domain {
						return d
					}

				}
				return nil
			},
		}).ParseFS(tmplFS, "views/"+tmpl.Name(), "views/layouts/*html")
		if err != nil {
			return nil, err
		}

		templates[tmpl.Name()] = pt
	}
	return &RenderEngine{
		templates:     templates,
		articles:      articles,
		outputDir:     outputDir,
		articleEngine: articleEngine,
		db:            db,
	}, nil
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
		return engine.articles[i].PublishedAt > engine.articles[j].PublishedAt
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

	// for _, article := range engine.articles {
	// 	err := engine.articlePage(article)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

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
		fmt.Println(err)
	}

	return nil
}

func (engine *RenderEngine) getFilePath(file string) string {
	return filepath.Join(engine.outputDir, file)

}
