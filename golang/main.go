package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	taskChan       = make(chan Task, 10)
	doneTasks      = make(chan Task, 10)
	undoneTasks    = make(chan Task, 10)
	wg             sync.WaitGroup
	allDoneTasks   []Task
	allUndoneTasks []Task
)

const (
	PRINT_INTERVAL        = 3 * time.Second
	TASK_CREATION_TIMEOUT = 10 * time.Second
)

type Task struct {
	id             int
	creationTime   time.Time
	processionTime time.Time
	succeed        bool
	err            error
}

// Генерирует таски 10 секунд и закрывает taskChan
func taskCreator(taskChan chan<- Task) {
	defer close(taskChan)
	defer wg.Done()

	// Запускаем таймер
	timer := time.NewTimer(TASK_CREATION_TIMEOUT)
	// id теперь просто идущие подряд числа
	for id := 1; ; id++ {
		select {
		// Если прошло 10 секунд, выходим из функции
		case <-timer.C:
			return
		default:
			currentTime := time.Now()
			task := Task{
				id:           id,
				creationTime: currentTime,
			}
			// Удаляем из числа наносекунд незначащие нули справа
			nanoSecStr := strings.TrimRight(strconv.Itoa(currentTime.Nanosecond()), "0")
			nanoSec, _ := strconv.Atoi(nanoSecStr)
			if nanoSec%2 > 0 {
				task.err = fmt.Errorf("random error occurred")
			}
			// Задержка для замедления генерации
			time.Sleep(time.Second)

			taskChan <- task
		}
	}
}

// Обрабатаывает таски, отправляет в соответсвующий канал и завершает работу при закрытии taskChan
func taskWorker(taskChan <-chan Task, doneTasks, undoneTasks chan<- Task, quit chan<- struct{}) {
	defer wg.Done()
	// Сообщаем о завершении работы
	defer close(quit)

	for task := range taskChan {
		task.processionTime = time.Now()
		if task.err == nil {
			task.succeed = true
			doneTasks <- task
		} else {
			task.succeed = false
			undoneTasks <- task
		}
		time.Sleep(time.Millisecond * 150)
	}
}

// Выводит обработанные таски каждые *PRINT_INTERVAL* секунд и завершает работу при закрытии taskChan
func taskSorter(doneTasks, undoneTasks <-chan Task, quit <-chan struct{}) {
	defer wg.Done()

	ticker := time.NewTicker(PRINT_INTERVAL)
	defer ticker.Stop()
	for {
		select {
		case <-quit:
			return
		case task := <-doneTasks:
			allDoneTasks = append(allDoneTasks, task)
		case task := <-undoneTasks:
			allUndoneTasks = append(allUndoneTasks, task)
		case <-ticker.C:
			fmt.Printf("Succeed tasks at %s\n", time.Now().Format("15:04:05"))
			taskPrinter(allDoneTasks)
			fmt.Printf("Unsucceed tasks at %s\n", time.Now().Format("15:04:05"))
			taskPrinter(allUndoneTasks)
			fmt.Println("----------------------------------------------")

		}
	}

}

func taskPrinter(tasks []Task) {
	for _, task := range tasks {
		fmt.Printf("id: %d, created: %s, processed: %s\n", task.id, task.creationTime.Format("15:04:05"), task.processionTime.Format("15:04:05"))
	}
}

func main() {
	// Канал, который сообщит о завершении работы
	quit := make(chan struct{})

	wg.Add(3)
	go taskCreator(taskChan)
	go taskWorker(taskChan, doneTasks, undoneTasks, quit)
	go taskSorter(doneTasks, undoneTasks, quit)
	wg.Wait()
}
