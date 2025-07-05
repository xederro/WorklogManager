package main

import (
	"context"
	"fmt"
	"github.com/apple/pkl-go/pkl"
	"github.com/xederro/WorklogManager/gen/config"
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

	fmt.Println("Configuration loaded successfully:", cfg)
}
