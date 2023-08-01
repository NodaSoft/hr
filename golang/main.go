package main

import (
	"context"
	"errors"
	"sync"
	"time"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнены остальных тасков

// A Task represents a meaninglessness of our life
type Task struct {
	id              int
	createdTime     time.Time
	executedTime    time.Time
	taskResultBytes []byte
	failed          bool
}

type TaskResult struct {
	task *Task
	data []byte
	err  error
}

const (
	maxTasksInQueue = 10

	taskWorkerSleepInterval = time.Millisecond * 150
	taskBehindInterval      = -20 * time.Second

	taskDistributorPoolCount = 2
)

var (
	somethingWasOccuredErr = errors.New("Some error occured")
	taskInThePastErr       = errors.New("task in the future")

	globalNow = time.Now()
)

// пример на случай разделенных данных

// type TaskErrors struct {
// 	errors []error
// 	mutex  sync.RWMutex
// }

// func (te *TaskErrors) add(err error) {
// 	te.mutex.Lock()
// 	defer te.mutex.Unlock()

// 	te.errors = append(te.errors, err)
// }

// func (te *TaskErrors) read(f func(errs []error)) {
// 	te.mutex.RLock()
// 	defer te.mutex.RUnlock()

// 	f(te.errors)
// }

// type TaskMapResults struct {
// 	results map[int]TaskResult
// 	mutex   sync.RWMutex
// }

// func (tm *TaskMapResults) read(f func(map[int]TaskResult)) {
// 	tm.mutex.RLock()
// 	defer tm.mutex.RUnlock()

// 	f(tm.results)
// }

// func (tm *TaskMapResults) add(tr TaskResult) {
// 	tm.mutex.Lock()
// 	defer tm.mutex.Unlock()

// 	tm.results[tr.task.id] = tr
// }

func runWorkerPoolAsync(ctx context.Context, wg *sync.WaitGroup, taskInput chan Task, w func(Task) TaskResult) chan TaskResult {
	results := make(chan TaskResult)
	closingWorkersWg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		// wait until all workers are finish writing and then close the channel
		closingWorkersWg.Wait()

		close(results)
	}()

	// make the distribution more controllable
	for i := 0; i < taskDistributorPoolCount; i++ {
		wg.Add(1)
		closingWorkersWg.Add(1)
		go func() {
			defer closingWorkersWg.Done()
			defer wg.Done()
			for {
				select {
				case task := <-taskInput:
					taskResult := w(task)
					select {
					case results <- taskResult:
					case <-ctx.Done():
						return
					}

				case <-ctx.Done():
					return
				}
			}
		}()
	}

	return results
}

func main() {
	wg := sync.WaitGroup{}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tasksChan := runTaskCreatorAsync(ctx, &wg)
	taskResults := runWorkerPoolAsync(ctx, &wg, tasksChan, taskWorker)

	result := map[int]TaskResult{}
	errs := []error{}

	doneTasks, undoneTasks := runDistributorAsync(ctx, &wg, taskResults)

	// тут контекст не нужен
	wg.Add(1)
	go func() {
		defer wg.Done()
		for err := range undoneTasks {
			errs = append(errs, err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for tr := range doneTasks {
			result[tr.task.id] = tr
		}
	}()

	time.Sleep(time.Second * 3)
	cancel()  // заканчиваем генерацию и обработку заданий
	wg.Wait() // ждем окончания всех

	println("Errors:")

	for _, r := range errs {
		println(r.Error())
	}

	println("Done tasks ids:")
	for id := range result {
		println(id)
	}
}

func runDistributorAsync(ctx context.Context, wg *sync.WaitGroup, taskResults chan TaskResult) (
	done chan TaskResult,
	undone chan error,
) {
	wg.Add(1)
	done = make(chan TaskResult)
	undone = make(chan error)

	// distributor reading task results
	go func() {
		defer func() {
			close(done)
			close(undone)
		}()
		defer wg.Done()

		for tr := range taskResults {
			if tr.err == nil {
				select {
				case done <- tr:
				case <-ctx.Done():
					return
				}

			} else {
				select {
				case undone <- tr.err:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return
}

func runTaskCreatorAsync(ctx context.Context, wg *sync.WaitGroup) chan Task {
	wg.Add(1)
	tasks := make(chan Task, maxTasksInQueue)
	var id = 0
	go func() {
		defer close(tasks)
		defer wg.Done()

		for {
			now := time.Now()
			failed := false
			if time.Since(globalNow).Nanoseconds()%2 > 0 { // вот такое условие появления ошибочных тасков
				failed = true
			}
			// возможно заменить на UUID
			id += 1
			select {
			case tasks <- Task{createdTime: now, id: id, failed: failed}: // передаем таск на выполнение
			case <-ctx.Done():
				return
			}
		}
	}()

	return tasks
}

func taskWorker(a Task) (taskResult TaskResult) {
	taskResult.task = &a
	if a.failed {
		taskResult.err = somethingWasOccuredErr
	} else {
		if a.createdTime.After(time.Now().Add(taskBehindInterval)) {
			taskResult.data = []byte("task has been succeeded")
		} else {
			taskResult.err = taskInThePastErr
		}
	}

	taskResult.task.executedTime = time.Now()

	time.Sleep(taskWorkerSleepInterval)

	return
}
