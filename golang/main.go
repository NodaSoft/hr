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
	creationTimeFormat    = time.RFC3339
	executionTimeFormat   = time.RFC3339Nano
	tasksCreatingDuration = 10 * time.Second
	taskCreatingDelay     = 150 * time.Millisecond
	maxTaskCreatingDelay  = 20 * time.Second
	tasksPrintingDelay    = 3 * time.Second
)

var (
	idCounter = 0
	mu1       sync.Mutex
	mu2       sync.Mutex
	mu3       sync.Mutex
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

	mu1.Lock()
	taskId := idCounter
	idCounter++
	failCondition := time.Now().Nanosecond()%2 > 0 // вот такое условие появления ошибочных тасков
	time.Sleep(taskCreatingDelay)
	mu1.Unlock()

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
	fmt.Println("Handled tasks:")
	for _, t := range tasks {
		fmt.Printf(
			"taskId: %d\t"+
				"creationTime: %s\t"+
				"executionTime: %s\t"+
				"logs: %s\n",
			t.id, t.creationTime.Format(creationTimeFormat), t.executionTime.Format(executionTimeFormat),
			strings.Join(t.logs, ", "))
	}
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
			if time.Since(startTime) >= tasksCreatingDuration {
				close(stopChan)
				return
			}
			go func() {
				task := createTask()
				handleTask(task)
				mu2.Lock()
				handledTasks = append(handledTasks, task)
				mu2.Unlock()
			}()
		}
	}()

	go func() {
		defer wg.Done()
		ticker := time.NewTicker(tasksPrintingDelay)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				mu3.Lock()
				printTasks(handledTasks)
				mu3.Unlock()
			case <-stopChan:
				return
			}
		}
	}()

	wg.Wait()

	printTasks(handledTasks)
}
