package worklogList

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/xederro/WorklogManager/internal/config"
	"time"
)

type WorklogItem struct {
	Issue     *jira.Issue
	Stopwatch *stopwatch.Model
	LogText   string
	UseAi     int
}

func NewItem(issue jira.Issue) list.Item {
	s := stopwatch.NewWithInterval(time.Second)
	return &WorklogItem{
		Issue:     &issue,
		Stopwatch: &s,
		LogText:   config.Conf.Jira.DefaultWorklogComment,
	}
}

func (i *WorklogItem) Title() string {
	if i.Issue.Fields.Summary != "" {
		return fmt.Sprintf("%s - %s", i.Issue.Key, i.Issue.Fields.Summary)
	}
	return fmt.Sprintf("%s", i.Issue.Key)
}
func (i *WorklogItem) Description() string            { return i.Stopwatch.View() }
func (i *WorklogItem) FilterValue() string            { return i.Title() }
func (i *WorklogItem) GetLogText() *string            { return &i.LogText }
func (i *WorklogItem) GetStopwatch() *stopwatch.Model { return i.Stopwatch }
func (i *WorklogItem) UpdateStopwatch(msg tea.Msg) tea.Cmd {
	m, cmd := i.Stopwatch.Update(msg)
	i.Stopwatch = &m
	return cmd
}
func (i *WorklogItem) GetIssue() *jira.Issue {
	return i.Issue
}
