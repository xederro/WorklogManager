package main

import (
	"context"
	"fmt"
	"github.com/xederro/WorklogManager/internal/config"
	"google.golang.org/genai"
	"log"
)

func main() {
	config.Init()

	u, _, err := config.JiraClient.User.GetSelfWithContext(context.Background())
	if err != nil {
		panic(fmt.Errorf("failed to get user: %w", err))
	}

	fmt.Printf("Email: %v\n", u.EmailAddress)
	fmt.Println("Success!")

	if config.Conf.UseAi {
		result, err := config.GoogleClient.Models.GenerateContent(
			context.Background(),
			config.Conf.GoogleAi.DefaultModel,
			genai.Text(config.Conf.GoogleAi.DefaultPrompt+" "+u.EmailAddress+" did a lot of work today :3"),
			&genai.GenerateContentConfig{
				ThinkingConfig: &genai.ThinkingConfig{
					ThinkingBudget:  nil,
					IncludeThoughts: false,
				},
			},
		)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result.Text())
	}
}
