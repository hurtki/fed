package cli_ui

import (
	"bufio"
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
	// Создаем спиннер сразу с привязкой к переданному io.Writer
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond, spinner.WithWriter(w))

	return &CLI{
		s: s,
		w: w,
	}
}

func (c *CLI) Status(message string) {
	c.m.Lock()
	defer c.m.Unlock()

	c.s.Stop()

	// Используем цветной вывод с учетом нашего io.Writer
	color.New(color.FgGreen).Fprintln(c.w, message)

	c.s.Start()
}

func (c *CLI) Log(message string) {
	c.m.Lock()
	defer c.m.Unlock()

	c.s.Stop()
	color.New(color.FgWhite).Fprintf(c.w, "Log: %s\n", message)
	c.s.Start()
}

func (c *CLI) Plan(plan domain.Plan) {
	c.m.Lock()
	defer c.m.Unlock()

	c.s.Stop()
	color.New(color.FgBlue).Fprintln(c.w, "Presented a plan, not implemented in internal/ui/cli")
	c.s.Start() // Не забываем вернуть спиннер, если нужно
}

func (c *CLI) Result(success bool, message string) {
	c.m.Lock()
	defer c.m.Unlock()

	c.s.Stop()

	if success {
		if message != "" {
			color.New(color.FgGreen).Fprintf(c.w, "success: %s\n", message)
		} else {
			color.New(color.FgGreen).Fprintln(c.w, "success")
		}
	} else {
		if message != "" {
			color.New(color.FgRed).Fprintf(c.w, "failure: %s\n", message)
		} else {
			color.New(color.FgRed).Fprintln(c.w, "failure")
		}
	}
}

func (c *CLI) Approve(message string) bool {
	// 1. Сначала останавливаем спиннер под мьютексом
	c.m.Lock()
	c.s.Stop()
	color.New(color.FgRed).Fprintf(c.w, "%s\nApprove y/n: ", message)
	c.m.Unlock() // Обязательно отпускаем мьютекс ПЕРЕД чтением из консоли!

	// 2. Спокойно ждем ввода от пользователя (мьютекс свободен, другие горутины не зависнут)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	// 3. Более надежная валидация ответа
	input = strings.TrimSpace(strings.ToLower(input))
	return input == "y" || input == "yes"
}
