package main

import (
	_ "embed"
	"html/template"
	"log"
	"net/http"
)

//go:embed index.html.tmpl
var indexHtml string

func indexHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("indexHandler has been called")

	s := initiateLoginIfRequired(w, req)
	if s == nil {
		return
	}

	indexHtmlTmpl, err := template.New("indexhtml").Parse(indexHtml)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = indexHtmlTmpl.Execute(w, sessionsStore[s.ID])
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
