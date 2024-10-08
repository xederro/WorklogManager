package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/stopwatch"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
	"github.com/xederro/WorklogManager/jira"
	"github.com/xederro/WorklogManager/state"
	"github.com/xederro/WorklogManager/tui/issueList"
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
	affected *issueList.ListItem
}

type Model struct {
	list         list.Model
	log          *huh.Form
	delegateKeys *issueList.DelegateKeyMap
	state        *state.State
}

func NewModel() Model {
	var delegateKeys = issueList.NewDelegateKeyMap()

	// Setup issueList
	delegate := issueList.NewItemDelegate(delegateKeys)
	issues := list.New(nil, delegate, 0, 0)
	issues.Title = "Issues"
	issues.Styles.Title = titleStyle

	// Make issueList of items
	var items []list.Item
	issuesList, err := jiraClient.GetIssues()
	if err != nil {
		panic(err)
	}
	for _, issue := range issuesList.Issues {
		s := stopwatch.NewWithInterval(time.Second)
		items = append(items, &issueList.ListItem{
			Issue:     issue,
			Stopwatch: &s,
			LogText:   "",
		})
	}
	issues.SetItems(items)

	return Model{
		list:         issues,
		delegateKeys: delegateKeys,
		state:        state.NewState(),
	}
}

func (m Model) Init() tea.Cmd {
	tw, th, _ := term.GetSize(os.Stdout.Fd())
	h, v := appStyle.GetFrameSize()
	m.list.SetSize(tw-h, th-v)
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch m.state.GetState() {
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
			case key.Matches(msg, m.delegateKeys.Choose):
				status := ""
				if m.list.SelectedItem().(*issueList.ListItem).GetStopwatch().Running() {
					status = fmt.Sprintf("Stopped %s Stopwatch", *m.list.SelectedItem().(*issueList.ListItem).Issue.Key)
				} else {
					status = fmt.Sprintf("Started %s Stopwatch", *m.list.SelectedItem().(*issueList.ListItem).Issue.Key)
				}

				cmds = append(cmds, m.list.NewStatusMessage(statusMessageStyle(status)))
				cmds = append(cmds, m.list.SelectedItem().(*issueList.ListItem).GetStopwatch().Toggle())
				break
			case key.Matches(msg, m.delegateKeys.StopAll):
				cmds = append(cmds, m.list.NewStatusMessage(
					statusMessageStyle("Stopped All Stopwatches"),
				))
				for _, item := range m.list.Items() {
					cmds = append(cmds, item.(*issueList.ListItem).GetStopwatch().Stop())
				}
				break
			case key.Matches(msg, m.delegateKeys.Worklog):
				cmds = append(cmds, m.list.NewStatusMessage(
					statusMessageStyle(
						fmt.Sprintf("Sending %s Worklog", *m.list.SelectedItem().(*issueList.ListItem).Issue.Key),
					),
				))

				cmds = append(cmds, m.list.SelectedItem().(*issueList.ListItem).GetStopwatch().Stop())
				m.log = huh.NewForm(
					huh.NewGroup(
						huh.NewText().
							Title(fmt.Sprintf(
								"%s @ %s",
								m.list.SelectedItem().(*issueList.ListItem).GetStopwatch().View(),
								m.list.SelectedItem().(*issueList.ListItem).Title()),
							).
							Value(m.list.SelectedItem().(*issueList.ListItem).GetLogText()),
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
				t := int(m.list.SelectedItem().(*issueList.ListItem).GetStopwatch().Elapsed().Seconds())
				w := jira.Worklog{
					Comment:          m.list.SelectedItem().(*issueList.ListItem).GetLogText(),
					TimeSpentSeconds: &t,
					Started:          &startTime,
				}

				i := m.list.SelectedItem().(*issueList.ListItem)
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
					worklogResp.affected.GetStopwatch().Reset(),
					m.list.NewStatusMessage(
						statusMessageStyle(
							fmt.Sprintf("Worklog sent to %s", *worklogResp.affected.GetIssue().Key),
						),
					),
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
		cmds = append(cmds, item.(*issueList.ListItem).UpdateStopwatch(msg))
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	switch m.state.GetState() {
	case state.TICKETS:
		return formStyle(m.list.View())
	case state.WORKLOG:
		return formStyle(m.log.View())
	}
	panic("unreachable")
}
