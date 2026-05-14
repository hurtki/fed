package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/hurtki/fed/internal/config"
	"github.com/hurtki/fed/internal/infrastructure/ollama"
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

	ollamaCfg, err := config.LoadOllamaConfig()
	if err != nil {
		logger.Error("can't initialize ollama config", "err", err)
		return
	}

	cl := ollama.NewOllamaClient(ollamaCfg)
	res, err := cl.Generate(context.Background(), "hi")
	logger.Info("cfg", "model", ollamaCfg.Model)
	logger.Info("generate", "res", res, "err", err)
}
