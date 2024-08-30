package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
	"github.com/xederro/WorklogManager/jira"
	"github.com/xederro/WorklogManager/state"
	"os"
	"time"
)

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)

type customItem struct {
	issue     *jira.Issue
	stopwatch *stopwatch.Model
}

func (i *customItem) Title() string {
	return fmt.Sprintf("%s ~ %s", *i.issue.ID, i.stopwatch.View())
}
func (i *customItem) Description() string            { return *i.issue.ID }
func (i *customItem) FilterValue() string            { return *i.issue.ID }
func (i *customItem) GetStopwatch() *stopwatch.Model { return i.stopwatch }
func (i *customItem) UpdateStopwatch(msg tea.Msg) tea.Cmd {
	m, cmd := i.stopwatch.Update(msg)
	i.stopwatch = &m
	return cmd
}

type model struct {
	list         list.Model
	login        *huh.Form
	delegateKeys *delegateKeyMap
	state        state.State
}

func newModel() model {
	var delegateKeys = newDelegateKeyMap()
	var isToken bool
	var token string
	var login string
	var pass string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Do you have token?").
				Value(&isToken),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Login").
				Prompt("> ").
				Value(&login),
			huh.NewInput().
				Title("Password").
				Prompt("> ").
				EchoMode(huh.EchoModePassword).
				Value(&pass).
				Validate(validator(&login)),
		).WithHideFunc(func() bool {
			return isToken
		}),
		huh.NewGroup(
			huh.NewInput().
				Title("Token").
				Prompt("> ").
				Value(&token),
		).WithHideFunc(func() bool {
			return !isToken
		}),
	)

	// Setup list
	delegate := newItemDelegate(delegateKeys)
	groceryList := list.New(nil, delegate, 0, 0)
	groceryList.Title = "Issues"
	groceryList.Styles.Title = titleStyle

	return model{
		list:         groceryList,
		delegateKeys: delegateKeys,
		login:        form,
	}
}

func (m model) Init() tea.Cmd {
	return m.login.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch m.state.GetState() {
	case state.LOGIN:
		form, cmd := m.login.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.login = f
		}
		cmds = append(cmds, cmd)

		if m.login.State == huh.StateCompleted {
			// Make list of items
			var items []list.Item
			issues, err := jira.GetIssues()
			if err != nil {

			}
			for _, issue := range issues.Issues {
				s := stopwatch.NewWithInterval(time.Second)
				items = append(items, &customItem{
					issue:     issue,
					stopwatch: &s,
				})
			}
			m.list.SetItems(items)
			m.state.Login()
			tw, th, _ := term.GetSize(os.Stdout.Fd())
			h, v := appStyle.GetFrameSize()
			m.list.SetSize(tw-h, th-v)
		}
		break
	case state.TICKETS:
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			h, v := appStyle.GetFrameSize()
			m.list.SetSize(msg.Width-h, msg.Height-v)

		case tea.KeyMsg:
			// Don't match any of the keys below if we're actively filtering.
			if m.list.FilterState() == list.Filtering {
				break
			}
			switch {
			case key.Matches(msg, m.delegateKeys.choose):
				cmds = append(cmds, m.list.SelectedItem().(*customItem).GetStopwatch().Toggle())
			}
		}

		// This will also call our delegate's update function.
		for _, item := range m.list.Items() {
			cmds = append(cmds, item.(*customItem).UpdateStopwatch(msg))
		}
		newListModel, cmd := m.list.Update(msg)
		m.list = newListModel
		cmds = append(cmds, cmd)
		break
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	switch m.state.GetState() {
	case state.LOGIN:
		return appStyle.Render(m.login.View())
	case state.TICKETS:
		return appStyle.Render(m.list.View())
	}
	panic("unreachable")
}

func main() {
	if _, err := tea.NewProgram(newModel(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func validator(login *string) func(string) error {
	return func(pass string) error {
		return jira.Jira{}.SetBasicAuth(*login, pass)
	}
}
