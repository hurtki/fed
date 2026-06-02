package domain

type Step struct {
	Description string
	Files       []ProjectFile
}

func NewStep(description string, files []ProjectFile) (Step, error) {
	return Step{
		Description: description,
		Files:       files,
	}, nil
}
