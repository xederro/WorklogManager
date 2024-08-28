package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/xederro/WorklogManager/jira"
	"github.com/xederro/WorklogManager/state"
)

type model struct {
	stopwatch stopwatch.Model
	form      *huh.Form
	writeLog  *huh.Form
	state     *state.State
	selected  *jira.Issue
	text      *string
	keymap    keymap
	help      help.Model
	time      time.Time
}

type keymap struct {
	start key.Binding
	stop  key.Binding
	reset key.Binding
	quit  key.Binding
}

func (m model) helpView() string {
	return "\n" + m.help.ShortHelpView([]key.Binding{
		m.keymap.start,
		m.keymap.stop,
		m.keymap.reset,
		m.keymap.quit,
	})
}

func initialModel() model {
	issues := jira.GetIssues()
	var s *jira.Issue
	var choises []huh.Option[jira.Issue]
	if issues == nil {
		log.Println("cannot optain issues")
	} else {
		for _, v := range issues.Issues {
			if s == nil {
				s = v
			}
			choises = append(choises, huh.NewOption(
				fmt.Sprintf("%s | %s | %s", *v.ID, *v.Fields.Status.Name, *v.Fields.Summary), *v,
			))
		}
	}
	t := ""
	m := model{
		state:     state.NewState(),
		selected:  s,
		text:      &t,
		stopwatch: stopwatch.NewWithInterval(time.Minute),
		keymap: keymap{
			start: key.NewBinding(
				key.WithKeys("s"),
				key.WithHelp("s", "start"),
			),
			stop: key.NewBinding(
				key.WithKeys("s"),
				key.WithHelp("s", "stop"),
			),
			reset: key.NewBinding(
				key.WithKeys("r"),
				key.WithHelp("r", "reset"),
			),
			quit: key.NewBinding(
				key.WithKeys("ctrl+c", "q"),
				key.WithHelp("q", "quit"),
			),
		},
		help: help.New(),
	}

	m.keymap.start.SetEnabled(false)

	m.form = huh.NewForm(huh.NewGroup(
		huh.NewSelect[jira.Issue]().
			Options(choises...).
			Value(m.selected),
	))
	m.writeLog = huh.NewForm(huh.NewGroup(
		huh.NewText().
			Title("Worklog").
			Value(m.text),
	))
	return m
}

func (m model) Init() tea.Cmd {
	return m.form.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state.State {
	case state.CHOOSING_TICKET:
		form, cmd := m.form.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
		}
		return m, cmd
	case state.MEASURING_TIME:
		if m.state.IsNew() {
			m.time = time.Now()
			return m, m.stopwatch.Start()
		}
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, m.keymap.quit):
				m.state.Log()
				return m, m.stopwatch.Stop()
			case key.Matches(msg, m.keymap.reset):
				return m, m.stopwatch.Reset()
			case key.Matches(msg, m.keymap.start, m.keymap.stop):
				m.keymap.stop.SetEnabled(!m.stopwatch.Running())
				m.keymap.start.SetEnabled(m.stopwatch.Running())
				return m, m.stopwatch.Toggle()
			}
		}
		w, cmd := m.stopwatch.Update(msg)
		m.stopwatch = w
		return m, cmd
	case state.WRITING_WORKLOG:
		if m.state.IsNew() {
			return m, m.writeLog.Init()
		}
		form, cmd := m.writeLog.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.form = f
		}
		return m, cmd
	case state.SENDING:
		if m.state.IsNew() {
			t := int(m.stopwatch.Elapsed().Seconds())
			s := m.time.
				Format("2006-01-02T15:04:05.000-0700")
			issue := jira.GetIssues().Issues

			fmt.Println(*issue[0].Self)

			(&jira.Worklog{
				Comment:          m.text,
				TimeSpentSeconds: &t,
				Started:          &s,
			}).AddToIssue(m.selected)
			return m, tea.Quit
		}
		panic("wrong update")
	}
	panic("wrong update")
}

func (m model) View() string {
	switch m.state.State {
	case state.CHOOSING_TICKET:
		if m.form.State != huh.StateCompleted {
			return m.form.View()
		}
		m.state.Choose()
		return m.View()
	case state.MEASURING_TIME:
		s := m.stopwatch.View() + "\n"
		s = "Elapsed: " + s
		s += m.helpView()
		return s
	case state.WRITING_WORKLOG:
		if m.writeLog.State != huh.StateCompleted {
			return m.writeLog.View()
		}
		m.state.Send()
	}
	return ""
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
