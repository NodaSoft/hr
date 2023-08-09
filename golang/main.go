package main

import (
	"context"
	"main/config"
	"main/internal/task"
	st "main/internal/task/structs"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		panic("Could not load config: " + err.Error())
	}

	taskHandler := task.NewTaskHandler(cfg.Task)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.WorkTime)
	defer cancel()

	tasksResult := make(chan st.TasksResult)
	taskHandler.Run(ctx, tasksResult)

	(<-tasksResult).PrintResult()
}
