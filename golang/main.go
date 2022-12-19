package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнены остальных тасков

// A taskT represents a meaninglessness of our life
type taskT struct { // The trailing "T" is recommended here
	id            int64
	creationTime  time.Time // время создания
	executionTime time.Time // время выполнения
	result        string
}

type (
	taskC  = chan *taskT // Trailing letters is not necessarily my way of coding.
	taskWG = *sync.WaitGroup
)

const delta = 1

func isNanosecondEven(nanoSecond int) bool { // I hope you just want me to copy the condition
	return nanoSecond^delta == nanoSecond+delta // but not to do some fancy stuff
}

func startTaskCreator(ctx context.Context, c taskC, wg taskWG) {
	wg.Add(delta)

	go func() {
		log.Printf("INFO: Task creator started!")
		defer log.Printf("INFO: Task creator finished!")

		defer wg.Done()

		taskTicker := time.NewTicker(time.Second) // to prevent ID duplication
		defer taskTicker.Stop()

		defer close(c)

		for {
			select {
			case <-ctx.Done():
				return
			case tick := <-taskTicker.C:
				id := tick.Unix()

				if !isNanosecondEven(tick.Nanosecond()) {
					tick = time.Time{} // zero time.Time struct represents creation error

					taskTicker.Reset(time.Second) // to prevent nanosecond duplication
				}

				select {
				case <-ctx.Done():
					return
				case c <- &taskT{id: id, creationTime: tick}:
					if !tick.IsZero() {
						log.Printf("INFO: Task with ID %d created at %s", id, tick)

						continue
					}

					log.Printf("ERROR: Error creating task with id %d", id)
				}
			}
		}
	}()
}

func startTaskWorker(c taskC, wg taskWG) (finished taskC) {
	const jobDuration = time.Second * 2

	wg.Add(delta)

	fC := make(taskC)

	go func() {
		defer wg.Done()
		defer close(fC)

		finishedTasks := make([]*taskT, 0)

		for {
			task, ok := <-c
			if !ok {
				break
			}

			if task.creationTime.IsZero() {
				task.result = "job failed"
				finishedTasks = append(finishedTasks, task)

				continue
			}

			time.Sleep(jobDuration) // the job itself

			task.result = "job succeeded"
			task.executionTime = time.Now()

			finishedTasks = append(finishedTasks, task)
		}

		for _, task := range finishedTasks {
			fC <- task
		}
	}()

	return fC
}

func startTaskSorter(finishedTasksC []taskC, wg taskWG) {
	wg.Add(delta)

	go func() {
		succeededTasks := make([]*taskT, 0)
		failedTasks := make([]*taskT, 0)

		defer wg.Done()

		for _, taskChan := range finishedTasksC {
			for task := range taskChan {
				if !task.creationTime.IsZero() && !task.executionTime.IsZero() {
					succeededTasks = append(succeededTasks, task)

					continue
				}

				failedTasks = append(failedTasks, task)
			}
		}

		bldr := strings.Builder{}
		bldr.WriteString(fmt.Sprintf("INFO: Succeeded tasks:\n"))

		for _, task := range succeededTasks {
			bldr.WriteString(fmt.Sprintf("ID: %d, Result: %s, Creation Time: %s, Execution Time: %s\n",
				task.id, task.result, task.creationTime, task.executionTime))
		}

		log.Print(bldr.String())

		bldr.Reset()
		bldr.WriteString(fmt.Sprintf("INFO: Failed tasks:\n"))

		for _, task := range failedTasks {
			bldr.WriteString(fmt.Sprintf("ID: %d, Result: %s\n", task.id, task.result))
		}

		log.Print(bldr.String())
	}()
}

func main() {
	const (
		workers = 10
	)

	var (
		wg            sync.WaitGroup
		taskChan      = make(taskC) // is buffer here needed?
		finishedTasks = make([]taskC, 0, workers)
	)

	globalCtx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	for ctr := 0; ctr < workers; ctr++ {
		finishedTasks = append(finishedTasks, startTaskWorker(taskChan, &wg))
	}

	startTaskCreator(globalCtx, taskChan, &wg)

	startTaskSorter(finishedTasks, &wg)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	<-sig

	cancelFunc()

	wg.Wait()
}
