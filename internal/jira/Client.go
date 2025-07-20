package jira

import (
	"github.com/andygrunwald/go-jira"
	tea "github.com/charmbracelet/bubbletea"
	"net/http"
)

// Client wraps the go-jira client to provide additional functionality.
type Client struct {
	*jira.Client
}

// NewClient creates a new Jira client with the provided HTTP client and base URL.
func NewClient(httpClient *http.Client, baseURL string) (*Client, error) {
	client, err := jira.NewClient(httpClient, baseURL)
	if err != nil {
		return nil, err
	}
	return &Client{client}, nil
}

// WorklogItemsMsg is a message type used to communicate the results of fetching worklog items.
// Compatible with the Bubbletea tea.Cmd
type WorklogItemsMsg struct {
	Issues []jira.Issue
	Err    error
}

// GetItemsToUpdate fetches Jira issues based on the provided JQL query.
func (m Client) GetItemsToUpdate(jql string) tea.Cmd {
	return func() tea.Msg {
		// Fetch issues from Jira
		issues, _, err := m.Issue.Search(jql, nil)
		return WorklogItemsMsg{
			Issues: issues,
			Err:    err,
		}
	}
}
