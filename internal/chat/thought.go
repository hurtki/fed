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

func (c *Chat) Think(ctx context.Context, prompt string) (ThinkResult, error) {
	formattedPrompt := fmt.Sprintf(`You are a UNIX assistant.
Return ONLY a valid JSON object. No markdown, no triple backticks.
Structure:
{"short_text": "...", "shell_action": {"command": "..."}}

User input: %s`, prompt)

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
		// Полезно логировать и ошибку, и сырой ответ для отладки
		return ThinkResult{}, fmt.Errorf("can't unmarshal ai response: %w. Raw: %s", err, res)
	}

	return dto, nil
}
