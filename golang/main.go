package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Приложение эмулирует получение и обработку неких тасков. Пытается и получать, и обрабатывать в многопоточном режиме.
// После обработки тасков в течении 3 секунд приложение должно выводить накопленные к этому моменту успешные таски и отдельно ошибки обработки тасков.

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

type TaskManager struct {
	superChan   chan Ttype
	doneTasks   chan Ttype
	undoneTasks chan Ttype
	result      map[int]Ttype
	err         map[int]Ttype
	semaphore   chan struct{}
	mu          sync.Mutex
}

func NewTaskManager(seconds int) *TaskManager {
	taskManager := &TaskManager{
		superChan:   make(chan Ttype),
		doneTasks:   make(chan Ttype),
		undoneTasks: make(chan Ttype),
		result:      map[int]Ttype{},
		err:         map[int]Ttype{},
		semaphore:   make(chan struct{}, 5),
	}

	go taskManager.taskCreturer()
	go taskManager.taskStarter()
	taskManager.resultReader()

	time.Sleep(time.Second * time.Duration(seconds))

	return taskManager
}

func (t *TaskManager) taskCreturer() {
	for {
		t.semaphore <- struct{}{}

		var err []byte
		ft := time.Now().Format(time.RFC3339)
		if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
			err = []byte("some error occured")
		}
		t.superChan <- Ttype{id: int(time.Now().Unix()), cT: ft, taskRESULT: err} // передаем таск на выполнение

		<-t.semaphore
	}
}

func (t *TaskManager) taskStarter() {
	// получение тасков
	for task := range t.superChan {
		t.semaphore <- struct{}{}
		go func(a Ttype) {
			defer func() { <-t.semaphore }()
			t.taskSorter(taskWorker(a))
		}(task)
	}

	close(t.superChan)
}

func (t *TaskManager) resultReader() {
	//  собираем таски
	go func() {
		for r := range t.doneTasks {
			t.mu.Lock()
			t.result[r.id] = r
			t.mu.Unlock()
		}
		close(t.doneTasks)
	}()

	go func() {
		for r := range t.undoneTasks {
			t.mu.Lock()
			t.err[r.id] = r
			t.mu.Unlock()
		}
		close(t.undoneTasks)
	}()
}

func (t *TaskManager) taskSorter(a Ttype) {
	if string(a.taskRESULT[14:]) == "successed" {
		t.doneTasks <- a
	} else {
		t.undoneTasks <- a
	}
}

func taskWorker(a Ttype) Ttype {
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(800))) // имитация выполнения каких-то действий

	tt, _ := time.Parse(time.RFC3339, a.cT)

	if tt.Before(time.Now().Add(-10 * time.Millisecond)) {
		time.Sleep(time.Millisecond * 10)
	}

	if len(a.taskRESULT) == 0 {
		if tt.After(time.Now().Add(-2 * time.Second)) {
			a.taskRESULT = []byte("task has been successed")
		} else {
			a.taskRESULT = []byte("something went wrong")
		}
	}
	a.fT = time.Now().Format(time.RFC3339Nano)

	return a
}

func (t *TaskManager) Print() {
	// печатаем результаты
	go func() {
		t.semaphore <- struct{}{}
		defer func() { <-t.semaphore }()

		if len(t.err) != 0 {
			println("Errors:")
			printRes(t.err, &t.mu)
			fmt.Println()
		}

		if len(t.result) != 0 {
			println("Done tasks:")
			printRes(t.result, &t.mu)
			fmt.Println()
		}
	}()
}

func printRes(res map[int]Ttype, mu *sync.Mutex) {
	mu.Lock()

	for _, r := range res {
		println("Task id:", r.id, "create time:", r.cT,
			"finish time:", r.fT, string(r.taskRESULT))
	}
	mu.Unlock()
}

func main() {
	lifeCycleOnSeconds := 3

	for {
		t := NewTaskManager(lifeCycleOnSeconds) // передаем время жизненного цикла параметром
		t.Print()
	}
}
