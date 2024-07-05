package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// Task описывает задачу
type Task struct {
	id          int
	createdAt   time.Time // время создания задачи
	processedAt time.Time // время выполнения
	result      bool
}

// taskProducer создает новые таски и помещает их в output chan
func taskProducer(output chan *Task) {
	log.Printf("Task producer started\n")
	defer close(output) // закрываем канал после окончания работы функции
	wg := sync.WaitGroup{}
	done := time.After(10 * time.Second) // создаем таймер, тк таски должны генерится в течени 10 секунд
	var count atomic.Uint64              // счетчик созданных задач
	for {
		select {
		case <-done:
			log.Printf("Task producer stopped. Num of tasks generated: %v\n", count.Load())
			wg.Wait()
			return
		default:
			wg.Add(1)
			go func() {
				defer wg.Done()
				task := Task{createdAt: time.Now(), id: rand.Intn(1000)}
				if time.Now().Nanosecond()%2 > 0 { // Задача считается ошибочной, если текущее кол-во наносекунд является нечетным
					task.createdAt = time.Time{}
				}
				count.Add(1)
				output <- &task
			}()

			time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond) // рандомный слип
		}
	}
}

// taskWorker обрабатывает задачи и помещает в каналы в зависимости от того, ошибочная таска или нет
func taskWorker(task *Task, processedTasks chan *Task, failedTasks chan error, wg *sync.WaitGroup) {
	log.Println("task worker started\n")
	defer wg.Done()

	task.processedAt = time.Now()

	// проверяем, была ли задача создана в течение последних 20 секунд и не ошибочная ли она
	if !task.createdAt.IsZero() && time.Since(task.createdAt) <= 20*time.Second {
		task.result = true
		processedTasks <- task
	} else {
		task.result = false
		failedTasks <- fmt.Errorf("failed task id %d", task.id)
	}

	time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond) // рандомный слип для имитации обработки задачи

}

func main() {

	producerChan := make(chan *Task, 10) //канал для созданных задач
	go taskProducer(producerChan)        //генерация задач

	processedTasks := make(chan *Task) //канал для успешно обработанных задач
	failedTasks := make(chan error)    //канал для ошибочных задач

	// обработка созданных задач
	wg := sync.WaitGroup{}
	go func() {
		defer close(processedTasks)
		defer close(failedTasks)
		for task := range producerChan {
			wg.Add(1)
			go taskWorker(task, processedTasks, failedTasks, &wg)
		}
		wg.Wait()
	}()

	successTasks := make(map[int]*Task) // мапа для успешно обработанных задач
	mu := sync.Mutex{}
	go func() {
		for task := range processedTasks {
			mu.Lock()
			successTasks[task.id] = task
			mu.Unlock()
		}
	}()

	errorTasks := make([]error, 0) //канал для ошибочных задач
	mu2 := sync.Mutex{}
	go func() {
		for task := range failedTasks {
			mu2.Lock()
			errorTasks = append(errorTasks, task)
			mu2.Unlock()
		}
	}()

	// Запуск горутины для печати содержимого каналов каждые три секунды
	ticker := time.NewTicker(3 * time.Second)
	quit := make(chan struct{})
	tmpE := 0
	tmpS := 0
	wg.Add(1)
	go func() {
		for {
			select {
			case <-ticker.C:
				//проверяет увеличилось ли количество обработанных задач с прошлой итерации
				mu.Lock()
				mu2.Lock()
				if len(errorTasks) != tmpE || len(successTasks) != tmpS {

					fmt.Printf("Num of errors: %v\n", len(errorTasks))
					fmt.Printf("Num of succeeded tasks: %v\n", len(successTasks))

					fmt.Printf("Success: %v\n", successTasks)
					fmt.Printf("Errors: %v \n", errorTasks)

					tmpE, tmpS = len(errorTasks), len(successTasks) //увеличиваем значение обработанных задач
					mu.Unlock()
					mu2.Unlock()
				} else { // если количество задач осталось неизменным, значит все задачи уже напечатаны
					close(quit)
				}

			case <-quit:
				log.Println("end of printing tasks")
				ticker.Stop()
				wg.Done()
				return
			}
		}
	}()

	wg.Wait()
}
