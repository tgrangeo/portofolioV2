package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"portofolio/src"
	"time"

	"strings"

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
		// w.Header().Set("HX-Redirect", "/")
		tmpl.ExecuteTemplate(w, "home.html", nil)
	})

	http.HandleFunc("/user/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 *time.Millisecond)
		path := r.URL.Path
		userID := strings.TrimPrefix(path, "/user/")
		if userID == "" {
			http.Error(w, "User ID is required", http.StatusBadRequest)
			return
		}
		os.Setenv("USERNAME", userID)
		w.Header().Set("HX-Redirect", "/user/tgrangeo")
		err := tmpl.ExecuteTemplate(w, "index.html", nil)
		if err != nil {
			http.Error(w, "Template rendering error", http.StatusInternalServerError)
			log.Println("Template error:", err)
		}
	})
	http.HandleFunc("/contact", src.ContactHandler)
	http.HandleFunc("/submit", src.HandleSubmit)
	http.HandleFunc("/projects", src.ShowProjects)
	http.HandleFunc("/about", src.ShowAbout)
	http.HandleFunc("/readme/", src.ShowProjectReadme)
	http.HandleFunc("/profile-picture", src.ProfilePictureHandler)
	http.HandleFunc("/username", src.GetUsername)
	http.HandleFunc("/submit-username", src.SetUsername)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	log.Println("🚀 Starting up on port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
