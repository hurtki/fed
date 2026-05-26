package domain

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Project struct {
	BasePath string
}

func NewProject(basePath string) (Project, error) {
	return Project{BasePath: basePath}, nil
}

func (p Project) NewFile(relPath string) (ProjectFile, error) {
	relPath = filepath.Clean(relPath)
	if !strings.HasPrefix(relPath, "./") {
		return ProjectFile{}, fmt.Errorf("%s is not a relative path", relPath)
	}

	return ProjectFile{
		RelativePath: relPath,
		ProjectRoot:  p.BasePath,
	}, nil
}

type ProjectFile struct {
	RelativePath string
	ProjectRoot  string
}

func (pf ProjectFile) NewChange(find, replace string) FileChange {
	return FileChange{
		File:    pf,
		Find:    find,
		Replace: replace,
	}
}

func (pf ProjectFile) GetAbsPath() string {
	return filepath.Join(pf.ProjectRoot, pf.RelativePath)
}
