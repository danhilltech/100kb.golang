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
	"strings"
	"text/template"
	"time"

	"github.com/danhilltech/100kb.golang/pkg/article"
	"github.com/danhilltech/100kb.golang/pkg/domain"
	"github.com/danhilltech/100kb.golang/pkg/scorer"
)

var categories = []string{
	"technology",
	"life",
	"family",
	"science",
	"politics",
	"news",
	"programming",
	"food",
	"investing",
	"management",
	"nature",
}

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
	TagTitle string
	Data     []*article.Article
	Page     int
	PrevPage string
	NextPage string

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
	totalDomains := len(engine.domains)

	engine.log.Printf("Articles:\t%d\n", articleCount)
	engine.log.Printf("Page size:\t%d\n", pageSize)

	sort.Slice(goodArticles, func(i, j int) bool {
		return goodArticles[i].PublishedAt > goodArticles[j].PublishedAt
	})

	err := engine.buildCategory(goodArticles, "", articleCount, totalDomains)
	if err != nil {
		return err
	}

	for _, cat := range categories {

		catArticles := []*article.Article{}

		for _, a := range goodArticles {
			tgs := a.GetZeroShot()
			for _, t := range tgs {
				if t == cat {
					catArticles = append(catArticles, a)
				}
			}
		}

		sort.Slice(catArticles, func(i, j int) bool {
			return catArticles[i].PublishedAt > catArticles[j].PublishedAt
		})

		catArticleCount := len(catArticles)

		err := engine.buildCategory(catArticles, cat, catArticleCount, totalDomains)
		if err != nil {
			return err
		}
	}

	err = engine.aboutPage(0, goodArticles, articleCount, totalDomains)
	if err != nil {
		return err
	}

	return nil
}

func (engine *RenderEngine) buildCategory(articles []*article.Article, tag string, articleCount, totalDomains int) error {

	numPages := int(math.Ceil(float64(len(articles)) / float64(pageSize)))

	engine.log.Printf("Tag: %s Pages:\t%d\n", tag, numPages)
	k := ""
	pos := 1

	for _, a := range articles {
		d := time.Unix(a.PublishedAt, 0)

		if k != d.Format("2006-01-02") {
			pos = 1
			k = d.Format("2006-01-02")
		}

		a.DayPosition = pos

		pos++

	}

	for page := 0; page < numPages-1; page++ {
		start := page * pageSize
		end := (page + 1) * pageSize
		pageArticles := articles[start:end]

		err := engine.articleListsPage(page, tag, pageArticles, articleCount, totalDomains)
		if err != nil {
			return err
		}
	}
	return nil
}

func buildPagePath(tag string, page int) string {

	fullTag := "/"
	if tag != "" {
		fullTag = fmt.Sprintf("/%s", tag)
	}

	if page <= 0 {
		return "/"
	}

	pageSegment := "/"
	if page > 0 {
		pageSegment = "/page"
	}

	fullPage := fmt.Sprintf("/%d.html", page)
	if page == 0 {
		fullPage = ""
	}

	path := fmt.Sprintf("%s%s%s", fullTag, pageSegment, fullPage)

	path = strings.ReplaceAll(path, "//", "/")

	return path
}

func buildPageFilePath(tag string, page int) string {

	fullTag := "/"
	if tag != "" {
		fullTag = fmt.Sprintf("/%s", tag)
	}

	pageSegment := "/"
	if page > 0 {
		pageSegment = "/page"
	}

	fullPage := fmt.Sprintf("/%d.html", page)
	if page == 0 {
		fullPage = "/index.html"
	}

	path := fmt.Sprintf("%s%s%s", fullTag, pageSegment, fullPage)

	path = strings.ReplaceAll(path, "//", "/")
	path = strings.ReplaceAll(path, "//", "/")

	return path
}

func (engine *RenderEngine) articleListsPage(page int, tag string, articles []*article.Article, totalArticles int, totalDomains int) error {

	filePath := buildPageFilePath(tag, page)

	err := os.MkdirAll(filepath.Join(engine.outputDir, filepath.Dir(filePath)), os.ModePerm)
	if err != nil {
		return fmt.Errorf("mkdir %s %w", filepath.Join(engine.outputDir, filepath.Dir(filePath)), err)
	}

	f, err := os.Create(engine.getFilePath(filePath))
	if err != nil {
		return fmt.Errorf("create %s %w", engine.getFilePath(filePath), err)
	}
	defer f.Close()

	title := fmt.Sprintf("Page %d ~ 100kb", page)
	if tag != "" {
		title = fmt.Sprintf("Writing about %s. Page %d ~ 100kb", tag, page)
	}

	tagTitle := ""
	if tag != "" {
		tagTitle = fmt.Sprintf("Writing about %s", tag)
	}

	prevPage := buildPagePath(tag, page-1)

	nextPage := buildPagePath(tag, page+1)

	pageData := ArticleListData{
		Title:         title,
		TagTitle:      tagTitle,
		Data:          articles,
		Page:          page,
		NextPage:      nextPage,
		PrevPage:      prevPage,
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
		NextPage:      "",
		PrevPage:      "",
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
