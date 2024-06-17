package main

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
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

const (
	currentTimeFormat = time.RFC3339Nano
)

const (
	printInterval   = 3 * time.Second
	workingInterval = 10 * time.Second
)

const (
	taskLifeTime = 20 * time.Second
)

const (
	sleepInterval = 150 * time.Millisecond
)

var (
	badTaskError     = errors.New("it very bad time")
	veryOldTaskError = errors.New("task very old")
)

type Task struct {
	ID            int64
	Error         error
	CreateTime    string
	CompletedTime string
}

type App struct {
	newTaskQueue chan *Task
	doneTasks    chan *Task
	undoneTasks  chan *Task
	stopTimer    *time.Timer
	run          atomic.Bool
	wg           sync.WaitGroup
}

func main() {
	a := &App{}
	a.Run()
}

func (a *App) Run() {
	a.wg = sync.WaitGroup{}
	a.newTaskQueue = make(chan *Task, 10)
	a.doneTasks = make(chan *Task)
	a.undoneTasks = make(chan *Task)

	a.run.Store(true)

	a.startTaskProvider()
	a.startTaskWorker()
	a.startTaskPrinter()

	a.stopTimer = time.NewTimer(workingInterval)
	for a.run.Load() {
		select {
		case <-a.stopTimer.C:
			a.run.Store(false)
		}
	}

	a.wg.Wait()
}

func (a *App) startTaskPrinter() {
	a.wg.Add(1)
	go func() { // todo возможно, нужно было разделить вывод ошибок и успехов в разные горутины...
		ticker := time.NewTicker(printInterval)
		doneChannelClosed := false
		undoneChannelClosed := false
		for !doneChannelClosed || !undoneChannelClosed {
			doneTasks := make([]*Task, 0, 10)
			undoneTasks := make([]*Task, 0, 10)

			collectTasks := true
			for collectTasks {
				select {
				case doneTask, ok := <-a.doneTasks:
					if !ok {
						doneChannelClosed = true
						break
					}

					doneTasks = append(doneTasks, doneTask)
				case undoneTask, ok := <-a.undoneTasks:
					if !ok {
						undoneChannelClosed = true
						break
					}

					undoneTasks = append(undoneTasks, undoneTask)
				case <-ticker.C:
					collectTasks = false
				}
			}

			if len(doneTasks) == 0 && len(undoneTasks) == 0 {
				continue
			}

			fmt.Printf("------------%s------------\n", time.Now().Format(currentTimeFormat))

			if len(undoneTasks) != 0 {
				fmt.Println("Errors: ")
				for _, t := range undoneTasks {
					fmt.Printf("ID: %v, create time: %s, error: %s\n", t.ID, t.CreateTime, t.Error)
				}

				fmt.Println("")
			}

			if len(doneTasks) != 0 {
				fmt.Println("Done tasks: ")
				for _, t := range doneTasks {
					fmt.Printf("ID: %v, create time: %s, completed time: %s\n", t.ID, t.CreateTime, t.CompletedTime)
				}
			}

			fmt.Printf("------------------------------------\n\n")
		}

		ticker.Stop()
		a.wg.Done()
	}()
}

func (a *App) startTaskWorker() {
	a.wg.Add(1)
	go func() {
		for {
			task, ok := <-a.newTaskQueue
			if !ok {
				break
			}

			a.handleTask(task)
			a.sortTask(task)
		}

		close(a.doneTasks)
		close(a.undoneTasks)

		a.wg.Done()
	}()
}

func (a *App) startTaskProvider() {
	a.wg.Add(1)
	go func() {
		var nextTaskID int64 = 0
		for a.run.Load() {
			var err error
			createTime := time.Now().Format(currentTimeFormat)
			if time.Now().Nanosecond()%2 > 0 {
				err = badTaskError
			}

			task := Task{
				ID:         nextTaskID,
				Error:      err,
				CreateTime: createTime,
			}

			a.newTaskQueue <- &task
			nextTaskID++
		}

		close(a.newTaskQueue)
		a.wg.Done()
	}()
}

func (a *App) handleTask(task *Task) {
	if task.Error != nil {
		return
	}

	createTime, _ := time.Parse(currentTimeFormat, task.CreateTime)
	if time.Since(createTime) > taskLifeTime {
		task.Error = veryOldTaskError
	}
	task.CompletedTime = time.Now().Format(currentTimeFormat)

	time.Sleep(sleepInterval)
}

func (a *App) sortTask(task *Task) {
	if task.Error != nil {
		a.undoneTasks <- task
	} else {
		a.doneTasks <- task
	}
}
