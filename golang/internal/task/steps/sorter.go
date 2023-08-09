package steps

import (
	"context"
	"fmt"
	st "main/internal/task/structs"
	"strings"
)

func RunSort(ctx context.Context, taskInput <-chan st.Task, doneTasks chan<- st.Task, undoneTasks chan<- error) {
	for {
		select {
		case <-ctx.Done():
			close(doneTasks)
			close(undoneTasks)
			return
		case task, ok := <-taskInput:
			if !ok {
				return
			}
			sort(task, doneTasks, undoneTasks)
		}
	}
}

func sort(task st.Task, doneTasks chan<- st.Task, undoneTasks chan<- error) {
	if strings.Contains(string(task.Result), "completed") {
		doneTasks <- task
	} else {
		undoneTasks <- fmt.Errorf("task id %d time %s, error %s", task.Id, task.CreateTime, task.Result)
	}
}
