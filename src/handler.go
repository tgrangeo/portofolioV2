package src

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

var tmpl *template.Template

func GetUsername(w http.ResponseWriter, r *http.Request) {
	username := os.Getenv("USERNAME")
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(username))
}

func SetUsername(w http.ResponseWriter, r *http.Request) {
	os.Setenv("USERNAME","edepauw")
	log.Println(os.Getenv("USERNAME"))
	tmpl.ExecuteTemplate(w, "index.html", nil)
}

func ShowAbout(w http.ResponseWriter, r *http.Request) {
	username := os.Getenv("USERNAME")
	readme, err := fetchUserReadme(username)
	if err != nil {
		log.Printf("Error fetching user README: %v", err)
		http.Error(w, fmt.Sprintf("Error fetching README: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("<div class=\"content-readme\""))
	w.Write([]byte(readme))
	w.Write([]byte("</div>"))
}

func ShowProjects(w http.ResponseWriter, r *http.Request) {
	username := os.Getenv("USERNAME")
	repos, err := getRepos(username)
	if err != nil {
		log.Printf("Error fetching repos : %v", err)
		http.Error(w, fmt.Sprintf("Error fetching repos: %v", err), http.StatusInternalServerError)
		return
	}

	// Generate HTML for each repo
	var builder strings.Builder
	builder.WriteString("<div class=\"projects\"><button class=\"side-project-button\" onClick=\"openProjectSide()\">Select a project</button>")



	builder.WriteString("<div id=\"side\" class=\"side\"><ul>")
	for _, repo := range repos {
		builder.WriteString(fmt.Sprintf("<li><button onclick=\"closeProjectSide()\" hx-get=\"/readme/%s\" hx-target=\"#project-readme\" hx-swap=\"innerHTML\">%s</button></li>", repo, repo))
	}
	builder.WriteString("</ul></div>")
	builder.WriteString("<div id=\"project-readme\">select a project to begin</div>")
	builder.WriteString("</div>")

	// Write HTML to response
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(builder.String()))
}

func ShowProjectReadme(w http.ResponseWriter, r *http.Request) {
	username := os.Getenv("USERNAME")
	// Extract the path parameter manually
	path := r.URL.Path
	segments := strings.Split(path, "/")
	if len(segments) < 3 || segments[2] == "" {
		http.Error(w, "Name parameter is required", http.StatusBadRequest)
		return
	}
	name := segments[2]

	readme, err := fetchRepoReadme(username, name)
	if err != nil {
		log.Printf("Error fetching repo README: %v", err)
		http.Error(w, fmt.Sprintf("Error fetching README: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(readme))
}

func ProfilePictureHandler(w http.ResponseWriter, r *http.Request) {
	username := os.Getenv("USERNAME")

	avatarURL, err := GetProfilePicture(username)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching profile picture: %v", err), http.StatusInternalServerError)
		return
	}
	var builder strings.Builder
	w.Header().Set("Content-Type", "text/html")
	builder.WriteString(fmt.Sprintf("<img class=\"profile-picture\" id=\"profile-picture\" width=\"70px\" src=\"%s\" alt=\"Profile\" disabled>", avatarURL))
	w.Write([]byte(builder.String()))
}

func ContactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(contactTemplate))
}

func HandleSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		form := SubmitForm{
			name: r.FormValue("name"),
			mail: r.FormValue("mail"),
			message: r.FormValue("message"),
		}
		w.Header().Set("Content-Type", "text/html")
		if ContactSend(form){
			w.Write([]byte(contactTemplateSubmitted))
		} else{
			w.Write([]byte(contactTemplateFalse))
		}
		return
	}
	http.Redirect(w, r, "/contact", http.StatusSeeOther)
}	
