package tools

import "github.com/hurtki/fed/internal/domain"

type DiffFormatter interface {
	FormatDiff(path, find, replace string) string
}

type Approver interface {
	Approve(msg string) bool
	DiffFormatter
}

type FileRightsStorage interface {
	EligibleForEdit(pf domain.ProjectFile) (bool, error)
	EligibleForRead(pf domain.ProjectFile) (bool, error)

	SetEligibleForEdit(pf domain.ProjectFile) error
	SetEligibleForRead(pf domain.ProjectFile) error
}

type ToolChain struct {
	approver          Approver
	fileRightsStorage FileRightsStorage
}

func NewToolChain(approver Approver, fileRightsStorage FileRightsStorage) *ToolChain {
	return &ToolChain{
		approver:          approver,
		fileRightsStorage: fileRightsStorage,
	}
}
