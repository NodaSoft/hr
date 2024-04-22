package main

import (
	"fmt"
	"time"
)

const (
	DateTimeFormat   = time.RFC3339
	DateTimeNano     = time.RFC3339Nano
	TaskResultErrMsg = "Some error occurred"
	TaskResultOK     = "Task has been created"
)

// Task represents a task
type Task struct {
	ID          int
	CreatedTime string // creation time
	FinishTime  string // finish time
	Result      string // execution result
}

func main() {
	taskCreator := func(tasks chan<- Task) {
		for {
			createdTime := time.Now().Format(DateTimeFormat)
			result := TaskResultOK
			if time.Now().Nanosecond()%2 > 0 {
				result = TaskResultErrMsg
			}
			task := Task{
				CreatedTime: createdTime,
				ID:          int(time.Now().Unix()),
				Result:      result,
			}
			tasks <- task
		}
	}

	taskChannel := make(chan Task, 10)

	go taskCreator(taskChannel)

	taskWorker := func(task Task) Task {
		task.FinishTime = time.Now().Format(DateTimeNano)
		time.Sleep(time.Millisecond * 150)

		switch task.Result {
		case TaskResultOK:
			task.Result = "Task has been succeeded"
		default:
			task.Result = "Something went wrong"
		}

		return task
	}

	doneTasks := make(chan Task)
	undoneTasks := make(chan Task)

	taskSorter := func(task Task) {
		if task.Result == "Task has been succeeded" {
			doneTasks <- task
		} else {
			undoneTasks <- task
		}
	}

	go func() {
		for task := range taskChannel {
			task := taskWorker(task)
			go taskSorter(task)
		}
		close(taskChannel)
	}()

	done := make(chan struct{})
	go func() {
		for {
			select {
			case task := <-doneTasks:
				fmt.Printf("Task ID: %d, Finish Time: %s, Result: %s\n", task.ID, task.FinishTime, task.Result)
			case task := <-undoneTasks:
				fmt.Printf("Task ID: %d, Finish Time: %s, Result: %s\n", task.ID, task.FinishTime, task.Result)
			case <-time.After(time.Second * 3):
				close(doneTasks)
				close(undoneTasks)
				close(done)
				return
			}
		}
	}()

	<-done
}
