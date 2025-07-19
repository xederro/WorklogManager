package worklogList

import (
	tea "github.com/charmbracelet/bubbletea"
)

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
