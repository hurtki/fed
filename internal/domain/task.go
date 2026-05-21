package domain

type TaskTest struct {
	Description string
}

func NewTaskTest(description string) (TaskTest, error) {
	return TaskTest{}, nil
}

type Task struct {
	Description string
	Solved      bool
	Files       []ProjectFile
	Tests       []TaskTest
}

func NewTask(description string, files []ProjectFile, tests []TaskTest) (Task, error) {
	return Task{
		Description: description,
		Solved:      false,
		Files:       files,
		Tests:       tests,
	}, nil
}
