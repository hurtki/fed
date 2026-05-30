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

	geminiCfg, err := config.LoadGeminiConfig()
	if err != nil {
		logger.Error("can't initialize ollama config", "err", err)
		return
	}

	cl, err := gemini.NewGeminiAI(geminiCfg.Token, geminiCfg.Model)
	if err != nil {
		logger.Error("can't initialize gemini", "err", err)
	}

	reporter := cli_reporter.NewCLI(os.Stdout)

	absPath, _ := filepath.Abs("./")

	proj, err := domain.NewProject(absPath)

	fileRightsStorage := storage.NewMemoryFileRightsStorage()

	toolchain := tools.NewToolChain(reporter, fileRightsStorage)

	a := agent.NewAgent(cl, reporter, proj, toolchain)

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("->")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		err = a.Prompt(context.Background(), input, nil)
		logger.Info("Prompt executed", "err", err)
	}

}
