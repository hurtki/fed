package agent

import (
	"context"
	"fmt"
)

func (a *Agent) FixResponse(ctx context.Context, res, prompt, problem string) (string, error) {
	aiPrompt := fmt.Sprintf(`
<prvious_prompt>
%s
</prvious_prompt>

<ai_answer>
%s
</ai_answer>

<ai_asnwer_problem>
%s
</ai_asnwer_problem>


<task_definition>
AI made a mistake when answering
fix it to exactly match what was asked in prevoius_prompt
</task_definition>`, prompt, res, problem)
	return a.ai.GenerateJSON(ctx, aiPrompt)
}
