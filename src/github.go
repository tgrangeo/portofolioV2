package src

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type GitHubFileContent struct {
	Content string `json:"content"`
}

type GitHubRepo struct {
	Name     string `json:"name"`
	PushedAt string `json:"pushed_at"`
}

// RepoNameAndDate holds the repo name and parsed date for sorting
type RepoNameAndDate struct {
	Name     string
	PushedAt time.Time
}

type GitHubUser struct {
	AvatarURL string `json:"avatar_url"`
}

func GetProfilePicture(username string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s", username)
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
		return "", fmt.Errorf("error fetching user data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var user GitHubUser
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}
	return user.AvatarURL, nil
}

func getRepos(user string) ([]string, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/repos", user)
	log.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		req.Header.Set("Authorization", "token "+token)
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

	var repoNamesAndDates []RepoNameAndDate
	var repoNamesAndDatesMutex sync.Mutex
	var wg sync.WaitGroup

	for _, repo := range repos {
		wg.Add(1)
		go func(repo GitHubRepo) {
			defer wg.Done()
			if hasReadme(user, repo.Name) {
				pushedAt, err := time.Parse(time.RFC3339, repo.PushedAt)
				if err != nil {
					log.Printf("Error parsing time for repo %s: %v", repo.Name, err)
					return
				}
				repoNamesAndDatesMutex.Lock()
				repoNamesAndDates = append(repoNamesAndDates, RepoNameAndDate{
					Name:     repo.Name,
					PushedAt: pushedAt,
				})
				repoNamesAndDatesMutex.Unlock()
			}
		}(repo)
	}

	wg.Wait()

	sort.Slice(repoNamesAndDates, func(i, j int) bool {
		return repoNamesAndDates[i].PushedAt.After(repoNamesAndDates[j].PushedAt)
	})

	repoNames := make([]string, len(repoNamesAndDates))
	for i, repo := range repoNamesAndDates {
		repoNames[i] = repo.Name
	}

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
