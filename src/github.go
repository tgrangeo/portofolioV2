package src

import (
	"bytes"
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
	Lang     string `json:"language"`
	Color    string `json:"color"`
	Desc     string `json:"description"`
	PushedAt string `json:"pushed_at"`
}

type RepoNameAndDate struct {
	Name     string
	PushedAt time.Time
}

type GitHubUser struct {
	AvatarURL string `json:"avatar_url"`
}

func getRepos(user string) ([]Project, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/repos?per_page=60&page=1", user)
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
				repoNamesAndDates = append(repoNamesAndDates, RepoNameAndDate{
					Name:     repo.Name,
					PushedAt: pushedAt,
				})
			}
		}(repo)
	}
	wg.Wait()

	sort.Slice(repoNamesAndDates, func(i, j int) bool {
		return repoNamesAndDates[i].PushedAt.After(repoNamesAndDates[j].PushedAt)
	})

	var projects []Project
	for _, repoNameAndDate := range repoNamesAndDates {
		for _, repo := range repos {
			if repo.Name == repoNameAndDate.Name {
				langs, err := getTopLanguages(repo.Name)
				if err != nil {
					fmt.Println(err)
				}
				project := Project{
					Title: repo.Name,
					Lang:  langs,
					Desc:  repo.Desc,
				}
				projects = append(projects, project)
				break
			}
		}
	}
	return projects, nil
}

func getTopLanguages(repoName string) ([]string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/languages",username, repoName)
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
		return nil, fmt.Errorf("erreur HTTP: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la lecture de la réponse: %v", err)
	}
	languages := make(map[string]int)
	err = json.Unmarshal(body, &languages)
	if err != nil {
		return nil, fmt.Errorf("erreur lors du parsing JSON: %v", err)
	}
	type languageUsage struct {
		Name  string
		Bytes int
	}

	var sortedLanguages []languageUsage
	for lang, bytes := range languages {
		sortedLanguages = append(sortedLanguages, languageUsage{Name: lang, Bytes: bytes})
	}

	sort.Slice(sortedLanguages, func(i, j int) bool {
		return sortedLanguages[i].Bytes > sortedLanguages[j].Bytes
	})
	topLanguages := []string{}
	for i, lang := range sortedLanguages {
		if i >= 3 {
			break
		}
		topLanguages = append(topLanguages, lang.Name)
	}

	return topLanguages, nil
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

func hasReadme(username, repo string) bool {
	if username != repo{
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
	return false
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

type GraphQLRequest struct {
	Query string `json:"query"`
}

type GraphQLResponse struct {
	Data struct {
		Repository struct {
			OpenGraphImageUrl string `json:"openGraphImageUrl"`
		} `json:"repository"`
	} `json:"data"`
}

func GetRepoImage(name string) string {
	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		log.Fatal("Please set your GitHub token in the GITHUB_TOKEN environment variable")
	}

	query := `
	query {
		repository(owner: "`+ username +`", name: "` + name + `") {
			openGraphImageUrl
		}
	}
`
	requestBody, err := json.Marshal(GraphQLRequest{Query: query})
	if err != nil {
		log.Fatalf("Failed to marshal request: %v", err)
	}

	req, err := http.NewRequest("POST", "https://api.github.com/graphql", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+githubToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("Request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var graphQLResp GraphQLResponse
	if err := json.NewDecoder(resp.Body).Decode(&graphQLResp); err != nil {
		log.Fatalf("Failed to decode response: %v", err)
	}
	return graphQLResp.Data.Repository.OpenGraphImageUrl
}
