package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"portofolio/src"
	"time"

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

type Data struct {
	Content template.HTML
}

func middleware(next func(http.ResponseWriter, *http.Request) string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", time.Unix(0, 0).Format(http.TimeFormat))
		if r.Header.Get("HX-Request") != "true" {
			data := Data{
				Content: template.HTML(next(w, r)),
			}
			tmpl.ExecuteTemplate(w, "index.html", data)
		} else {
			w.Write([]byte(next(w, r)))
		}
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/browse", http.StatusSeeOther)
	})
	http.HandleFunc("/about", middleware(src.ShowAbout))
	http.HandleFunc("/contact", middleware(src.ContactHandler))
	http.HandleFunc("/submit", middleware(src.HandleSubmit))
	http.HandleFunc("/projects", middleware(src.ShowProjects))
	http.HandleFunc("/browse", middleware(src.ShowBrowse))
	http.HandleFunc("/readme/", middleware(src.ShowProjectReadme))
	http.HandleFunc("/blog", middleware(src.ShowBlog))
	http.HandleFunc("/article/", middleware(src.ShowArticle))
	http.HandleFunc("/profile-picture", src.ProfilePictureHandler)
	http.HandleFunc("/username", src.GetUsername)
	http.HandleFunc("/new-username/", src.SetUsername)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	log.Println("🚀 Starting up on port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
