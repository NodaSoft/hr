package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type Task struct {
	ID          int
	CreatedAt   string
	ProcessedAt string
}

type CreateTaskRes struct {
	Task      Task
	CreateErr error
}

type RunTaskRes struct {
	Task Task
	Err  error
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	processRunResults(taskRunner(ctx, taskCreator(ctx)))
}

func taskCreator(ctx context.Context) <-chan CreateTaskRes {
	out := make(chan CreateTaskRes)
	create := func() CreateTaskRes {
		res := CreateTaskRes{}
		now := time.Now()
		if now.Nanosecond()%2 > 0 {
			res.CreateErr = errors.New("some error occured")
		} else {
			res.Task = Task{
				ID:        int(now.Unix()),
				CreatedAt: now.Format(time.RFC3339),
			}
		}
		return res
	}
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
			case out <- create():
			}
		}
	}()
	return out
}

func taskRunner(ctx context.Context, maybeTasks <-chan CreateTaskRes) <-chan RunTaskRes {
	out := make(chan RunTaskRes)
	run := func(maybeTask CreateTaskRes) RunTaskRes {
		if maybeTask.CreateErr != nil {
			return RunTaskRes{Err: fmt.Errorf("failed to run task: %w", maybeTask.CreateErr)}
		}
		task := maybeTask.Task
		createdAt, err := time.Parse(time.RFC3339, task.CreatedAt)
		if err != nil {
			return RunTaskRes{
				Err: fmt.Errorf(
					"failed to run task ID=%d, CreatedAt: %s, err: %w",
					task.ID,
					task.CreatedAt,
					err,
				),
			}
		}

		if !createdAt.After(time.Now().Add(-20 * time.Second)) {
			return RunTaskRes{
				Err: fmt.Errorf("failed to run task ID=%d, CreatedAt: %s", task.ID, task.CreatedAt),
			}
		}
		task.ProcessedAt = time.Now().Format(time.RFC3339)
		time.Sleep(150 * time.Millisecond)

		return RunTaskRes{Task: task}
	}
	go func() {
		defer close(out)
		for maybeTask := range maybeTasks {
			select {
			case <-ctx.Done():
				return
			case out <- run(maybeTask):
			}
		}
	}()
	return out
}

func processRunResults(runResults <-chan RunTaskRes) {
	errors := make([]error, 0)
	executedTasks := make([]Task, 0)
	for runRes := range runResults {
		if runRes.Err != nil {
			errors = append(errors, runRes.Err)
		} else {
			executedTasks = append(executedTasks, runRes.Task)
		}
	}

	fmt.Println("Errors:")
	for _, err := range errors {
		fmt.Println(err)
	}
	fmt.Println("Done tasks:")
	for _, task := range executedTasks {
		fmt.Println(task)
	}
}
