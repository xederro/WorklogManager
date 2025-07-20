package config

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	gojira "github.com/andygrunwald/go-jira"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-co-op/gocron/v2"
	"github.com/xederro/WorklogManager/db"
	"github.com/xederro/WorklogManager/internal/gen/config"
	"github.com/xederro/WorklogManager/internal/gen/sqlc"
	"github.com/xederro/WorklogManager/internal/jira"
	"google.golang.org/genai"
	_ "modernc.org/sqlite"
	"net/http"
	"os"
)

var (
	Conf         *config.Config
	JiraClient   *jira.Client
	GoogleClient *genai.Client
	Schedule     gocron.Scheduler
	Queries      *sqlc.Queries
	Ch           chan tea.Cmd
	database     *sql.DB
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

	sqlite, err := sql.Open("sqlite", cfg.DbPath)
	if err != nil {
		fmt.Printf("Failed to start db %v\n", err)
		panic(err)
	}
	if _, err = sqlite.ExecContext(context.Background(), db.WorklogsSchema); err != nil {
		fmt.Printf("Failed to start db %v\n", err)
		panic(err)
	}

	Queries = sqlc.New(sqlite)
	database = sqlite

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
func Shutdown() {
	if Conf == nil {
		return
	}

	_ = Schedule.Shutdown()

	close(Ch)

	_ = database.Close()
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
