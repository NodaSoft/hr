package main

import (
	"fmt"
	"math/rand"
	"sync"
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

// Task - структура задач
type Task struct {
	ID         int
	CreateTime time.Time
	FinishTime time.Time
	Result     string
}

var (
	successTasks []Task
	failedTasks  []Task
	mu           sync.Mutex
)

// generateTasks - генерация тасков
func generateTasks(taskChan chan<- Task, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case taskChan <- Task{
			ID:         rand.Intn(1000),
			CreateTime: time.Now(),
		}:
		case <-time.After(10 * time.Millisecond):
			return
		}
	}
}

// processTask - обработка тасков
func processTask(taskChan <-chan Task, resultChan chan<- Task, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range taskChan {
		time.Sleep(150 * time.Millisecond)
		task.FinishTime = time.Now()
		if task.CreateTime.Nanosecond()%2 > 0 {
			task.Result = "error"
		} else {
			task.Result = "success"
		}
		resultChan <- task
	}
}

// collectStatistic - сбор статистики
func collectStatistic(resultChan <-chan Task, doneChan chan<- struct{}) {
	for task := range resultChan {
		mu.Lock()
		if task.Result == "success" {
			successTasks = append(successTasks, task)
		} else {
			failedTasks = append(failedTasks, task)
		}
		mu.Unlock()
	}
	close(doneChan)
}

// printResults - вывод результатов по таскам
func printResults() {
	mu.Lock()
	defer mu.Unlock()

	fmt.Println("Successful tasks:")
	for _, task := range successTasks {
		fmt.Printf("ID: %d\nCreate Time: %s\nFinish Time: %s\n",
			task.ID, task.CreateTime.Format(time.RFC3339),
			task.FinishTime.Format(time.RFC3339))
	}

	fmt.Println("Failed tasks:")
	for _, task := range failedTasks {
		fmt.Printf("ID: %d\nCreate Time: %s\nFinish Time: %s\n",
			task.ID, task.CreateTime.Format(time.RFC3339),
			task.FinishTime.Format(time.RFC3339))
	}
	fmt.Println()
}

// ticker - вывод результатов в момент времени
func ticker(timer *time.Ticker, doneChan chan struct{}) {
	for {
		select {
		case <-timer.C:
			printResults()
		case <-doneChan:
			return
		}
	}
}

func main() {
	taskChan := make(chan Task)
	resultChan := make(chan Task)
	doneChan := make(chan struct{})
	mu = sync.Mutex{}

	var wg sync.WaitGroup

	wg.Add(1)
	go generateTasks(taskChan, &wg)

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go processTask(taskChan, resultChan, &wg)
	}

	go collectStatistic(resultChan, doneChan)

	timer := time.NewTicker(3 * time.Second)
	go ticker(timer, doneChan)

	time.Sleep(10 * time.Second)
	close(taskChan)

	wg.Wait()
	close(resultChan)
	<-doneChan
	timer.Stop()

	printResults()
}
