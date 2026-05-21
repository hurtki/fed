package agent

import (
	"context"
)

type AI interface {
	Generate(ctx context.Context, prompt string) (string, error)
}

type Agent struct {
	ai AI
}

func NewAgent(ai AI) Agent {
	return Agent{
		ai: ai,
	}
}
