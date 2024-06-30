package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/joho/godotenv"
)

type GitHubFileContent struct {
	Content string `json:"content"`
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var repos []struct {
		Name string `json:"name"`
	}
	err = json.Unmarshal(body, &repos)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response: %w", err)
	}

	var repoNames []string
	for _, repo := range repos {
		repoNames = append(repoNames, repo.Name)
	}

	return repoNames, nil
}

func fetchRepoReadme(username string, repo string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/README.md", username, repo)
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
	w.Write([]byte(readme))
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
	builder.WriteString("<ul>")
	for i, repo := range repos {
		builder.WriteString(fmt.Sprintf("<li>%d <button hx-get=\"/readme/%s\">%s<button></li>", i, repo, repo))
	}
	builder.WriteString("</ul>")

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
	http.HandleFunc("/readme/", ShowProjectReadme) // Changed to handle parameter

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("🚀 Starting up on port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
