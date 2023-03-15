package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

const (
	tmplPath   = "public/templates"
	tmplSuffix = ".html.tmpl"
)

func main() {
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))
	http.HandleFunc("/", serveTemplate)
	log.Print("Listening on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println(err)
	}
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Redirect(w, r, "home", http.StatusMovedPermanently)
		return
	}

	basePath := filepath.Join(tmplPath, "base"+tmplSuffix)
	sidebarPath := filepath.Join(tmplPath, "sidebar"+tmplSuffix)
	filePath := filepath.Join(tmplPath, filepath.Clean(r.URL.Path+tmplSuffix))

	// Return a 404 if the template doesn't exist
	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
	}

	// Return a 404 if the request is for a directory
	if info.IsDir() {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles(basePath, sidebarPath, filePath)
	if err != nil {
		log.Println(err)
	}

	tmpl.ExecuteTemplate(w, "base", nil)
}
