package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

// Приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме.
// В конце приложение должно выводить успешные таски и ошибки выполнения остальных тасков.

const TaskProcessTimeout = time.Millisecond * 150

type Task struct {
	ID         ulid.ULID
	CreatedAt  time.Time
	FinishedAt time.Time
	Msg        string
	Error      error
}

type DoneTasks struct {
	sync.Mutex
	Tasks *[]Task
}

func (m *DoneTasks) Add(t Task) {
	m.Lock()
	defer m.Unlock()

	*m.Tasks = append(*m.Tasks, t)
}

func (m *DoneTasks) Drain() []Task {
	m.Lock()
	defer m.Unlock()

	stored := m.Tasks
	m.Tasks = &[]Task{}
	return *stored
}

type UndoneTasks struct {
	sync.Mutex
	Tasks *[]Task
}

func (m *UndoneTasks) Add(t Task) {
	m.Lock()
	defer m.Unlock()

	*m.Tasks = append(*m.Tasks, t)
}

func (m *UndoneTasks) Drain() []Task {
	m.Lock()
	defer m.Unlock()

	stored := m.Tasks
	m.Tasks = &[]Task{}
	return *stored
}

func NewTask() Task {
	task := Task{ID: ulid.Make(), CreatedAt: time.Now()}

	if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
		task.Error = errors.New("some error occurred")
	}

	task.FinishedAt = time.Now()

	return task
}

func checkCreatedTime(task Task) Task {
	if task.CreatedAt.After(time.Now().Add(-20 * time.Second)) {
		task.Msg = "task completed successfully"
	} else {
		task.Error = errors.New("something went wrong")
	}

	return task
}

func ProcessTask(task Task, f func(Task) Task, timeout time.Duration) Task {
	task = f(task)

	task.FinishedAt = time.Now()

	time.Sleep(timeout)

	return task
}

func filterTask(
	ctx context.Context,
	task Task,
	completed *DoneTasks,
	failed *UndoneTasks,
) {
	select {
	case <-ctx.Done():
		return
	default:
		if task.Error == nil {
			completed.Add(task)
		} else {
			failed.Add(task)
		}
	}
}

func generateTasks(ctx context.Context, taskCh chan<- Task) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			task := NewTask()
			taskCh <- task // передаем таск на выполнение
		}
	}
}

func retrieveTasks(
	ctx context.Context,
	taskCh <-chan Task,
	completed *DoneTasks,
	failed *UndoneTasks,
) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// получение тасков
			for task := range taskCh {
				task = ProcessTask(task, checkCreatedTime, TaskProcessTimeout)
				go filterTask(ctx, task, completed, failed)
			}
		}
	}
}

func main() {
	taskPool := make(chan Task, 10)
	ctx, ctxCancel := context.WithCancel(context.Background())

	go generateTasks(ctx, taskPool)

	done, undone := DoneTasks{Tasks: &[]Task{}}, UndoneTasks{Tasks: &[]Task{}}

	go retrieveTasks(ctx, taskPool, &done, &undone)

	time.Sleep(time.Second * 3)
	ctxCancel()

	completed, failed := done.Drain(), undone.Drain()

	println("Errors:")
	for _, fail := range failed {
		fmt.Printf("Task ID: %s, Error: %s\n", fail.ID.String(), fail.Error)
	}

	println("Done tasks:")
	for _, success := range completed {
		fmt.Printf(
			"Task ID: %s, Created: %s, \nFinished: %s, Msg: '%s'\n\n",
			success.ID.String(),
			success.CreatedAt,
			success.FinishedAt,
			success.Msg,
		)
	}
}
