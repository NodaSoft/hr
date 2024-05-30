package main

import (
	"fmt"
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

type taskStatus string

const (
	SuccesedTask        taskStatus = "task has been successed"
	badTaskCreationTime taskStatus = "bad task creation time"
	creationTaskError   taskStatus = "failted to create task"
)

const (
	creationTimeLimit   = 200                     // nanosecond, deault = 2
	creationTimeOffset  = -150 * time.Millisecond // default = 20 seconds
	workerSleepDuration = 150 * time.Millisecond  // default = 150 miliseconds
	appWorkTime         = 3 * time.Second         // default = 3 seconds
)

// A Task represents a meaninglessness of our life
type Task struct {
	id          int64
	createdTime time.Time // время создания
	finishTime  time.Time // время выполнения
	taskResult  taskStatus
}

func taskCreturer(taskChan chan Task) {
	go func() {
		for {
			timeCreation := time.Now()
			var task Task
			if time.Now().Nanosecond()%creationTimeLimit > 0 { // вот такое условие появления ошибочных тасков
				task.taskResult = creationTaskError // ft = "Some error occured"
			}
			task.id = timeCreation.UnixNano()
			task.createdTime = timeCreation
			// передаем таск на выполнение
			taskChan <- task
		}
	}()
}

func taskWorker(task Task) Task {
	if task.createdTime.After(time.Now().Add(creationTimeOffset)) {
		task.taskResult = SuccesedTask
	} else if task.taskResult == "" {
		task.taskResult = badTaskCreationTime
	}
	time.Sleep(workerSleepDuration)
	task.finishTime = time.Now() // set finish time after useful work (sleep)
	return task
}

func taskSorter(t Task, doneTasks chan Task, undoneTasks chan error) {
	if t.taskResult == SuccesedTask {
		doneTasks <- t
	} else {
		undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.createdTime, t.taskResult)
	}
}

func main() {
	superChan := make(chan Task)
	doneTasks := make(chan Task)
	undoneTasks := make(chan error)
	result := make(map[int64]Task)
	err := make([]error, 0)
	defer close(superChan)
	defer close(doneTasks)
	defer close(undoneTasks)

	taskCreturer(superChan)

	// tasks working and sort
	go func() {
		// получение тасков
		for t := range superChan {
			taskSorter(taskWorker(t), doneTasks, undoneTasks)
		}
	}()

	// fill results
	go func() {
		for r := range doneTasks {
			result[r.id] = r
		}
	}()
	go func() {
		for r := range undoneTasks {
			err = append(err, r)
		}
	}()

	time.Sleep(appWorkTime)

	// print results
	fmt.Println("Errors:")
	for _, r := range err {
		fmt.Println(r.Error())
	}
	fmt.Println("Done tasks:")
	for _, r := range result {
		fmt.Printf("Task id %d time %s, %s\n", r.id, r.createdTime.String(), r.taskResult)
	}
}
