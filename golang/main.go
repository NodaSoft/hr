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

// A Ttype represents a meaninglessness of our life
type Ttype struct {
	id         int
	cT         string // время создания
	fT         string // время выполнения
	taskRESULT []byte
}

func taskCreturer(a chan Ttype, wg *sync.WaitGroup) {
	defer wg.Done()
	start := time.Now()

	for {
		if time.Since(start) >= 10*time.Second {
			close(a)
			return
		}
		creationTime := time.Now().Format(time.RFC3339)
		if time.Now().Nanosecond()%2 > 0 {
			creationTime = "Some error occurred"
		}
		task := Ttype{id: int(time.Now().Unix()), cT: creationTime}
		a <- task
		time.Sleep(300 * time.Millisecond)
	}
}

func task_worker(a Ttype) Ttype {
	tt, _ := time.Parse(time.RFC3339, a.cT)
	if tt.After(time.Now().Add(-20 * time.Second)) {
		a.taskRESULT = []byte("task has been successed")
	} else {

		a.taskRESULT = []byte("something went wrong")
	}
	a.fT = time.Now().Format(time.RFC3339Nano)

	time.Sleep(time.Millisecond * 150)

	return a
}

func task_work(tasksChan <-chan Ttype, doneTasks chan<- Ttype, undoneTasks chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	// получение тасков
	for t := range tasksChan {
		t = task_worker(t)
		if string(t.taskRESULT[14:]) == "successed" {

			doneTasks <- t
		} else {
			undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT, t.taskRESULT)
		}
	}

	close(doneTasks)
	close(undoneTasks)
}

func tasksorter(t Ttype, doneTasks chan<- Ttype, undoneTasks chan<- error, wg *sync.WaitGroup) {
	if string(t.taskRESULT[14:]) == "successed" {
		doneTasks <- t
	} else {
		undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT, t.taskRESULT)
	}
}

func chanPrinter(doneTasks <-chan Ttype, undoneTasks <-chan error) {
	fmt.Println("Done tasks:")
	done := true
	for done {
		select {
		case task, ok := <-doneTasks:
			if ok {
				fmt.Println(task.id)
			} else {
				done = false
			}
		default:
			done = false
		}
	}
	fmt.Println("Error tasks:")
	errorDone := true
	for errorDone {
		select {
		case task, ok := <-undoneTasks:
			if ok {
				fmt.Println(task.Error())
			} else {
				errorDone = false
			}
		default:
			errorDone = false
		}
	}
}

func periodPrinter(doneTasks <-chan Ttype, undoneTasks <-chan error, wg *sync.WaitGroup) {
	defer wg.Done()

	//tickre := time.NewTicker(3 * time.Second)
	timer := time.NewTimer(3 * time.Second)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			chanPrinter(doneTasks, undoneTasks)
		case <-time.After(10*time.Second + 3*time.Second):
			chanPrinter(doneTasks, undoneTasks)
			return
		}
	}
}

func main() {

	var wg sync.WaitGroup

	superChan := make(chan Ttype)
	doneTasks := make(chan Ttype, 100)
	undoneTasks := make(chan error, 100)

	//Create Tasks by 10 seconds
	wg.Add(1)
	go taskCreturer(superChan, &wg)

	//Work width tasks
	wg.Add(1)
	go task_work(superChan, doneTasks, undoneTasks, &wg)

	//Print Tasks by 3 seconds
	wg.Add(1)
	go periodPrinter(doneTasks, undoneTasks, &wg)

	wg.Wait()
}
