package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
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
	ctx, cancel := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	taskCreturer := func() <-chan Ttype {
		out := make(chan Ttype, 10)
		go func() {
			defer close(out)

			t := time.NewTicker(time.Millisecond * 150)
			defer t.Stop()

			for {
				select {
				case <-ctx.Done():
					return
				case <-t.C:
					ft := time.Now().Format(time.RFC3339)
					if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
						ft = "Some error occured"
					}
					out <- Ttype{cT: ft, id: int(time.Now().Unix())} // передаем таск на выполнение
				}
			}
		}()
		return out
	}

	task_worker := func(in <-chan Ttype) <-chan Ttype {
		out := make(chan Ttype)
		go func() {
			for a := range in {
				tt, err := time.Parse(time.RFC3339, a.cT)
				if err == nil && tt.After(time.Now().Add(-20*time.Second)) {
					a.taskRESULT = []byte("task has been successed")
				} else {
					a.taskRESULT = []byte("something went wrong")
				}

				a.fT = time.Now().Format(time.RFC3339Nano)
				out <- a
			}
			close(out)
		}()
		return out
	}

	tasksorter := func(in <-chan Ttype) (<-chan Ttype, <-chan error) {
		doneTasks := make(chan Ttype)
		undoneTasks := make(chan error)
		go func() {
			for t := range in {
				if string(t.taskRESULT[14:]) == "successed" {
					doneTasks <- t
				} else {
					undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT, t.taskRESULT)
				}
			}
			close(doneTasks)
			close(undoneTasks)
		}()
		return doneTasks, undoneTasks
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	result := map[int]Ttype{}
	err := []error{}
	collectTasks := func(doneTasks <-chan Ttype, undoneTasks <-chan error) {
		go func() {
			defer wg.Done()
		for_exit:
			for {
				select {
				case r, ok := <-doneTasks:
					if !ok {
						break for_exit
					}
					result[r.id] = r

				case r, ok := <-undoneTasks:
					if !ok {
						break for_exit
					}
					err = append(err, r)
				}
			}
		}()
	}

	collectTasks(tasksorter(task_worker(taskCreturer())))

	time.Sleep(3 * time.Second)
	cancel()

	wg.Wait()

	fmt.Println("Errors:")
	for r, t := range err {
		fmt.Println(r, t)
	}

	fmt.Println("Done tasks:")
	for r, t := range result {
		fmt.Printf("%v, %v\n", r, string(t.taskRESULT))
	}
}
