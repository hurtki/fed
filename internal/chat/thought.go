package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type ThinkResult struct {
	ShortText   string       `json:"short_text"`
	ShellAction *ShellAction `json:"shell_action"`
}

func (c *Chat) Think(ctx context.Context, prompt string, chatContext Context) (ThinkResult, error) {
	formattedPrompt := fmt.Sprintf(`
You are a UNIX assistant.
Return ONLY a valid JSON object. No markdown, no triple backticks.

Structure if you want to execute some command:
{"short_text": "...", "shell_action": {"command": "..."}}

command field should be a bash script without shebang.
command will be executed using bash -c "command"

Structure if you don't want to execute command:
{"short_text": "..."}

short_text field should describe in at least 2 sentences what you think right now

Chat Context: %s
User input: %s`, chatContext.String(), prompt)

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
