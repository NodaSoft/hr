package main

import (
	"fmt"
	"sync"
	"time"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнены остальных тасков

type TaskResultStatus int

const (
	TaskResultStatusAccepted TaskResultStatus = iota
	TaskResultStatusRejected
)

// A Ttype represents a meaninglessness of our life
type Ttype struct {
	id               int
	cT               string // время создания
	fT               string // время выполнения
	err              string
	taskResultMsg    []byte
	taskResultStatus TaskResultStatus
}

type ConcurrentErrSlice struct {
	sync.RWMutex
	items []error
}

type ConcurrentErrSliceItem struct {
	Index int
	Value error
}

func (cs *ConcurrentErrSlice) Append(item error) {
	cs.Lock()
	defer cs.Unlock()

	cs.items = append(cs.items, item)
}

func (cs *ConcurrentErrSlice) Iter() <-chan ConcurrentErrSliceItem {
	c := make(chan ConcurrentErrSliceItem)

	f := func() {
		cs.Lock()
		defer cs.Unlock()
		for index, value := range cs.items {
			c <- ConcurrentErrSliceItem{index, value}
		}
		close(c)
	}
	go f()

	return c
}

func main() {
	superChan := make(chan Ttype, 10)
	wg := sync.WaitGroup{}
	superChanCancel := generateTasks(superChan, &wg)

	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)
	go func() {
		// получение тасков
		for t := range superChan {
			taskWorker(&t)
			taskSorter(t, doneTasks, undoneTasks)
		}
	}()

	result := sync.Map{}
	go handleDoneTasks(doneTasks, &wg, &result)

	err := ConcurrentErrSlice{}
	go handleUndoneTasks(undoneTasks, &wg, &err)

	executionTime := time.Second * 3
	time.Sleep(executionTime)
	superChanCancel()
	wg.Wait()
	close(doneTasks)
	close(undoneTasks)

	fmt.Println("Errors:")
	for item := range err.Iter() {
		fmt.Printf("%v\n", item.Value)
	}

	fmt.Println("Done tasks:")
	result.Range(func(key, value any) bool {
		fmt.Printf("%d\n", key.(int))
		return true
	})
}

func generateTasks(superChan chan Ttype, wg *sync.WaitGroup) func() {
	var cancel func()
	go func() {
		defer close(superChan)
		isCanceled := false
		cancel = func() {
			isCanceled = true
		}

		for !isCanceled {
			wg.Add(1)
			cT := time.Now().Format(time.RFC3339)
			var err string
			// TODO Возможно стоит заменить функцию Nanosecond на другую, например, UnixMicro,
			// 		так как функция Now не возвращает наносекунды.
			if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
				err = "Some error occured"
			}
			superChan <- Ttype{id: int(time.Now().Unix()), cT: cT, err: err} // передаем таск на выполнение
		}
	}()

	return cancel
}

func taskWorker(a *Ttype) {
	tt, _ := time.Parse(time.RFC3339, a.cT)
	if tt.After(time.Now().Add(-20*time.Second)) && len(a.err) == 0 {
		a.taskResultMsg = []byte("task has been successed")
		a.taskResultStatus = TaskResultStatusAccepted
	} else {
		a.taskResultMsg = []byte("something went wrong")
		a.taskResultStatus = TaskResultStatusRejected
	}
	a.fT = time.Now().Format(time.RFC3339Nano)

	// emulating of long-running code
	time.Sleep(time.Millisecond * 150)
}

func taskSorter(t Ttype, doneTasks chan Ttype, undoneTasks chan error) {
	if t.taskResultStatus == TaskResultStatusAccepted {
		doneTasks <- t
	} else {
		undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT, t.taskResultMsg)
	}
}

func handleDoneTasks(tasks chan Ttype, wg *sync.WaitGroup, result *sync.Map) {
	for task := range tasks {
		go func(task Ttype) {
			result.Store(task.id, task)
			wg.Done()
		}(task)
	}
}

func handleUndoneTasks(tasks chan error, wg *sync.WaitGroup, err *ConcurrentErrSlice) {
	for task := range tasks {
		go func(task error) {
			err.Append(task)
			wg.Done()
		}(task)
	}
}
