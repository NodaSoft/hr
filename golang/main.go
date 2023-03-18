package main

import (
	"fmt"
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
	id            int
	createTime    string
	executionTime string
	taskResult    []byte
}

var (
	taskChan    chan Task
	doneTasks   chan Task
	undoneTasks chan error
)

func init() {
	taskChan = make(chan Task, 10)
	doneTasks = make(chan Task)
	undoneTasks = make(chan error)
}

func main() {
	println("Please wait 3 seconds")

	go taskCreator(taskChan)

	go func() {
		// получение тасков
		for task := range taskChan {
			task = taskWorker(task)
			go taskSorter(task)
		}
		close(taskChan)
	}()

	result := map[int]Task{}
	err := []error{}
	go func() {
		for {
			select {
			case doneTask, ok := <-doneTasks:
				result[doneTask.id] = doneTask
				if !ok {
					doneTasks = nil
				}
			case undoneTask, ok := <-undoneTasks:
				err = append(err, undoneTask)
				if !ok {
					undoneTasks = nil
				}
			}

			if doneTasks == nil && undoneTasks == nil {
				break
			}
		}
	}()

	time.Sleep(time.Second * 3)
	close(doneTasks)
	close(undoneTasks)

	outputResult(result, err)
}

func taskCreator(taskChan chan Task) {
	for {
		createTime := time.Now().Format(time.RFC3339)
		if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
			createTime = "Some error occured"
		}
		taskChan <- Task{id: int(time.Now().Unix()), createTime: createTime} // передаем таск на выполнение
	}
}

func taskWorker(task Task) Task {
	taskTime, _ := time.Parse(time.RFC3339, task.createTime)
	if taskTime.After(time.Now().Add(-20 * time.Second)) {
		task.taskResult = []byte("task has been successed")
	} else {
		task.taskResult = []byte("something went wrong")
	}
	task.executionTime = time.Now().Format(time.RFC3339Nano)

	time.Sleep(time.Millisecond * 150)

	return task
}

func taskSorter(task Task) {
	if string(task.taskResult[14:]) == "successed" {
		doneTasks <- task
	} else {
		undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", task.id, task.createTime, task.taskResult)
	}
}

func outputResult(result map[int]Task, err []error) {
	println("Done tasks:")
	for val := range result {
		println(val)
	}

	println("Errors:")
	for _, val := range err {
		println(val.Error())
	}
}
