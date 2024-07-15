package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"sync"
	"time"
)

// Task contains information about a one task
type Task struct {
	Uuid       uuid.UUID
	CreateTime time.Time
	FinishTime time.Time
	Result     []byte
	Err        error
}

// IncorrectTaskTimeDivision is used to generate tasks with an error.
//
// If, when dividing the number of nanoseconds from the current time by this number, the remainder is greater than zero, then the task will be done with an error.
var IncorrectTaskTimeDivision int = 2

// CreateTaskDuration is the frequency with which tasks will be created.
var CreateTaskDuration time.Duration = time.Millisecond * 100

// ProcessTaskDuration is the duration of the pause during processing. It is used only for the correct task.
var ProcessTaskDuration time.Duration = time.Millisecond * 150

// ErrIncorrectTaskCreateTime is an error that will be returned if the task is created with an incorrect time.
var ErrIncorrectTaskCreateTime = errors.New("incorrect task create time")

// GenerateOneTask creates a new task.
//
// Use IncorrectTaskTimeDivision for simulating generate incorrect tasks.
// If task is incorrect - ErrIncorrectTaskCreateTime stores in Err field.
func GenerateOneTask(logger *zap.Logger) Task {
	logger = logger.Named("GenerateOneTask")

	currentTime := time.Now()
	task := Task{
		Uuid:       GenerateUuid(),
		CreateTime: currentTime,
	}

	if time.Now().Nanosecond()%IncorrectTaskTimeDivision > 0 {
		task.Err = ErrIncorrectTaskCreateTime
	}

	logger.Debug("Create a new one task", zap.Reflect("Task", task))

	return task
}

// FillTaskChannel filling provided channel with tasks.
func FillTaskChannel(
	tasksChannel chan<- Task,
	ctx context.Context,
	wg *sync.WaitGroup,
	logger *zap.Logger,
) {
	defer wg.Done()

	startTime := time.Now()
	logger = logger.Named("FillTaskChannel")

	defer func() {
		duration := zap.Duration("Duration", time.Since(startTime))

		if r := recover(); r != nil {
			logger.Error("Recovered!", zap.Any("Cause", r), duration)
		} else {
			logger.Info("Finish filling task channel", duration)
		}
	}()

	ticker := time.NewTicker(CreateTaskDuration) // ticker for simulating events for task creation

	for {
		select {
		case <-ctx.Done():
			logger.Info("Context is closed!")
			return
		case <-ticker.C:
			select {
			case <-ctx.Done():
				logger.Info("Context is closed while trying to insert a new one task!")
			case tasksChannel <- GenerateOneTask(logger):
				logger.Debug("Insert a one task")
			}
		}
	}
}

// ProcessOneTask processing task.
//
// Use ProcessTaskDuration for simulate processing but only if task is correct.
func ProcessOneTask(task Task, logger *zap.Logger) Task {
	logger = logger.Named("ProcessOneTask")

	if task.Err != nil {
		logger.Debug("Task is incorrect - can't be processed", zap.Error(task.Err)) // debug because there will be a lot of error messages
		task.Result = []byte(fmt.Sprintf("Task is incorrect - can't be processed: %s", task.Err))
	} else {
		time.Sleep(ProcessTaskDuration) // simulation of processing
		task.Result = []byte("OK")
		logger.Debug("Process a one task", zap.Reflect("Task", task))
	}

	task.FinishTime = time.Now()

	return task
}

// ProcessTaskChannel processing tasks channel.
func ProcessTaskChannel(
	unprocessedTasksChannel <-chan Task,
	processedTasksChannel chan<- Task,
	ctx context.Context,
	wg *sync.WaitGroup,
	logger *zap.Logger,
) {
	defer wg.Done()

	startTime := time.Now()
	logger = logger.Named("ProcessTaskChannel")

	defer func() {
		duration := zap.Duration("Duration", time.Since(startTime))

		if r := recover(); r != nil {
			logger.Error("Recovered!", zap.Any("Cause", r), duration)
		} else {
			logger.Info("Finish processing tasks channel", duration)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			logger.Info("Context is closed")
			return
		case task, isNotClosed := <-unprocessedTasksChannel:
			if !isNotClosed {
				logger.Info("Channel is closed while trying to read an unprocessed task!")
				return
			}

			select {
			case processedTasksChannel <- ProcessOneTask(task, logger):
				logger.Debug("Process a task", zap.String("Task uuid", task.Uuid.String()))
			case <-ctx.Done():
				logger.Info("Context is closed while adding a new one task!")
				return
			}
		}
	}
}

// HandleProcessedTasksChannel output results of processing tasks.
//
// Return processed and unprocessed tasks.
func HandleProcessedTasksChannel(
	processedTasksChannel <-chan Task,
	ctx context.Context,
	detailed bool,
	printPeriod time.Duration,
	wg *sync.WaitGroup,
	logger *zap.Logger,

) ([]Task, []Task) {
	defer wg.Done()

	processedTasks := make([]Task, 0)
	unprocessedTasks := make([]Task, 0)

	ticker := time.NewTicker(printPeriod)

	for {
		select {
		case <-ctx.Done():
			logger.Info("Context is closed")
			return processedTasks, unprocessedTasks
		case task, isNotClosed := <-processedTasksChannel:
			if !isNotClosed {
				logger.Info("Channel is closed while trying to read a processed task!")
				return processedTasks, unprocessedTasks
			}

			if task.Err != nil {
				unprocessedTasks = append(unprocessedTasks, task)
			} else {
				processedTasks = append(processedTasks, task)
			}
		case <-ticker.C:
			PrintTasksResult(processedTasks, unprocessedTasks, detailed)
		}
	}
}

// PrintTasksResult print tasks result.
//
// If detailed is true, besides count of processed and unprocessed tasks also print all tasks with all their fields.
func PrintTasksResult(processedTasks []Task, unprocessedTasks []Task, detailed bool) {
	fmt.Printf("Tasks result on %s:\n", time.Now().Format(time.RFC3339Nano))
	fmt.Printf("\tProcessed tasks count: %d\n", len(processedTasks))
	fmt.Printf("\tUnprocessed tasks count: %d\n", len(unprocessedTasks))

	if detailed {
		if len(processedTasks) > 0 {
			fmt.Println("\n\tAll processed tasks:")
			PrintTasksSlice(processedTasks)
		} else {
			fmt.Println("\n\tNo processed tasks")
		}

		if len(unprocessedTasks) > 0 {
			fmt.Println("\n\tAll unprocessed tasks:")
			PrintTasksSlice(unprocessedTasks)
		} else {
			fmt.Println("\n\tNo unprocessed tasks")
		}

		fmt.Println() // just add new line for separating next output
	}
}

// PrintTasksSlice print tasks slice.
// Each line is a single task in `+v` format.
func PrintTasksSlice(tasks []Task) {
	for _, task := range tasks {
		fmt.Printf("\t%+v\n", task)
	}
}
