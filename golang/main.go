package main

import (
	"context"
	"fmt"
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

// task represents the meaninglessness of our lives
type task struct {
	id             int
	creationTime   string // время создания
	completionTime string // время выполнения
	result         string
}

const (
	creatorDuration = 3 * time.Second
	inputBuffer     = 10
	numWorkers      = 10
)

func createTasks(ctx context.Context, ch chan *task) {
	// у time.Now() наносекунды всегда "000"
	start := time.Now()

	for {
		select {
		case <-ctx.Done():
			close(ch)

			return
		default:
			var ct string

			// Получаем ненулевые наносекунды
			unique := int(time.Since(start).Nanoseconds())

			if unique%2 > 0 { // вот такое условие появления ошибочных тасков
				// Я предположил, что на этом этапе мы не знаем об ошибке, и поймем это только при обработке
				ct = "Some error occurred"
			} else {
				ct = time.Now().Format(time.RFC3339)
			}

			ch <- &task{id: unique, creationTime: ct} // передаем таск на выполнение
		}
	}
}

func processTasks(inputCh, doneCh chan *task, failedCh chan error) {
	// ограничил макс кол-во параллельных воркеров
	sem := make(chan struct{}, numWorkers)
	var wg sync.WaitGroup

	for t := range inputCh {
		sem <- struct{}{}
		wg.Add(1)

		go func(t *task) {
			defer func() {
				wg.Done()
				<-sem
			}()

			processTask(t, doneCh, failedCh)
		}(t)
	}

	wg.Wait()

	close(doneCh)
	close(failedCh)
}

func processTask(t *task, doneCh chan *task, failCh chan error) {
	// По-хорошему такого не должно быть конечно, но для читаемости вывода можно оставить
	time.Sleep(time.Millisecond * 150)

	ct, err := time.Parse(time.RFC3339, t.creationTime)
	if err != nil {
		// Тут t.creationTime естественно не будет, потому что оно и не записалось правильно
		failCh <- fmt.Errorf("task id: %d, error: %w", t.id, err)

		return
	}

	if ct.Before(time.Now().Add(-20 * time.Second)) {
		t.result = "something went wrong"
		failCh <- fmt.Errorf("task id: %d time: %s, error: %s", t.id, t.creationTime, t.result)

		return
	}

	t.result = "task completed successfully"
	t.completionTime = time.Now().Format(time.RFC3339)
	doneCh <- t
}

func printTasks(doneCh chan *task, failCh chan error) {
	var result []*task
	var errs []*error

	resDone := make(chan struct{})

	go func() {
		for t := range doneCh {
			result = append(result, t)
		}

		resDone <- struct{}{}
	}()

	for e := range failCh {
		e := e
		errs = append(errs, &e)
	}

	<-resDone

	fmt.Println("Errors:")

	for _, e := range errs {
		fmt.Println(*e)
	}

	fmt.Println()
	fmt.Println("Completed tasks:")

	for _, r := range result {
		fmt.Println(*r)
	}

	fmt.Println()
	fmt.Printf("Completed %d tasks, failed %d tasks", len(result), len(errs))
}

func main() {
	inputCh := make(chan *task, inputBuffer)
	doneCh := make(chan *task)
	failCh := make(chan error)

	ctx, cancel := context.WithTimeout(context.Background(), creatorDuration)
	defer cancel()

	// создание тасков
	go createTasks(ctx, inputCh)

	// обработка тасков
	go processTasks(inputCh, doneCh, failCh)

	// вывод результатов
	printTasks(doneCh, failCh)
}
