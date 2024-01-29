package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Ttype struct {
	Id         int64
	CreatedAt  time.Time // время создания
	FinishedAt time.Time // время выполнения
	TaskResult string
	Error      error
}

const (
	TASK_RESULT_SUCCESS   = "task has been successed"
	TASK_RESULT_UNSUCCESS = "something went wrong"
	CHAN_SIZE             = 10
)

type TaskWorker struct {
	mu *sync.Mutex
	wg *sync.WaitGroup

	newTasks    chan Ttype
	doneTasks   chan Ttype
	undoneTasks chan Ttype

	completeTasks map[int64]Ttype
	errorTasks    map[int64]Ttype
}

func NewTaskWorker() *TaskWorker {
	return &TaskWorker{
		mu:          &sync.Mutex{},
		wg:          &sync.WaitGroup{},
		doneTasks:   make(chan Ttype),
		undoneTasks: make(chan Ttype),
		newTasks:    make(chan Ttype, CHAN_SIZE),
	}
}

func (tw *TaskWorker) taskCreturer(ctx context.Context) {
	go func() {
		for {
			select {
			default:
				tw.wg.Add(1)

				var err error
				createdAt := time.Now()

				// вместо time.Now().Nanoseconds()
				id := createRandomId()

				if id%2 > 0 { // вот такое условие появления ошибочных тасков
					err = fmt.Errorf("Some error occured")
				}
				tw.newTasks <- Ttype{CreatedAt: createdAt, Id: int64(id), Error: err} // передаем таск на выполнение

				// охладить пыл по созданию тасок
				time.Sleep(time.Millisecond * 50)

			case <-ctx.Done():
				return
			}

		}
	}()
}

// time.Now().Nanosecond() всегда генерировал числа которые делятся на 2 без остатка, поэтому решил немного их рандомить
func createRandomId() int {
	rand.NewSource(time.Now().UnixNano())
	min := 10
	max := 30
	r := rand.Intn(max-min+1) + min

	return time.Now().Nanosecond() + r
}

func (tw *TaskWorker) taskHandler(task Ttype) {
	if task.Error == nil && task.CreatedAt.After(time.Now().Add(-20*time.Second)) {
		task.TaskResult = TASK_RESULT_SUCCESS
		task.FinishedAt = time.Now()
		tw.doneTasks <- task
	} else {
		task.TaskResult = TASK_RESULT_UNSUCCESS
		task.FinishedAt = time.Now()
		log.Printf("Task id %d time %s, error %s", task.Id, task.CreatedAt, task.TaskResult)
		tw.undoneTasks <- task
	}

}

func (tw *TaskWorker) createTaskHandlers(handlerCount int64) {
	for i := 0; i < int(handlerCount); i++ {
		go func() {
			// получение тасков
			for t := range tw.newTasks {
				task := t
				go tw.taskHandler(task)
			}
		}()
	}
}

func (tw *TaskWorker) taskOutput() {
	go func() {
		for r := range tw.doneTasks {
			tw.mu.Lock()
			tw.completeTasks[r.Id] = r
			tw.mu.Unlock()

			tw.wg.Done()
		}

	}()

	go func() {
		for r := range tw.undoneTasks {
			tw.mu.Lock()
			tw.errorTasks[r.Id] = r
			tw.mu.Unlock()

			tw.wg.Done()
		}

	}()
}

func (tw *TaskWorker) Start(
	ctxTimeout context.Context,
	handlerCount int64,
) {
	ctx, stop := signal.NotifyContext(ctxTimeout, os.Interrupt, syscall.SIGTERM)
	defer stop()

	// запускаем процесс создания тасок
	tw.taskCreturer(ctx)

	// создаем n обработчиков тасок
	tw.createTaskHandlers(handlerCount)

	// запускаем процесс записи тасок в мапы
	tw.taskOutput()

	// подобие graceful shutdown
	func() {
		<-ctx.Done()
		fmt.Println("APP CLOSING")
		close(tw.newTasks)
		tw.wg.Wait()
	}()

	println("Error tasks:")
	for r := range tw.errorTasks {
		println(r)
	}

	println("Done tasks:")
	for r := range tw.completeTasks {
		println(r)
	}
}

func main() {
	// Таймаут через которое приложение закроется
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tw := NewTaskWorker()
	tw.Start(ctx, CHAN_SIZE)
}
