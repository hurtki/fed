package domain

type Project struct {
	BasePath string
}

func NewProject(basePath string) (Project, error) {
	return Project{
		BasePath: basePath,
	}, nil
}
