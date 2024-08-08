package main

import (
	"fmt"
	"strings"
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

const (
	tasksCreatingDuration = 10 * time.Second
	taskCreatingDelay     = 150 * time.Millisecond
	maxTaskCreatingDelay  = 20 * time.Second
	tasksPrintingDelay    = 3 * time.Second
)

var (
	idCounter   = 0
	creationMtx sync.Mutex
	handlingMtx sync.Mutex
	printingMtx sync.Mutex
)

// A Task represents a meaninglessness of our life
type Task struct {
	id            int
	creationTime  time.Time // время создания
	executionTime time.Time // время выполнения
	isSuccess     bool
	isCompleted   bool
	logs          []string
}

func createTask() *Task {
	creationTime := time.Now()

	creationMtx.Lock()
	taskId := idCounter
	idCounter++
	failCondition := time.Now().Nanosecond()%2 > 0 // Condition for the appearance of failed tasks
	time.Sleep(taskCreatingDelay)                  // Pretend that the task is working too long
	creationMtx.Unlock()

	if failCondition {
		return &Task{
			id:           taskId,
			creationTime: creationTime,
			isSuccess:    false,
			logs:         []string{"Some error occurred"},
		}
	}
	return &Task{
		id:           taskId,
		creationTime: creationTime,
		isSuccess:    true,
	}
}

func handleTask(t *Task) {
	duration := t.creationTime.Sub(time.Now())

	if duration <= maxTaskCreatingDelay {
		t.isCompleted = true
		t.logs = append(t.logs, "Task was completed successfully")
	} else {
		t.isCompleted = false
		t.logs = append(t.logs, "Task is working too long, something went wrong")
	}

	t.executionTime = time.Now()
}

func printTasks(tasks []*Task) {
	success := make([]*Task, 0, len(tasks))
	failed := make([]*Task, 0, len(tasks))

	for _, t := range tasks {
		if t.isSuccess && t.isCompleted {
			success = append(success, t)
			continue
		}
		failed = append(failed, t)
	}

	fmt.Println("\033[32mSuccess:") // Green color console
	for _, t := range success {
		fmt.Printf(
			"taskId: %d\t"+
				"creationTime: %s\t"+
				"executionTime: %s\n",
			t.id, t.creationTime.Format(time.RFC3339), t.executionTime.Format(time.RFC3339Nano))
	}

	fmt.Println("\033[31mFailed:") // Red color console
	for _, t := range failed {
		fmt.Printf(
			"taskId: %d\t"+
				"creationTime: %s\t"+
				"logs: %s\n",
			t.id, t.creationTime.Format(time.RFC3339), strings.Join(t.logs, ", "))
	}

	fmt.Println()
}

func main() {
	stopChan := make(chan struct{})
	handledTasks := make([]*Task, 0)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		startTime := time.Now()
		for {
			// If 10 seconds have passed
			if time.Since(startTime) >= tasksCreatingDuration {
				close(stopChan)
				return
			}
			// Creation and handling tasks in goroutines
			go func() {
				task := createTask()
				handleTask(task)
				// Appending to slice with handled tasks concurrently
				handlingMtx.Lock()
				handledTasks = append(handledTasks, task)
				handlingMtx.Unlock()
			}()
		}
	}()

	go func() {
		defer wg.Done()
		ticker := time.NewTicker(tasksPrintingDelay)
		defer ticker.Stop()
		for {
			select {
			// Every time 3 seconds have passed
			case <-ticker.C:
				// Printing tasks concurrently
				printingMtx.Lock()
				printTasks(handledTasks)
				printingMtx.Unlock()
				// Clearing slice concurrently
				handlingMtx.Lock()
				handledTasks = handledTasks[:0]
				handlingMtx.Unlock()
			case <-stopChan:
				return
			}
		}
	}()

	wg.Wait()

	printTasks(handledTasks)
}
