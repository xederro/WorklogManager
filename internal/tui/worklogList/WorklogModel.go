package worklogList

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/xederro/WorklogManager/internal/config"
	"github.com/xederro/WorklogManager/internal/gen/sqlc"
	"maps"
	"slices"
	"strings"
	"time"
)

// WorklogModel is a model for managing a list of worklogs.
type WorklogModel struct {
	list.Model
	items map[string]list.Item
}

func New(issues []jira.Issue, delegate list.ItemDelegate, width, height int) WorklogModel {
	m := WorklogModel{
		Model: list.New(nil, delegate, width, height),
		items: make(map[string]list.Item),
	}
	m.UpdateWorklogs(issues)

	return m
}

func (m WorklogModel) Update(msg tea.Msg) (WorklogModel, tea.Cmd) {
	var cmd tea.Cmd
	m.Model, cmd = m.Model.Update(msg)

	return m, cmd
}

func (m WorklogModel) UpdateWorklogs(issues []jira.Issue) (WorklogModel, tea.Cmd) {
	for _, issue := range issues {
		if _, exists := m.items[issue.Key]; exists {
			item := m.items[issue.Key].(*WorklogItem)
			item.Issue = &issue
			_ = config.Queries.UpdateWorklog(context.Background(), sqlc.UpdateWorklogParams{
				JiraData: getJiraIssueJson(issue),
				Duration: item.Stopwatch.Elapsed().Nanoseconds(),
				Running:  item.Stopwatch.Running(),
				LogText:  item.LogText,
				JiraKey:  issue.Key,
			})
		} else {
			m.items[issue.Key] = NewItem(issue)
			_ = config.Queries.CreateWorklog(context.Background(), sqlc.CreateWorklogParams{
				JiraKey:  issue.Key,
				JiraData: getJiraIssueJson(issue),
				Duration: 0,
				Running:  false,
				LogText:  "",
			})
		}
	}
	items := slices.Collect(maps.Values(m.items))
	slices.SortFunc(items, orderItems)
	return m.Update(m.Model.SetItems(items))
}

func (m WorklogModel) GetItem(key string) (list.Item, error) {
	k, exist := m.items[key]
	if !exist {
		return nil, fmt.Errorf("item with key %s not found", key)
	}
	return k, nil
}

func (m WorklogModel) Quit() {
	for _, item := range m.items {
		_ = config.Queries.UpdateWorklog(context.Background(), sqlc.UpdateWorklogParams{
			JiraData: getJiraIssueJson(*item.(*WorklogItem).Issue),
			Duration: item.(*WorklogItem).GetStopwatch().Elapsed().Nanoseconds(),
			Running:  item.(*WorklogItem).GetStopwatch().Running(),
			LogText:  item.(*WorklogItem).LogText,
			JiraKey:  item.(*WorklogItem).Issue.Key,
		})
	}
}

func orderItems(a, b list.Item) int {
	ai := a.(*WorklogItem)
	bi := b.(*WorklogItem)
	return strings.Compare(ai.Issue.Key, bi.Issue.Key)
}

func getJiraIssueJson(issue jira.Issue) string {
	bytes, err := json.Marshal(issue)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func getJiraIssueFromJson(data string) (jira.Issue, error) {
	var issue jira.Issue
	if err := json.Unmarshal([]byte(data), &issue); err != nil {
		return jira.Issue{}, fmt.Errorf("failed to unmarshal jira issue: %w", err)
	}
	return issue, nil
}

func (m WorklogModel) GetJiraFromDB() (WorklogModel, tea.Cmd) {
	var cmds []tea.Cmd
	now := time.Now()
	worklogs, err := config.Queries.ListWorklog(context.Background())
	if err != nil {
		return m, nil
	}

	for _, worklog := range worklogs {
		issue, err := getJiraIssueFromJson(worklog.JiraData)
		if err != nil {
			continue
		}
		item := NewItem(issue)
		item.(*WorklogItem).LogText = worklog.LogText
		if worklog.Running {
			cmds = append(cmds, item.(*WorklogItem).Stopwatch.Start())
			closed := now.Sub(worklog.UpdatedAt)
			item.(*WorklogItem).Stopwatch = item.(*WorklogItem).Stopwatch.SetDuration(
				((time.Duration(worklog.Duration) + closed) / 1e9) * 1e9,
			)
		} else {
			item.(*WorklogItem).Stopwatch = item.(*WorklogItem).Stopwatch.SetDuration(time.Duration(worklog.Duration))
		}
		m.items[issue.Key] = item
	}

	items := slices.Collect(maps.Values(m.items))
	slices.SortFunc(items, orderItems)
	var cmd tea.Cmd
	cmd = m.Model.SetItems(items)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}
