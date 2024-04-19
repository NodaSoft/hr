package main

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// * обновленный код отправить через pull-request.

// Приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме.
// Должно выводить успешные таски и ошибки по мере выполнения.
// Как видите, никаких привязок к внешним сервисам нет - полный карт-бланш на модификацию кода.

const (
	taskWorkerCount   = 4
	simulationSeconds = 3
)

// A task represents a meaninglessness of our life
type task struct {
	id          int
	createTime  time.Time
	processTime time.Duration
	err         error
}

// Abstract sequence of task identifiers
type taskIdSequence interface {
	NextId() int
}

// In-memory implementation of a sequence of task identifiers
type taskIdSequenceInMemory struct {
	id atomic.Int32
}

func newTaskIdSequence() taskIdSequence {
	// TODO: implement a persistent version or switch to uuid
	return &taskIdSequenceInMemory{}
}

func (s *taskIdSequenceInMemory) NextId() int {
	return int(s.id.Add(1))
}

// Task producer generates tasks and writes them to a channel
type taskProducer struct {
	tasks  chan<- task
	tidseq taskIdSequence
}

func newTaskProducer(tidseq taskIdSequence, tasks chan<- task) taskProducer {
	return taskProducer{
		tasks:  tasks,
		tidseq: tidseq,
	}
}

func (p *taskProducer) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			curtime := time.Now()
			if curtime.Nanosecond()%2 > 0 {
				// simulation of incorrect task
				curtime = time.Time{}
			}
			p.tasks <- task{
				id:         p.tidseq.NextId(),
				createTime: curtime,
			}
		}
	}
}

// Task worker processes tasks from a channel and writes them to another
type taskWorker struct {
	pendingTasks   <-chan task
	processedTasks chan<- task
}

func newTaskWorker(pendingTasks <-chan task, completedTasks chan<- task) taskWorker {
	return taskWorker{
		pendingTasks:   pendingTasks,
		processedTasks: completedTasks,
	}
}

func (c *taskWorker) run() {
	for task := range c.pendingTasks {
		if task.createTime.IsZero() {
			task.err = fmt.Errorf("task creation time is zero")
		} else {
			curtime := time.Now()
			if task.createTime.After(curtime.Add(-20 * time.Second)) {
				task.processTime = curtime.Sub(task.createTime)
			} else {
				task.err = fmt.Errorf("something went wrong")
			}
		}
		time.Sleep(time.Millisecond * 150)
		c.processedTasks <- task
	}
}

func main() {
	tidseq := newTaskIdSequence()
	tasks := make(chan task)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*simulationSeconds))
	defer cancel()
	go func() {
		producer := newTaskProducer(tidseq, tasks)
		producer.run(ctx)
		close(tasks)
	}()
	workerWaitGroup := sync.WaitGroup{}
	processedTasks := make(chan task)
	for i := 0; i < taskWorkerCount; i++ {
		workerWaitGroup.Add(1)
		go func() {
			defer workerWaitGroup.Done()
			worker := newTaskWorker(tasks, processedTasks)
			worker.run()
		}()
	}
	go func() {
		workerWaitGroup.Wait()
		close(processedTasks)
	}()
	successedTasks := make(chan task)
	failedTasks := make(chan task)
	taskSorterWaitGroup := sync.WaitGroup{}
	taskSorterWaitGroup.Add(1)
	go func() {
		defer taskSorterWaitGroup.Done()
		for t := range processedTasks {
			taskSorterWaitGroup.Add(1)
			go func(t task) {
				defer taskSorterWaitGroup.Done()
				if t.err == nil {
					successedTasks <- t
				} else {
					failedTasks <- t
				}
			}(t)
		}
	}()
	go func() {
		taskSorterWaitGroup.Wait()
		close(successedTasks)
		close(failedTasks)
	}()
	successedTasksMap := make(map[int]task)
	successedTasksMutex := sync.Mutex{}
	failedTasksArr := []task{}
	failedTasksMutex := sync.Mutex{}
	taskConsumerWaitGroup := sync.WaitGroup{}
outer:
	for {
		select {
		case t, ok := <-successedTasks:
			if !ok {
				break outer
			}
			taskConsumerWaitGroup.Add(1)
			go func(t task) {
				defer taskConsumerWaitGroup.Done()
				successedTasksMutex.Lock()
				successedTasksMap[t.id] = t
				successedTasksMutex.Unlock()
			}(t)
		case t, ok := <-failedTasks:
			if !ok {
				break outer
			}
			taskConsumerWaitGroup.Add(1)
			go func(t task) {
				defer taskConsumerWaitGroup.Done()
				failedTasksMutex.Lock()
				failedTasksArr = append(failedTasksArr, t)
				failedTasksMutex.Unlock()
			}(t)
		}
	}
	taskConsumerWaitGroup.Wait()
	println("Errors:")
	for _, t := range failedTasksArr {
		fmt.Printf("\tTask id %d time %s, error: %s\n", t.id, t.createTime.Format(time.RFC3339Nano), t.err)
	}

	println("Done tasks:")
	for t := range successedTasksMap {
		fmt.Printf("\t%v\n", t)
	}
}
