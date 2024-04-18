package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
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

// A Ttype represents a meaninglessness of our life
type Task struct {
	Id           int
	CreationTime string // время создания
	FullfullTime string // время выполнения
	Error        error
	Result       []byte
}

func main() {
	taskCreator := func(out chan<- Task, stop context.Context, wg *sync.WaitGroup) {
		defer wg.Done()

		var id int = 0

		for {
			select {
			case <-stop.Done():
				return
			default:
				var err error = nil
				ft := time.Now().Format(time.RFC3339)

				// time.Now().Nanosecond()%2 > 0, у меня всегда возвращал false
				if (time.Now().Nanosecond()/1000)%2 > 0 {
					err = errors.New("Some error occured")
				}

				id += 1
				out <- Task{
					CreationTime: ft,
					Error:        err,
					Id:           id,
				}
			}
		}
	}

	taskWorker := func(input <-chan Task, done chan<- Task, errChan chan<- error, stop context.Context, wg *sync.WaitGroup) {
		defer wg.Done()

		for {
			select {
			case <-stop.Done():
				return
			default:
				task := <-input

				time.Sleep(time.Millisecond * 150)

				if task.Error != nil {
					errChan <- fmt.Errorf("Task id [%d] time [%s], error [%s]", task.Id, task.CreationTime, task.Error.Error())
					continue
				}

				creationTime, err := time.Parse(time.RFC3339, task.CreationTime)
				if err != nil { // Немного бессмысленно, но пусть будет
					errChan <- fmt.Errorf("Task id [%d] time [%s], error [%s]", task.Id, task.CreationTime, task.Error.Error())
					continue
				}

				if creationTime.After(time.Now().Add(-20 * time.Second)) {
					task.Result = []byte("task has been successed")
				} else {
					task.Result = []byte("something went wrong")
				}
				task.FullfullTime = time.Now().Format(time.RFC3339Nano)

				done <- task
			}
		}
	}

	taskChan := make(chan Task, 10)
	doneTasks := make(chan Task, 10)
	undoneTasks := make(chan error, 10)

	var wg *sync.WaitGroup = &sync.WaitGroup{}

	workersContext, cancelWorkers := context.WithCancel(context.Background())

	wg.Add(2)

	go taskWorker(taskChan, doneTasks, undoneTasks, workersContext, wg)
	go taskCreator(taskChan, workersContext, wg)

	ticker := time.NewTicker(3 * time.Second)

	results := make(map[int]Task)
	errorList := make([]error, 0)

	func() { // Чтобы удобнее выходить из цикла
		for {
			select {
			case <-ticker.C:
				cancelWorkers()

				// TaskCraetor создает задачи без задержки, из-за чего большую часть времени будет в блокировке
				if cap(taskChan) == len(taskChan) {
					<-taskChan
				}
				wg.Wait()
				return
			case doneTask := <-doneTasks:
				fmt.Printf("Task id [%d] time [%s], result: [%s]\n", doneTask.Id, doneTask.CreationTime, doneTask.Result)
				results[doneTask.Id] = doneTask
			case err := <-undoneTasks:
				fmt.Println(err)
				errorList = append(errorList, err)
			}
		}
	}()

	close(doneTasks)
	close(undoneTasks)
	close(taskChan)

	println("Errors:")
	for _, v := range errorList {
		fmt.Println(v)
	}

	println("Done tasks:")
	for _, v := range results {
		fmt.Printf("%+v\n", v)
	}
}
