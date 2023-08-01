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

type TaskResult = []byte

var (
	TaskSuccess = TaskResult("task has been successed")
	TaskError   = TaskResult("something went wrong")
)

// A Ttype represents a meaninglessness of our life
// Методы обработки задачи и ее состояния хранятся в самой таске в текущей реализации.
type Ttype struct {
	id     int
	ct     string // время создания
	ft     string // время выполнения
	result TaskResult
	err    error
}

func NewTtype(ft string) Ttype {
	return Ttype{ct: ft, id: int(time.Now().UnixMilli())} // для уникальности в текущем контексте подойдет UnixMilli с интервалом отправки 150ms, но лучше uuid, serial
}

// Метод обработки и проверки состояния таски
func (t Ttype) Process(done func(Ttype)) {
	t.worker()
	t.sorter()
	done(t)
}

func (t *Ttype) worker() {
	tt, err := time.Parse(time.RFC3339, t.ct)
	if err == nil && tt.After(time.Now().Add(-20*time.Second)) {
		t.result = TaskSuccess
	} else {
		t.result = TaskError
	}

	t.ft = time.Now().Format(time.RFC3339Nano)
}

func (t *Ttype) sorter() {
	if string(t.result) != string(TaskSuccess) {
		t.err = fmt.Errorf("Task id %d time %s, error %s", t.id, t.ct, t.result)
	}
}

func (t Ttype) Error() error {
	return t.err
}

// Метод для генерации задач. Возвращаем канал на чтение куда с интервалом времени отправляются задачи.
func taskCreturer(ctx context.Context) <-chan Ttype {
	ch := make(chan Ttype, 10)

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(ch)
				return
			default:
				ft := time.Now().Format(time.RFC3339)
				if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
					ft = "Some error occured"
				}
				ch <- NewTtype(ft)                 // передаем таск на выполнение
				time.Sleep(time.Millisecond * 150) // отправляем таск в канал через указанный интервал
			}
		}
	}()

	return ch
}

// Метод для распараллеливания обработки задач.
func concurrencyPipe(ch <-chan Ttype) (<-chan Ttype, <-chan error) {
	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)

	go func() {
		wg := sync.WaitGroup{}
		done := func(t Ttype) {
			defer wg.Done()

			if err := t.Error(); err != nil {
				undoneTasks <- err
			} else {
				doneTasks <- t
			}
		}

		// получение тасков
		for t := range ch {
			wg.Add(1)
			go t.Process(done)
		}

		wg.Wait()
		close(doneTasks)
		close(undoneTasks)
	}()

	return doneTasks, undoneTasks
}

func main() {
	var (
		result = map[int]Ttype{}
		errs   = []error{}
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	superChan := taskCreturer(ctx)
	doneTasks, undoneTasks := concurrencyPipe(superChan)

	for loop := true; loop; {
		select {
		case <-ctx.Done():
			loop = false
		case t, ok := <-doneTasks:
			if loop = ok; ok {
				result[t.id] = t
			}
		case err, ok := <-undoneTasks:
			if loop = ok; ok {
				errs = append(errs, err)
			}
		}
	}

	println("Errors:")
	for _, t := range errs {
		println(t.Error())
	}

	println("Done tasks:")
	for r := range result {
		println(fmt.Sprintf("Task id %d successed", r))
	}
}
