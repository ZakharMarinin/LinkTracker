package domain

import "time"

type Link struct {
	ID          int       `json:"link_id"`
	URL         string    `json:"url"`
	Alias       string    `json:"alias"`
	Desc        string    `json:"desc"`
	ChatID      int64     `json:"chat_id"`
	LastUpdated time.Time `json:"last_updated"`
}

type Response struct {
	ChatIDs []int64 `json:"chat_ids"`
	Link    *Link   `json:"link"`
}

type AddLinkRequest struct {
	ChatID int64  `json:"chat_id"`
	URL    string `json:"url"`
	Desc   string `json:"desc"`
}

type DeleteLinkRequest struct {
	ChatID int64  `json:"chat_id"`
	Alias  string `json:"alias"`
}

type GitHubContent struct {
	Issues       []GitHubIssue `json:"issues"`
	PullRequests []GitHubIssue `json:"pull_requests"`
}

type GitHubIssue struct {
	Title       string    `json:"title"`
	UpdatedAt   time.Time `json:"updated_at"`
	User        User      `json:"user"`
	Body        string    `json:"body"`
	PullRequest *struct{} `json:"pull_request,omitempty"`
}

type User struct {
	Login string `json:"login"`
}
