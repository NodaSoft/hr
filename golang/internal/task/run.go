package task

import (
	"context"
	"main/config"
	"main/internal/task/steps"
	st "main/internal/task/structs"
)

type TaskHandler struct {
	cfg config.TaskConfig
}

func NewTaskHandler(cfg config.TaskConfig) *TaskHandler {
	return &TaskHandler{cfg: cfg}
}

func (th *TaskHandler) Run(ctx context.Context, resultCh chan<- st.TasksResult) {
	tasksToProcess := make(chan st.Task, 10)
	tasksForSorting := make(chan st.Task)
	doneTasks := make(chan st.Task)
	undoneTasks := make(chan error)

	go steps.RunCreation(ctx, tasksToProcess)
	go steps.RunProcessing(ctx, tasksToProcess, tasksForSorting, th.cfg.TaskRelevanceTimDuration)
	go steps.RunSort(ctx, tasksForSorting, doneTasks, undoneTasks)
	go steps.CollectResult(doneTasks, undoneTasks, resultCh)
}
