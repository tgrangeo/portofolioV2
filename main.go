package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/joho/godotenv"
)

type GitHubFileContent struct {
	Content string `json:"content"`
}

type GitHubRepo struct {
	Name string `json:"name"`
}

func getRepos(user string) ([]string, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/repos", user)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching repos: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var repos []GitHubRepo
	err = json.Unmarshal(body, &repos)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	var wg sync.WaitGroup
	repoNames := make([]string, 0)
	repoNamesMutex := &sync.Mutex{}

	for _, repo := range repos {
		wg.Add(1)
		go func(repo GitHubRepo) {
			defer wg.Done()
			if hasReadme(user, repo.Name) {
				repoNamesMutex.Lock()
				repoNames = append(repoNames, repo.Name)
				repoNamesMutex.Unlock()
			}
		}(repo)
	}

	wg.Wait()
	return repoNames, nil
}

func hasReadme(username, repo string) bool {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/README.md", username, repo)
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		log.Printf("error creating request: %v", err)
		return false
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		req.Header.Set("Authorization", "token "+token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error fetching repo README: %v", err)
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func fetchRepoReadme(username string, repo string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/readme", username, repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		req.Header.Set("Authorization", "token "+token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error fetching repo README: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch repo README: %s", resp.Status)
	}

	var fileContent GitHubFileContent
	if err := json.NewDecoder(resp.Body).Decode(&fileContent); err != nil {
		return "", fmt.Errorf("error decoding repo README: %w", err)
	}

	decodedContent, err := base64.StdEncoding.DecodeString(fileContent.Content)
	if err != nil {
		return "", fmt.Errorf("error decoding base64 content: %w", err)
	}

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	mdParser := parser.NewWithExtensions(extensions)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	renderer := html.NewRenderer(html.RendererOptions{Flags: htmlFlags})

	md := markdown.ToHTML(decodedContent, mdParser, renderer)

	return string(md), nil
}

func fetchUserReadme(username string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/readme", username, username)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		req.Header.Set("Authorization", "token "+token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error fetching user README: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch user README: %s", resp.Status)
	}

	var fileContent GitHubFileContent
	if err := json.NewDecoder(resp.Body).Decode(&fileContent); err != nil {
		return "", fmt.Errorf("error decoding user README: %w", err)
	}

	decodedContent, err := base64.StdEncoding.DecodeString(fileContent.Content)
	if err != nil {
		return "", fmt.Errorf("error decoding base64 content: %w", err)
	}

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	mdParser := parser.NewWithExtensions(extensions)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	renderer := html.NewRenderer(html.RendererOptions{Flags: htmlFlags})

	md := markdown.ToHTML(decodedContent, mdParser, renderer)

	return string(md), nil
}

var tmpl *template.Template

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
		log.Printf("Error fetching repos: %v", err)
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
	builder.WriteString("<div id=\"project-readme\"></div>")
	builder.WriteString("</div>")

	// Write HTML to response
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(builder.String()))
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	tmpl = template.Must(template.ParseGlob("views/*.html"))
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

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "index.html", nil)
	})
	http.HandleFunc("/projects", ShowProjects)
	http.HandleFunc("/about", ShowAbout)
	http.HandleFunc("/readme/", ShowProjectReadme)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("🚀 Starting up on port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
