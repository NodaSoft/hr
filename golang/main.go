package main

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

// Приложение эмулирует получение и обработку неких тасков. Пытается и получать, и обрабатывать в многопоточном режиме.
// Приложение должно генерировать таски 10 сек. Каждые 3 секунды должно выводить в консоль результат всех обработанных к этому моменту тасков (отдельно успешные и отдельно с ошибками).

// ЗАДАНИЕ: сделать из плохого кода хороший и рабочий - as best as you can.
// Важно сохранить логику появления ошибочных тасков.
// Важно оставить асинхронные генерацию и обработку тасков.
// Сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через pull-request в github
// Как видите, никаких привязок к внешним сервисам нет - полный карт-бланш на модификацию кода.

const (
	timeToCreate     = 10
	printReportEvery = 3
	successStatus    = "task has been successfully completed"
	nonSuccessStatus = "something went wrong"
)

// Task представляет задачу с необходимыми атрибутами.
type Task struct {
	id         int64
	createdAt  time.Time
	finishedAt time.Time
	result     string
}

// taskCreator эмулирует процесс создания задач.
func taskCreator(ctx context.Context, tasks chan<- Task) {
	defer close(tasks)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			now := time.Now()
			id := now.Unix()
			if now.Nanosecond()%2 > 0 {
				id = 0
			}
			task := Task{
				createdAt: now,
				id:        id,
			}
			tasks <- task
		}
	}
}

// taskWorker обрабатывает задачи из канала tasks, проставляет статусы и отправляет в соответствующие каналы.
func taskWorker(wg *sync.WaitGroup, tasks <-chan Task, doneTasks chan<- Task, undoneTasks chan<- Task) {
	defer wg.Done()
	for task := range tasks {
		time.Sleep(150 * time.Millisecond)
		task.finishedAt = time.Now()
		if task.id == 0 {
			task.result = nonSuccessStatus
			undoneTasks <- task
		} else {
			task.result = successStatus
			doneTasks <- task
		}
	}
}

// taskLogger выводит в консоль задачи из каналов.
func taskLogger(doneTasks <-chan Task, undoneTasks <-chan Task) {
	logTasks := func(ch <-chan Task, taskType string) {
		for {
			select {
			case task, ok := <-ch:
				if !ok {
					return
				}

				fmt.Printf("[%s] task created: %s, finished: %s, result: %s\n",
					taskType,
					task.createdAt.Format(time.RFC3339),
					task.finishedAt.Format(time.RFC3339),
					task.result)
			default:
				return
			}
		}
	}

	logTasks(doneTasks, "Done")
	logTasks(undoneTasks, "Undone")
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), timeToCreate*time.Second)
	defer cancel()

	tasks := make(chan Task)
	doneTasks := make(chan Task)
	undoneTasks := make(chan Task)

	go taskCreator(ctx, tasks)

	var wg sync.WaitGroup
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go taskWorker(&wg, tasks, doneTasks, undoneTasks)
	}

	go func() {
		wg.Wait()
		close(doneTasks)
		close(undoneTasks)
	}()

	ticker := time.NewTicker(printReportEvery * time.Second)
	defer ticker.Stop()

	done := false
	for !done {
		select {
		case <-ctx.Done():
			done = true
		case <-ticker.C:
			taskLogger(doneTasks, undoneTasks)
		}
	}
}
