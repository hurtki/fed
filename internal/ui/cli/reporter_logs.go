package cli_ui

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hurtki/fed/internal/domain"
)

type CLILogs struct {
	w io.Writer
}

func NewCLILogs(w io.Writer) *CLILogs {
	return &CLILogs{
		w: w,
	}
}

func (c *CLILogs) Status(message string) {
	fmt.Fprintf(c.w, "Thought: %s\n", message)
}

func (c *CLILogs) Log(message string) {
	fmt.Fprintf(c.w, "Log: %s\n", message)
}

func (c *CLILogs) Plan(plan domain.Plan) {
	fmt.Fprintf(c.w, "Presented a plan")
}

func (c *CLILogs) Result(success bool, message string) {
	fmt.Fprintf(c.w, "result, success: %t", success)
}

func (c *CLILogs) Approve(message string) bool {
	fmt.Fprintf(c.w, "%s\nApprove y/n:", message)

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')

	return strings.Contains(input, "y")
}
