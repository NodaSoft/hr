package app

import (
	"context"
	"os"
	"testtask/internal/config"
	"testtask/internal/task"
)

const (
	doneTasksBufferSize = 1000
)

func Run(cfg config.Config) {
	taskProducer := task.NewTaskProducer(cfg.TaskExpirationDuration)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		cfg.TaskGenerationDuration,
	)
	defer cancel()
	toDoTasks := taskProducer.ProduceTasks(ctx)

	taskExecutor := task.NewTaskExecutor(
		cfg.TaskExecutorsLimit,
		doneTasksBufferSize,
	)
	doneTasks := taskExecutor.ExecuteTasks(toDoTasks)

	task.
		NewTaskReporter(
			doneTasks,
			os.Stdout,
			cfg.TaskResultReportingPeriod,
		).
		ReportTaskResults()
}
