package sorter

import (
	"fmt"
	"taskConcurrency/internal/domain/task"
)

type Sorter struct{}

func (s *Sorter) Sort(tasks <-chan task.Task,
	doneTasks chan<- task.Task, undoneTasks chan<- error) {
	go func() {
		for task := range tasks {
			go s.sortTask(task, doneTasks, undoneTasks)
		}
		close(doneTasks)
		close(undoneTasks)
	}()
}

func (s *Sorter) sortTask(task task.Task,
	doneTasks chan<- task.Task, undoneTasks chan<- error) {
	if len(task.TaskResult) == 0 {
		undoneTasks <- fmt.Errorf("task id: %d, error: %s", task.Id, task.TaskResult)
		return
	}
	if string(task.TaskResult[14:]) == "successed" {
		doneTasks <- task
	} else {
		undoneTasks <- fmt.Errorf("task id: %d, error: %s", task.Id, task.TaskResult)
	}
}
