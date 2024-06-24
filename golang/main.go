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

type Task struct {
	id         int
	createTime time.Time
	finishTime time.Time
	result     string
	err        bool
}

func taskCreate(a chan Task) {
	timer := time.NewTimer(time.Second * 10)
	defer timer.Stop()
	defer close(a)
	for {
		select {
		case <-timer.C:
			return
		default:
			ft := time.Now()
			a <- Task{
				id:         int(time.Now().Unix()),
				createTime: ft,
				err:        ft.Nanosecond()%2 > 0,
			} // передаем таск на выполнение
		}
	}
}

func taskWorker(a Task) Task {
	if time.Since(a.createTime) > time.Second*20 || a.err {
		a.err = true
		a.result = "something went wrong"
	} else {
		a.result = "task has been successed"
	}
	a.finishTime = time.Now()

	time.Sleep(time.Millisecond * 150)
	return a
}

func taskSorter(a Task, done, undone chan Task) {
	if a.err {
		undone <- a
	} else {
		done <- a
	}
}

func main() {

	superChan := make(chan Task, 10)
	go taskCreate(superChan)
	doneTasks := make(chan Task)
	undoneTasks := make(chan Task)
	defer close(doneTasks)
	defer close(undoneTasks)
	result := make([]Task, 0)
	errors := make([]Task, 0)
	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			fmt.Println("\nErrors:")
			for _, r := range errors {
				fmt.Printf("%v ", r.id) // печатаем ошибки(r)
			}
			println("\nDone tasks:")
			for _, r := range result {
				fmt.Printf("%v ", r.id)
			}
		case task, ok := <-superChan:
			task1 := taskWorker(task)
			go taskSorter(task1, doneTasks, undoneTasks)
			if !ok {
				return
			}
		case suc := <-doneTasks:
			result = append(result, suc)
		case err := <-undoneTasks:
			errors = append(errors, err)
		}
	}
}
