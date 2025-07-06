package main

import (
	"flag"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/xederro/WorklogManager/internal/config"
	"github.com/xederro/WorklogManager/internal/tui"
	"os"
)

func main() {
	path := flag.String("config", "config.pkl", "The path to the config file")
	flag.Parse()
	if *path == "" {
		fmt.Println("Using default config file: config.pkl")
		*path = "config.pkl"
	}

	if _, err := os.Stat(*path); os.IsNotExist(err) {
		fmt.Printf("Config file not found at %s. Please provide a valid path.\n", *path)
		os.Exit(1)
	}

	config.Init(*path)

	if _, err := tea.NewProgram(tui.NewModel(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
