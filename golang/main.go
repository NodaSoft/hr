package main

import (
	"fmt"
	"time"
)

// Были изменен нейминг сущностей на кэмел кейс, каждое название теперь отражает предназначение
// В структуре кода был разделен по областям код, описывающий функции, код иницилизации каналов и код запуска горутин

// Task represents task that contains result
type Task struct {
	ID int
	// Были убраны поля FinishedAt(fT) и CreatedAt(cT), т.к. в этих полях нет необходимости для реализации задачи
	Result string // изменен тип поля на string, т.к. фактическое его предназначение - хранение строки
}

func main() {

	// taskCreator создает задачи, симулируя задержку между их появлением
	taskCreator := func(taskChan chan<- Task) {
		go func() {
			for {
				createdAt := time.Now()

				task := Task{
					ID: int(createdAt.Unix()),
				}

				taskChan <- task
				time.Sleep(time.Second) // Искусственная задержка между созданиями задач
			}
		}()
	}

	// taskWorker обрабатывает задачи и записывает в них результат
	taskWorker := func(task Task) Task {
		if time.Now().Nanosecond()%2 > 0 { // Симуляция ошибочного кейса
			task.Result = "task has encountered errors"
		} else {
			task.Result = "task has been succeed"
		}
		time.Sleep(time.Millisecond * 100)

		return task
	}

	//taskSorter
	taskSorter := func(task Task, doneTasks chan<- Task, errorTasks chan<- Task) {
		if task.Result == "task has been succeed" {
			doneTasks <- task
		} else {
			errorTasks <- task
		}
	}

	tasksChan := make(chan Task)
	doneTasksChan := make(chan Task)
	errorTasksChan := make(chan Task)

	go taskCreator(tasksChan)

	//
	go func() {
		for task := range tasksChan {
			go func(task Task) {
				task = taskWorker(task)
				taskSorter(task, doneTasksChan, errorTasksChan)
			}(task)
		}
	}()

	// Была добавлена горутина для вывода тасков, чтобы результат отображался сразу после выполнения
	go func() {
		for {
			select {
			case task := <-doneTasksChan:
				fmt.Printf("Done Task: ID %d, Result: %s\n", task.ID, task.Result)
			case task := <-errorTasksChan:
				fmt.Printf("Error Task: ID %d, Result: %s\n", task.ID, task.Result)
			}
		}
	}()

	time.Sleep(time.Second * 10)
}
