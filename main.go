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
	token := os.Getenv("JIRA_PAT")
	if server == "" {
		panic("JIRA_URL environment variable not set. Example: https://jira.test.server.com/rest/api/2")
	}
	if token == "" {
		panic("JIRA_PAT environment variable not set. Add PAT token")
	}
	j := jira.Jira{}
	j.SetUrlBase(server)
	err := j.SetAuth(token)
	if err != nil {
		panic(fmt.Sprintf("There was a problem with checking your PAT: %s", err.Error()))
	}
	if _, err := tea.NewProgram(tui.NewModel(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
