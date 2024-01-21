package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func (engine *RenderEngine) postHumanClassification(w http.ResponseWriter, r *http.Request) {

	url := r.FormValue("url")
	classification := r.FormValue("classification")
	// redir := r.FormValue("redir")

	class, err := strconv.ParseInt(classification, 10, 32)
	if err != nil {
		fmt.Println(err)
	}

	_, err = engine.db.Exec("UPDATE articles SET humanClassification = ? WHERE url = ?", class, url)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Updated %s to %d\n", url, class)

	// http.Redirect(w, r, redir, http.StatusFound)

}

func (engine *RenderEngine) runHttp() {
	http.HandleFunc("/human", engine.postHumanClassification)

	fs := http.FileServer(http.Dir("./output"))
	http.Handle("/", fs)

	http.ListenAndServe(":8080", nil)
}
