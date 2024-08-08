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
	failCondition := time.Now().Nanosecond()%2 > 0 // вот такое условие появления ошибочных тасков
	time.Sleep(taskCreatingDelay)
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

	fmt.Println("\033[32mSuccess:")
	for _, t := range success {
		fmt.Printf(
			"taskId: %d\t"+
				"creationTime: %s\t"+
				"executionTime: %s\n",
			t.id, t.creationTime.Format(time.RFC3339), t.executionTime.Format(time.RFC3339Nano))
	}

	fmt.Println("\033[31mFailed:")
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
			if time.Since(startTime) >= tasksCreatingDuration {
				close(stopChan)
				return
			}
			go func() {
				task := createTask()
				handleTask(task)
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
			case <-ticker.C:
				printingMtx.Lock()
				printTasks(handledTasks)
				printingMtx.Unlock()

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
