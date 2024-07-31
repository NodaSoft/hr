package main

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"
)

// TaskReporter separately reports error tasks and done tasks into a provided io.Writer output.
type TaskReporter struct {
	out               io.Writer
	okTasks, errTasks []*Task
}

// NewTaskReporter returns new TaskReporter.
func NewTaskReporter(out io.Writer) *TaskReporter {
	return &TaskReporter{
		out: out,
	}
}

// Record records a task. If the task result is successful,
// it records it into a group of done tasks.
// Otherwise, task is recorded as error-task.
func (r *TaskReporter) Record(task *Task) {
	if task.Result().IsSuccessful() {
		r.okTasks = append(r.okTasks, task)
	} else {
		r.errTasks = append(r.errTasks, task)
	}
}

// Report reports all the recorded tasks into the io.Writer output. It panics if io.Writer returns an error.
func (r *TaskReporter) Report() {
	var sb strings.Builder
	sb.WriteString("Errors:\n")
	for _, t := range r.errTasks {
		sb.WriteString(
			fmt.Sprintf("Task id=%d time=%s error=%s\n",
				t.ID(),
				t.CreationTime().Format(time.RFC3339),
				t.Result().Message(", ")))
	}

	sb.WriteString("Done tasks:\n")
	for _, t := range r.okTasks {
		sb.WriteString(fmt.Sprintf("%+v\n", t))
	}

	_, err := r.out.Write([]byte(sb.String()))
	if err != nil {
		panic(err)
	}
}

// StartReporting starts an asynchronous operation of
// Report function invocations once in a provided interval.
func (r *TaskReporter) StartReporting(ctx context.Context, interval time.Duration) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				break
			default:
				time.Sleep(interval)
			}

			r.Report()
		}
	}()
}

// NewRecorderPipe is a Pipe that records all the input tasks into a reporter TaskReporter.
func NewRecorderPipe(reporter *TaskReporter) Pipe[*Task] {
	return func(ctx context.Context, in <-chan *Task) <-chan *Task {
		out := make(chan *Task, cap(in))

		go func() {
			defer close(out)

			for t := range in {
				if err := ctx.Err(); err != nil {
					break
				}

				reporter.Record(t)

				out <- t
			}
		}()

		return out
	}
}
