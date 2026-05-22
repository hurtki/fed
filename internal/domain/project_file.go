package domain

import (
	"fmt"
	"path/filepath"
	"strings"
)

type ProjectFile struct {
	Path    string
	Project *Project
}

func NewProjectFile(filePath string, project *Project) (ProjectFile, error) {
	if !strings.HasPrefix(filePath, "./") {
		return ProjectFile{}, fmt.Errorf("not project file")
	}

	return ProjectFile{
		Path:    filePath,
		Project: project,
	}, nil
}

func (pf ProjectFile) GetAbsPath() string {
	return filepath.Join(pf.Project.BasePath, pf.Path)
}
