package issueList

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
)

type ListItem struct {
	Issue     *jira.Issue
	Stopwatch *stopwatch.Model
	LogText   string
}

func (i *ListItem) Title() string {
	if i.Issue.Fields.Summary != "" {
		return fmt.Sprintf("%s - %s", i.Issue.Key, i.Issue.Fields.Summary)
	}
	return fmt.Sprintf("%s", i.Issue.Key)
}
func (i *ListItem) Description() string            { return i.Stopwatch.View() }
func (i *ListItem) FilterValue() string            { return i.Title() }
func (i *ListItem) GetLogText() *string            { return &i.LogText }
func (i *ListItem) GetStopwatch() *stopwatch.Model { return i.Stopwatch }
func (i *ListItem) UpdateStopwatch(msg tea.Msg) tea.Cmd {
	m, cmd := i.Stopwatch.Update(msg)
	i.Stopwatch = &m
	return cmd
}
func (i *ListItem) GetIssue() *jira.Issue {
	return i.Issue
}
