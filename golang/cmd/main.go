package main

import (
	"hr/internal/task"
)

func main() {
	tasks := make(chan task.Task, task.MaxExecutionCount)
	completedTasks := make(chan task.Task, task.MaxExecutionCount)

	taskBuilder := task.Builder{Tasks: tasks}
	taskWorker := task.Worker{NewTasks: tasks, CompletedTasks: completedTasks}
	taskLogger := task.Logger{Tasks: completedTasks}

	go taskBuilder.Start()
	go taskWorker.Start()

	taskLogger.Start()
}
