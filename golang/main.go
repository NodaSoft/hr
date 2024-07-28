package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Приложение эмулирует получение и обработку неких тасков. Пытается и получать, и обрабатывать в многопоточном режиме.
// Приложение должно генерировать таски 10 сек.
// Каждые 3 секунды должно выводить в консоль результат всех обработанных к этому моменту тасков
// (отдельно успешные и отдельно с ошибками).

// ЗАДАНИЕ: сделать из плохого кода хороший и рабочий - as best as you can.
// Важно сохранить логику появления ошибочных тасков.
// Важно оставить асинхронные генерацию и обработку тасков.
// Сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через pull-request в github
// Как видите, никаких привязок к внешним сервисам нет - полный карт-бланш на модификацию кода.

// A Task represents a meaninglessness of our life

const (
	CreatingDurationSeconds = 10
	WorkerPoolSize          = 2
	ReportPeriodSeconds     = 3
	WorkSimulationMs        = 150
)

var (
	ErrTaskWasFailed = errors.New("task was failed")
)

type StatusType int

const (
	UndefinedStatus StatusType = iota
	Success
	Failure
)

type FailureReasonType int

const (
	UndefinedFailureReason FailureReasonType = iota
	ErrCreating
	ErrOverdue
)

type Task struct {
	id            int
	status        StatusType
	failureReason FailureReasonType
	createdAt     time.Time // время создания
	finishedAt    time.Time // время выполнения
}

func (t *Task) GetCreatedAt() time.Time {
	return t.createdAt
}

func (t *Task) SetFailure(reason FailureReasonType) {
	t.finishedAt = time.Now()
	t.status = Failure
	t.failureReason = reason
}

func (t *Task) SetSuccess() error {
	if t.status == Failure {
		return ErrTaskWasFailed
	}
	t.finishedAt = time.Now()
	t.status = Success

	return nil
}

func (t *Task) GetStatus() StatusType {
	return t.status
}

func taskCreator(ctx context.Context, a chan Task) {
	defer close(a)

	for {
		select {
		case <-ctx.Done():
			log.Printf("creating cancelled: %v\n", ctx.Err())
			return
		default:
			currentTime := time.Now()
			task := Task{createdAt: currentTime, id: int(currentTime.Unix())}

			if currentTime.Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
				task.SetFailure(ErrCreating)
			}
			a <- task // передаем таск на выполнение
		}
	}
}

func taskWorker(
	ctx context.Context,
	wg *sync.WaitGroup,
	inputChan <-chan Task,
	outputChan chan<- Task,
) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			log.Printf("worker: %v\n", ctx.Err())
			return
		case task, ok := <-inputChan:
			if !ok {
				return
			}

			if task.GetStatus() != Failure && task.GetCreatedAt().After(time.Now().Add(-20*time.Second)) {
				err := task.SetSuccess()
				if err != nil {
					log.Printf("%v\n", err)
				}
			} else {
				task.SetFailure(ErrOverdue)
			}

			// work simulation
			time.Sleep(time.Millisecond * WorkSimulationMs)
			outputChan <- task
		}
	}
}

func main() {
	ctx, cancelMain := context.WithCancel(context.Background())
	defer cancelMain()

	newTasksCh := make(chan Task)

	ctxWithTimeout, cancelWithTimeout := context.WithTimeout(ctx, time.Second*CreatingDurationSeconds)
	defer cancelWithTimeout()

	go taskCreator(ctxWithTimeout, newTasksCh)

	wg := &sync.WaitGroup{}
	processedTasksCh := make(chan Task)

	go func() {
		for i := 0; i < WorkerPoolSize; i++ {
			wg.Add(1)
			go taskWorker(ctx, wg, newTasksCh, processedTasksCh)
		}

		wg.Wait()
		close(processedTasksCh)
	}()

	var successTasks []Task
	var failureTasks []Task
	var undefinedTasks []Task

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
		<-quit
		cancelMain()
	}()

	ticker := time.NewTicker(time.Second * ReportPeriodSeconds)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			PrintReport(successTasks, failureTasks, undefinedTasks)
		case task, ok := <-processedTasksCh:
			if !ok {
				fmt.Println("Final report:")
				PrintReport(successTasks, failureTasks, undefinedTasks)
				fmt.Println("All tasks have been processed")

				return
			}

			switch task.GetStatus() {
			case Success:
				successTasks = append(successTasks, task)
			case Failure:
				failureTasks = append(failureTasks, task)
			default:
				undefinedTasks = append(undefinedTasks, task)
			}
		}
	}
}

func PrintReport(successTasks, failureTasks, undefinedTasks []Task) {
	fmt.Printf("successTasks: %v\n", len(successTasks))
	fmt.Printf("failureTasks: %v\n", len(failureTasks))
	fmt.Printf("undefinedTasks: %v\n", len(undefinedTasks))
	fmt.Printf("_______\n")
}
