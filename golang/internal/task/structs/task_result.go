package structs

type TasksResult struct {
	DoneTasks   map[int]Task
	UndoneTasks []error
}

func (t TasksResult) PrintResult() {
	println("Errors:")
	for _, err := range t.UndoneTasks {
		println(err.Error())
	}

	println("Done tasks:")
	for id := range t.DoneTasks {
		println(id)
	}
}
