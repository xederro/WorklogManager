package worklogList

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/xederro/WorklogManager/internal/config"
	"github.com/xederro/WorklogManager/internal/tui/worklogTimer"
)

type WorklogItem struct {
	Issue     *jira.Issue
	Stopwatch worklogTimer.WorklogTimer
	LogText   string
}

func NewItem(issue jira.Issue) list.Item {
	return &WorklogItem{
		Issue:     &issue,
		Stopwatch: worklogTimer.New(),
		LogText:   config.Conf.Jira.DefaultWorklogComment,
	}
}

func (i *WorklogItem) Title() string {
	if i.Issue.Fields.Summary != "" {
		return fmt.Sprintf("%s - %s", i.Issue.Key, i.Issue.Fields.Summary)
	}
	return fmt.Sprintf("%s", i.Issue.Key)
}
func (i *WorklogItem) Description() string                     { return i.Stopwatch.View() }
func (i *WorklogItem) FilterValue() string                     { return i.Title() }
func (i *WorklogItem) GetLogText() *string                     { return &i.LogText }
func (i *WorklogItem) GetStopwatch() worklogTimer.WorklogTimer { return i.Stopwatch }
func (i *WorklogItem) UpdateStopwatch(msg tea.Msg) tea.Cmd {
	m, cmd := i.Stopwatch.Update(msg)
	i.Stopwatch = m
	return cmd
}
func (i *WorklogItem) GetIssue() *jira.Issue {
	return i.Issue
}

// WorklogResponse is the message type returned by the worklog actions.
type WorklogResponse struct {
	Err      error
	Affected *WorklogItem
}

// ReturnCmd represents a worklog item.
func ReturnCmd(err error, affected *WorklogItem) tea.Cmd {
	return func() tea.Msg {
		return WorklogResponse{
			Err:      err,
			Affected: affected,
		}
	}
}
