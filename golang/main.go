package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Приложение эмулирует получение и обработку неких тасков. Пытается и получать, и обрабатывать в многопоточном режиме.
// После обработки тасков в течении 3 секунд приложение должно выводить накопленные к этому моменту успешные таски и отдельно ошибки обработки тасков.

// TODO: сделать из плохого кода хороший и рабочий - as best as you can.
// Важно сохранить логику появления ошибочных тасков.
// Важно оставить асинхронные генерацию и обработку тасков.
// Сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через pull-request в github
// Как видите, никаких привязок к внешним сервисам нет - полный карт-бланш на модификацию кода.

// Task представляет собой задачу
type Task struct {
	id            int
	creationTime  time.Time // Время создания
	executionTime time.Time // Время выполнения
	result        []byte    // Результат выполнения
	err           error     // Ошибка (если есть)
}

const timeFormat = time.RFC3339Nano

// TaskCreator генерирует задачи и отправляет их в канал out
func TaskCreator(out chan<- Task, ctx context.Context) {
	defer close(out)
	for {
		select {
		case <-ctx.Done(): // если контекст отменен - завершаем генерацию
			return
		default: // иначе продолжаем отправлять задачи в out
			var task Task
			task.creationTime = time.Now()
			if time.Now().Nanosecond()%2 > 0 { // Условие для появления ошибочных задач
				task.err = errors.New("some error occurred")
			}
			task.id = int(time.Now().Unix())

			out <- task // Передача задачи на выполнение
		}
	}
}

// TaskHandler обрабатывает/выполняет задачи
func TaskHandler(task Task) Task {
	if task.creationTime.After(time.Now().Add((-20)*time.Second)) && task.err == nil {
		task.result = []byte("task has been successful")
	} else {
		task.result = []byte("something went wrong")
	}
	task.executionTime = time.Now()

	time.Sleep(time.Millisecond * 150)

	return task
}

// TaskSorter сортирует задачи на завершенные и незавершенные (с ошибкой)
func TaskSorter(task Task, doneTasks, undoneTasks chan<- Task) {
	if task.err == nil {
		doneTasks <- task
	} else {
		undoneTasks <- task
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	tasksChan := make(chan Task, 10)
	doneTasks := make(chan Task)
	undoneTasks := make(chan Task)

	// Отслеживание сигналов завершения приложения
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		cancel() // Отмена контекста при получении сигнала завершения (Ctrl + C)
	}()

	go TaskCreator(tasksChan, ctx) // генерируем задачи и отправляем их в tasksChan

	go func() {
		// Закрываем каналы после завершения
		defer close(doneTasks)
		defer close(undoneTasks)
		// Получение задач
		for {
			select {
			case t, ok := <-tasksChan: // получили какую-то задачу
				if !ok { // если канал был закрыт
					return
				}
				go TaskSorter(TaskHandler(t), doneTasks, undoneTasks)
			}
		}
	}()

	completedTasks := make(map[int]Task)
	var taskErrors []error

	go func() { // сохраняем выполненные таски
		for task := range doneTasks {
			completedTasks[task.id] = task
		}
	}()
	go func() { // сохраняем ошибки тасков
		for task := range undoneTasks {
			err := fmt.Errorf("task id: %d, creation time: %s, execution time: %s, result: %s, error: %w",
				task.id, task.creationTime.Format(timeFormat), task.executionTime.Format(timeFormat), task.result, task.err)
			taskErrors = append(taskErrors, err)
		}
	}()

	time.Sleep(3 * time.Second)
	cancel()

	fmt.Println("Errors:")
	for _, err := range taskErrors {
		fmt.Println(err.Error())
	}

	fmt.Println("Done tasks:")
	for taskID, task := range completedTasks {
		fmt.Printf(
			"ID: %d, creation time: %s, execution time: %s, result %s\n",
			taskID, task.creationTime.Format(timeFormat), task.executionTime.Format(timeFormat), task.result,
		)
	}
}
