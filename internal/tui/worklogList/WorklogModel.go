package worklogList

import (
	"github.com/andygrunwald/go-jira"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"maps"
	"slices"
)

type WorklogModel struct {
	list.Model
	items map[string]list.Item
}

func New(issues []jira.Issue, delegate list.ItemDelegate, width, height int) WorklogModel {
	m := WorklogModel{
		Model: list.New(nil, delegate, width, height),
	}
	m.SetWorklogs(issues)

	return m
}

func (m WorklogModel) Update(msg tea.Msg) (WorklogModel, tea.Cmd) {
	var cmd tea.Cmd
	m.Model, cmd = m.Model.Update(msg)

	return m, cmd
}

func (m WorklogModel) SetWorklogs(issues []jira.Issue) (WorklogModel, tea.Cmd) {
	items := make(map[string]list.Item)
	for _, issue := range issues {
		items[issue.Key] = NewItem(issue)
	}

	return m.Update(m.Model.SetItems(slices.Collect(maps.Values(items))))
}

func (m WorklogModel) UpdateWorklogs(issues []jira.Issue) {
	for _, issue := range issues {
		if _, exists := m.items[issue.Key]; exists {
			m.items[issue.Key].(*WorklogItem).Issue = &issue
		} else {
			m.items[issue.Key] = NewItem(issue)
		}
	}
}
