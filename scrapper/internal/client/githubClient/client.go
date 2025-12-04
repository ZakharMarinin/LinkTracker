package githubClient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"scrapper/internal/domain"
	"strings"
	"time"
)

type GithubClient struct {
	httpClient *http.Client
	token      string
	baseURL    string
	log        *slog.Logger
}

func NewGithubClient(token string, log *slog.Logger) *GithubClient {
	return &GithubClient{
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
		token:   token,
		baseURL: "https://api.github.com",
		log:     log,
	}
}

func (gh *GithubClient) DoRequest(ctx context.Context, method, path string) ([]byte, error) {
	url := gh.baseURL + path

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		gh.log.Error("Error creating request", "error", err)
		return nil, err
	}

	req.Header.Set("Authorization", "token "+gh.token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := gh.httpClient.Do(req)
	if err != nil {
		gh.log.Error("Error executing request", "error", err, "url", url)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		gh.log.Error("Error executing request", "status", resp.Status)
		return nil, fmt.Errorf(resp.Status)
	}

	return io.ReadAll(resp.Body)
}

func (gh *GithubClient) ParseGitHubURL(url string) (owner, repo string, err error) {
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "github.com/")

	parts := strings.Split(url, "/")

	if len(url) < 2 {
		return "", "", nil
	}

	return parts[0], parts[1], nil
}

func (gh *GithubClient) GetUpdates(ctx context.Context, link *domain.Link) (update *domain.GitHubContent, err error) {
	owner, repo, err := gh.ParseGitHubURL(link.URL)
	if err != nil {
		gh.log.Error("Error parsing GitHub URL", "error", err)
		return nil, err
	}

	repoPR := fmt.Sprintf("/repos/%s/%s/pulls?state=open", owner, repo)
	repoIssues := fmt.Sprintf("/repos/%s/%s/issues?state=open", owner, repo)

	fmt.Println(repoIssues)

	prData, err := gh.DoRequest(ctx, "GET", repoPR)
	if err != nil {
		return nil, err
	}

	issuesData, err := gh.DoRequest(ctx, "GET", repoIssues)
	if err != nil {
		return nil, err
	}

	var allIssues domain.GitHubContent
	if err := json.Unmarshal(prData, &allIssues.PullRequests); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(issuesData, &allIssues.Issues); err != nil {
		return nil, err
	}

	var updates domain.GitHubContent

	for _, issue := range allIssues.Issues {
		if issue.UpdatedAt.After(link.LastUpdated) {
			updates.Issues = append(updates.Issues, issue)
		}
	}

	for _, pr := range allIssues.PullRequests {
		if pr.UpdatedAt.After(link.LastUpdated) {
			updates.PullRequests = append(updates.PullRequests, pr)
		}
	}

	gh.log.Info("got update", "update", updates)

	return &updates, nil
}
