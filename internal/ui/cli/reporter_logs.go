package cli_ui

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/hurtki/fed/internal/domain"
)

type CLILogs struct {
	logger *slog.Logger
}

func NewCLILogs(logger *slog.Logger) *CLILogs {
	return &CLILogs{
		logger: logger,
	}
}

func (c *CLILogs) Status(message string) {
	c.logger.Info("Status", "message", message)
}

func (c *CLILogs) Log(message string) {
	c.logger.Info("Log", "message", message)
}

func (c *CLILogs) Plan(plan domain.Plan) {
	c.logger.Info("Plan", "plan", plan, "status", "not implemented domain model")
}

func (c *CLILogs) Result(success bool, message string) {
	c.logger.Info("rseult", "success", success, "message", message)
}

func (c *CLILogs) Approve(message string) bool {
	c.logger.Info("Approve", "message", message)
	fmt.Print("Approve y/n:")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')

	return strings.Contains(input, "y")
}
