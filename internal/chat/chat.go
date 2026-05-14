package chat

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type AI interface {
	Generate(ctx context.Context, prompt string) (string, error)
}

type Chat struct {
	ai AI
}

func NewChat(ai AI) *Chat {
	return &Chat{
		ai: ai,
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

	res, err := c.Think(ctx, input)
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return
	}
	fmt.Printf("Thought: %s\n", res.ShortText)
	fmt.Printf("Want to execute: %s?\n", res.ShellAction.Command)
	reader = bufio.NewReader(os.Stdin)
	input, _ = reader.ReadString('\n')

	cmd := exec.Command(
		strings.Split(res.ShellAction.Command, " ")[0],
		strings.Split(res.ShellAction.Command, " ")[1:]...,
	)

	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error when executing: %s\n", err)
		return
	}

	fmt.Println(string(out))
}
