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

var username = "tgrangeo"

func GetUsername(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(username))
}

func SetUsername(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	n := strings.Split(path, "/")[2]
	username = n
}

func ShowBrowse(w http.ResponseWriter, r *http.Request) string {
	htmlFile, err := os.ReadFile("views/browse.html")
	if err != nil {
		http.Error(w, "Unable to read HTML file", http.StatusInternalServerError)
		return ""
	}
	w.Header().Set("Content-Type", "text/html")
	return string(htmlFile)
}

func ShowAbout(w http.ResponseWriter, r *http.Request) string {
	readme, err := fetchUserReadme(username)
	if err != nil {
		log.Printf("Error fetching user README: %v", err)
		http.Error(w, fmt.Sprintf("Error fetching README: %v", err), http.StatusInternalServerError)
		return ""
	}
	return string("<div class=\"content-readme\"" + readme + "</div>")
}

func ShowProjects(w http.ResponseWriter, r *http.Request) string {
	repos, err := getRepos(username)
	if err != nil {
		log.Printf("Error fetching repos : %v", err)
		http.Error(w, fmt.Sprintf("Error fetching repos: %v", err), http.StatusInternalServerError)
		return ""
	}
	var builder strings.Builder
	builder.WriteString("<div class=\"projects\"><button class=\"side-project-button\" onClick=\"openProjectSide()\">Select a project</button>")

	builder.WriteString("<div id=\"side\" class=\"side\"><ul>")
	for _, repo := range repos {
		builder.WriteString(fmt.Sprintf(`<li><button onclick="closeProjectSide()" hx-get="/readme/%s" hx-target="#project-readme" hx-swap="innerHTML">%s</button><button onclick="linkToGithub('%s')"><img class="github-icon" src="../static/github-black.png" width="18px" /></button></li>`, repo, repo, repo))
	}
	builder.WriteString("</ul></div>")
	builder.WriteString("<div id=\"project-readme\">select a project to begin</div>")
	builder.WriteString("</div>")
	w.Header().Set("Content-Type", "text/html")
	return builder.String()
}

func ShowProjectReadme(w http.ResponseWriter, r *http.Request) string{
	path := r.URL.Path
	segments := strings.Split(path, "/")
	if len(segments) < 3 || segments[2] == "" {
		http.Error(w, "Name parameter is required", http.StatusBadRequest)
		return ""
	}
	name := segments[2]

	readme, err := fetchRepoReadme(username, name)
	if err != nil {
		log.Printf("Error fetching repo README: %v", err)
		http.Error(w, fmt.Sprintf("Error fetching README: %v", err), http.StatusInternalServerError)
		return ""
	}
	w.Header().Set("Content-Type", "text/html")
	return readme
}

func ProfilePictureHandler(w http.ResponseWriter, r *http.Request) {
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