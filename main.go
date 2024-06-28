package main

import (
	"encoding/base64"
	"encoding/json"
	"html/template"

	// "io"
	"log"
	"net/http"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type GitHubFileContent struct {
	Content string `json:"content"`
}

func fetchRepoReadme(username string,repo string) (string, error) {
	url := "https://api.github.com/repos/" + username + "/" + repo + "/contents" + "/README.md"
	log.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()


	var fileContent GitHubFileContent
	if err := json.NewDecoder(resp.Body).Decode(&fileContent); err != nil {
		return "", err
	}

	decodedContent, err := base64.StdEncoding.DecodeString(fileContent.Content)
	if err != nil {
		return "", err
	}

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	mdParser := parser.NewWithExtensions(extensions)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	renderer := html.NewRenderer(html.RendererOptions{Flags: htmlFlags})

	md := markdown.ToHTML(decodedContent, mdParser, renderer)

	return string(md), nil
}

func fetchUserReadme(username string) (string, error) {
	url := "https://api.github.com/repos/" + username + "/" + username + "/readme"
	log.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var fileContent GitHubFileContent
	if err := json.NewDecoder(resp.Body).Decode(&fileContent); err != nil {
		return "", err
	}

	decodedContent, err := base64.StdEncoding.DecodeString(fileContent.Content)
	if err != nil {
		return "", err
	}

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	mdParser := parser.NewWithExtensions(extensions)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	renderer := html.NewRenderer(html.RendererOptions{Flags: htmlFlags})

	md := markdown.ToHTML(decodedContent, mdParser, renderer)

	return string(md), nil
}

var tmpl *template.Template

func ShowHomePage(w http.ResponseWriter, r *http.Request) {
	// readme, _ := fetchUserReadme("tgrangeo")
	readme, err := fetchUserReadme("tgrangeo")
	//TODO: error handling here
	if err != nil {
		log.Fatalln(err)
		// return tmpl.HTML(http.StatusInternalServerError, fmt.Sprintf("Error fetching README: %v", err))
	}
	data := map[string]interface{}{
		"Readme": template.HTML(readme),
	}

	tmpl.ExecuteTemplate(w, "index.html", data)
}

func ShowProjectPage(w http.ResponseWriter, r *http.Request) {
	readme, err := fetchRepoReadme("tgrangeo", "meteor")
	log.Println(readme)
	//TODO: error handling here
	if err != nil {
		log.Fatalln(err)
		// return tmpl.HTML(http.StatusInternalServerError, fmt.Sprintf("Error fetching README: %v", err))
	}
	data := map[string]interface{}{
		"Readme": template.HTML(readme),
	}

	tmpl.ExecuteTemplate(w, "index.html", data)
}

func init() {
	if tmpl == nil {
		if tmpl == nil {
			tmpl = template.Must(tmpl.ParseGlob("views/*.html"))
			template.Must(tmpl.ParseGlob("views/*.html"))
		}
	}
}

func main() {
	http.HandleFunc("/", ShowHomePage)
	http.HandleFunc("/projects", ShowProjectPage)


	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("🚀 Starting up on port 3000")

	log.Fatal(http.ListenAndServe(":3000", nil))
}
