package output

import (
	"fmt"
	"net/http"
)

// var trackFile = "scoring/scored.csv"

type ScoreRequest struct {
	URL   string
	Score int
}

func (engine *RenderEngine) RunHttp() {

	fmt.Println("Starting output http server...")

	fs := http.FileServer(http.Dir("./output-train"))
	http.Handle("/", fs)

	http.ListenAndServe(":8081", nil)
}

// func (engine *RenderEngine) handleScore(w http.ResponseWriter, r *http.Request) {
// 	err := r.ParseForm()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	score := r.PostFormValue("score")
// 	url := r.PostFormValue("url")

// 	f, err := os.OpenFile(trackFile,
// 		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	defer f.Close()
// 	if _, err := f.WriteString(fmt.Sprintf("%s,%s\n", url, score)); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// articlesFiltered := []*article.Article{}

// 	// for _, a := range engine.articles {
// 	// 	if a.FirstPersonRatio > 0.03 && a.WordCount > 200 && a.WordCount < 2000 && a.BadCount < 50 {
// 	// 		articlesFiltered = append(articlesFiltered, a)
// 	// 	}
// 	// }

// 	// rand.Seed(time.Now().Unix())
// 	// article := articlesFiltered[rand.Intn(len(articlesFiltered))]
// 	article := engine.articles[rand.Intn(len(engine.articles))]

// 	http.Redirect(w, r, fmt.Sprintf("/article/%s", article.GetSlug()), http.StatusMovedPermanently)
// }
