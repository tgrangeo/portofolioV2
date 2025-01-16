package src

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"
)

type Project struct{
	Title string
	Lang string
	Desc string
}

type Projects_Data struct {
	Projects []Project
}

func ShowProjects(w http.ResponseWriter, r *http.Request) string {
	repos, err := getRepos(username)
	if err != nil {
		log.Printf("Error fetching repos : %v", err)
		http.Error(w, fmt.Sprintf("Error fetching repos: %v", err), http.StatusInternalServerError)
		return ""
	}
	data := Projects_Data{
		Projects: repos,
	}
	tmpl, err := template.ParseFiles("views/projects.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return ""
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return ""
	}
	return buf.String()
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
	return readme
}
