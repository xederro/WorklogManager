package worklogList

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"io"
)

// WorklogDelegate is a custom delegate for the worklog list.
type WorklogDelegate struct {
	list.DefaultDelegate
}

func NewWorklogDelegate(keys *DelegateKeyMap) WorklogDelegate {
	d := WorklogDelegate{
		list.NewDefaultDelegate(),
	}

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		return nil
	}

	help := []key.Binding{keys.Choose, keys.Worklog, keys.StopAll}

	d.ShortHelpFunc = func() []key.Binding {
		return help
	}

	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}

	return d
}

type DelegateKeyMap struct {
	Choose  key.Binding
	Worklog key.Binding
	StopAll key.Binding
	Refetch key.Binding
}

func (d DelegateKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		d.Choose,
		d.Worklog,
		d.StopAll,
		d.Refetch,
	}
}

func (d DelegateKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			d.Choose,
			d.Worklog,
			d.StopAll,
			d.Refetch,
		},
	}
}

func NewDelegateKeyMap() *DelegateKeyMap {
	return &DelegateKeyMap{
		Choose: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "choose"),
		),
		Worklog: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "worklog"),
		),
		StopAll: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "stop all"),
		),
		Refetch: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "refetch"),
		),
	}
}

// Render prints an listItem.
func (d WorklogDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	var (
		title, desc  string
		matchedRunes []int
		s            = &d.Styles
	)

	if i, ok := item.(*WorklogItem); ok {
		title = i.Title()
		desc = i.Description()
	} else {
		return
	}

	// Conditions
	var (
		isSelected  = index == m.Index()
		emptyFilter = m.FilterState() == list.Filtering && m.FilterValue() == ""
		isFiltered  = m.FilterState() == list.Filtering || m.FilterState() == list.FilterApplied
	)

	if emptyFilter {
		title = s.DimmedTitle.Render(title)
		desc = s.DimmedDesc.Render(desc)
	} else if isSelected && m.FilterState() != list.Filtering {
		if isFiltered {
			// Highlight matches
			unmatched := s.SelectedTitle.Inline(true)
			matched := unmatched.Inherit(s.FilterMatch)
			title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
		}
		title = s.SelectedTitle.Render(title)
		desc = s.SelectedDesc.Render(desc)
	} else {
		if isFiltered {
			// Highlight matches
			unmatched := s.NormalTitle.Inline(true)
			matched := unmatched.Inherit(s.FilterMatch)
			title = lipgloss.StyleRunes(title, matchedRunes, matched, unmatched)
		}
		title = s.NormalTitle.Render(title)
		desc = s.NormalDesc.Render(desc)
	}

	if d.ShowDescription {
		_, _ = fmt.Fprintf(w, "%s\n%s", title, desc)
		return
	}
	_, _ = fmt.Fprintf(w, "%s", title)
}
