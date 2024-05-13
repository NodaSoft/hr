package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Приложение эмулирует получение и обработку неких тасков. Пытается и получать, и обрабатывать в многопоточном режиме.
// После обработки тасков в течении 3 секунд приложение должно выводить накопленные к этому моменту успешные таски и отдельно ошибки обработки тасков.

// ЗАДАНИЕ: сделать из плохого кода хороший и рабочий - as best as you can.
// Важно сохранить логику появления ошибочных тасков.
// Важно оставить асинхронные генерацию и обработку тасков.
// Сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через pull-request в github
// Как видите, никаких привязок к внешним сервисам нет - полный карт-бланш на модификацию кода.

// A Task represents a meaninglessness of our life
type Task struct { // Изменяем название структуры и полей на более читаемый
	ID         int
	CreatedAt  string // время создания
	FinishedAt string // время выполнения
	Result     []byte
}

type TaskError struct {
	TaskID  int
	Message string
}

func NewTaskError(taskID int, message string) *TaskError {
	return &TaskError{
		TaskID:  taskID,
		Message: message,
	}
}

func (e *TaskError) Error() string {
	return fmt.Sprintf("Task id %d, error %s", e.TaskID, e.Message)
}

func createTasks(ctx context.Context, taskChan chan<- Task, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			ft := time.Now().Format(time.RFC3339)
			if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
				ft = "Some error occured"
			}
			task := Task{CreatedAt: ft, ID: int(time.Now().Unix())} // передаем таск на выполнение
			taskChan <- task
		}
	}
}

func processTask(task Task) (Task, error) {
	tt, _ := time.Parse(time.RFC3339, task.CreatedAt)
	if tt.After(time.Now().Add(-20 * time.Second)) {
		task.Result = []byte("task has been successed")
	} else {
		return task, NewTaskError(task.ID, "task has been failed")
	}
	task.FinishedAt = time.Now().Format(time.RFC3339Nano)
	time.Sleep(time.Millisecond * 150)
	return task, nil
}

func sortTasks(ctx context.Context, taskChan <-chan Task, doneTasks chan<- Task, errorTasks chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range taskChan {
		select {
		case <-ctx.Done():
			return
		default:
			task, err := processTask(task)
			if err != nil {
				errorTasks <- err
			} else if string(task.Result) == "task has been successed" {
				doneTasks <- task
			}
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	taskChan := make(chan Task)
	doneTasks := make(chan Task)
	errorTasks := make(chan error)

	var wg sync.WaitGroup
	wg.Add(2)

	go createTasks(ctx, taskChan, &wg)
	go sortTasks(ctx, taskChan, doneTasks, errorTasks, &wg)

	go func() {
		wg.Wait()
		close(taskChan)
		close(doneTasks)
		close(errorTasks)
	}()

	fmt.Println("Errors Tasks:")
	for err := range errorTasks {
		fmt.Println(err)
	}

	fmt.Println("Done Tasks:")
	for task := range doneTasks {
		fmt.Printf("Task ID: %d, Created At: %s, Finished At: %s, Result: %s\n", task.ID, task.CreatedAt, task.FinishedAt, task.Result)
	}
}
