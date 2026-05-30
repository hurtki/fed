package agent

import (
	"context"

	"github.com/hurtki/fed/internal/domain"
)

type AI interface {
	GenerateJSON(ctx context.Context, prompt string) (string, error)
}

type ToolChain interface {
	RunFileChange(ch domain.FileChange) error
	ReadFile(f domain.ProjectFile) ([]byte, error)
}

type AgentReporter interface {
	Status(message string)
	Log(message string)
	Plan(plan domain.Plan)
	Result(success bool, message string)
}

type Agent struct {
	proj domain.Project

	ai        AI
	reporter  AgentReporter
	toolchain ToolChain
}

func NewAgent(ai AI, reporter AgentReporter, proj domain.Project, toolChain ToolChain) Agent {
	return Agent{
		ai:       ai,
		reporter: reporter,
		proj:     proj,
	}
}
