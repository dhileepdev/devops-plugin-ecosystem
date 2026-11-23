package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Issue represents a GitHub issue used as a health signal.
type Issue struct {
	Number    int       `json:"number"`
	Title     string    `json:"title"`
	State     string    `json:"state"`
	Labels    []Label   `json:"labels"`
	HTMLURL   string    `json:"html_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
}

// Label is a GitHub issue label.
type Label struct {
	Name string `json:"name"`
}

// Client wraps the GitHub REST API.
type Client struct {
	BaseURL string
	Token   string
	HTTP    *http.Client
}

// NewClient creates a GitHub API client. Token may be empty for public repos.
func NewClient(token string) *Client {
	return &Client{
		BaseURL: "https://api.github.com",
		Token:   token,
		HTTP:    &http.Client{Timeout: 15 * time.Second},
	}
}

// ListIssues fetches issues for the given owner/repo. state is "open", "closed", or "all".
func (c *Client) ListIssues(owner, repo, state string) ([]Issue, error) {
	if state == "" {
		state = "open"
	}
	url := fmt.Sprintf("%s/repos/%s/%s/issues?state=%s&per_page=100", c.BaseURL, owner, repo, state)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("github api %d: %s", resp.StatusCode, string(body))
	}
	var issues []Issue
	if err := json.NewDecoder(resp.Body).Decode(&issues); err != nil {
		return nil, err
	}
	return issues, nil
}

// TerraformSDKRepo is the canonical repo for Terraform provider SDK issues.
const TerraformSDKRepo = "hashicorp/terraform-plugin-sdk"

// FetchTerraformSDKIssues fetches open issues from hashicorp/terraform-plugin-sdk.
func (c *Client) FetchTerraformSDKIssues() ([]Issue, error) {
	return c.ListIssues("hashicorp", "terraform-plugin-sdk", "open")
}
