package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// * обновленный код отправить через pull-request.

// Приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме.
// Должно выводить успешные таски и ошибки по мере выполнения.
// Как видите, никаких привязок к внешним сервисам нет - полный карт-бланш на модификацию кода.

const TasksCount = 10

var (
	tasksCreated   = 0
	tasksCompleted = 0
	tasksFails     = 0
)

type task struct {
	id          uuid.UUID
	create      time.Time
	fulfillment time.Time
	result      *error
}

func createTask() task {
	return task{
		id:     uuid.New(),
		create: time.Now(),
	}
}

func executeTask(ctx context.Context, res *task) {
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer func() {
		res.fulfillment = time.Now()
		cancel()
	}()

	var err error
	errCh := make(chan error)

	go func() {
		if time.Now().Nanosecond()%3 > 0 { // %2 выдает 0
			errCh <- errors.New("some error occured")
			return
		}

		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		err = errors.New("task timeout")
	case err = <-errCh:

	}

	if err != nil {
		res.result = &err
		tasksFails++
		log.Print(err)
	}

	tasksCompleted++
}

func main() {
	ctx := context.Background()

	tasks := make(chan *task)
	for i := 1; i <= TasksCount; i++ {
		tasksCreated++

		go func() {
			task := createTask()
			tasks <- &task
		}()
	}

	for {
		select {
		case t := <-tasks:
			go func() {
				executeTask(ctx, t)
			}()
		default:
			if tasksCreated > 0 && tasksCreated == tasksCompleted {
				fmt.Printf("total: %d\r\n", tasksCompleted)
				fmt.Printf("successed: %d\r\n", tasksCompleted-tasksFails)
				fmt.Printf("fails: %d\r\n", tasksFails)
				return
			}
			continue
		}
	}
}
