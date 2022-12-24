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

// A Task represents a meaninglessness of our life
type Task struct {
	Id          int
	CreateTime  string // время создания
	ProcessTime string // время выполнения
	TaskResult  []byte
}

func taskCreturer(ctx context.Context) <-chan Task {
	superChan := make(chan Task, 10)
	go func() {
		defer close(superChan)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				cT := time.Now().Format(time.RFC3339)
				if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
					cT = "Some error occurred"
				}
				superChan <- Task{CreateTime: cT, Id: int(time.Now().Unix())} // передаем таск на выполнение
				time.Sleep(1 * time.Microsecond)
			}
		}
	}()
	return superChan
}

func taskWorker(t Task) Task {
	tt, err := time.Parse(time.RFC3339, t.CreateTime)

	if err != nil {
		t.TaskResult = []byte("invalid parse task create time")
		return t
	}
	if tt.After(time.Now().Add(-20 * time.Second)) {
		t.TaskResult = []byte("task has been successed")
	} else {
		t.TaskResult = []byte("invalid task create time")
	}

	t.ProcessTime = time.Now().Format(time.RFC3339Nano)

	time.Sleep(time.Millisecond * 150)

	return t
}

func taskSorter(t Task, done chan<- Task, undone chan<- error) {

	if string(t.TaskResult[14:]) == "successed" {
		done <- t
	} else {
		undone <- fmt.Errorf("Task id %d time %s, error %s", t.Id, t.CreateTime, string(t.TaskResult))
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	doneTasks := make(chan Task, 100)
	undoneTasks := make(chan error, 100)
	defer close(doneTasks)
	defer close(undoneTasks)

	superChan := taskCreturer(ctx)

	wg := sync.WaitGroup{}

	go func() {
		// получение тасков
		for {
			select {
			case <-ctx.Done():
				return
			case task := <-superChan:
				go func() {
					wg.Add(1)
					processed := taskWorker(task)
					taskSorter(processed, doneTasks, undoneTasks)
					wg.Done()
				}()

			}
		}

	}()

	result := map[int]Task{}
	err := make([]error, 0, 16)
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case r := <-doneTasks:
				result[r.Id] = r
			case er := <-undoneTasks:
				err = append(err, er)
			}
		}
	}(ctx)

	time.Sleep(time.Second * 3)
	cancel()
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
