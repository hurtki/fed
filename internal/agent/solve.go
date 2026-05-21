package agent

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hurtki/fed/internal/domain"
)

func (a *Agent) SolveTask(ctx context.Context, task domain.Task) error {
	var filesContext strings.Builder

	for _, pf := range task.Files {
		absPath := pf.GetAbsPath()

		fData, err := os.ReadFile(absPath)
		if err != nil {
			return fmt.Errorf("one of task files is not reachable: %w", err)
		}

		fmt.Fprintf(&filesContext, `
		<%s>
		%s
		</%s>
		`, pf.Path, string(fData), pf.Path)
	}

	res, err := a.ai.Generate(ctx, fmt.Sprintf(`
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
	"changes": [{"path": "./path", find": "", "replace": ""}, {"find": "", "replace": ""}]
}

in changes, "path" should be one of file paths in given files_context
in changes, "find" key should be minimum 3-4 lines of exact code as in file_context in file
in changes, "replace" is what you want to insert instead of "find" part
</task_definition>

<task_description>
%s
</task_description>
	`, filesContext.String(), task.Description))

	if err != nil {
		return fmt.Errorf("can't generate using ai: %w", err)
	}

	cleanRes := strings.TrimSpace(res)
	cleanRes = strings.TrimPrefix(cleanRes, "```json")
	cleanRes = strings.TrimPrefix(cleanRes, "```")
	cleanRes = strings.TrimSuffix(cleanRes, "```")
	cleanRes = strings.TrimSpace(cleanRes)

	resDto := TaskSolveResponse{}

	err = json.Unmarshal([]byte(cleanRes), &resDto)
	if err != nil {
		fmt.Println(cleanRes)
		return fmt.Errorf("not valid response from ai")
	}

	for _, c := range resDto.Changes {
		fmt.Printf(`
>>>>>> %s
%s
====== %s
%s
`, c.Path, c.Find, c.Path, c.Replace)
		fmt.Printf("Apply?y/n:")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')

		if strings.Contains(input, "y") {
			err = applyTaskChange(*task.Files[0].Project, c)
			if err == nil {
				fmt.Println("\napplied change")
			} else {
				fmt.Printf("\nerr: %s\n", err.Error())
			}
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

func applyTaskChange(proj domain.Project, ch TaskChange) error {
	path := filepath.Join(proj.BasePath, ch.Path)

	f, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("can't read file in task change: %w", err)
	}

	findCount := strings.Count(string(f), ch.Find)

	if findCount != 1 {
		return fmt.Errorf("%d find blocks in given change", findCount)
	}

	replacedF := strings.Replace(string(f), ch.Find, ch.Replace, 1)

	return overwriteFile(path, []byte(replacedF))
}

func overwriteFile(filepath string, newContent []byte) error {
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("can't open file: %w", err)
	}
	defer file.Close()

	_, err = file.Write(newContent)
	if err != nil {
		return fmt.Errorf("write error: %w", err)
	}

	return nil
}
