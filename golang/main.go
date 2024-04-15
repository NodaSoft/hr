package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

const (
	numWorkers = 10
	timeoutSec = 10
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// * обновленный код отправить через pull-request.

// Приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме.
// Должно выводить успешные таски и ошибки по мере выполнения.
// Как видите, никаких привязок к внешним сервисам нет - полный карт-бланш на модификацию кода.

// A Task represents a task to be processed.
type Task struct {
	id         int
	createdAt  string // время создания
	finishedAt string // время выполнения
	result     []byte
	err        error
}

func (t Task) String() string {
	return fmt.Sprintf(
		"task id: %d, created at: %s, finished at: %s, result: %s, err: %v",
		t.id, t.createdAt, t.finishedAt, string(t.result), t.err,
	)
}

func main() {
	/* let's organize the task processing in a pipeline manner */
	ctx, cancel := context.WithTimeout(context.Background(), timeoutSec*time.Second)
	defer cancel()
	tasksDoneCh, tasksUndoneCh := dispatcher(ctx, runWorkerPool(ctx, generator(ctx), numWorkers))

	/* read and print results */
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for task := range tasksDoneCh {
			fmt.Printf("done: %v\n", task)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for task := range tasksUndoneCh {
			fmt.Printf("undone: %v\n", task)
		}
	}()
	wg.Wait()
}

// generator is in charge of creating tasks.
func generator(ctx context.Context) <-chan Task {
	tasksCh := make(chan Task, 10)
	go func() {
		defer close(tasksCh)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				id := int(time.Now().Unix())
				createdAt := time.Now().Format(time.RFC3339)
				// эту логику оставил нетронутой
				if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
					createdAt = "bug: wrong time format"
				}
				tasksCh <- Task{
					id:        id,
					createdAt: createdAt,
				} // передаем таск на выполнение
			}
		}
	}()
	return tasksCh
}

func worker(ctx context.Context, inTasksCh <-chan Task) <-chan Task {
	outTasksCh := make(chan Task)

	go func() {
		defer close(outTasksCh)
		for task := range inTasksCh {
			select {
			case <-ctx.Done(): // the done-guard pattern
				return
			default:
			}

			createdAt, err := time.Parse(time.RFC3339, task.createdAt)
			if err != nil {
				task.err = err
				select {
				case <-ctx.Done():
					return
				case outTasksCh <- task:
					continue
				}
			}

			if createdAt.After(time.Now().Add(-20 * time.Second)) {
				task.result = []byte("task has been succeed")
			} else {
				task.err = errors.New("something went wrong")
			}
			time.Sleep(time.Millisecond * 150)
			task.finishedAt = time.Now().Format(time.RFC3339Nano) // переместил точку окончания работу после sleep
			select {
			case <-ctx.Done():
				return
			case outTasksCh <- task:
			}
		}
	}()

	return outTasksCh
}

// runWorkerPool is a simple helper to run multiple workers concurrently.
func runWorkerPool(ctx context.Context, tasksCh <-chan Task, numWorkers int) <-chan Task {
	workerChannels := make([]<-chan Task, numWorkers)
	for i := range numWorkers {
		workerChannels[i] = worker(ctx, tasksCh)
	}
	return fanIn(ctx, workerChannels...)
}

// fanIn implements fan-in pattern.
func fanIn(ctx context.Context, inChs ...<-chan Task) <-chan Task {
	var wg sync.WaitGroup
	multiplexedCh := make(chan Task)

	wg.Add(len(inChs))
	for _, ch := range inChs {
		go func(ch <-chan Task) {
			defer wg.Done()
			for task := range ch {
				select {
				case <-ctx.Done():
					return
				case multiplexedCh <- task:
				}
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(multiplexedCh)
	}()

	return multiplexedCh
}

// dispatcher distributes the task flow through two channels - done tasks and undone tasks.
func dispatcher(ctx context.Context, tasks <-chan Task) (<-chan Task, <-chan Task) {
	tasksDoneCh := make(chan Task)
	tasksUndoneCh := make(chan Task)
	go func() {
		defer close(tasksDoneCh)
		defer close(tasksUndoneCh)
		for task := range tasks {
			select {
			case <-ctx.Done(): // the done-guard pattern
				return
			default:
				if task.err != nil {
					task.err = fmt.Errorf("task id: %d, time: %s, error: %w", task.id, task.createdAt, task.err)
					select {
					case <-ctx.Done():
						return
					case tasksUndoneCh <- task:
					}
				} else {
					select {
					case <-ctx.Done():
						return
					case tasksDoneCh <- task:
					}

				}
			}
		}
	}()
	return tasksDoneCh, tasksUndoneCh
}
