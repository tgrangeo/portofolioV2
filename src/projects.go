package src

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type Project struct{
	Title string
	Lang []string
	Desc string
}

type Projects_Data struct {
	Projects []Project
}

func GetProjects(w http.ResponseWriter, r *http.Request) {
    repos, err := getRepos(username)
    if err != nil {
        log.Printf("Error fetching repos: %v", err)
        http.Error(w, "Error fetching repos", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(repos)
}

func ShowProjects(w http.ResponseWriter, r *http.Request) string {
	htmlFile, err := os.ReadFile("views/projects.html")
	if err != nil {
			http.Error(w, "Unable to read HTML file", http.StatusInternalServerError)
			return ""
		}
		w.Header().Set("Content-Type", "text/html")
		return string(htmlFile)
	}

func ShowProjectReadme(w http.ResponseWriter, r *http.Request) string {
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
	return string("<div class=\"content-readme\"" + readme + "</div>")
}
