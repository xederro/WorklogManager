package worklogList

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"maps"
	"slices"
	"strings"
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
			m.items[issue.Key].(*WorklogItem).Issue = &issue
		} else {
			m.items[issue.Key] = NewItem(issue)
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

func orderItems(a, b list.Item) int {
	ai := a.(*WorklogItem)
	bi := b.(*WorklogItem)
	return strings.Compare(ai.Issue.Key, bi.Issue.Key)
}
