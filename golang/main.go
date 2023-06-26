package main

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	id         int64
	cT         string
	fT         string
	taskRESULT []byte
}

func taskCreator(taskChan chan<- Task) {
	for {
		createTime := time.Now().Format(time.RFC3339)
		if time.Now().Nanosecond()%2 > 0 { // условие появления ошибочных заданий
			createTime = "Some error occurred"
		}
		taskChan <- Task{cT: createTime, id: time.Now().Unix()} // передаем задание на выполнение
		time.Sleep(time.Second) // добавляем задержку, чтобы избежать бесконечного цикла
	}
}

func taskWorker(task Task) Task {
	taskTime, _ := time.Parse(time.RFC3339, task.cT)
	if taskTime.After(time.Now().Add(-20 * time.Second)) {
		task.taskRESULT = []byte("task has been successed")
	} else {
		task.taskRESULT = []byte("something went wrong")
	}
	task.fT = time.Now().Format(time.RFC3339Nano)
	time.Sleep(time.Millisecond * 150)
	return task
}

func taskSorter(doneTasks chan<- Task, undoneTasks chan<- error, task Task) {
	if string(task.taskRESULT[14:]) == "successed" {
		doneTasks <- task
	} else {
		undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", task.id, task.cT, string(task.taskRESULT))
	}
}

func main() {
	taskChan := make(chan Task, 10)
	doneTasks := make(chan Task)
	undoneTasks := make(chan error)

	go taskCreator(taskChan)

	var wg sync.WaitGroup

	go func() {
		for task := range taskChan {
			wg.Add(1)
			go func(t Task) {
				defer wg.Done()
				t = taskWorker(t)
				taskSorter(doneTasks, undoneTasks, t)
			}(task)
		}
	}()

	go func() {
		wg.Wait()
		close(doneTasks)
		close(undoneTasks)
	}()

	result := map[int64]Task{}
	err := []error{}

	for r := range doneTasks {
		result[r.id] = r
	}
	for r := range undoneTasks {
		err = append(err, r)
	}

	fmt.Println("Errors:")
	for _, e := range err {
		fmt.Println(e)
	}

	fmt.Println("Done tasks:")
	for _, r := range result {
		fmt.Printf("Task id: %d, create time: %s, finish time: %s, result: %s\n", r.id, r.cT, r.fT, string(r.taskRESULT))
	}
}
