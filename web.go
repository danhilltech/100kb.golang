package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
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

	// rand.Seed(time.Now().Unix())
	article := engine.articles[rand.Intn(len(engine.articles))]

	http.Redirect(w, r, fmt.Sprintf("/article/%s", article.GetSlug()), http.StatusMovedPermanently)
}
