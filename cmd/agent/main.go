package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/hurtki/fed/internal/agent"
	"github.com/hurtki/fed/internal/config"
	"github.com/hurtki/fed/internal/domain"
	"github.com/hurtki/fed/internal/infrastructure/gemini"
	"github.com/hurtki/fed/internal/infrastructure/ollama"
	"github.com/hurtki/fed/internal/storage"
	"github.com/hurtki/fed/internal/tools"
	cli_ui "github.com/hurtki/fed/internal/ui/cli"
)

func main() {
	ui := cli_ui.NewCLI(os.Stdout, cli_ui.NewNeonCLIPalette())

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

	llmName := ui.RequestMulLines("Chose llm to use(gemini,ollama):")

	switch llmName {
	case "gemini":
		geminiCfg, err := config.LoadGeminiConfigFromEnvFile(envSrc)
		if err != nil {
			return
		}

		ai, err = gemini.NewGeminiAI(geminiCfg.Token, geminiCfg.Model)
		if err != nil {
		}
	case "ollama":
		ollamaCfg, err := config.LoadOllamaConfigFromEnvFile(envSrc)
		if err != nil {
			return
		}

		ai = ollama.NewOllamaClient(ollamaCfg)
	default:
		return
	}

	absPath, _ := filepath.Abs("./")

	proj, err := domain.NewProject(absPath)

	fileRightsStorage := storage.NewMemoryFileRightsStorage()

	toolchain := tools.NewToolChain(ui, fileRightsStorage)

	a := agent.NewAgent(ai, ui, proj, toolchain)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		os.Exit(0)
	}()

	for {
		input := ui.RequestMulLines("->")
		if input == "" {
			return
		}

		err = a.Prompt(context.Background(), input, nil)
		if err != nil {
			ui.Result(false, fmt.Sprintf("error occured: %s", err.Error()))
		} else {
			ui.Result(true, "")
		}
	}

}
