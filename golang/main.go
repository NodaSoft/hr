package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type Task struct {
	ID        int           `json:"id"`
	CreatedAt time.Time     `json:"-"`
	RunTime   time.Duration `json:"-"`
	Result    string        `json:"result,omitempty"`
	ErrMsg    string        `json:"error,omitempty"`
}

func New() *Task {
	return &Task{
		// INFO: there is collisions without external lib
		ID:        time.Now().Nanosecond(),
		CreatedAt: time.Now(),
	}
}

func (t *Task) recordRunTime() func() {
	return func() {
		t.RunTime = time.Since(t.CreatedAt)
	}
}

func (t *Task) Do() *Task {
	defer t.recordRunTime()()

	// INFO: error condition
	if time.Now().Nanosecond()%2 > 0 {
		t.ErrMsg = "something went wrong"
	} else {
		t.Result = "Success"
	}

	return t
}

func producer(ctx context.Context) <-chan *Task {
	tasks := make(chan *Task)

	go func() {
		defer close(tasks)

		for {
			go func() {
				select {
				case <-ctx.Done():
					return
				default:
					tasks <- New()
				}
			}()
		}
	}()

	return tasks
}

func worker(tasks <-chan *Task) <-chan *Task {
	finishedTasks := make(chan *Task)

	go func() {
		defer close(finishedTasks)

		for task := range tasks {
			go func(task *Task) {
				finishedTasks <- task.Do()
			}(task)
		}
	}()

	return finishedTasks
}

func parse(successful, failed []*Task) {
	s := struct {
		Successful any
		Failed     any
	}{
		Successful: successful,
		Failed:     failed,
	}

	bytes, _ := json.Marshal(&s)

	fmt.Println(string(bytes))
}

func show(ctx context.Context, tasks <-chan *Task) {
	successful, failed := make([]*Task, 0), make([]*Task, 0)
	ticker := time.NewTicker(3 * time.Second)

	for {
		select {
		case <-ctx.Done():
			parse(successful, failed)
			return
		case <-ticker.C:
			parse(successful, failed)
			successful, failed = nil, nil
		case task := <-tasks:
			if task.ErrMsg != "" {
				failed = append(failed, task)
				continue
			}

			successful = append(successful, task)
		}
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	show(ctx, worker(producer(ctx)))
}
