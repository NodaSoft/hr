package steps

import (
	st "main/internal/task/structs"
)

func CollectResult(doneTasks <-chan st.Task, undoneTasks <-chan error, resultCh chan<- st.TasksResult) {
	result := st.TasksResult{
		DoneTasks:   make(map[int]st.Task, 0),
		UndoneTasks: make([]error, 0),
	}

	done := false
	undone := false

	for !done || !undone {
		select {
		case task, ok := <-doneTasks:
			if !ok {
				done = true
				continue
			}
			result.DoneTasks[task.Id] = task

		case err, ok := <-undoneTasks:
			if !ok {
				undone = true
				continue
			}
			result.UndoneTasks = append(result.UndoneTasks, err)
		}
	}

	resultCh <- result
}
