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

// A Ttype represents a meaninglessness of our life
type Ttype struct {
	id         int
	cT         string // время создания
	fT         string // время выполнения
	taskRESULT []byte
}

const WorkersCount = 70000
const Timout = 3 * time.Second

func main() {
	taskCreturer := func(a chan Ttype) {
		go func() {
			for {
				ft := time.Now().Format(time.RFC3339)
				if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
					ft = "Some error occured"
				}
				a <- Ttype{cT: ft, id: int(time.Now().Unix())} // передаем таск на выполнение
			}
		}()
	}

	superChan := make(chan Ttype, 10)

	go taskCreturer(superChan)

	task_worker := func(a Ttype) Ttype {
		tt, _ := time.Parse(time.RFC3339, a.cT)
		if tt.After(time.Now().Add(-20 * time.Second)) {
			a.taskRESULT = []byte("task has been successed")
		} else {
			a.taskRESULT = []byte("something went wrong")
		}
		a.fT = time.Now().Format(time.RFC3339Nano)

		time.Sleep(time.Millisecond * 150)

		return a
	}

	handledTasks := make(chan Ttype)
	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)

	tasksorter := func(t Ttype) {
		if string(t.taskRESULT[14:]) == "successed" {
			doneTasks <- t
		} else {
			undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT, t.taskRESULT)
		}
	}
	ctx, _ := context.WithTimeout(context.Background(), Timout)

	workerWg := &sync.WaitGroup{}
	workerWg.Add(WorkersCount)
	for i := 0; i < WorkersCount; i++ {
		go func(ctx context.Context) {
			defer workerWg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case t, ok := <-superChan:
					if !ok {
						return
					}
					handledTasks <- task_worker(t)
				}
			}
		}(ctx)
	}
	go func() {
		workerWg.Wait()
		close(handledTasks)
	}()

	go func() {
		defer close(doneTasks)
		defer close(undoneTasks)
		for t := range handledTasks {
			tasksorter(t)
		}
	}()

	result := map[int]Ttype{}
	err := []error{}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for r := range doneTasks {
			result[r.id] = r
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for r := range undoneTasks {
			err = append(err, r)
		}
	}()

	<-ctx.Done()
	wg.Wait()

	println("Errors:")
	for r := range err {
		println(r)
	}

	println("Done tasks:")
	for r := range result {
		println(r)
	}
}
