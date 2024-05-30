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
	Id         int
	CreateTime time.Time // время создания
	FinishTime time.Time // время выполнения
	Error      error
}

// эмуляция бесконечного получения(создания), обработки и вывода тасков
func main() {
	// канал полученных тасок
	rawTask := make(chan *Task)
	// канал обработанных тасок
	completedTask := make(chan *Task)

	// создатель тасок
	go taskCreator(rawTask)
	// обработчик тасок
	go taskHandler(rawTask, completedTask)
	// вывод обработанных тасок / ошибок
	go completedTaskPrinter(completedTask)

	for {
	}
}

func taskCreator(res chan *Task) {
	// имитация получения тасок каждые 0.5 секунд
	for {
		time.Sleep(500 * time.Millisecond)
		// для уникальности Id можно использовать пакет "github.com/google/uuid"
		task := &Task{CreateTime: time.Now(), Id: int(time.Now().Unix())}
		if time.Now().Nanosecond()%2 > 0 { // условие появления ошибочных тасков
			task.Error = fmt.Errorf("Some error occurred")
		}
		// кладем таску в канал для дальнейшей обработки
		res <- task
	}
}

func taskHandler(rawTask chan *Task, completedTask chan *Task) {
	for t := range rawTask {
		// не перезаписываем ошибку, если она была
		if !t.CreateTime.After(time.Now().Add(-20*time.Second)) && t.Error == nil {
			t.Error = fmt.Errorf("something went wrong")
		}
		t.FinishTime = time.Now()
		completedTask <- t
	}
}

func completedTaskPrinter(resultTask chan *Task) {
	// сигнальный канал для начала вывода значений
	signal := time.Tick(3 * time.Second)
	// слайсы для хранения значений
	errors := make([]string, 0)
	result := make([]*Task, 0)

	for {
		select {
		case <-signal:
			// вывод значений в отдельной горутине
			go printResults(result, errors)
			// очистка слайсов
			errors = make([]string, 0)
			result = make([]*Task, 0)
		case t := <-resultTask:
			// у ошибочных тасок записываем только ошибки в отдельный слайс
			if t.Error != nil {
				errors = append(errors, t.Error.Error())
			} else {
				result = append(result, t)
			}
		}
	}
}

func printResults(res []*Task, errors []string) {
	for _, v := range errors {
		fmt.Println(v)
	}
	for _, v := range res {
		fmt.Println(v.Id)
	}
}
