package main

import (
	"html/template"
	"log"
	"net/http"
	"portofolio/src"

	"github.com/joho/godotenv"
)

var tmpl *template.Template

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	tmpl = template.Must(template.ParseGlob("views/*.html"))
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "index.html", nil)
	})
	http.HandleFunc("/projects", src.ShowProjects)
	http.HandleFunc("/about", src.ShowAbout)
	http.HandleFunc("/readme/", src.ShowProjectReadme)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("🚀 Starting up on port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
