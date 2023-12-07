package main

import (
	"context"
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

// A Ttype represents a meaninglessness of our life

const (
	mainChanSize = 10
	execTime     = 3
)

type tTask struct {
	id                int
	createTime        time.Time // время создания
	finishTime        time.Time // время выполнения
	taskResult        string
	createTimeInvalid bool
}

type chanExch struct {
	inCh      chan tTask
	successCh chan tTask
	errorCh   chan tTask
	eTask     []tTask
	sTask     []tTask
}

func (ce *chanExch) CreateTasks(ctx context.Context) {

	for {

		select {
		case <-ctx.Done():
			close(ce.inCh)
			return

		default:
			//fmt.Println(time.Now())
			//time.Sleep(time.Second)
			ce.inCh <- tTask{
				id:                int(time.Now().Unix()),
				createTime:        time.Now(),
				taskResult:        "",
				createTimeInvalid: time.Now().Nanosecond()%2 > 0,
			}
		}
	}
}

func (ce *chanExch) workTasks() {
	defer close(ce.successCh)
	defer close(ce.errorCh)
	for msg := range ce.inCh {

		msg.taskResult = "something went wrong"
		resultChan := ce.errorCh
		// проверка условия msg.createTime.After(time.Now().Add(-20*time.Second)) убрана, т.к. по условию программа выполняется всего 3 секунды,
		// сообщений в канале старше 20 секунд появится не может в принципе при данных условиях, всегда возвращает true
		if msg.createTimeInvalid == false {
			msg.taskResult = "task has been successed"
			resultChan = ce.successCh

		}
		msg.finishTime = time.Now()
		time.Sleep(time.Millisecond * 150)
		resultChan <- msg

	}
}

func (ce *chanExch) successTasks(wg *sync.WaitGroup) {
	defer wg.Done()
	for msg := range ce.successCh {
		ce.sTask = append(ce.sTask, msg)
	}
}

func (ce *chanExch) errorTasks(wg *sync.WaitGroup) {
	defer wg.Done()
	for msg := range ce.errorCh {
		ce.eTask = append(ce.eTask, msg)
	}
}

func main() {

	channels := chanExch{
		inCh:      make(chan tTask, mainChanSize),
		successCh: make(chan tTask),
		errorCh:   make(chan tTask),
	}

	ctx, cancel := context.WithTimeout(context.Background(), execTime*time.Second)
	defer cancel()

	go channels.CreateTasks(ctx)
	go channels.workTasks()

	var wg sync.WaitGroup
	wg.Add(1)
	go channels.successTasks(&wg)
	wg.Add(1)
	go channels.errorTasks(&wg)
	wg.Wait()

	fmt.Println("Success task:")
	for _, val := range channels.sTask {
		fmt.Println(val.id, val.createTime.String(), val.finishTime.String(), val.taskResult, val.createTimeInvalid)
	}

	fmt.Println("Error task:")
	for _, val := range channels.eTask {
		fmt.Println(val.id, val.createTime.String(), val.finishTime.String(), val.taskResult, val.createTimeInvalid)
	}

	return
}
