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

// A Ttype represents a meaninglessness of our life
type Task struct {
	ID         int
	CreateTime time.Time
	FinishTime time.Time
	Result     string
}

func main() {
	taskChan := make(chan Task)
	resultChan := make(chan Task)
	doneChan := make(chan struct{})

	var wg sync.WaitGroup
	
	// Start task generator
	wg.Add(1)
	go generateTasks(taskChan, &wg)

	// Start task processors
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go processTask(taskChan, resultChan, &wg)
	}

	// Start result collector
	go collectResults(resultChan, doneChan)

	// Periodically print results
	ticker := time.NewTicker(3 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				printResults()
			case <-doneChan:
				return
			}
		}
	}()

	// Run for 10 seconds
	time.Sleep(10 * time.Second)
	close(taskChan)

	wg.Wait()
	close(resultChan)
	<-doneChan
	ticker.Stop()

	// Final results
	printResults()
}

var (
	successTasks []Task
	failedTasks  []Task
	mu           sync.Mutex
)

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

func collectResults(resultChan <-chan Task, doneChan chan<- struct{}) {
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

func printResults() {
	mu.Lock()
	defer mu.Unlock()
	
	fmt.Println("Successful tasks:")
	for _, task := range successTasks {
		fmt.Printf("  ID: %d, Create Time: %s, Finish Time: %s\n",
			task.ID, task.CreateTime.Format(time.RFC3339),
			task.FinishTime.Format(time.RFC3339))
	}
	
	fmt.Println("Failed tasks:")
	for _, task := range failedTasks {
		fmt.Printf("  ID: %d, Create Time: %s, Finish Time: %s\n",
			task.ID, task.CreateTime.Format(time.RFC3339),
			task.FinishTime.Format(time.RFC3339))
	}
	fmt.Println()
}
