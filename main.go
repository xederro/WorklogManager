package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/xederro/WorklogManager/jira"
	"github.com/xederro/WorklogManager/tui"
	"os"
)

func main() {
	server := os.Getenv("JIRA_URL")
	if server == "" {
		panic("JIRA_URL environment variable not set. Example: https://jira.test.server.com/rest/api/2/")
	}
	j := jira.Jira{}
	j.SetUrlBase(server)
	if _, err := tea.NewProgram(tui.NewModel(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
