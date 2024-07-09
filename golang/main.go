package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Task struct {
	ID         int
	CreatedAt  time.Time
	FinishedAt time.Time
	Result     string
	Error      error
}

const (
	taskDuration       = 10 * time.Second
	taskCount          = 10
	readDuration       = 3 * time.Second
	taskProcessingTime = 150 * time.Millisecond
	programLifetime    = time.Minute
	workersCount       = 4
)

func createNewTask(tasksStorage chan<- Task) {
	now := time.Now().UTC()

	task := Task{
		ID:        int(now.Unix()),
		CreatedAt: now,
		Error:     nil,
	}
	if time.Now().Second() > 30 {
		task.Error = errors.New("some error occured")
	}

	tasksStorage <- task
}

func createTasks(tasksStorage chan<- Task, shutdown <-chan bool) {
	ticker := time.NewTicker(taskDuration)
	defer ticker.Stop()

	for {
		select {
		case <-shutdown:
			return
		case <-ticker.C:
			for range taskCount {
				go createNewTask(tasksStorage)
			}
		}
	}
}

type WorkerChanels struct {
	TaskStorage    <-chan Task
	Shutdown       <-chan bool
	ComplitedTasks chan<- Task
}

func worker(chanels WorkerChanels) {
	for {
		select {
		case <-chanels.Shutdown:
			return
		case task := <-chanels.TaskStorage:
			task.FinishedAt = time.Now().UTC()
			if task.Error != nil {
				task.Result = "something went wrong"
			} else {
				task.Result = "task has been successed"
			}

			chanels.ComplitedTasks <- task

			time.Sleep(taskProcessingTime)
		}
	}
}

func runWorkers(chanels WorkerChanels) {
	for range workersCount {
		go worker(chanels)
	}
}

func printTasks(complitedTasks <-chan Task, shutdown <-chan bool) {
	ticker := time.NewTicker(readDuration)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			var errTasks []Task
			var successTasks []Task

		Loop:
			for task := range complitedTasks {
				if task.Error != nil {
					errTasks = append(errTasks, task)
				} else {
					successTasks = append(successTasks, task)
				}

				if len(complitedTasks) == 0 {
					break Loop
				}
			}

			fmt.Println("Error tasks: ")
			for _, task := range errTasks {
				fmt.Printf("%+v\n", task)
			}

			fmt.Println("Success tasks: ")
			for _, task := range successTasks {
				fmt.Printf("%+v\n", task)
			}
		case <-shutdown:
			return
		}
	}
}

func gracefulShutdown(
	shutdownCreateTasks chan bool,
	shutdownWorkers chan bool,
	shutdownPrintTasks chan bool,
	taskStorage chan Task,
	complitedTasks chan Task,
) {
	shutdownCreateTasks <- true
	for len(taskStorage) != 0 {
		time.Sleep(time.Second)
	}

	for range workersCount {
		shutdownWorkers <- true
	}

	for len(complitedTasks) != 0 {
		time.Sleep(time.Second)
	}
	shutdownPrintTasks <- true
}

func main() {
	taskStorage := make(chan Task, 100)
	complitedTasks := make(chan Task, 100)
	shutdownCreateTasks := make(chan bool)
	shutdownPrintTasks := make(chan bool)
	shutdownWorkers := make(chan bool, workersCount)
	defer close(shutdownWorkers)
	defer close(shutdownCreateTasks)
	defer close(shutdownPrintTasks)
	defer close(complitedTasks)
	defer close(taskStorage)

	go createTasks(taskStorage, shutdownCreateTasks)
	runWorkers(
		WorkerChanels{
			TaskStorage:    taskStorage,
			Shutdown:       shutdownWorkers,
			ComplitedTasks: complitedTasks,
		},
	)
	go printTasks(complitedTasks, shutdownPrintTasks)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case stopSignal := <-stop:
			fmt.Printf("Stopping program with signal: %s\n", stopSignal.String())
			gracefulShutdown(shutdownCreateTasks, shutdownWorkers, shutdownPrintTasks, taskStorage, complitedTasks)
			return
		case <-time.After(programLifetime):
			fmt.Println("Program lifetime exceeded, stopping...")
			gracefulShutdown(shutdownCreateTasks, shutdownWorkers, shutdownPrintTasks, taskStorage, complitedTasks)
			return
		}
	}
}
