package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/hurtki/fed/internal/domain"
	"github.com/hurtki/fed/internal/tools"
)

func (a *Agent) EditProject(ctx context.Context, msg string, initProjectFiles []domain.ProjectFile) error {
	var filesContext strings.Builder
	for _, pf := range initProjectFiles {
		fileData, err := a.toolchain.ReadFile(pf)
		if err != nil {
			if errors.Is(err, tools.ErrUserDenied) {
				fmt.Fprintf(&filesContext, "<file src=%s editable=false>user denied to access to this file</file>", pf.RelativePath)
			} else {
				fmt.Fprintf(&filesContext, "<file src=%s editable=false>can't read file, error: %s</file>", pf.RelativePath, err.Error())
			}
			continue
		}
		if bytes.IndexByte(fileData[:min(512, len(fileData))], 0) != -1 {
			fmt.Fprintf(&filesContext, "<file src=%s editable=false>file contains null byte (0x00) in first 512 bytes, likely binary</file>", pf.RelativePath)
			continue
		}
		fmt.Fprintf(&filesContext, "<file src=%s editable=true>%s</file>", pf.RelativePath, string(fileData))
	}

	aiPrompt := fmt.Sprintf(`
<files_context>
%s
</files_context>

<task_definition>
Create a plan for project editing splited in steps
Initial description lays in initial_description tag
Description of
Every step should be one strict task, that can be completed in 1-5 patches of files.
Example 1 of one step: "change loop iterations count in "filename" from 1000 to 98
Example 2 of one step: "delete current implementation of PostgresSQL repo located in internal/storage/postgres.go with structure name UsersRepo"

Your output should be valid JSON with this structure

{
	"steps": [
		{"description": "description of tasks that should include clear changes that will need to be made and in what files should those changes be made, will be used by other agent ot generate patches"}
	]
}
</task_definition>

<initial_description>
%s
</initial_description>
	`, filesContext.String(), msg)

	a.reporter.Status("Planning")
	res, err := a.ai.GenerateJSON(ctx, aiPrompt)

	if err != nil {
		return fmt.Errorf("error from AI: %w", err)
	}

	resDto := EditingPlan{}
	err = json.Unmarshal([]byte(res), &resDto)
	if err != nil {
		for i := range fixTries {
			a.reporter.Status(fmt.Sprintf("Fixing JSON issue ( try %d )", i+1))
			res, err = a.FixResponse(ctx, res, aiPrompt, err.Error())
			if err != nil {
				return err
			}
			resDto = EditingPlan{}
			err = json.Unmarshal([]byte(res), &resDto)
			if err == nil {
				break
			}
		}
	}
	if err != nil {
		return fmt.Errorf("AI couldn't recover from json unmarshaling issue")
	}

	for i, s := range resDto.Steps {
		step, err := domain.NewStep(s.Description, initProjectFiles)
		if err != nil {
			return fmt.Errorf("can't create step: %w", err)
		}
		a.reporter.Status(fmt.Sprintf("working on step %d: %s", i+1, step.Description))
		err = a.editStep(ctx, step)
		if err != nil {
			return fmt.Errorf("step exited with error: %w", err)
		}
	}

	return nil
}

type EditingPlan struct {
	Steps []struct {
		Description string
	}
}
