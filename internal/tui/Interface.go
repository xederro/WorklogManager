package tui

import (
	"context"
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
	"github.com/xederro/WorklogManager/internal/config"
	"github.com/xederro/WorklogManager/internal/state"
	"github.com/xederro/WorklogManager/internal/tui/worklogList"
	"google.golang.org/genai"
	"os"
	"time"
)

var (
	startTime   = time.Now()
	Ch          = make(chan tea.Cmd, 2)
	currSending = 0

	appStyle = lipgloss.NewStyle().Padding(1, 2)

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

type Model struct {
	list         worklogList.WorklogModel
	log          *huh.Form
	delegateKeys *worklogList.DelegateKeyMap
	state        *state.State
}

func NewModel() Model {
	var delegateKeys = worklogList.NewDelegateKeyMap()

	// Setup worklogList
	delegate := worklogList.NewWorklogDelegate(delegateKeys)
	issues := worklogList.New(nil, delegate, 0, 0)
	issues.Title = "Issues"
	issues.Styles.Title = titleStyle

	// Make a worklogList of items
	go func() {
		Ch <- worklogList.GetItemsToUpdate()
	}()

	return Model{
		list:         issues,
		delegateKeys: delegateKeys,
		state:        state.New(),
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
				if m.list.SelectedItem().(*worklogList.WorklogItem).GetStopwatch().Running() {
					status = fmt.Sprintf("Stopped %s Stopwatch", m.list.SelectedItem().(*worklogList.WorklogItem).Issue.Key)
				} else {
					status = fmt.Sprintf("Started %s Stopwatch", m.list.SelectedItem().(*worklogList.WorklogItem).Issue.Key)
				}

				cmds = append(cmds, m.list.NewStatusMessage(statusMessageStyle(status)))
				cmds = append(cmds, m.list.SelectedItem().(*worklogList.WorklogItem).GetStopwatch().Toggle())
				break
			case key.Matches(msg, m.delegateKeys.StopAll):
				cmds = append(cmds, m.list.NewStatusMessage(
					statusMessageStyle("Stopped All Stopwatches"),
				))
				for _, item := range m.list.Items() {
					cmds = append(cmds, item.(*worklogList.WorklogItem).GetStopwatch().Stop())
				}
				break
			case key.Matches(msg, m.delegateKeys.Worklog):
				cmds = append(cmds, m.list.NewStatusMessage(
					statusMessageStyle(
						fmt.Sprintf("Sending %s Worklog", m.list.SelectedItem().(*worklogList.WorklogItem).Issue.Key),
					),
				))

				cmds = append(cmds, m.list.SelectedItem().(*worklogList.WorklogItem).GetStopwatch().Stop())

				if config.Conf.UseAi {
					m.log = huh.NewForm(
						huh.NewGroup(
							huh.NewText().
								Title(fmt.Sprintf(
									"%s @ %s",
									m.list.SelectedItem().(*worklogList.WorklogItem).GetStopwatch().View(),
									m.list.SelectedItem().(*worklogList.WorklogItem).Title()),
								).
								Value(m.list.SelectedItem().(*worklogList.WorklogItem).GetLogText()),
							huh.NewSelect[int]().
								Options(
									huh.NewOption("NAAAH!", 0).Selected(true),
									huh.NewOption("YIS!", 1),
								).Title("Use AI?").
								Value(&m.list.SelectedItem().(*worklogList.WorklogItem).UseAi),
						),
					)
				} else {
					m.log = huh.NewForm(
						huh.NewGroup(
							huh.NewText().
								Title(fmt.Sprintf(
									"%s @ %s",
									m.list.SelectedItem().(*worklogList.WorklogItem).GetStopwatch().View(),
									m.list.SelectedItem().(*worklogList.WorklogItem).Title()),
								).
								Value(m.list.SelectedItem().(*worklogList.WorklogItem).GetLogText()),
						),
					)
				}

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
			go func(ch chan tea.Cmd) {
				t := int(m.list.SelectedItem().(*worklogList.WorklogItem).GetStopwatch().Elapsed().Seconds())

				comment := *m.list.SelectedItem().(*worklogList.WorklogItem).GetLogText()
				if m.list.SelectedItem().(*worklogList.WorklogItem).UseAi == 1 {
					result, err := config.GoogleClient.Models.GenerateContent(
						context.Background(),
						config.Conf.GoogleAi.DefaultModel,
						genai.Text(config.Conf.GoogleAi.DefaultPrompt+" "+comment),
						&genai.GenerateContentConfig{
							ThinkingConfig: &genai.ThinkingConfig{
								ThinkingBudget:  nil,
								IncludeThoughts: false,
							},
						},
					)
					if err != nil {
						ch <- worklogList.ReturnCmd(err, nil)
						return
					}
					comment = result.Text()
				}

				w := jira.WorklogRecord{
					Comment:          comment,
					TimeSpentSeconds: t,
					Started:          (*jira.Time)(&startTime),
				}

				i := m.list.SelectedItem().(*worklogList.WorklogItem)
				_, _, err := config.JiraClient.Issue.AddWorklogRecord(i.Issue.ID, &w)
				ch <- worklogList.ReturnCmd(err, i)
			}(Ch)
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

	switch msg.(type) {
	case worklogList.WorklogItemsMsg:
		msg := msg.(worklogList.WorklogItemsMsg)
		if msg.Err != nil {
			cmds = append(cmds, m.list.NewStatusMessage(
				statusMessageStyle(
					fmt.Sprintf("Error fetching issues: %s", msg.Err.Error()),
				),
			))
		} else {
			var cmd tea.Cmd
			m.list, cmd = m.list.UpdateWorklogs(msg.Issues)
			cmds = append(
				cmds,
				cmd,
				m.list.NewStatusMessage(
					statusMessageStyle(
						fmt.Sprintf("Fetched %d issues", len(msg.Issues)),
					),
				),
			)
		}
	case worklogList.WorklogResponse:
		msg := msg.(worklogList.WorklogResponse)
		currSending--
		if msg.Err != nil {
			cmds = append(cmds, m.list.NewStatusMessage(
				statusMessageStyle(
					msg.Err.Error(),
				),
			))
		} else {
			cmds = append(
				cmds,
				msg.Affected.GetStopwatch().Reset(),
				m.list.NewStatusMessage(
					statusMessageStyle(
						fmt.Sprintf("Worklog sent to %s", msg.Affected.GetIssue().Key),
					),
				),
			)
		}
		if currSending == 0 {
			m.list.StopSpinner()
		}
	}

	select {
	case worklogResp := <-Ch:
		cmds = append(cmds, worklogResp)
	default:
		break
	}

	for _, item := range m.list.Items() {
		cmds = append(cmds, item.(*worklogList.WorklogItem).UpdateStopwatch(msg))
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	switch m.state.GetState() {
	case state.TICKETS:
		return formStyle(m.list.View())
	case state.WORKLOG:
		return formStyle(m.log.View())
	default:
		panic("unreachable")
	}
}
