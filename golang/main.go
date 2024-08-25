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

// Мы даем тестовое задание чтобы:
// * уменьшить время технического собеседования - лучше вы потратите пару часов в спокойной домашней обстановке, чем будете волноваться, решая задачи под взором наших ребят;
// * увеличить вероятность прохождения испытательного срока - видя сразу стиль и качество кода, мы можем быть больше уверены в выборе;
// * снизить число коротких собеседований, когда мы отказываем сразу же.

// Выполнение тестового задания не гарантирует приглашение на собеседование, т.к. кроме качества выполнения тестового задания, оцениваются и другие показатели вас как кандидата.

// Мы не даем комментариев по результатам тестового задания. Если в случае отказа вам нужен наш комментарий по результатам тестового задания, то просим об этом написать вместе с откликом.

// A Ttype represents a meaninglessness of our life
type Ttype struct {
	id     int
	cT     string // время создания
	fT     string // время выполнения
	result []byte
}

func main() {
	tasks := make(chan Ttype, 10)
	defer close(tasks)

	go CreateTasks(tasks)

	results := make(chan Ttype, 10)
	defer close(results)

	n := 2 //кол-во горутин, занимающихся выполнением таска
	for i := 0; i < n; i++ {
		go WorkTask(tasks, results)
	}

	doneTasks := make(chan Ttype)
	defer close(doneTasks)

	undoneTasks := make(chan error)
	defer close(undoneTasks)

	go SortResult(results, doneTasks, undoneTasks)

	result := make(map[int]Ttype)
	errors := make([]error, 0)

	go WriteResults(doneTasks, undoneTasks, result, &errors)

	time.Sleep(time.Second * 3)

	println("Errors:")
	for r := range errors {
		println(r)
	}

	println("Done tasks:")
	for r := range result {
		println(r)
	}
}

func CreateTasks(a chan<- Ttype) {
	go func() {
		for {
			ft := time.Now().Format(time.RFC3339)
			if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
				ft = "Some error occured"
			}
			a <- Ttype{cT: ft, id: int(time.Now().Unix())} // передаем таск на выполнение
		}
	}()
}

func WorkTask(tasks <-chan Ttype, results chan<- Ttype) {
	for task := range tasks {
		tt, _ := time.Parse(time.RFC3339, task.cT)

		if tt.After(time.Now().Add(-20 * time.Second)) {
			task.result = []byte("task has been successed")
		} else {
			task.result = []byte("something went wrong")
		}

		task.fT = time.Now().Format(time.RFC3339Nano)

		time.Sleep(time.Millisecond * 150)

		results <- task
	}
}

func SortResult(results <-chan Ttype, doneTasks chan<- Ttype, undoneTasks chan<- error) {
	for res := range results {
		if string(res.result[14:]) == "successed" {
			doneTasks <- res
		} else {
			undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", res.id, res.cT, res.result)
		}
	}
}

func WriteResults(doneTasks <-chan Ttype, undoneTasks <-chan error, result map[int]Ttype, errors *[]error) {
	for {
		select {
		case doneTask := <-doneTasks:
			result[doneTask.id] = doneTask
		case undoneTask := <-undoneTasks:
			*errors = append(*errors, undoneTask)
		}
	}
}
