package main

import (
	"fmt"
	"net/http"
)

func (engine *RenderEngine) runHttp() {

	fmt.Println("Starting output http server...")

	fs := http.FileServer(http.Dir("./output"))
	http.Handle("/", fs)

	http.ListenAndServe(":8080", nil)
}
