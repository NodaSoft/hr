package main

import (
	"fmt"
	"sync"
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

type Ttype struct {
	id         int
	cT         string // время создания
	fT         string // время выполнения
	taskRESULT []byte
}

func main() {
	superChan := make(chan Ttype, 10)
	doneTasks := make(chan Ttype, 10)
	undoneTasks := make(chan error, 10)
	done := make(chan struct{})

	var wg sync.WaitGroup

	go taskCreator(superChan)

	go func() {
		for t := range superChan {
			wg.Add(1)
			go func(task Ttype) {
				defer wg.Done()
				processedTask := taskWorker(task)
				taskSorter(processedTask, doneTasks, undoneTasks)
			}(t)
		}

		wg.Wait()
		close(doneTasks)
		close(undoneTasks)
		close(done)
	}()

	go func() {
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				fmt.Println("Errors:")
				errorsProcessed := false
				for !errorsProcessed {
					select {
					case undoneTask, ok := <-undoneTasks:
						if !ok {
							errorsProcessed = true
						} else {
							fmt.Println(undoneTask.Error())
						}
					default:
						errorsProcessed = true
					}
				}

				fmt.Println("Done tasks:")
				doneTasksProcessed := false
				for !doneTasksProcessed {
					select {
					case doneTask, ok := <-doneTasks:
						if !ok {
							doneTasksProcessed = true
						} else {
							fmt.Printf("Task id %d creation time %s, finished at %s\n", doneTask.id, doneTask.cT, doneTask.fT)
						}
					default:
						doneTasksProcessed = true
					}
				}
			case <-done:
				return
			}
		}
	}()

	<-done
}

func taskCreator(a chan Ttype) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	stop := time.After(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			ft := time.Now().Format(time.RFC3339)
			if time.Now().Nanosecond()%2 > 0 { // условие для генерации ошибок
				ft = "Some error occurred"
			}
			a <- Ttype{cT: ft, id: int(time.Now().Unix())} // передаем таск на выполнение
		case <-stop:
			close(a)
			return
		}
	}
}

func taskSorter(t Ttype, doneTasks chan Ttype, undoneTasks chan error) {
	if string(t.taskRESULT) == "task has been succeeded" {
		doneTasks <- t
	} else {
		undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT, t.taskRESULT)
	}
}

func taskWorker(a Ttype) Ttype {
	tt, err := time.Parse(time.RFC3339, a.cT)
	if err != nil || tt.After(time.Now().Add(-20*time.Second)) {
		a.taskRESULT = []byte("task has been succeeded")
	} else {
		a.taskRESULT = []byte("something went wrong")
	}
	a.fT = time.Now().Format(time.RFC3339Nano)

	time.Sleep(time.Millisecond * 150)

	return a
}
