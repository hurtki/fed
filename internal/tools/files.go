package tools

import (
	"fmt"
	"os"
	"strings"

	"github.com/hurtki/fed/internal/domain"
)

func (t *ToolChain) RunFileChange(ch domain.FileChange) error {
	f, err := os.ReadFile(ch.File.GetAbsPath())
	if err != nil {
		return fmt.Errorf("can't read file in task change: %w", err)
	}

	findCount := strings.Count(string(f), ch.Find)

	if findCount > 1 {
		return ErrTooManyFindBlocks
	} else if findCount < 1 {
		return ErrNoFindBlocks
	}

	replacedF := strings.Replace(string(f), ch.Find, ch.Replace, 1)

	approved := t.approver.Approve(fmt.Sprintf(
		`>>>>>> %s
%s
====== %s
%s`, ch.File.RelativePath, ch.Find, ch.File.RelativePath, ch.Replace))

	if !approved {
		return ErrUserDenied
	}

	return overwriteFile(ch.File.GetAbsPath(), []byte(replacedF))
}

func overwriteFile(filepath string, newContent []byte) error {
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("can't open file: %w", err)
	}
	defer file.Close()

	_, err = file.Write(newContent)
	if err != nil {
		return fmt.Errorf("write error: %w", err)
	}

	return nil
}
