package src

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func ShowAbout(w http.ResponseWriter, r *http.Request) {
	readme, err := fetchUserReadme("tgrangeo")
	if err != nil {
		log.Printf("Error fetching user README: %v", err)
		http.Error(w, fmt.Sprintf("Error fetching README: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("<div style=\"width:60vw\">"))
	w.Write([]byte(readme))
	w.Write([]byte("</div>"))
}

func ShowProjects(w http.ResponseWriter, r *http.Request) {
	repos, err := getRepos("tgrangeo")
	if err != nil {
		log.Printf("Error fetching repos : %v", err)
		http.Error(w, fmt.Sprintf("Error fetching repos: %v", err), http.StatusInternalServerError)
		return
	}

	// Generate HTML for each repo
	var builder strings.Builder
	builder.WriteString("<div class=\"projects\">")
	builder.WriteString("<div class=\"side\"><ul>")
	for _, repo := range repos {
		builder.WriteString(fmt.Sprintf("<li><button hx-get=\"/readme/%s\" hx-target=\"#project-readme\" hx-swap=\"innerHTML\">%s</button></li>", repo, repo))
	}
	builder.WriteString("</ul></div>")
	builder.WriteString("<div id=\"project-readme\"><div style=\"display:flex;height:100%;align-items:center;justify-content:center\">select a project to begin</div></div>")
	builder.WriteString("</div>")

	// Write HTML to response
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(builder.String()))
}

func ShowProjectReadme(w http.ResponseWriter, r *http.Request) {
	// Extract the path parameter manually
	path := r.URL.Path
	segments := strings.Split(path, "/")
	if len(segments) < 3 || segments[2] == "" {
		http.Error(w, "Name parameter is required", http.StatusBadRequest)
		return
	}
	name := segments[2]

	readme, err := fetchRepoReadme("tgrangeo", name)
	if err != nil {
		log.Printf("Error fetching repo README: %v", err)
		http.Error(w, fmt.Sprintf("Error fetching README: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(readme))
}
