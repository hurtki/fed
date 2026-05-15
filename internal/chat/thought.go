package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type ThinkResult struct {
	Text        string       `json:"response"`
	ShellAction *ShellAction `json:"shell_action"`
}

func (c *Chat) Think(ctx context.Context, prompt string, chatContext Context) (ThinkResult, error) {
	formattedPrompt := fmt.Sprintf(`
======
Chat Context: %s
======

You are an agent who can run bash scripts to help the user complete their task.
Return ONLY a valid JSON object. No markdown, no triple backticks.
Don't use line breaks, instead use \n symbols.

Structure if you want to execute some script from current working directory of the user
{"response": "...", "shell_action": {"bash_script": "..."}}

bash_script field should be a bash script without shebang.
It can be couple lines of bash code
DON'T USE SUDO COMMAND
DON'T write the same shell script in your text response

Structure if user is not asking for something or you need to get details about his request
You can ask him about some details and then you will recieve his answer so you can
create a better bash script
{"response": "..."}

User request: %s
`, chatContext.String(), prompt)

	res, err := c.ai.Generate(ctx, formattedPrompt)
	if err != nil {
		return ThinkResult{}, err
	}

	cleanRes := strings.TrimSpace(res)
	cleanRes = strings.TrimPrefix(cleanRes, "```json")
	cleanRes = strings.TrimPrefix(cleanRes, "```")
	cleanRes = strings.TrimSuffix(cleanRes, "```")
	cleanRes = strings.TrimSpace(cleanRes)

	var dto ThinkResult
	err = json.Unmarshal([]byte(cleanRes), &dto)
	if err != nil {
		return ThinkResult{}, fmt.Errorf("can't unmarshal ai response: %w. Raw: %s", err, res)
	}

	return dto, nil
}
