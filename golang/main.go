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

//
// A Ttype represents a meaninglessness of our life
//
type LifeTask struct {
	id         int
	createTime string // время создания
	leadTime   string // время выполнения
	taskRESULT []byte
}

func main() {

	var createdNumber int
	var donedNumber int
	var undonedNumber int

	taskCreator := func(a chan LifeTask) {
		go func() {
			for {
				currentTimeNanoSec := time.Now().Format(time.RFC3339Nano)

				if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
					currentTimeNanoSec = "Some error occured"
				}

				a <- LifeTask{
					createTime: currentTimeNanoSec,
					id:         int(time.Now().Nanosecond()),
				} // передаем таск на выполнение

				createdNumber++
			}
		}()
	}

	superChannel := make(chan LifeTask, 10)

	go taskCreator(superChannel)

	task_worker := func(task LifeTask) LifeTask {

		createTime, _ := time.Parse(time.RFC3339, task.createTime)

		if createTime.After(time.Now().Add(-20 * time.Second)) {
			task.taskRESULT = []byte("task has been successed")
		} else {
			task.taskRESULT = []byte("something went wrong")
		}

		task.leadTime = time.Now().Format(time.RFC3339Nano)

		time.Sleep(time.Millisecond * 150)

		return task
	}

	result := map[int]LifeTask{}
	errors := []error{}

	doneTasks := make(chan LifeTask)
	undoneTasks := make(chan error)

	go func(task_channel chan LifeTask) {
		for task := range task_channel {
			result[task.id] = task
			donedNumber++
		}
	}(doneTasks)

	go func(task_channel chan error) {
		for task := range task_channel {
			errors = append(errors, task)
			undonedNumber++
		}
	}(undoneTasks)

	tasksorter := func(t LifeTask) {
		if string(t.taskRESULT[14:]) == "successed" {
			doneTasks <- t
		} else {
			undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.createTime, t.taskRESULT)
		}
	}

	go func() {
		// получение тасков
		for t := range superChannel {
			t = task_worker(t)
			go tasksorter(t)
		}

		close(superChannel)
	}()

	time.Sleep(time.Second * 3)

	fmt.Println("\nCreated - ", createdNumber)

	fmt.Printf("\nErrors (%d):\n", undonedNumber)
	for _, error := range errors {
		println(error)
	}

	fmt.Printf("\nDone tasks (%d): \n", donedNumber)
	for id := range result {
		println(id)
	}

}
