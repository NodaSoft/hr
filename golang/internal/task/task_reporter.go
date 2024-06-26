package task

import (
	"fmt"
	"io"
	"sync"
	"time"
)

type TaskReporter struct {
	doneTasks                 <-chan Task
	output                    io.Writer
	taskResultReportingPeriod time.Duration
}

func NewTaskReporter(
	doneTasks <-chan Task,
	output io.Writer,
	taskResultReportingPeriod time.Duration) *TaskReporter {
	return &TaskReporter{
		doneTasks:                 doneTasks,
		output:                    output,
		taskResultReportingPeriod: taskResultReportingPeriod,
	}
}

func (tr *TaskReporter) ReportTaskResults() {
	successfulTasks := []Task{}
	failedTasks := []Task{}
	var successfulTasksMutex sync.Mutex
	var failedTasksMutex sync.Mutex

	doneReading := make(chan struct{})
	go func() {
		for task := range tr.doneTasks {
			if task.Error == nil {
				successfulTasksMutex.Lock()
				successfulTasks = append(successfulTasks, task)
				successfulTasksMutex.Unlock()
			} else {
				failedTasksMutex.Lock()
				failedTasks = append(failedTasks, task)
				failedTasksMutex.Unlock()
			}
		}
		doneReading <- struct{}{}
	}()

	printReport := func() {
		fmt.Fprintln(tr.output, "Successful tasks:")
		successfulTasksMutex.Lock()
		for _, task := range successfulTasks {
			fmt.Fprintln(tr.output, task)
		}
		successfulTasksMutex.Unlock()

		fmt.Fprintln(tr.output, "Failed tasks:")
		failedTasksMutex.Lock()
		for _, err := range failedTasks {
			fmt.Fprintln(tr.output, err)
		}
		failedTasksMutex.Unlock()

		fmt.Fprintln(tr.output, "")
	}

	timeToReportTicker := time.NewTicker(tr.taskResultReportingPeriod)
	defer timeToReportTicker.Stop()
	for {
		select {
		case <-timeToReportTicker.C:
			printReport()
		case <-doneReading:
			printReport()
			return
		}
	}
}
