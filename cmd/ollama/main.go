package main

import (
	"log/slog"
	"os"

	"github.com/hurtki/fed/internal/chat"
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
	logger.Info("cfg", "model", ollamaCfg.Model)

	chat := chat.NewChat(cl)
	chat.Loop()
}
