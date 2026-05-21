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
	task := domain.Task{
		Description: "delete http server logic, devide class should stay",
		Solved:      false,
		Files: []domain.ProjectFile{
			{Path: "./main.py", Project: &domain.Project{BasePath: "/Users/hurtki/tmp/"}},
		},
	}
	err = a.SolveTask(context.Background(), task)
	logger.Info("solved task", "err", err)
}
