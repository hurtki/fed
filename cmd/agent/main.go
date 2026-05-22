package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/hurtki/fed/internal/agent"
	"github.com/hurtki/fed/internal/config"
	"github.com/hurtki/fed/internal/domain"
	"github.com/hurtki/fed/internal/infrastructure/gemini"
)

func main() {
	logger := slog.New(
		slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		),
	)

	geminiCfg, err := config.LoadGeminiConfig()
	if err != nil {
		logger.Error("can't initialize ollama config", "err", err)
		return
	}

	cl, err := gemini.NewGeminiAI(geminiCfg.Token, geminiCfg.Model)
	if err != nil {
		logger.Error("can't initialize gemini", "err", err)
	}

	a := agent.NewAgent(cl)

	proj := &domain.Project{BasePath: "/Users/hurtki/Projects/test/github-fetcher/"}

	task := domain.Task{
		Description: "separate graphQL API calling logic from main.go to graphql.go. And add showcase in main.go of its usage",
		Solved:      false,
		Files: []domain.ProjectFile{
			{Path: "./main.go", Project: proj},
			{Path: "./graphql.go", Project: proj},
		},
	}

	err = a.SolveTask(context.Background(), task)

	logger.Info("solved task", "err", err)
}
