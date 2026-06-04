package tools

import (
	"fmt"
	"os/exec"
)

func (t *ToolChain) RunBashScript(script string) (string, bool) {
	approved := t.approver.Approve(fmt.Sprintf("Want to execute?\n===\n %s\n===", script))
	if !approved {
		return "", false
	}

	cmd := exec.Command(
		"bash",
		"-c",
		script,
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

	return out, true
}
