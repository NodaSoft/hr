package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"task_service/internal/dispatcher"
	"task_service/internal/domain"
	"task_service/internal/generator"
	"task_service/internal/receivers"
	"task_service/internal/repository"
	"task_service/internal/workers"
	"time"
)

const superChSize = 6
const taskChSize = 3
const failedTaskChSize = 3

const timeOut = time.Second * 5

func main() {
	ctx := context.Background()

	superCh := make(chan domain.Task, superChSize)
	taskCh := make(chan domain.Task, taskChSize)
	failedTaskCh := make(chan domain.Task, failedTaskChSize)

	generator := generator.New(superCh)

	workerManager := workers.NewWorkerManager(superCh, taskCh, failedTaskCh)
	worker := workers.NewSimpleWorker()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	taskRepo := repository.NewTaskRepository()
	errorRepo := repository.NewErrorRepository()

	recevier := receivers.New(taskCh, failedTaskCh, taskRepo, errorRepo)

	ctx, _ = context.WithTimeout(ctx, timeOut)
	dispatcher := dispatcher.NewDispatcher(taskCh, failedTaskCh, shutdown)

	go recevier.Run()
	go workerManager.Run(worker)
	go generator.Run()
	dispatcher.Dispatch(ctx)

	doneTasks := taskRepo.List()
	undoneTasks := errorRepo.List()

	fmt.Println("Done tasks:")
	for id := range doneTasks {
		fmt.Println(id)
	}

	fmt.Println("Errors:")
	for id := range undoneTasks {
		fmt.Println(id)
	}
}
