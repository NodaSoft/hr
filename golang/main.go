package main

import (
	"context"
	"fmt"
	"log"
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

// A Ttype represents a meaninglessness of our life
type Ttype struct {
	id         int
	cT         string // время создания
	fT         string // время выполнения
	taskRESULT []byte
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	superChan := make(chan Ttype, 10)
	go taskCreator(ctx, superChan)

	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)

	go taskProcessor(ctx, superChan, doneTasks, undoneTasks)

	displayResults(ctx, doneTasks, undoneTasks)

	time.Sleep(time.Second * 3)
	cancel()
}

func taskCreator(ctx context.Context, tasks chan<- Ttype) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			ft := time.Now().Format(time.RFC3339)
			if time.Now().Nanosecond()%2 > 0 {
				ft = "Some error occured"
			}
			tasks <- Ttype{cT: ft, id: int(time.Now().Unix())}
			time.Sleep(100 * time.Millisecond) // Добавлено для предотвращения чрезмерной нагрузки
		}
	}
}

func taskProcessor(ctx context.Context, tasks <-chan Ttype, doneTasks chan<- Ttype, undoneTasks chan<- error) {
	for {
		select {
		case <-ctx.Done():
			return
		case t := <-tasks:
			t = processTask(t)
			if string(t.taskRESULT) == "task has been successed" {
				doneTasks <- t
			} else {
				undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT, t.taskRESULT)
			}
		}
	}
}

func processTask(t Ttype) Ttype {
	tt, err := time.Parse(time.RFC3339, t.cT)
	if err != nil || tt.After(time.Now().Add(-20*time.Second)) {
		t.taskRESULT = []byte("task has been successed")
	} else {
		t.taskRESULT = []byte("something went wrong")
	}
	t.fT = time.Now().Format(time.RFC3339Nano)
	time.Sleep(time.Millisecond * 150)
	return t
}

func displayResults(ctx context.Context, doneTasks <-chan Ttype, undoneTasks <-chan error) {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case task := <-doneTasks:
				fmt.Printf("Done Task: %v\n", task)
			}
		}
	}()

	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case err := <-undoneTasks:
				log.Printf("Error: %v\n", err)
			}
		}
	}()

	wg.Wait()
}
