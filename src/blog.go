package src

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type Article struct {
	Title           string
	PublicationDate time.Time
	FilePath        string
}

type Blog_Data struct {
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

func initArticles() []Article {
	res := []Article{}
	res = append(res, NewArticle("gitignore", "../static/blog/gitignore.md", time.Now()))
	res = append(res, NewArticle("htmx_go", "../static/blog/htmx_go.md", time.Now()))
	return res
}

func ShowBlog(w http.ResponseWriter, r *http.Request) string {
	if r.Header.Get("HX-Request") != "true" {
		err := tmpl.ExecuteTemplate(w, "index.html", nil)
		if err != nil {
			http.Error(w, "Error rendering index template", http.StatusInternalServerError)
			return ""
		}
		return ""
	}
	data := Blog_Data{
		Articles: initArticles(),
	}
	tmpl, err := template.ParseFiles("views/blog.html")
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


func ShowArticle(w http.ResponseWriter, r *http.Request) string {
	path := strings.TrimPrefix(r.URL.Path, "/article/")
	title := strings.TrimSpace(path)
	title = filepath.Clean(title)
	filePath := filepath.Join("./static/articles/", title+".md")
	art, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, "Article not found", http.StatusNotFound)
		fmt.Println("Error opening file:", err)
		return ""
	}
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	mdParser := parser.NewWithExtensions(extensions)
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	renderer := html.NewRenderer(html.RendererOptions{Flags: htmlFlags})
	md := markdown.ToHTML(art, mdParser, renderer)
	return string(`<div class="content-readme">` + string(md) + `</div>`)
}
