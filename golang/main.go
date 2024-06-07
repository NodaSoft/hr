package main

import (
	"fmt"
	"strings"
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