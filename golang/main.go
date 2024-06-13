package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

// A Ttype represents a meaninglessness of our life
type Ttype struct {
	id         int
	cT         string // время создания
	fT         string // время выполнения
	taskRESULT []byte
}

type Task struct {
	id      int
	created time.Time // время создания
	handled time.Time // время выполнения
	result  []byte
	err     error
}

var ErrOnCreateTask = errors.New("err on create task")
var ErrOnHandleTask = errors.New("err on handle task")

func NewTaskCreator() (<-chan Task, func(ctx context.Context)) {
	tasksChan := make(chan Task)
	started := atomic.Bool{}
	started.Store(false)
	return tasksChan, func(ctx context.Context) {
		if started.Swap(true) || ctx.Err() != nil {
			log.Printf("[ERROR] attempt to start creator again")
			return
		}
		for {
			t := Task{
				id:      int(time.Now().Unix()),
				created: time.Now(),
			}
			if time.Now().Nanosecond()%2 > 0 { // фейлятся таски в нечетные наносекунды
				t.err = ErrOnCreateTask
			}
			select {
			case <-ctx.Done():
				return
			case tasksChan <- t:
			}
		}
	}
}

type TaskHandler struct {
	tasksCh   <-chan Task
	doneTasks map[int]Task
	errors    []error
	mu        sync.Mutex
}

func NewTaskHandler(tasksCh <-chan Task) *TaskHandler {
	th := &TaskHandler{
		tasksCh:   tasksCh,
		doneTasks: make(map[int]Task),
	}
	go th.handleTasks()
	return th
}

func handleTask(task Task) Task {
	if task.created.After(time.Now().Add(-20 * time.Second)) {
		task.result = []byte("task has been successed")
	} else {
		task.result = []byte("something went wrong")
		task.err = ErrOnHandleTask
	}
	task.handled = time.Now()
	time.Sleep(time.Millisecond * 150)
	return task
}

func (h *TaskHandler) handleTasks() {
	for t := range h.tasksCh {
		t = handleTask(t)
		func() {
			h.mu.Lock()
			defer h.mu.Unlock()
			if t.err != nil {
				err := fmt.Errorf("handle task id=%d created=%s: %w", t.id, t.created.Format(time.RFC3339), t.err)
				h.errors = append(h.errors, err)
			} else {
				h.doneTasks[t.id] = t
			}
		}()
	}
}

type TaskHandleResult struct {
	DoneTasks map[int]Task
	Errors    []error
}

func (h *TaskHandler) Peek() TaskHandleResult {
	h.mu.Lock()
	defer h.mu.Unlock()
	r := TaskHandleResult{
		DoneTasks: make(map[int]Task, len(h.doneTasks)),
		Errors:    make([]error, len(h.errors)),
	}
	copy(r.Errors, h.errors)
	for k, v := range h.doneTasks {
		r.DoneTasks[k] = v
	}
	return r
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	tasksChan, runCreator := NewTaskCreator()
	taskHandler := NewTaskHandler(tasksChan)
	go runCreator(ctx)
	time.Sleep(time.Second * 3)
	cancel()
	result := taskHandler.Peek()
	fmt.Println("Errors:")
	for _, err := range result.Errors {
		fmt.Printf("%v\n", err)
	}
	fmt.Println("Done tasksChan:")
	for id, task := range result.DoneTasks {
		fmt.Printf("id=%v, created=%v\n", id, task.created.Format(time.RFC3339))
	}
}
