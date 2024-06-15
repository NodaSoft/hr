package main

import (
	"context"
	"fmt"
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

// A TaskType represents a meaninglessness of our life
type TaskType struct {
	id         int
	createdAt  string // время создания
	executedAt string // время выполнения
	taskResult string
	err        error
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	doneTasks := make(chan TaskType, 10)
	undoneTasks := make(chan TaskType, 10)

	var wg sync.WaitGroup

	createChan := taskCreator(ctx)
	workerChan := taskWorker(createChan)
	taskSorter(workerChan, doneTasks, undoneTasks, &wg)

	go func() {
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				printResults(doneTasks, undoneTasks)
			case <-ctx.Done():
				return
			}
		}
	}()

	wg.Wait()
}

func taskCreator(ctx context.Context) <-chan TaskType {
	out := make(chan TaskType)
	go func() {
		defer close(out)
		id := 1

		for {
			select {
			case <-ctx.Done():
				return
			default:
				timeStamp := time.Now().Format(time.RFC3339)
				if time.Now().Nanosecond()%2 > 0 {
					timeStamp = "пол шестого"
				}
				out <- TaskType{createdAt: timeStamp, id: id}
				id++
			}
		}
	}()
	return out
}

func taskWorker(in <-chan TaskType) <-chan TaskType {
	out := make(chan TaskType)
	go func() {
		defer close(out)
		for t := range in {
			timeStamp, err := time.Parse(t.createdAt, time.RFC3339)
			if err != nil {
				t.taskResult = "task failed"
				t.err = err
			} else {
				t.taskResult = fmt.Sprintf("task has been completed at %s", timeStamp)
			}
			t.executedAt = time.Now().Format(time.RFC3339Nano)
			time.Sleep(time.Millisecond * 150)
			out <- t
		}
	}()
	return out
}

func taskSorter(in <-chan TaskType, doneChan, errChan chan TaskType, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for t := range in {
			if t.err != nil {
				errChan <- t
			} else {
				doneChan <- t
			}
		}
	}()
}

func printResults(doneTasks, undoneTasks <-chan TaskType) {
	fmt.Println("Done tasks:")
	for len(doneTasks) > 0 {
		t := <-doneTasks
		fmt.Printf("Task id %d created at %s executed at %s message: %s\n", t.id, t.createdAt, t.executedAt, t.taskResult)
	}
	fmt.Println("Errors:")
	for len(undoneTasks) > 0 {
		t := <-undoneTasks
		fmt.Printf("Task id %d created at %s executed at %s error: %s\n", t.id, t.createdAt, t.executedAt, t.err.Error())
	}
}
