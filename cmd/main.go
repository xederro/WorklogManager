package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/xederro/WorklogManager/internal/config"
	"github.com/xederro/WorklogManager/internal/tui"
	"os"
)

func main() {
	config.Init()
	if _, err := tea.NewProgram(tui.NewModel(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
