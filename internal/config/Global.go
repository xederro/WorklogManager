package config

import (
	"context"
	"fmt"
	gojira "github.com/andygrunwald/go-jira"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-co-op/gocron/v2"
	"github.com/xederro/WorklogManager/internal/gen/config"
	"github.com/xederro/WorklogManager/internal/jira"
	"google.golang.org/genai"
	"net/http"
	"os"
)

var (
	Conf         *config.Config
	JiraClient   *jira.Client
	GoogleClient *genai.Client
	Schedule     gocron.Scheduler
	Ch           chan tea.Cmd
)

// Init initializes the global configs.
func Init(path string) {
	cfg, err := config.LoadFromPath(context.Background(), path)
	if err != nil {
		panic(err)
	}

	Conf = cfg

	var client *http.Client

	switch cfg.Jira.ServerType {
	case "cloud":
		transport := gojira.BasicAuthTransport{
			Username: cfg.Jira.CloudConfig.Email,
			Password: cfg.Jira.CloudConfig.ApiToken,
		}
		client = transport.Client()
		break
	case "on-premise":
		transport := &gojira.BearerAuthTransport{
			Token: cfg.Jira.OnPremiseConfig.Pat,
		}
		client = transport.Client()
		break
	default:
		panic("Unknown server type: " + cfg.Jira.ServerType)
	}

	jiraClient, err := jira.NewClient(client, cfg.Jira.Url)
	if err != nil {
		panic(err)
	}
	JiraClient = jiraClient

	if cfg.UseAi {
		ctx := context.Background()
		cc := &genai.ClientConfig{
			APIKey: cfg.GoogleAi.APIKey,
		}
		googleClient, err := genai.NewClient(ctx, cc)
		if err != nil {
			panic(err)
		}
		GoogleClient = googleClient
	}

	Ch = make(chan tea.Cmd, 10)

	s, err := gocron.NewScheduler()
	if err != nil {
		fmt.Printf("Failed to create scheduler: %v\n", err)
		os.Exit(2)
	}
	for _, request := range Conf.Jira.Requests {
		_, err = s.NewJob(
			gocron.DurationJob(
				request.RefetchInterval.GoDuration(),
			),
			gocron.NewTask(
				func() {
					Ch <- JiraClient.GetItemsToUpdate(request.Jql)
				},
			),
		)
		if err != nil {
			fmt.Printf("Failed to start task %v\n", err)
			panic(err)
		}
	}
	s.Start()
	Schedule = s
}

// Shutdown gracefully shuts down the global configs.
func Shutdown() error {
	if Conf == nil {
		return nil
	}

	err := Schedule.Shutdown()
	if err != nil {
		return fmt.Errorf("failed to shutdown scheduler: %w", err)
	}

	close(Ch)

	return nil
}

// TriggerUpdate runs all jobs immediately.
func TriggerUpdate() {
	if Conf == nil {
		return
	}

	for _, job := range Schedule.Jobs() {
		err := job.RunNow()
		if err != nil {
			panic(fmt.Errorf("failed to run job now: %w", err))
		}
	}
}
