package config

import (
	"context"
	"github.com/andygrunwald/go-jira"
	"github.com/xederro/WorklogManager/internal/gen/config"
	"google.golang.org/genai"
	"net/http"
)

var Conf *config.Config
var JiraClient *jira.Client
var GoogleClient *genai.Client

func Init(path string) {
	cfg, err := config.LoadFromPath(context.Background(), path)
	if err != nil {
		panic(err)
	}

	Conf = cfg

	var client *http.Client

	switch cfg.Jira.ServerType {
	case "cloud":
		transport := jira.BasicAuthTransport{
			Username: cfg.Jira.CloudConfig.Email,
			Password: cfg.Jira.CloudConfig.ApiToken,
		}
		client = transport.Client()
		break
	case "on-premise":
		transport := &jira.BearerAuthTransport{
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
}
