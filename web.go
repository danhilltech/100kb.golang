package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"

	"github.com/danhilltech/100kb.golang/pkg/article"
)

var trackFile = "output/scored.csv"

type ScoreRequest struct {
	URL   string
	Score int
}

func (engine *RenderEngine) runHttp() {

	fmt.Println("Starting output http server...")

	http.HandleFunc("/score", engine.handleScore)

	fs := http.FileServer(http.Dir("./output"))
	http.Handle("/", fs)

	http.ListenAndServe(":8080", nil)
}

func (engine *RenderEngine) handleScore(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	score := r.PostFormValue("score")
	url := r.PostFormValue("url")

	f, err := os.OpenFile(trackFile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()
	if _, err := f.WriteString(fmt.Sprintf("%s,%s\n", url, score)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	articlesFiltered := []*article.Article{}

	for _, a := range engine.articles {
		if a.FirstPersonRatio > 0.03 && a.WordCount > 200 && a.WordCount < 2000 && a.BadCount < 50 {
			articlesFiltered = append(articlesFiltered, a)
		}
	}

	// rand.Seed(time.Now().Unix())
	article := articlesFiltered[rand.Intn(len(articlesFiltered))]

	http.Redirect(w, r, fmt.Sprintf("/article/%s", article.GetSlug()), http.StatusMovedPermanently)
}
