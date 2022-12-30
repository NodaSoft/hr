package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"runtime"
	"sync"
	"time"
)

type Result struct {
	success bool
	error   error
}

type Task struct {
	id         string
	createdTs  string
	executedTs string
	result     Result
}

type doneTasks map[string]Task

type WorkResults struct {
	tasks  doneTasks
	errors []error
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var results []<-chan Task
	for i := 0; i < runtime.NumCPU(); i++ {
		results = append(results, taskWorker(taskCreator(ctx)))
	}
	wr := WorkResults{
		tasks: map[string]Task{},
	}
	for task := range merge(results...) {
		if task.result.success {
			wr.tasks[task.id] = task
		} else {
			wr.errors = append(wr.errors, fmt.Errorf("task id %s time %s, error %s", task.id, task.createdTs, task.result.error))
		}
	}

	println("Errors:")
	for _, err := range wr.errors {
		println(err.Error())
	}

	println("Done tasks:")
	for id, _ := range wr.tasks {
		println(id)
	}
}

func taskCreator(ctx context.Context) <-chan Task {
	out := make(chan Task)
	go func() {
		defer close(out)
		for {
			createdTs := time.Now().Format(time.RFC3339)
			if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
				// тут не вполне ясно, если это условность для эмуляции битых данных то ок, а вообще ошибку в таймштампе передавать это дичь
				createdTs = "Some error occurred"
			}
			select {
			case out <- Task{createdTs: createdTs, id: uuid.New().String()}:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

func taskWorker(in <-chan Task) <-chan Task {
	out := make(chan Task)
	go func() {
		for task := range in {
			createdTs, _ := time.Parse(time.RFC3339, task.createdTs)
			if createdTs.After(time.Now().Add(-20 * time.Second)) {
				task.result = Result{
					success: true,
				}
			} else {
				task.result = Result{
					success: false,
					error:   errors.New("something went wrong"),
				}
			}
			task.executedTs = time.Now().Format(time.RFC3339Nano)
			time.Sleep(time.Millisecond * 150)
			out <- task
		}
		close(out)
	}()
	return out
}

func merge(cs ...<-chan Task) <-chan Task {
	var wg sync.WaitGroup
	out := make(chan Task)

	output := func(c <-chan Task) {
		defer wg.Done()
		for n := range c {
			out <- n
		}
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
