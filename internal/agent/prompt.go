package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	maxProjectsFiles = 2500
	fixTries         = 2
)

type PromptOptions struct {
}

func (a *Agent) Prompt(ctx context.Context, msg string, opts *PromptOptions) error {
	projectFiles := getFiles(a.proj.BasePath)
	if len(projectFiles) > maxProjectsFiles {
		a.reporter.Log(
			fmt.Sprintf("too many files in project: %d ( should be lower than %d)",
				len(projectFiles), maxProjectsFiles),
		)
		a.reporter.Status("Thinking without project context")
	} else {
		a.reporter.Status("Thinking")
	}
	projectFilesText := strings.Join(projectFiles, "\n")

	a.reporter.Log(fmt.Sprintf("got %d files in file tree of the project injected into context with their text length of %d symbols", len(projectFiles), len(projectFilesText)))

	aiPrompt := fmt.Sprintf(`
<files_in_project>
%s
</files_in_project>

<task_definition>
You need to determine type of message that user requested
Types of messages:
1. Not connected to project, just discuss
your output: {"not_project": true} response
2. Connected to project but without editing
your output: {"discuss_project": true, "wanted_files": ["", ""]}
wanted files is what files you think you want to have in context to answer user's request
3. Edit Project
your output: {"edit_project": true, "wanted_files": ["", ""]}
wanted files is what files you think you want to have in context to process user's request
</task_definition>

<user_request>
%s
</user_request>`, projectFilesText, msg)

	res, err := a.ai.GenerateJSON(ctx, aiPrompt)
	if err != nil {
		return fmt.Errorf("error from AI: %w", err)
	}

	resDto := PromptTypeResolving{}
	err = json.Unmarshal([]byte(res), &resDto)
	if err != nil {
		for i := range fixTries {
			a.reporter.Status(fmt.Sprintf("Fixing JSON issue ( try %d )", i+1))
			res, err = a.FixResponse(ctx, res, aiPrompt, err.Error())
			if err != nil {
				return err
			}
			resDto = PromptTypeResolving{}
			err = json.Unmarshal([]byte(res), &resDto)
			if err == nil {
				break
			}
		}
	}
	if err != nil {
		return fmt.Errorf("AI couldn't recover from json unmarshaling issue")
	}

	a.reporter.Log(fmt.Sprintf(`
NotProject: %t
DiscussProject: %t
EditProject: %t

Wanted files: %v
	`, resDto.NotProject, resDto.DiscussProject, resDto.EditProject, resDto.WantedFiles))

	return nil
}

func getFiles(root string) (res []string) {
	_ = filepath.Walk(root, func(p string, info os.FileInfo, _ error) error {
		if info != nil && !info.IsDir() {
			res = append(res, p)
		}
		return nil
	})
	return
}
