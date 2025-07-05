package main

import (
	"context"
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/apple/pkl-go/pkl"
	"github.com/xederro/WorklogManager/gen/config"
	"net/http"
)

func main() {
	evaluator, err := pkl.NewProjectEvaluator(context.Background(), "pkl/", pkl.PreconfiguredOptions)
	if err != nil {
		panic(err)
	}

	cfg, err := config.Load(context.Background(), evaluator, pkl.FileSource("pkl/dev/config.pkl"))
	if err != nil {
		panic(err)
	}

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

	u, _, err := jiraClient.User.GetSelfWithContext(context.Background())
	if err != nil {
		panic(fmt.Errorf("failed to get user: %w", err))
	}

	fmt.Printf("Email: %v\n", u.EmailAddress)
	fmt.Println("Success!")
}
