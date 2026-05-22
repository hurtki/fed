package chat

import (
	"context"
	"encoding/json"
	"fmt"
)

type ThinkResult struct {
	Text        string       `json:"response"`
	ShellAction *ShellAction `json:"shell_action"`
}

func (c *Chat) Think(ctx context.Context, prompt string, chatContext Context) (ThinkResult, error) {
	formattedPrompt := fmt.Sprintf(`
<context>
Chat Context: %s
</context>

<task_definition>
Output ONLY a valid JSON object. Do not include any thinking, explanations, or text outside of the JSON structure.

You've analyzed the context above.
Now move on to executing the user's actual request.
Execution criteria:
1. Focus exclusively on the task below.
2. Use data from  only to extract facts about chat history with user or information from bash previous bash scripts execution

You are an agent who can run bash scripts to help the user complete their task.
Return ONLY a valid JSON object. No markdown, no triple backticks.
Don't use line breaks or \n
Use ; to separate commands

Structure if you want to execute some script from current working directory of the user
{"response": "...", "shell_action": {"bash_script": "..."}}

bash_script field should be a bash script without shebang.
DON'T USE SUDO COMMAND
DON'T write the same shell script in your text response

Structure if user is not asking for something or you need to get details about his request
You can ask him about some details and then you will recieve his answer so you can
create a better bash script

{"response": "..."}
</task_definition>

<user_query>
%s
</user_query>
`, chatContext.String(), prompt)

	res, err := c.ai.GenerateJSON(ctx, formattedPrompt)
	if err != nil {
		return ThinkResult{}, err
	}

	var dto ThinkResult
	err = json.Unmarshal([]byte(res), &dto)
	if err != nil {
		return ThinkResult{}, fmt.Errorf("can't unmarshal ai response: %w. Raw: %s", err, res)
	}

	return dto, nil
}
