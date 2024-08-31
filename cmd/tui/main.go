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
	startTime   = time.Now().Format("2006-01-02T15:04:05.000-0700")
	ch          = make(chan tea.Msg, 2)
	currSending = 0

	jiraClient = jira.Jira{}
	appStyle   = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render

	formStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
			Render
)

type worklogResponse struct {
	err      error
	affected *customItem
}

type customItem struct {
	issue     *jira.Issue
	stopwatch *stopwatch.Model
	logText   string
}

func (i *customItem) Title() string                  { return *i.issue.ID }
func (i *customItem) Description() string            { return i.stopwatch.View() }
func (i *customItem) FilterValue() string            { return *i.issue.ID }
func (i *customItem) GetLogText() *string            { return &i.logText }
func (i *customItem) GetStopwatch() *stopwatch.Model { return i.stopwatch }
func (i *customItem) UpdateStopwatch(msg tea.Msg) tea.Cmd {
	m, cmd := i.stopwatch.Update(msg)
	i.stopwatch = &m
	return cmd
}
func (i *customItem) GetIssue() *jira.Issue {
	return i.issue
}

type model struct {
	list         list.Model
	login        *huh.Form
	log          *huh.Form
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
				Value(&token).
				Validate(jiraClient.SetTokenAuth),
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
			issues, err := jiraClient.GetIssues()
			if err != nil {
				panic(err)
			}
			for _, issue := range issues.Issues {
				s := stopwatch.NewWithInterval(time.Second)
				items = append(items, &customItem{
					issue:     issue,
					stopwatch: &s,
					logText:   "",
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
			if m.list.FilterState() == list.Filtering {
				break
			}
			switch {
			case key.Matches(msg, m.delegateKeys.choose):
				status := ""
				if m.list.SelectedItem().(*customItem).GetStopwatch().Running() {
					status = fmt.Sprintf("Stopped %s Stopwatch", *m.list.SelectedItem().(*customItem).issue.Key)
				} else {
					status = fmt.Sprintf("Started %s Stopwatch", *m.list.SelectedItem().(*customItem).issue.Key)
				}

				cmds = append(cmds, m.list.NewStatusMessage(statusMessageStyle(status)))
				cmds = append(cmds, m.list.SelectedItem().(*customItem).GetStopwatch().Toggle())
				break
			case key.Matches(msg, m.delegateKeys.stopAll):
				cmds = append(cmds, m.list.NewStatusMessage(
					statusMessageStyle("Stopped All Stopwatches"),
				))
				for _, item := range m.list.Items() {
					cmds = append(cmds, item.(*customItem).GetStopwatch().Stop())
				}
				break
			case key.Matches(msg, m.delegateKeys.worklog):
				cmds = append(cmds, m.list.NewStatusMessage(
					statusMessageStyle(
						fmt.Sprintf("Sending %s Worklog", *m.list.SelectedItem().(*customItem).issue.Key),
					),
				))

				cmds = append(cmds, m.list.SelectedItem().(*customItem).GetStopwatch().Stop())
				m.log = huh.NewForm(
					huh.NewGroup(
						huh.NewText().
							Title(fmt.Sprintf(
								"%s @ %s",
								m.list.SelectedItem().(*customItem).GetStopwatch().View(),
								m.list.SelectedItem().(*customItem).Title()),
							).
							Value(m.list.SelectedItem().(*customItem).GetLogText()),
					),
				)
				cmds = append(cmds, m.log.Init())
				m.state.LogWork()
				break
			}
		}
		newListModel, cmd := m.list.Update(msg)
		m.list = newListModel
		cmds = append(cmds, cmd)
		break
	case state.WORKLOG:
		form, cmd := m.log.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.log = f
		}
		cmds = append(cmds, cmd)

		if m.log.State == huh.StateCompleted {
			m.state.Logged()
			cmds = append(cmds, m.list.StartSpinner())
			currSending++
			go func(ch chan tea.Msg) {
				t := int(m.list.SelectedItem().(*customItem).GetStopwatch().Elapsed().Seconds())
				w := jira.Worklog{
					Comment:          m.list.SelectedItem().(*customItem).GetLogText(),
					TimeSpentSeconds: &t,
					Started:          &startTime,
				}

				i := m.list.SelectedItem().(*customItem)
				err := jiraClient.AddWorklogToIssue(&w, i.GetIssue())
				ch <- worklogResponse{
					err:      err,
					affected: i,
				}
			}(ch)
			m.log = nil
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, key.NewBinding(
				key.WithKeys("esc"), //TODO: move to separate object for better visibility
			)):
				m.state.Logged()
				m.log = nil
				break
			}
		}
		break
	}

	select {
	case worklogResp := <-ch:
		if worklogResp, ok := worklogResp.(worklogResponse); ok {
			currSending--
			if worklogResp.err != nil {
				cmds = append(cmds, m.list.NewStatusMessage(
					statusMessageStyle(
						worklogResp.err.Error(),
					),
				))
			} else {
				cmds = append(
					cmds,
					m.list.NewStatusMessage(
						statusMessageStyle(
							fmt.Sprintf("Worklog sent to %s", *worklogResp.affected.GetIssue().Key),
						),
					),
					worklogResp.affected.GetStopwatch().Reset(),
				)
			}
			if currSending == 0 {
				m.list.StopSpinner()
			}
		}
	default:
		break
	}

	for _, item := range m.list.Items() {
		cmds = append(cmds, item.(*customItem).UpdateStopwatch(msg))
	}
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	switch m.state.GetState() {
	case state.LOGIN:
		return formStyle(m.login.View())
	case state.TICKETS:
		return formStyle(m.list.View())
	case state.WORKLOG:
		return formStyle(m.log.View())
	}
	panic("unreachable")
}

func main() {
	server := os.Getenv("JIRA_URL")
	if server == "" {
		panic("JIRA_URL environment variable not set. Example: https://jira.test.server.com/rest/api/2/")
	}
	jiraClient.SetUrlBase(server)
	if _, err := tea.NewProgram(newModel(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func validator(login *string) func(string) error {
	return func(pass string) error {
		return jiraClient.SetBasicAuth(*login, pass)
	}
}
