package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Приложение эмулирует получение и обработку неких тасков. Пытается и получать, и обрабатывать в многопоточном режиме.
// Приложение должно генерировать таски 10 сек. Каждые 3 секунды должно выводить в консоль результат всех обработанных к этому моменту тасков (отдельно успешные и отдельно с ошибками).

// ЗАДАНИЕ: сделать из плохого кода хороший и рабочий - as best as you can.
// Важно сохранить логику появления ошибочных тасков.
// Важно оставить асинхронные генерацию и обработку тасков.
// Сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через pull-request в github
// Как видите, никаких привязок к внешним сервисам нет - полный карт-бланш на модификацию кода.

type Task struct {
	id         int
	createTime string
	finishTime string
	result     []byte
}

func (t Task) String() string {
	return fmt.Sprintf("task id: %d, create time: %s, finish time: %s", t.id, t.createTime, t.finishTime)
}

func taskMaker(tasks chan Task) {
	for {
		createTime := time.Now().Format(time.RFC3339)
		if rand.Intn(100)%2 > 0 {
			createTime = "Some error occured"
		}
		tasks <- Task{id: int(time.Now().Unix()), createTime: createTime}
	}
}

func taskWorker(tasks chan Task, completedTasks chan Task) {
	for task := range tasks {
		taskTime, _ := time.Parse(time.RFC3339, task.createTime)
		if taskTime.After(time.Now().Add(-20 * time.Second)) {
			task.result = []byte("task has been successed")
		} else {
			task.result = []byte("something went wrong")
		}
		task.finishTime = time.Now().Format(time.RFC3339Nano)
		time.Sleep(150 * time.Millisecond)
		completedTasks <- task
	}
	close(completedTasks)
}

func taskSorter(tasks chan Task, doneTasks chan Task, undoneTasks chan Task) {
	for task := range tasks {
		if string(task.result[14:]) == "successed" {
			doneTasks <- task
		} else {
			undoneTasks <- task
		}
	}
	close(doneTasks)
	close(undoneTasks)
}

func main() {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	done := map[int]Task{}
	undone := map[int]Task{}

	go func() {
		taskChannel := make(chan Task, 10)
		completedTasks := make(chan Task)
		doneTasks := make(chan Task)
		undoneTasks := make(chan Task)

		go taskMaker(taskChannel)
		go taskWorker(taskChannel, completedTasks)
		go taskSorter(completedTasks, doneTasks, undoneTasks)

		go func() {
			for t := range doneTasks {
				done[t.id] = t
			}
		}()
		go func() {
			for t := range undoneTasks {
				undone[t.id] = t
			}
		}()
	}()

	go func() {
		for range ticker.C {
			fmt.Println("Done tasks:")
			for _, t := range done {
				fmt.Println(t)
			}

			fmt.Println("Undone tasks:")
			for _, t := range undone {
				fmt.Println(t)
			}
		}
	}()

	time.Sleep(10 * time.Second)
}
