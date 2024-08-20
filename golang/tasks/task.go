package tasks

import (
	"context"
	"errors"
	"time"
)

type Task struct {
	ID         int
	CreatedAt  time.Time
	FinishedAt time.Time
	Result     []byte
	Err        error
}

// Генератор тасок не должен знать заранее какая таска ошибочная
func RunCreator(ctx context.Context) <-chan *Task {
	ch := make(chan *Task)

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(ch)
				return
			default:
				createdTime := time.Now()
				task := &Task{ID: int(createdTime.Unix()), CreatedAt: createdTime}

				// Сохраняю условие из исходного кода
				if task.CreatedAt.Nanosecond()%2 > 0 {
					task.Err = ErrTaskFailed
				}

				ch <- task
			}
		}
	}()

	return ch
}

// Имитация работы таски
func work(task *Task) {
	// Добавляю ещё одно условие из-за возможной таски с изначальной ошибкой
	// Остальные ветвления оставил без изменения (почти)
	if task.Err != nil {
		task.Result = []byte("initial error")
	} else if task.CreatedAt.After(time.Now().Add(-20 * time.Second)) {
		task.Result = []byte("task has been successed")
	} else {
		task.Err = errors.New("new error")
		task.Result = []byte("something went wrong")
	}

	task.FinishedAt = time.Now()
	time.Sleep(time.Millisecond * 150)
}

// Обработчик тасок
func RunWorker(ctx context.Context, ch <-chan *Task) <-chan *Task {
	handled := make(chan *Task)

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(handled)
				return
			case task := <-ch:
				work(task)
				handled <- task
			}
		}
	}()

	return handled
}
