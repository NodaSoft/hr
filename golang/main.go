// README
// Приложение эмулирует получение и обработку неких тасков - получает и обрабатывает в многопоточном режиме.
// При запуске можно указать флаги: 
// --time для указания продолжительности работы приложения в секундах (default = 10)
// --frequency для определения частоты создания задач в милисекундах (default 1/150)
// Приложение генерирует таски и каждые 3 секунды выводит в консоль результат всех обработанных к этому моменту задач (отдельно успешные и отдельно с ошибками).

package main

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"sync"
	"time"
)

var (
	execution_time int
	frequency      int
)

func init() {
	flag.IntVar(&execution_time, "time", 10, "Task execution time")
	flag.IntVar(&frequency, "frequency", 150, "How often the task is created: 1 per frequency in milliseconds")
	flag.Parse()
}

const (
	not_successfully = "something went wrong"
	successfully     = "task has been successed"
)

// A Tasker represents a meaninglessness of our life
type Tasker struct {
	id           int
	created_at   string
	completed_in string
	task_broken  bool
	description  string
}

// Closing function
func intSeq() func() int {
	i := 0
	return func() int {
		i++
		return i
	}
}
// The closure is used as a sequence number for identification.
var nextID = intSeq()

func main() {

	created_tasks := make(chan Tasker)
	doneTasks := make(chan Tasker)
	undoneTasks := make(chan error)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(execution_time))
	defer cancel()

	fmt.Println("Start processing tasks...")

	go TaskCreturer(ctx, created_tasks)

	go TaskProcessing(created_tasks, doneTasks, undoneTasks)

	Result(doneTasks, undoneTasks)

	fmt.Println("\nEnd of task processing...")
}

// TaskCreturer имитирует создание рабочих и ошибочных задач и отправляет результат в канал. Параметр frequency определяет частоту создания тасок.
func TaskCreturer(ctx context.Context, created_tasks chan Tasker) {
	for {
		select {
		case <-ctx.Done():
			close(created_tasks)
			return
		default:
			var broken bool
			time_now := time.Now().Format(time.RFC3339)
			if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
				broken = true
			}

			time.Sleep(time.Millisecond * time.Duration(frequency))
			created_tasks <- Tasker{id: nextID(), created_at: time_now, task_broken: broken} // передаем таск на выполнение
		}
	}
}

// TaskProcessing является оберткой над обработкой и сортировкой задач с проверкой на безопасность при многопоточном исполнении.
func TaskProcessing(created_tasks chan Tasker, doneTasks chan Tasker, undoneTasks chan error) {
	processed_tasks := make(chan Tasker)
	wg := sync.WaitGroup{}

	for task := range created_tasks {
		wg.Add(1)
		Processing(task, processed_tasks)
		Sorter(processed_tasks, doneTasks, undoneTasks, &wg)

	}
	defer func() {
		wg.Wait()
		close(processed_tasks)
		close(doneTasks)
		close(undoneTasks)
	}()
}

// Processing асинхронно имитирует обработку задач (с задержкой) и заполняет описание обработки.
func Processing(tsk Tasker, processed_tasks chan Tasker) {

	go func() {
		if tsk.task_broken {
			tsk.description = not_successfully
		} else {
			tsk.description = successfully
		}
		tsk.completed_in = time.Now().Format(time.RFC3339Nano)

		time.Sleep(time.Second * 2)

		processed_tasks <- tsk
	}()
}

// Sorter асинхронно сортирует результаты обработки по каналам.
func Sorter(processed_tasks chan Tasker, doneTasks chan Tasker, undoneTasks chan error, wg *sync.WaitGroup) {

	go func() {
		t := <-processed_tasks

		if strings.HasSuffix(t.description, "successed") {
			doneTasks <- t
		} else {
			undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.created_at, t.description)
		}
		wg.Done()
	}()
}

// Result забирает результаты после сортировки и отправляет на печать с периодичностью 3 сек.
func Result(doneTasks chan Tasker, undoneTasks chan error) {

	var processed = make(map[int]Tasker)
	var unprocessed = make([]error, 0)

	completion_done := make(chan struct{})
	completion_undone := make(chan struct{})
	completion_full := make(chan struct{})

	go func() {
		for r := range doneTasks {
			processed[r.id] = r
		}
		completion_done <- struct{}{}
		close(completion_done)
	}()
	go func() {
		for r := range undoneTasks {
			unprocessed = append(unprocessed, r)
		}
		completion_undone <- struct{}{}
		close(completion_undone)
	}()

	go func() {
		<-completion_done
		<-completion_undone
		completion_full <- struct{}{}
	}()

	var iteration int
	for {
		time.Sleep(time.Second * 3)
		iteration++
		fmt.Println("\nIteration:", iteration)
		fmt.Println("*****************************************************************")

		fmt.Println("Done tasks:")
		for _, p := range processed {
			fmt.Println(p)
		}
		clear(processed)

		fmt.Println("Errors:")
		for _, e := range unprocessed {
			fmt.Println(e)
		}
		unprocessed = unprocessed[:0]

		select {
		case <-completion_full:
			return
		default:
		}
	}
}
