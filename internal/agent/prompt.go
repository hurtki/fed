package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hurtki/fed/internal/domain"
)

const (
	maxProjectsFiles = 2500
	fixTries         = 2
)

type PromptOptions struct {
}

func (a *Agent) Prompt(ctx context.Context, msg string, opts *PromptOptions) error {
	projectFiles := getFiles(a.proj.BasePath)
	a.reporter.Status("Thinking")

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
Don't be afraid to take a lot of files into context, the more files, the better the answer.
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

	projFiles := []domain.ProjectFile{}
	for _, filePath := range resDto.WantedFiles {
		pf, err := a.proj.NewFile(filePath)
		if err != nil {
			a.reporter.Log("not existing file from ai, skipping")
			continue
		}
		projFiles = append(projFiles, pf)
	}

	switch {
	case resDto.NotProject:
		a.reporter.Log("Not connected to project request( not implemented agent logic )")
	case resDto.DiscussProject:
		a.reporter.Log("Connected to project 'discuss' request ( not implemented agent logic )")
	case resDto.EditProject:
		return a.EditProject(ctx, msg, projFiles)
	}
	a.reporter.Log(fmt.Sprint(resDto.NotProject, resDto.DiscussProject, resDto.EditProject))

	return nil
}

func getFiles(root string) (res []string) {
	_ = filepath.Walk(root, func(p string, info os.FileInfo, _ error) error {
		if info != nil && !info.IsDir() {
			// cut from abs path to relative path
			p = strings.TrimPrefix(p, root+"/")

			if strings.HasPrefix(p, ".git") {
				return nil
			}
			res = append(res, p)
		}
		return nil
	})
	return
}
