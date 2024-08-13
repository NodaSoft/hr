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
	id  int
	cT  string // время создания
	fT  string // время выполнения
	res []byte
}

func main() {
	taskCreator := func(taskChan chan Task) {
		defer close(taskChan)
		stop := time.NewTimer(10 * time.Second)

		for {
			select {
			case <-stop.C:
				return
			default:
				now := time.Now().Format(time.RFC3339)
				if time.Now().Nanosecond()/100%2 > 0 { // вот такое условие появления ошибочных тасков
					now = "Some error occurred"
				}
				taskChan <- Task{cT: now, id: int(time.Now().Unix())} // передаем таск на выполнение
			}
		}
	}

	taskChan := make(chan Task)
	go taskCreator(taskChan)

	taskWorker := func(task Task) Task {
		tt, err := time.Parse(time.RFC3339, task.cT)
		if err != nil || !tt.After(time.Now().Add(-10*time.Second)) {
			task.res = []byte("something went wrong")
		} else {
			task.res = []byte("task has been success")
		}

		task.fT = time.Now().Format(time.RFC3339Nano)
		time.Sleep(time.Millisecond * 150)

		return task
	}

	doneTasks := make(chan Task)
	undoneTasks := make(chan error)

	taskSorter := func(task Task) {
		if string(task.res) == "task has been success" {
			doneTasks <- task
		} else {
			undoneTasks <- fmt.Errorf("Task id: %d, time: %s, error: %s", task.id, task.cT, task.res)
		}
	}

	go func() {
		// получение тасков
		defer close(doneTasks)
		defer close(undoneTasks)

		for t := range taskChan {
			processedTask := taskWorker(t)
			go taskSorter(processedTask)
		}
	}()

	var results []Task
	go func() {
		for task := range doneTasks {
			results = append(results, task)
		}
	}()

	var arrOfErr []error
	go func() {
		for task := range undoneTasks {
			arrOfErr = append(arrOfErr, task)
		}
	}()

	done := make(chan struct{})
	go func() {
		defer close(done)

		stop := time.NewTimer(time.Second * 10)

		tick := time.NewTicker(time.Second * 3)
		defer tick.Stop()

		for {
			select {
			case <-stop.C:
				return
			case <-tick.C:
				fmt.Println("Errors:")
				for _, err := range arrOfErr {
					fmt.Println(err)
				}

				fmt.Println("Done tasks:")
				for _, res := range results {
					fmt.Printf("id: %d, ct: %s, ft: %s, res: %s\n", res.id, res.cT, res.fT, res.res)
				}
			}
		}
	}()
	<-done
}
