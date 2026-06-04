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
	"github.com/hurtki/fed/internal/domain"
	"golang.org/x/term"
)

type CLI struct {
	m sync.Mutex
	s *spinner.Spinner
	w io.Writer

	palette CLIPalette
}

func NewCLI(w io.Writer, palette CLIPalette) *CLI {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond, spinner.WithWriter(w))

	return &CLI{
		s:       s,
		w:       w,
		palette: palette,
	}
}

func (c *CLI) Status(message string) {
	c.m.Lock()
	defer c.m.Unlock()

	if c.s.Active() {
		c.s.Stop()
	}

	c.palette.Status.Fprintf(c.w, "└->%s\n", message)

	c.s.Start()
}

func (c *CLI) Log(message string) {
	c.m.Lock()
	defer c.m.Unlock()

	wasActive := c.s.Active()
	if wasActive {
		c.s.Stop()
	}

	c.palette.Log.Fprintf(c.w, "%s\n", message)

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

	c.palette.Plan.Fprintln(c.w, "Presented a plan, not implemented in internal/ui/cli")
}

func (c *CLI) Result(success bool, message string) {
	c.m.Lock()
	defer c.m.Unlock()

	if c.s.Active() {
		c.s.Stop()
	}

	if success {
		if message != "" {
			c.palette.ResultSuccess.Fprintf(c.w, "✓ success: %s\n", message)
		} else {
			c.palette.ResultSuccess.Fprintln(c.w, "✓ success")
		}
	} else {
		if message != "" {
			c.palette.ResultFailure.Fprintf(c.w, "✗ failure: %s\n", message)
		} else {
			c.palette.ResultFailure.Fprintln(c.w, "✗ failure")
		}
	}
}

func (c *CLI) Approve(message string) bool {
	c.m.Lock()
	defer c.m.Unlock()
	c.s.Stop()

	c.palette.Approve.Fprintf(c.w, "⚠  %s \n[y/N]->: ", message)

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	input = strings.TrimSpace(strings.ToLower(input))
	isApproved := input == "y" || input == "yes"

	c.deleteLastLines(strings.Count(message, "\n") + 2)

	c.s.Start()
	return isApproved
}

func (c *CLI) deleteLastLines(n int) {
	if n <= 0 {
		return
	}
	seq := strings.Repeat("\x1b[1A\x1b[2K", n)
	fmt.Fprint(c.w, seq)
}
func (c *CLI) RequestMulLines(prompt string) string {
	c.m.Lock()
	defer c.m.Unlock()

	fd := int(os.Stdin.Fd())
	if !term.IsTerminal(fd) {
		c.m.Unlock()
		return ""
	}
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		c.m.Unlock()
		return ""
	}

	defer func() {
		term.Restore(fd, oldState)
	}()

	c.palette.RequestMulLines.Fprint(c.w, prompt)

	buf := make([]byte, 1)
	var res strings.Builder

	timesEnterInRow := 0

	for {
		_, err := os.Stdin.Read(buf)
		if err != nil {
			return ""
		}

		char := buf[0]

		if char == 3 {
			return ""
		}

		if char == '\t' {
			break
		}

		if char == 127 || char == 8 {
			str := res.String()
			if len(str) > 0 {
				res.Reset()
				res.WriteString(str[:len(str)-1])

				fmt.Print("\b \b")
			}
			continue
		}

		if char == '\r' || char == '\n' {
			if timesEnterInRow > 0 {
				res := res.String()
				return res[:len(res)-1]
			}
			res.WriteByte('\n')
			fmt.Print("\r\n")
			timesEnterInRow++
			continue
		}

		res.WriteByte(char)
		c.palette.RequestMulLines.Fprint(c.w, string(char))
		timesEnterInRow = 0
	}

	fmt.Print("\r\n")
	return res.String()
}
