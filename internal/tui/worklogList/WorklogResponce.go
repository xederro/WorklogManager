package worklogList

import tea "github.com/charmbracelet/bubbletea"

type WorklogResponse struct {
	Err      error
	Affected *WorklogItem
}

func ReturnCmd(err error, affected *WorklogItem) tea.Cmd {
	return func() tea.Msg {
		return WorklogResponse{
			Err:      err,
			Affected: affected,
		}
	}
}
