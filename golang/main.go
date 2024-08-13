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

// A Task represents a meaninglessness of our life
type Task struct {
	id           int
	creationTime string // время создания
	finishTime   string // время выполнения
	result       []byte
	successed    bool
}

func main() {
	tasks := make(chan Task, 10)
	doneTasks := make(chan Task)
	undoneTasks := make(chan error)

	// Генерация тасков
	go GenerateTasksRoutine(tasks)

	// Обработка тасков
	go ProcessTasksRoutine(tasks, doneTasks, undoneTasks)

	// Вывод тасков в консоль
	go PrintTasksRoutine(doneTasks, undoneTasks)

	// Предотвращение от завершения
	fmt.Scanln()

	// Закрытие каналов
	close(doneTasks)
	close(undoneTasks)
	close(tasks)
}

func GenerateTasksRoutine(tasks chan<- Task) {
	routineStartTime := time.Now()
	fmt.Println("Generation started")
	for {
		creationTime := time.Now().Format(time.RFC3339)
		successed := true
		if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
			// Примечание выполняющего: условие не сработает, так как Nanosecond() всегда
			// заканчивается на "00" (у меня, по крайней мере).
			// Соответственно, и остаток от деления на 2 всегда будет выводить 0.
			// По заданию нельзя менять логику появления ошибок, так что оставил как есть.
			creationTime = "Some error occured"
			successed = false
		}
		tasks <- Task{creationTime: creationTime, id: int(time.Now().Unix()), successed: successed} // передаем таск на выполнение

		if time.Since(routineStartTime).Seconds() > 10 { // Закончить генерацию тасков спустя 10 секунд
			fmt.Println("Generation ended")
			break
		}
	}
}

func ProcessTasksRoutine(tasks chan Task, doneTasks chan<- Task, undoneTasks chan<- error) {
	for t := range tasks {
		t = ProcessTask(t)
		go SortTask(t, doneTasks, undoneTasks)
	}
}

func PrintTasksRoutine(doneTasks <-chan Task, undoneTasks <-chan error) {
	result := map[int]Task{}
	errors := []error{}
	go func() {
		for r := range doneTasks {
			result[r.id] = r
		}
	}()
	go func() {
		for r := range undoneTasks {
			errors = append(errors, r)
		}
	}()

	// Печать в консоль
	for {
		time.Sleep(time.Second * 3)

		fmt.Println("Errors:")
		for _, e := range errors {
			fmt.Println(e)
		}

		fmt.Println("Done tasks:")
		for _, v := range result {
			fmt.Printf("Task id %d time %s, %s\n", v.id, v.creationTime, v.result)
		}
	}
}

func ProcessTask(task Task) Task {
	tt, _ := time.Parse(time.RFC3339, task.creationTime)
	if tt.After(time.Now().Add(-20 * time.Second)) {
		task.result = []byte("task has been successed")
	} else {
		task.result = []byte("something went wrong")
	}
	task.finishTime = time.Now().Format(time.RFC3339Nano)

	time.Sleep(time.Millisecond * 150)

	return task
}

func SortTask(t Task, doneTasks chan<- Task, undoneTasks chan<- error) {
	if t.successed {
		doneTasks <- t
	} else {
		undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.creationTime, t.result)
	}
}
