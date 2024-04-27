package main

import (
	"context"
	"fmt"
	"time"
)

// Приложение эмулирует получение и обработку неких тасков. Пытается и получать, и обрабатывать в многопоточном режиме.
// Должно выводить успешные таски и ошибки обработки тасков по мере выполнения.
// ЗАДАНИЕ: сделать из плохого кода хороший и рабочий - as best as you can.
// Важно сохранить логику появления ошибочных тасков.
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

func taskCreator(ctx context.Context, out chan<- Ttype)  {
	for {
		select {
		case <-ctx.Done():
			close(out)
			return
		default:
			ft := time.Now().Format(time.RFC3339)
			if time.Now().Nanosecond() % 2 > 0 { // вот такое условие появления ошибочных тасков
				ft = "Some error occured"
			}
			out <- Ttype{cT: ft, id: int(time.Now().Unix())}
		}
	}
}

func taskWorker(t Ttype, d chan<- Ttype, u chan<- error) {
	tt, _ := time.Parse(time.RFC3339, t.cT)
	t.fT = time.Now().Format(time.RFC3339Nano)
	if tt.After(time.Now().Add(-20 * time.Second)) {
		t.taskRESULT = []byte("task has been successed")
		d <- t
	} else {
		t.taskRESULT = []byte("something went wrong")
		u <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT, t.taskRESULT)
	}
}

func main() {
	workTime := time.Millisecond * 1
	ctx, _ := context.WithTimeout(context.Background(), workTime)
	createdTasks := make(chan Ttype, 10)

	go taskCreator(ctx, createdTasks)

	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)
	result := []Ttype{}
	err := []error{}
WORKERS:
	for  {
		select {
		case <-ctx.Done():
			break WORKERS
		case task := <-createdTasks:
			go taskWorker(task, doneTasks, undoneTasks)
		}
	}

RESULT:
	for  {
		select {
		case dTask := <-doneTasks:
			result = append(result, dTask)
		case erTask := <-undoneTasks:
			err = append(err, erTask)
		default:
			break RESULT
		}

	}

	println("Errors:")
	for e := range err {
		println(e)
	}

	println("Done tasks:")
	for r := range result {
		println(r)
	}
}
