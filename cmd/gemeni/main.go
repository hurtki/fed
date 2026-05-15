package main

import (
	"log/slog"
	"os"

	"github.com/hurtki/fed/internal/chat"
	"github.com/hurtki/fed/internal/config"
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
	chat := chat.NewChat(cl)
	chat.Loop()
}
