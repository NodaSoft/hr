package main

import (
	"errors"
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

// A Task represents a meaninglessness of our life
type Task struct {
	id          int
	time_start  time.Time // время создания
	time_finish time.Time // время выполнения
	err         error
}

func taskCreator(tasks chan Task) {

	// параметры генератора
	const duration_secs = 10
	const interval_secs = 0.5

	start := time.Now()
	id := 0
	for time.Since(start).Seconds() < duration_secs {
		var task Task

		t := time.Now()
		task.id = id
		task.time_start = t
		if t.Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
			task.err = errors.New("Task creator error")
		}

		tasks <- task
		id++
		time.Sleep(time.Duration(float64(time.Second) * interval_secs))
	}
	close(tasks)
}

func taskLogger(successful_tasks *[]Task, failed_tasks *[]Task, stop_logger chan struct{}) {

	var log = func() {

		fmt.Println("Successful tasks: ", len(*successful_tasks))
		for _, t := range *successful_tasks {
			fmt.Printf("%-4d start: %s finish: %s\n", t.id, t.time_start.Format(time.RFC3339), t.time_finish.Format(time.RFC3339))
		}
		fmt.Println()

		fmt.Println("Failed tasks: ", len(*failed_tasks))
		for _, t := range *failed_tasks {
			fmt.Printf("%-4d start: %s error: %s\n", t.id, t.time_start.Format(time.RFC3339), t.err.Error())
		}
	}

	t := time.Now()
	for {
		select {
		case <-stop_logger:
			return
		default:
			if time.Since(t) > time.Second*3 {
				log()
				t = time.Now()
			}
		}
	}
}

func taskWorker(tasks chan Task, successful_tasks *[]Task, failed_tasks *[]Task, finished chan struct{}) {

	stop_logger := make(chan struct{})
	go taskLogger(successful_tasks, failed_tasks, stop_logger)

	for t := range tasks {
		if t.err != nil {
			*failed_tasks = append(*failed_tasks, t)
		} else {
			if rand.Intn(2) != 0 {
				t.err = errors.New("Task worker error") // условие ошибки в обработчике
				*failed_tasks = append(*failed_tasks, t)
			} else {
				t.time_finish = time.Now()
				*successful_tasks = append(*successful_tasks, t)
			}

			time.Sleep(time.Millisecond * 150)
		}
	}
	stop_logger <- struct{}{}
	time.Sleep(time.Millisecond * 100)

	finished <- struct{}{}
}

func main() {

	tasks := make(chan Task, 10)
	successful_tasks := []Task{}
	failed_tasks := []Task{}
	finished := make(chan struct{})

	// генератор задач
	go taskCreator(tasks)

	// обрабочик задач
	go taskWorker(tasks, &successful_tasks, &failed_tasks, finished)

	// ожидание завершения обработчика
	<-finished
}
