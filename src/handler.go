package src

import (
	"encoding/json"
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
	_, err := GetProfilePicture(strings.Split(r.URL.Path, "/")[2])
	if err != nil {
		w.WriteHeader(404)
	} else {
		username = strings.Split(r.URL.Path, "/")[2]
	}
}

func GetImage(w http.ResponseWriter, r *http.Request) {
	url := GetRepoImage(strings.Split(r.URL.Path, "/")[2])
	response := map[string]string{"imageUrl": url}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
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

func ShowHome(w http.ResponseWriter, r *http.Request) string {
	htmlFile, err := os.ReadFile("views/home.html")
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
		default_user_readme, _ := os.ReadFile("views/default_user_readme.html")
		return string(default_user_readme)
	}
	return string("<div class=\"content-readme\"" + readme + "</div>")
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
