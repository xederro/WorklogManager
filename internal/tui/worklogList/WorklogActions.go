package worklogList

import (
	"github.com/andygrunwald/go-jira"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/xederro/WorklogManager/internal/config"
)

type WorklogResponse struct {
	Err      error
	Affected *WorklogItem
}

func ReturnCmd(err error, affected *WorklogItem) tea.Cmd {
	return func() tea.Msg {
		return WorklogResponse{
			Err:      err,
			Affected: affected,
		}
	}
}

type WorklogItemsMsg struct {
	Issues []jira.Issue
	Err    error
}

func GetItemsToUpdate() tea.Cmd {
	return func() tea.Msg {
		// Fetch issues from Jira
		issues, _, err := config.JiraClient.Issue.Search("assignee = currentUser() AND statusCategory != Done ORDER BY updated DESC", nil)
		return WorklogItemsMsg{
			Issues: issues,
			Err:    err,
		}
	}
}
