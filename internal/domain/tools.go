package domain

type FileChange struct {
	File    ProjectFile
	Find    string
	Replace string
}
