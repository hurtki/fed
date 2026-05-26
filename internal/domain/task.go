package domain

type TaskTest struct {
	Description string
}

func NewTaskTest(description string) (TaskTest, error) {
	return TaskTest{}, nil
}

type Task struct {
	Description string
	Files       []ProjectFile
}

func NewTask(description string, files []ProjectFile) (Task, error) {
	return Task{
		Description: description,
		Files:       files,
	}, nil
}
