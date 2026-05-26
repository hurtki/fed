package cli_reporter

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hurtki/fed/internal/domain"
)

type CLI struct {
	w io.Writer
}

func NewCLI(w io.Writer) *CLI {
	return &CLI{
		w: w,
	}
}

func (c *CLI) Status(message string) {
	fmt.Fprintf(c.w, "Thought: %s\n", message)
}

func (c *CLI) Log(message string) {
	fmt.Fprintf(c.w, "Log: %s\n", message)
}

func (c *CLI) Plan(plan domain.Plan) {
	fmt.Fprintf(c.w, "Presented a plan")
}

func (c *CLI) Result(success bool, message string) {
	fmt.Fprintf(c.w, "result, success: %t", success)
}

func (c *CLI) Approve(message string) bool {
	fmt.Fprintf(c.w, "%s\nApprove y/n:", message)

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')

	return strings.Contains(input, "y")
}
