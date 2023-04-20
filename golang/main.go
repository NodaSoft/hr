package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
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

const (
	AppLifetime  = time.Second * 3
	TaskDuration = time.Millisecond * 150
	MaxTaskCount = AppLifetime/TaskDuration + 1
)

// A Task represents a meaninglessness of our life
type Task struct {
	id       uint32
	created  time.Time // время создания
	finished time.Time // время выполнения
	result   string
	hasError bool
}

type TaskManager struct {
	count  atomic.Uint32
	result []*Task
	errors []error
}

func (manager *TaskManager) NewTask() *Task {
	task := &Task{
		id:      manager.count.Add(1),
		created: time.Now(),
	}

	if task.created.Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
		task.hasError = true
	}

	return task
}

func (manager *TaskManager) PrintResult() {
	fmt.Println("Errors:")
	for err := range manager.errors {
		fmt.Println(err)
	}

	fmt.Println("Done tasks:")
	for task := range manager.result {
		fmt.Println(task)
	}
}

// главный метод, отрефакторенный бывший main()
func (manager *TaskManager) Run() {
	superChan := make(chan *Task, 10)
	sorterChan := make(chan *Task, 10)
	stopGenerate := make(chan bool, 1)
	wg := sync.WaitGroup{}

	// получение тасков
	wg.Add(1)
	go func() {
		defer close(sorterChan)
		defer wg.Done()

		for {
			task, ok := <-superChan
			if ok {
				manager.work(task)
				sorterChan <- task
			} else {
				return
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			task, ok := <-sorterChan
			if ok {
				manager.sort(task)
			} else {
				return
			}
		}
	}()

	//генерация новых тасков
	go func() {
		defer close(superChan)

		for {
			select {
			case <-stopGenerate:
				return

			default:
				// передаем таск на выполнение
				superChan <- manager.NewTask()
			}
		}
	}()

	// ждем таймаута или системного прерывания
	manager.waitForClose()
	//останавливаем генерацию тасков
	stopGenerate <- true
	// ждем остановки воркеров
	wg.Wait()
}

func (manager *TaskManager) work(task *Task) {
	// оставил оригинальное условие, хотя оно никогда не выполнится. Смотри выше task.hasError
	if time.Since(task.created) > 20*time.Second {
		task.result = "task has been successed"
	} else {
		task.result = "something went wrong"
	}

	task.finished = time.Now()
	time.Sleep(TaskDuration)
}

func (manager *TaskManager) sort(task *Task) {
	if task.hasError {
		err := fmt.Errorf("Task id %d time %s, error %s", task.id, task.created.Format(time.RFC3339), task.result)
		manager.errors = append(manager.errors, err)
	} else {
		manager.result = append(manager.result, task)
	}
}

func (manager *TaskManager) waitForClose() {
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Wait for interrupt signal to gracefully shutdown the application with
	// a timeout of 3 seconds.
	timer := time.NewTimer(AppLifetime)

	select {
	case <-quit:
		timer.Stop()

	case <-timer.C:
	}
}

func main() {
	taskManager := &TaskManager{
		result: make([]*Task, 0, int(MaxTaskCount)),
		errors: make([]error, 0, int(MaxTaskCount/10)),
	}

	taskManager.Run()
	taskManager.PrintResult()
}
