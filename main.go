package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		fmt.Println("Parsing template once ...")

		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})

	fmt.Println("Executing template ...")
	t.templ.Execute(w, nil)
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.Handle("/", &templateHandler{filename: "intro.html"})
	http.HandleFunc("/contact", func(w http.ResponseWriter, r *http.Request) {

		// gather form data
		//
	})

	fmt.Println("Starting web server ...")
	if err := http.ListenAndServe(":3030", nil); err != nil {
		log.Fatal("Server problems: ", err)
	}
}
