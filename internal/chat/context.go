package chat

import (
	"fmt"
	"strings"
)

type Context struct {
	Messages []string
}

func (c *Context) String() string {
	if len(c.Messages) < 1 {
		return "CONTEXT IS CLEAR"
	}
	return strings.Join(c.Messages, "\n")
}

func (c *Context) AddUserMessage(msg string) {
	c.Messages = append(c.Messages, "USER:"+msg)
}

func (c *Context) AddAgentThought(msg string) {
	c.Messages = append(c.Messages, "AGENT:"+msg)
}

func (c *Context) AddShellOutput(cmd string, out string) {
	c.Messages = append(c.Messages,
		fmt.Sprintf("SHELL CMD EXECUTION: command: %s, output: %s", cmd, out),
	)
}
