package output

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"text/template"
	"time"

	"github.com/danhilltech/100kb.golang/pkg/article"
	"github.com/danhilltech/100kb.golang/pkg/domain"
	"github.com/danhilltech/100kb.golang/pkg/scorer"
)

type RenderEngine struct {
	templates     map[string]*template.Template
	outputDir     string
	db            *sql.DB
	articleEngine *article.Engine
	log           *log.Logger

	articles []*article.Article
	domains  []*domain.Domain
	model    *scorer.LogisticModel
}

type ArticleListData struct {
	Title    string
	Data     []*article.Article
	Page     int
	PrevPage int
	NextPage int

	TotalArticles int
	TotalDomains  int

	GenDate string
}

var pageSize = 25

//go:embed views/*.html views/layouts/*.html
var tmplFS embed.FS

//go:embed views/static/*
var staticFS embed.FS

func NewRenderEngine(log *log.Logger, outputDir string, articles []*article.Article, domains []*domain.Domain, model *scorer.LogisticModel, db *sql.DB, articleEngine *article.Engine) (*RenderEngine, error) {

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
		templates:     templates,
		articles:      articles,
		domains:       domains,
		outputDir:     outputDir,
		articleEngine: articleEngine,
		model:         model,
		db:            db,
		log:           log,
	}, nil
}

func (engine *RenderEngine) Prepare() error {
	domainScores := make(map[string]float64, len(engine.domains))

	for _, a := range engine.articles {
		for _, d := range engine.domains {
			if d.Domain == a.Domain {
				if d.Articles == nil {
					d.Articles = []*article.Article{}
				}
				d.Articles = append(d.Articles, a)
			}
		}
	}

	// Prepare

	names := engine.domains[0].GetFloatFeatureNames()

	for _, d := range engine.domains {
		fts := d.GetFloatFeatures()
		score := engine.model.Predict([][]float64{fts}, names)
		if len(d.Articles) == 0 {
			d.LiveScore = 0
		} else {
			d.LiveScore = score[0]
		}

		domainScores[d.Domain] = d.LiveScore
	}

	for _, a := range engine.articles {
		a.DomainScore = domainScores[a.Domain]
	}

	return nil
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

	var goodArticles []*article.Article

	now := time.Now().Unix()

	for _, a := range engine.articles {
		if a.DomainScore > 0.5 && a.PublishedAt < now {
			goodArticles = append(goodArticles, a)
		}
	}

	articleCount := len(goodArticles)

	numPages := int(math.Ceil(float64(articleCount) / float64(pageSize)))
	engine.log.Printf("Articles:\t%d\n", articleCount)
	engine.log.Printf("Page size:\t%d\n", pageSize)
	engine.log.Printf("Pages:\t%d\n", numPages)

	sort.Slice(goodArticles, func(i, j int) bool {
		return goodArticles[i].PublishedAt > goodArticles[j].PublishedAt
	})

	k := ""
	pos := 1
	for _, a := range goodArticles {
		d := time.Unix(a.PublishedAt, 0)

		a.DayPosition = pos

		if k != d.Format("2006-01-02") {
			pos = 1
			k = d.Format("2006-01-02")
		}

		pos++
	}

	totalDomains := len(engine.domains)

	for page := 0; page < numPages-1; page++ {
		start := page * pageSize
		end := (page + 1) * pageSize
		pageArticles := goodArticles[start:end]

		err := engine.articleListsPage(page, pageArticles, articleCount, totalDomains)
		if err != nil {
			return err
		}
	}

	err := engine.aboutPage(0, goodArticles, articleCount, totalDomains)
	if err != nil {
		return err
	}

	return nil
}

func (engine *RenderEngine) articleListsPage(page int, articles []*article.Article, totalArticles int, totalDomains int) error {
	err := os.MkdirAll(filepath.Join(engine.outputDir, "page"), os.ModePerm)
	if err != nil {
		return err
	}

	f, err := os.Create(engine.getFilePath(fmt.Sprintf("page/%d.html", page)))
	if err != nil {
		return err
	}
	defer f.Close()

	pageData := ArticleListData{
		Title:         fmt.Sprintf("Page %d ~ 100kb", page),
		Data:          articles,
		Page:          page,
		NextPage:      page + 1,
		PrevPage:      page - 1,
		TotalArticles: totalArticles,
		TotalDomains:  totalDomains,
		GenDate:       time.Now().Format(time.RFC1123),
	}

	err = engine.templates["articleList.html"].Execute(f, pageData)
	if err != nil {
		// return err
		engine.log.Println(err)
	}

	return nil
}

func (engine *RenderEngine) aboutPage(page int, articles []*article.Article, totalArticles int, totalDomains int) error {

	f, err := os.Create(engine.getFilePath("about.html"))
	if err != nil {
		return err
	}
	defer f.Close()

	pageData := ArticleListData{
		Title:         "About",
		Data:          articles,
		Page:          page,
		NextPage:      page + 1,
		PrevPage:      page - 1,
		TotalArticles: totalArticles,
		TotalDomains:  totalDomains,
		GenDate:       time.Now().Format(time.RFC1123),
	}

	err = engine.templates["about.html"].Execute(f, pageData)
	if err != nil {
		// return err
		engine.log.Println(err)
	}

	return nil
}

func (engine *RenderEngine) getFilePath(file string) string {
	return filepath.Join(engine.outputDir, file)

}
