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

func (a *Agent) editStep(ctx context.Context, step domain.Step) error {
	var filesContext strings.Builder

	pfs := make(map[string]domain.ProjectFile)

	for _, pf := range step.Files {
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
		pfs[pf.RelativePath] = pf
		fmt.Fprintf(&filesContext, "<file src=%s editable=true>%s</file>", pf.RelativePath, string(fileData))
	}

	aiPrompt := fmt.Sprintf(`
<files_context>
%s
</files_context>

<task_definition>
Solve task described below by editing files in files context using search/replace output
YOUR OUTPUT SHOULD BE A VALID JSON OBJECT
Return ONLY a valid JSON object. No markdown, no triple backticks.

Do not include any thinking, explanations, or text outside of the JSON structure.
You've analyzed the context above.
Now move on to executing the user's actual request.
Execution criteria:
1. Focus exclusively on the task below.
2. Use data from  only to extract facts about chat history with user or information from bash previous bash scripts execution


Structure you need to give as output:
{
	"changes": [{"path": "", "find": "", "replace": ""}, {"path": "", "find": "", "replace": ""} ...]
}

in changes, "path" should be one of file paths, in given files_context it mentioned in src field
in changes, "find" key should be minimum 3-4 lines of exact code as in file_context in file
in changes, "replace" is what you want to insert instead of "find" part
</task_definition>

<task_description>
%s
</task_description>
	`, filesContext.String(), step.Description)

	res, err := a.ai.GenerateJSON(ctx, aiPrompt)

	if err != nil {
		return fmt.Errorf("can't generate using ai: %w", err)
	}

	resDto := TaskSolveResponse{}

	err = json.Unmarshal([]byte(res), &resDto)

	if err != nil {
		for i := range fixTries {
			a.reporter.Status(fmt.Sprintf("Fixing JSON issue ( try %d )", i+1))
			res, err = a.FixResponse(ctx, res, aiPrompt, err.Error())
			if err != nil {
				return err
			}
			resDto = TaskSolveResponse{}
			err = json.Unmarshal([]byte(res), &resDto)
			if err == nil {
				break
			}
		}
	}

	var changes []domain.FileChange
	for _, c := range resDto.Changes {
		changes = append(changes, pfs[c.Path].NewChange(c.Find, c.Replace))
	}

	err = a.toolchain.RunFileChanges(changes)
	if err != nil {
		if errors.Is(err, tools.ErrUserDenied) {
			return fmt.Errorf("user decieded not to edit file")
		}
		return err
	}

	for _, change := range changes {
		findLinesCount := strings.Count(change.Find, "\n")
		replaceLinesCount := strings.Count(change.Replace, "\n")
		switch {
		case findLinesCount > replaceLinesCount:
			a.reporter.Log(fmt.Sprintf("edited %s, %d lines deleted", change.File.RelativePath, findLinesCount-replaceLinesCount))
		case replaceLinesCount > findLinesCount:
			a.reporter.Log(fmt.Sprintf("edited %s, %d lines added", change.File.RelativePath, replaceLinesCount-findLinesCount))
		case replaceLinesCount == findLinesCount:
			a.reporter.Log(fmt.Sprintf("edited %s, %d lines edited", change.File.RelativePath, replaceLinesCount))
		}
	}

	return nil
}

type TaskSolveResponse struct {
	Changes []TaskChange
}

type TaskChange struct {
	Path    string `json:"path"`
	Find    string `json:"find"`
	Replace string `json:"replace"`
}
