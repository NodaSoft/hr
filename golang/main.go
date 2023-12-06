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

const (
	taskTTL         = 20 * time.Second
	taskChanSize    = 10
	serverSleepTime = 3 * time.Second
)

// A Task represents a meaninglessness of our life
type Task struct {
	id         int
	createTime string // время создания
	finishTime string // время выполнения
	successful bool
	resMessage string
}

func createTasks(ctx context.Context, taskChan chan Task) {
	id := 0
	for {
		select {
		case <-ctx.Done():
			close(taskChan)
			return
		default:
			curTime := time.Now().Format(time.RFC3339)
			if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
				continue
			}
			taskChan <- Task{createTime: curTime, id: id} // передаем таск на выполнение
			id++
		}
	}
}

func (task *Task) work() {
	taskCreateTime, err := time.Parse(time.RFC3339, task.createTime)
	if err != nil {
		return
	}
	if taskCreateTime.After(time.Now().Add(-1 * taskTTL)) {
		task.successful = true
	} else {
		task.successful = false
		task.resMessage = "something went wrong"
	}
	task.finishTime = time.Now().Format(time.RFC3339Nano)

	time.Sleep(time.Millisecond * 150)
}

func sortTasks(superChan, doneTasks chan Task, undoneTasks chan error) {
	taskSortersWG := sync.WaitGroup{}
	for t := range superChan {
		t.work()
		taskSortersWG.Add(1)
		go func(t Task) {
			if t.successful {
				doneTasks <- t
			} else {
				undoneTasks <- fmt.Errorf("task id %d time %s, error %s", t.id, t.createTime, t.resMessage)
			}
			taskSortersWG.Done()
		}(t)
	}
	taskSortersWG.Wait()
	close(doneTasks)
	close(undoneTasks)
}

func TestTasksWorker(ctx context.Context, result *[]Task, err *[]error) {
	superChan := make(chan Task, taskChanSize)

	go createTasks(ctx, superChan)

	doneTasks := make(chan Task)
	undoneTasks := make(chan error)

	go sortTasks(superChan, doneTasks, undoneTasks)

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		for r := range doneTasks {
			*result = append(*result, r)
		}
		wg.Done()
	}()
	go func() {
		for r := range undoneTasks {
			*err = append(*err, r)
		}
		wg.Done()
	}()
	wg.Wait()
}

func main() {
	ctx, ctxCanclel := context.WithCancel(context.Background())

	var result []Task
	var err []error
	go TestTasksWorker(ctx, &result, &err)
	time.Sleep(serverSleepTime)
	ctxCanclel()

	println("Errors:")
	for _, r := range err {
		fmt.Println(r)
	}

	println("Done tasks:")
	for _, r := range result {
		fmt.Println(r)
	}
}
