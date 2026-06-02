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

	eligible, err := t.fileRightsStorage.EligibleForEdit(ch.File)

	if err != nil {
		fmt.Printf("error when getting rights for file in rights storage: %s\n", err.Error())
	}

	if !eligible {
		approved := t.approver.Approve(fmt.Sprintf(
			`>>>>>> %s
%s
====== %s
%s`, ch.File.RelativePath, ch.Find, ch.File.RelativePath, ch.Replace))

		if !approved {
			return ErrUserDenied
		}

		err = t.fileRightsStorage.SetEligibleForEdit(ch.File)
		if err != nil {
			fmt.Printf("error when setting rights for file %s: %s", ch.File.GetAbsPath(), err.Error())
		}
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

func (t *ToolChain) ReadFile(pf domain.ProjectFile) ([]byte, error) {

	eligible, _ := t.fileRightsStorage.EligibleForRead(pf)
	if !eligible {
		approved := t.approver.Approve(fmt.Sprintf("Read file: %s?", pf.GetAbsPath()))
		if !approved {
			return nil, ErrUserDenied
		}
		err := t.fileRightsStorage.SetEligibleForRead(pf)
		if err != nil {
			fmt.Printf("error when setting rights for file %s: %s", pf.GetAbsPath(), err.Error())
		}
	}

	return os.ReadFile(pf.GetAbsPath())
}
