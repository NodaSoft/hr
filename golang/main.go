package main

import (
	"fmt"
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

// Мы даем тестовое задание чтобы:
// * уменьшить время технического собеседования - лучше вы потратите пару часов в спокойной домашней обстановке, чем будете волноваться, решая задачи под взором наших ребят;
// * увеличить вероятность прохождения испытательного срока - видя сразу стиль и качество кода, мы можем быть больше уверены в выборе;
// * снизить число коротких собеседований, когда мы отказываем сразу же.

// Выполнение тестового задания не гарантирует приглашение на собеседование, т.к. кроме качества выполнения тестового задания, оцениваются и другие показатели вас как кандидата.

// Мы не даем комментариев по результатам тестового задания. Если в случае отказа вам нужен наш комментарий по результатам тестового задания, то просим об этом написать вместе с откликом.

// A Ttype represents a meaninglessness of our life
type Ttype struct {
	id         int
	cT         string // время создания
	fT         string // время выполнения
	taskRESULT []byte
}

func main() {
	taskChannel := make(chan Ttype, 10)
	doneTasks := make(chan Ttype, 10)
	errorTasks := make(chan Ttype, 10)
	stopTaskGeneration := make(chan struct{})
	stopPrinting := make(chan struct{})

	var wg sync.WaitGroup

	// Start task generation
	go taskGenerator(taskChannel, stopTaskGeneration)

	// Start task processing
	wg.Add(1)
	go taskProcessor(taskChannel, doneTasks, errorTasks, &wg)

	// Start printing results every 3 seconds
	wg.Add(1)
	go printResults(doneTasks, errorTasks, stopPrinting)

	// Generate tasks for 10 seconds
	time.Sleep(10 * time.Second)
	close(stopTaskGeneration)
	close(taskChannel)

	// Wait for task processing to complete
	wg.Wait()
	close(doneTasks)
	close(errorTasks)
	close(stopPrinting)

	fmt.Println("All tasks are processed.")
}

// taskGenerator generates tasks and sends them to the taskChannel
func taskGenerator(taskChannel chan<- Ttype, stop <-chan struct{}) {
	for {
		select {
		case <-stop:
			return
		default:
			cT := time.Now().Format(time.RFC3339)
			if time.Now().Nanosecond()%2 > 0 { // condition for erroneous tasks
				cT = "Some error occurred"
			}
			task := Ttype{cT: cT, id: int(time.Now().UnixNano())}
			taskChannel <- task
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// taskProcessor processes tasks from the taskChannel and sends results to doneTasks and errorTasks
func taskProcessor(taskChannel <-chan Ttype, doneTasks chan<- Ttype, errorTasks chan<- Ttype, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range taskChannel {
		processTask(task, doneTasks, errorTasks)
	}
}

// processTask handles individual task processing
func processTask(t Ttype, doneTasks chan<- Ttype, errorTasks chan<- Ttype) {
	tt, err := time.Parse(time.RFC3339, t.cT)
	if err != nil || tt.After(time.Now().Add(-20*time.Second)) {
		t.taskRESULT = []byte("task has been successed")
		doneTasks <- t
	} else {
		t.taskRESULT = []byte("something went wrong")
		errorTasks <- t
	}
	t.fT = time.Now().Format(time.RFC3339Nano)
	time.Sleep(150 * time.Millisecond)
}

// printResults prints the results from doneTasks and errorTasks every 3 seconds
func printResults(doneTasks <-chan Ttype, errorTasks <-chan Ttype, stop <-chan struct{}) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("Done tasks:")
			printDoneTasks(doneTasks)

			fmt.Println("Errors:")
			printErrorTasks(errorTasks)
		case <-stop:
			return
		}
	}
}

// printDoneTasks prints all completed tasks
func printDoneTasks(doneTasks <-chan Ttype) {
	for {
		select {
		case task, ok := <-doneTasks:
			if !ok {
				return
			}
			fmt.Printf("Task ID: %d, Creation Time: %s, Finish Time: %s, Result: %s\n", task.id, task.cT, task.fT, string(task.taskRESULT))
		default:
			return
		}
	}
}

// printErrorTasks prints all tasks with errors
func printErrorTasks(errorTasks <-chan Ttype) {
	for {
		select {
		case task, ok := <-errorTasks:
			if !ok {
				return
			}
			fmt.Printf("Task ID: %d, Creation Time: %s, Error: %s\n", task.id, task.cT, string(task.taskRESULT))
		default:
			return
		}
	}
}
