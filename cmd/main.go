package main

import (
	"flag"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-co-op/gocron/v2"
	"github.com/xederro/WorklogManager/internal/config"
	"github.com/xederro/WorklogManager/internal/tui"
	"github.com/xederro/WorklogManager/internal/tui/worklogList"
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

	s, err := gocron.NewScheduler()
	if err != nil {
		fmt.Printf("Failed to create scheduler: %v\n", err)
		os.Exit(2)
	}

	_, err = s.NewJob(
		gocron.DurationJob(
			config.Conf.Jira.RefetchInterval.GoDuration(),
		),
		gocron.NewTask(
			func() {
				tui.Ch <- worklogList.GetItemsToUpdate()
			},
		),
	)
	if err != nil {
		fmt.Printf("Failed to start task %v\n", err)
		os.Exit(4)
	}

	s.Start()
	if _, err := tea.NewProgram(tui.NewModel(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(3)
	}
	_ = s.Shutdown()
}
