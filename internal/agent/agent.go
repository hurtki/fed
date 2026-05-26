package agent

import (
	"context"

	"github.com/hurtki/fed/internal/domain"
)

type AI interface {
	GenerateJSON(ctx context.Context, prompt string) (string, error)
}

type AgentReporter interface {
	Status(message string)
	Log(message string)
	Plan(plan domain.Plan)
	Result(success bool, message string)
}

type Agent struct {
	proj domain.Project

	ai       AI
	reporter AgentReporter
}

func NewAgent(ai AI, reporter AgentReporter, proj domain.Project) Agent {
	return Agent{
		ai:       ai,
		reporter: reporter,
		proj:     proj,
	}
}
