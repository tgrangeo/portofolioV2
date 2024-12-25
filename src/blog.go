package src

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	// "os"
	"text/template"
	"time"
)

type Article struct {
	Title           string
	PublicationDate time.Time
	FilePath        string
}

type Data struct {
	Articles []Article
}

func NewArticle(title, filePath string, publicationDate time.Time) Article {
	return Article{
		Title:           title,
		PublicationDate: publicationDate,
		FilePath:        filePath,
	}
}

func (a Article) String() string {
	return "Title: " + a.Title + "\nDate: " + a.PublicationDate.Format("2006-01-02") + "\nPath: " + a.FilePath
}

func initArticles() []Article{
	res := []Article{}
	res = append(res, NewArticle("gitignore", "../static/blog/gitignore.md", time.Now()))
	res = append(res, NewArticle("htmx_go", "../static/blog/htmx_go.md", time.Now()))
	return res
}

func ShowBlog(w http.ResponseWriter, r *http.Request) {
	data := Data{
		Articles: initArticles(),
	}
	tmpl, err := template.ParseFiles("views/blog.html")
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

func ShowArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("article")
	
	// Extract and sanitize the article title
	path := strings.TrimPrefix(r.URL.Path, "/article/")
	title := strings.TrimSpace(path)
	title = filepath.Clean(title) // Prevent directory traversal

	// Build the file path
	filePath := filepath.Join("./static/articles/", title + ".md")

	// Open the file
	art, err := os.ReadFile(filePath) // Read entire file content
	if err != nil {
		http.Error(w, "Article not found", http.StatusNotFound)
		fmt.Println("Error opening file:", err)
		return
	}

	// Write the response
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<div class="content-readme">`))
	w.Write([]byte(art))
	w.Write([]byte(`</div>`))
}
