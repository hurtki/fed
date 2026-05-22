package agent

import (
	"context"
)

type AI interface {
	GenerateJSON(ctx context.Context, prompt string) (string, error)
}

type Agent struct {
	ai AI
}

func NewAgent(ai AI) Agent {
	return Agent{
		ai: ai,
	}
}
