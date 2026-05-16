package chat

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type AI interface {
	Generate(ctx context.Context, prompt string) (string, error)
}

type Chat struct {
	ai AI
	Context
}

func NewChat(ai AI) *Chat {
	return &Chat{
		ai: ai,
		Context: Context{
			Messages: []string{},
		},
	}
}

func (c *Chat) Loop() {
	for {
		c.entry(context.Background())
	}
}

func (c *Chat) entry(ctx context.Context) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("->")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	start := time.Now()

	res, err := c.Think(ctx, input, c.Context)
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return
	}
	c.Context.AddUserMessage(input)
	fmt.Printf("%s: %s\n", time.Since(start).String(), res.Text)

	c.Context.AddAgentThought(res.Text)

	if res.ShellAction != nil {
		fmt.Printf("Want to execute?\n===\n %s\n===\ny/n:", res.ShellAction.Command)
		reader = bufio.NewReader(os.Stdin)
		input, _ = reader.ReadString('\n')

		if strings.Contains(input, "y") {
			cmd := exec.Command(
				"bash",
				"-c",
				res.ShellAction.Command,
			)

			outBytes, err := cmd.Output()
			out := string(outBytes)

			if err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					exitCode := exitErr.ExitCode()
					out += fmt.Sprintf(" exit code: %d", exitCode)
				} else {
					out = err.Error()
				}
			}
			fmt.Printf("\noutput:===\n%s\n===\n", out)
			c.Context.AddShellOutput(res.ShellAction.Command, out)
		}

	}
}
