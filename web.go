package main

import (
	"net/http"
)

func (engine *RenderEngine) runHttp() {

	fs := http.FileServer(http.Dir("./output"))
	http.Handle("/", fs)

	http.ListenAndServe(":8080", nil)
}
