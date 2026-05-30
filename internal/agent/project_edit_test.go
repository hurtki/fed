package agent_test

import (
	"context"
	"testing"

	"github.com/hurtki/fed/internal/agent"
	"github.com/hurtki/fed/internal/domain"
)

func TestAgent_EditProject(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		ai        agent.AI
		reporter  agent.AgentReporter
		proj      domain.Project
		toolChain agent.ToolChain
		// Named input parameters for target function.
		msg              string
		initProjectFiles []domain.ProjectFile
		wantErr          bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := agent.NewAgent(tt.ai, tt.reporter, tt.proj, tt.toolChain)
			gotErr := a.EditProject(context.Background(), tt.msg, tt.initProjectFiles)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("EditProject() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("EditProject() succeeded unexpectedly")
			}
		})
	}
}
