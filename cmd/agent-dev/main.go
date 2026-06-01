package main

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/hurtki/fed/internal/agent"
	"github.com/hurtki/fed/internal/config"
	"github.com/hurtki/fed/internal/domain"
	"github.com/hurtki/fed/internal/infrastructure/gemini"
	"github.com/hurtki/fed/internal/infrastructure/ollama"
	cli_reporter "github.com/hurtki/fed/internal/reporter/cli"
	"github.com/hurtki/fed/internal/storage"
	"github.com/hurtki/fed/internal/tools"
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

	var ai agent.AI

	cmdArgs := os.Args[1:]

	envSrc := ""
	switch len(cmdArgs) {
	case 0:
		envSrc = ".env"
	case 1:
		envSrc = cmdArgs[0]
	default:
		fmt.Println("too many args")
		return
	}

	fmt.Print("Chose llm to use(gemini,ollama):")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	switch input {
	case "gemini":
		logger.Info("trying to init gemini config from .env file")

		geminiCfg, err := config.LoadGeminiConfigFromEnvFile(envSrc)
		if err != nil {
			logger.Error("can't initialize gemini config", "err", err)
			return
		}

		ai, err = gemini.NewGeminiAI(geminiCfg.Token, geminiCfg.Model)
		if err != nil {
			logger.Error("can't initialize gemini", "err", err)
		}
	case "ollama":
		logger.Info("trying to init gemini config from .env file")

		ollamaCfg, err := config.LoadOllamaConfigFromEnvFile(envSrc)
		if err != nil {
			logger.Error("can't initialize ollama config", "err", err)
			return
		}

		ai = ollama.NewOllamaClient(ollamaCfg)
	default:
		logger.Info("not available llm source")
		return
	}

	reporter := cli_reporter.NewCLI(os.Stdout)

	absPath, _ := filepath.Abs("./")

	proj, err := domain.NewProject(absPath)

	fileRightsStorage := storage.NewMemoryFileRightsStorage()

	toolchain := tools.NewToolChain(reporter, fileRightsStorage)

	a := agent.NewAgent(ai, reporter, proj, toolchain)

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("->")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		err = a.Prompt(context.Background(), input, nil)
		logger.Info("Prompt executed", "err", err)
	}

}
