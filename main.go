package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"portofolio/src"

	"github.com/joho/godotenv"
)

var tmpl *template.Template

func init() {
	env := os.Getenv("ENV")
	log.Printf("Current ENV variable: %s", env)
	
	if env != "PROD" {
		err := godotenv.Load()
		if err != nil {
			log.Printf("Error loading .env file: %v", err)
		}
	}
	tmpl = template.Must(template.ParseGlob("views/*.html"))
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "index.html", nil)
	})
	http.HandleFunc("/contact", src.ContactHandler)
	http.HandleFunc("/submit", src.HandleSubmit)
	http.HandleFunc("/projects", src.ShowProjects)
	http.HandleFunc("/about", src.ShowAbout)
	http.HandleFunc("/browse", src.ShowBrowse)
	http.HandleFunc("/readme/", src.ShowProjectReadme)
	http.HandleFunc("/profile-picture", src.ProfilePictureHandler)
	http.HandleFunc("/username", src.GetUsername)
	http.HandleFunc("/new-username/", src.SetUsername)
	http.HandleFunc("/blog/", src.ShowBlog)
	http.HandleFunc("/article/", src.ShowArticle)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	log.Println("🚀 Starting up on port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
