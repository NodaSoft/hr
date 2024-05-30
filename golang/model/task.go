package model

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type Task struct {
	id         int       // идентификатор задачи
	craetedAt  time.Time // время создания задачи
	finishedAt time.Time // время выполнения задачи
	result     []byte    // результат обработки задачи
	err        error     // Выполнена ли задача
}

func (t Task) New() Task {
	return Task{id: int(time.Now().Unix()), craetedAt: time.Now()}
}

func (t *Task) SetFinishedAt(finishedAt time.Time) Task {
	t.finishedAt = finishedAt
	return *t
}

func (t *Task) SetResult(result []byte) Task {
	t.result = result
	return *t
}

func (t *Task) SetErr(err error) Task {
	t.err = err
	return *t
}

func (t *Task) GetId() int {
	return t.id
}

func (t *Task) GetCraetedAt() time.Time {
	return t.craetedAt
}

func (t *Task) GetErr() error {
	return t.err
}

func (t *Task) MetDeadline(delay time.Duration) bool {
	return t.finishedAt.Before(t.finishedAt.Add(-20 * time.Second))
}

func (t Task) Create(ctx context.Context, taskCh chan<- Task) {
	defer close(taskCh)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			t = t.New()
			if t.GetCraetedAt().Nanosecond()%2 > 0 {
				t.SetErr(errors.New("some error occured"))
			}
			taskCh <- t
			time.Sleep(time.Millisecond * 150) // Ограничиваем число созданных задач
		}
	}
}

func (t *Task) Process() Task {
	if t.GetErr() != nil {
		return *t
	}

	t.SetFinishedAt(time.Now())
	if t.MetDeadline(-20 * time.Second) {
		t.SetErr(errors.New("something went wrong"))
		return *t
	}

	t.SetResult([]byte("task has been successful"))

	return *t
}

func (t Task) Route(taskCh <-chan Task, doneCh chan<- Task, errCh chan<- error) {
	defer close(doneCh)
	defer close(errCh)

	for t := range taskCh {
		t = t.Process()
		Sort(doneCh, errCh, t)
	}
}

func Sort(doneCh chan<- Task, errCh chan<- error, t Task) {
	if t.GetErr() != nil {
		errCh <- fmt.Errorf(t.PrintErr())
	} else {
		doneCh <- t
	}
}

func (t *Task) Print() string {
	return fmt.Sprintf(
		"id: %d, result: %s, cratedAt: %s, finishedAt: %s",
		t.id,
		t.result,
		t.craetedAt.Format("2006-01-02 03:04:05"),
		t.finishedAt.Format("2006-01-02 03:04:05"))
}

func (t *Task) PrintErr() string {
	return fmt.Sprintf(
		"id: %d, err: %s, cratedAt: %s",
		t.id,
		t.err,
		t.craetedAt.Format("2006-01-02 03:04:05"),
	)
}

type Tasks map[int]Task

func (ts Tasks) Print() {
	fmt.Println("Done tasks:")
	for _, t := range ts {
		fmt.Println(t.Print())
	}
}
