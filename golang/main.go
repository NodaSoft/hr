package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// Приложение эмулирует получение и обработку неких тасков. Пытается и получать, и обрабатывать в многопоточном режиме.
// После обработки тасков в течении 3 секунд приложение должно выводить накопленные к этому моменту успешные таски и отдельно ошибки обработки тасков.

// ЗАДАНИЕ: сделать из плохого кода хороший и рабочий - as best as you can.
// Важно сохранить логику появления ошибочных тасков.
// Важно оставить асинхронные генерацию и обработку тасков.
// Сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через pull-request в github
// Как видите, никаких привязок к внешним сервисам нет - полный карт-бланш на модификацию кода.

// A Ttype represents a meaninglessness of our life
type Ttype struct {
	id         int
	cT         string // время создания
	fT         string // время выполнения
	taskRESULT []byte
}

func main() {
	var wg sync.WaitGroup
	var mu sync.Mutex
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	taskCreator := func(ctx context.Context, a chan Ttype) {
		go func() {
			for {
				select {
				case <-ctx.Done():
					close(a)
					return
				default:
					ft := time.Now().Format(time.RFC3339)
					simulatedError := time.Now().Nanosecond()%2 == 0
					if simulatedError {
						ft = "Some error occured"
					}
					a <- Ttype{cT: ft, id: int(time.Now().Unix())}
					time.Sleep(100 * time.Millisecond)
				}
			}
		}()
	}

	superChan := make(chan Ttype, 10)
	go taskCreator(ctx, superChan)

	taskWorker := func(ctx context.Context, t Ttype) (Ttype, error) {
		tt, err := time.Parse(time.RFC3339, t.cT)
		if err != nil {
			return Ttype{}, fmt.Errorf("error parsing task creation time: %w", err)
		}

		if tt.After(time.Now().Add(-20 * time.Second)) {
			t.taskRESULT = []byte("task has been successed")
		} else {
			t.taskRESULT = []byte("something went wrong")
		}
		t.fT = time.Now().Format(time.RFC3339Nano)

		select {
		case <-ctx.Done():
			return Ttype{}, context.Canceled
		case <-time.After(150 * time.Millisecond):
			return t, nil
		}
	}

	doneTasks := make(chan Ttype, 10)
	undoneTasks := make(chan error, 10)

	tasksorter := func(ctx context.Context, t Ttype) {
		if err := ctx.Err(); err != nil {
			return
		}
		if string(t.taskRESULT) == "task has been successed" {
			doneTasks <- t
		} else {
			undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT, string(t.taskRESULT))
		}
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for t := range superChan {
			task, err := taskWorker(ctx, t)
			if err != nil && !errors.Is(err, context.Canceled) {
				undoneTasks <- err
				continue
			}
			tasksorter(ctx, task)
		}
		close(doneTasks)
		close(undoneTasks)
	}()

	result := map[int]Ttype{}
	err := []error{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for r := range doneTasks {
			mu.Lock()
			result[r.id] = r
			mu.Unlock()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for r := range undoneTasks {
			mu.Lock()
			err = append(err, r)
			mu.Unlock()
		}
	}()

	wg.Wait()

	println("Errors:")
	for _, r := range err {
		fmt.Println(r)
	}

	println("Done tasks:")
	for _, r := range result {
		fmt.Println(r)
	}
}
