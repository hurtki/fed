package cli_ui

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/hurtki/fed/internal/domain"
)

type CLI struct {
	m sync.Mutex
	s *spinner.Spinner
	w io.Writer
}

func NewCLI(w io.Writer) *CLI {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond, spinner.WithWriter(w))

	return &CLI{
		s: s,
		w: w,
	}
}

func (c *CLI) Status(message string) {
	c.m.Lock()
	defer c.m.Unlock()

	if c.s.Active() {
		c.s.Stop()
	}

	color.RGB(95, 180, 156).Fprintf(c.w, "└->%s\n", message)

	c.s.Start()
}

func (c *CLI) Log(message string) {
	c.m.Lock()
	defer c.m.Unlock()

	wasActive := c.s.Active()
	if wasActive {
		c.s.Stop()
	}

	color.RGB(222, 239, 183).Fprintf(c.w, "%s\n", message)

	if wasActive {
		c.s.Start()
	}
}

func (c *CLI) Plan(plan domain.Plan) {
	c.m.Lock()
	defer c.m.Unlock()

	if c.s.Active() {
		c.s.Stop()
	}

	color.New(color.FgBlue).Fprintln(c.w, "Presented a plan, not implemented in internal/ui/cli")
}

func (c *CLI) Result(success bool, message string) {
	c.m.Lock()
	defer c.m.Unlock()

	if c.s.Active() {
		c.s.Stop()
	}

	if success {
		if message != "" {
			color.RGB(65, 66, 136).Fprintf(c.w, "✓ success: %s\n", message)
		} else {
			color.RGB(65, 66, 136).Fprintln(c.w, "✓ success")
		}
	} else {
		if message != "" {
			color.RGB(104, 45, 99).Fprintf(c.w, "✗ failure: %s\n", message)
		} else {
			color.RGB(104, 45, 99).Fprintln(c.w, "✗ failure")
		}
	}
}

func (c *CLI) Approve(message string) bool {
	c.s.Stop()

	color.New(color.FgYellow).Fprintf(c.w, "⚠  %s \n[y/N]->: ", message)

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	input = strings.TrimSpace(strings.ToLower(input))
	isApproved := input == "y" || input == "yes"

	c.DeleteLastLines(strings.Count(message, "\n") + 2)

	c.s.Start()
	return isApproved
}

func (c *CLI) DeleteLastLines(n int) {
	if n <= 0 {
		return
	}
	seq := strings.Repeat("\x1b[1A\x1b[2K", n)
	fmt.Fprint(c.w, seq)
}
