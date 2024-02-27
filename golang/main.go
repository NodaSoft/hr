package main

import (
	"fmt"
	"strings"
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
type taskType struct {
	id               int
	taskCreationTime string // время создания
	taskFailedTime   string // время выполнения
	taskResult       []byte
}

func main() {
	task := make(chan taskType, 10)

	go taskCreator(task)

	doneTasks := make(chan taskType)
	undoneTasks := make(chan error)

	result := make(map[int]taskType)
	errs := []error{}
	go func() {
        for {
            select {
            case r, ok := <-doneTasks:
                if !ok {
                    return
                }
                result[r.id] = r
            case err, ok := <-undoneTasks:
                if !ok {
                    return
                }
                errs = append(errs, err)
            }
        }
    }()

	// получение тасков
	for t := range task {
		t = taskWorker(t)
		taskSorter(t, doneTasks, undoneTasks)
	}
	
	close(doneTasks)
    close(undoneTasks)

	fmt.Println("Done tasks:")
	for id, result := range result {
		fmt.Printf("Task ID: %d, Result: %+v\n", id, result)
	}

	fmt.Println("Errors:")
	for _, err := range errs {
		fmt.Println(err)
	}
}

func taskCreator(a chan taskType) {
	for {
		creationTime := time.Now().Format(time.RFC3339)
		if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
			creationTime = "Some error occured"
		}
		a <- taskType{taskCreationTime: creationTime, id: int(time.Now().Unix())} // передаем таск на выполнение
	}
}

func taskWorker(a taskType) taskType {
	tt, err := time.Parse(time.RFC3339, a.taskCreationTime)
	if err != nil {
		a.taskResult = []byte("something went wrong")
	}

	if tt.After(time.Now().Add(-20*time.Second)) && a.taskResult == nil {
		a.taskResult = []byte("task has been successed")
	} else {
		a.taskResult = []byte("something went wrong")
	}
	a.taskFailedTime = time.Now().Format(time.RFC3339Nano)

	time.Sleep(time.Millisecond * 150)

	return a
}

func taskSorter(t taskType, doneTasks chan taskType, undoneTasks chan error) {
	if strings.Contains(string(t.taskResult), "successed") {
		doneTasks <- t
	} else {
		undoneTasks <- fmt.Errorf("Task id: %d, Creation time: %s, Result: %s", t.id, t.taskCreationTime, t.taskResult)
	}
}