package jira

import (
	"github.com/andygrunwald/go-jira"
	tea "github.com/charmbracelet/bubbletea"
	"net/http"
)

type Client struct {
	*jira.Client
}

func NewClient(httpClient *http.Client, baseURL string) (*Client, error) {
	client, err := jira.NewClient(httpClient, baseURL)
	if err != nil {
		return nil, err
	}
	return &Client{client}, nil
}

type WorklogItemsMsg struct {
	Issues []jira.Issue
	Err    error
}

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
