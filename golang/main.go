package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

type TaskID string

type Task struct {
	id                  TaskID
	createdAt           time.Time // Момент создания задачи.
	executionFinishedAt time.Time // Момент исполнения.
	mustExecutionFail   bool      // Подсказка о том, что должно случиться с задачей.
}

func (t Task) String() string {
	return fmt.Sprintf("Task (id=%v)", t.id)
}

func (t *Task) Execute() error {
	// Имитируем вычисления.
	time.Sleep(500 * time.Millisecond)
	t.executionFinishedAt = time.Now().UTC()
	if t.mustExecutionFail {
		return fmt.Errorf("%s errored out at %s", t, t.executionFinishedAt)
	}
	return nil
}

func NewTask() Task {
	t := Task{
		id:        <-id,
		createdAt: time.Now().UTC(),
	}
	// Условием тестового задания было сохранить логику возникновения ошибок:
	if time.Now().Nanosecond()%2 > 0 {
		t.mustExecutionFail = true
	}
	return t
}

var (
	wg sync.WaitGroup
	// Канал-источник читабельных идентификаторов задач (A, B, C, ..., A, B, C...)
	id chan TaskID
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	id = make(chan TaskID)

	go func() {
		start := "A"[0]
		end := "Z"[0]
		current := start
		for {
			id <- TaskID(current)
			current += 1
			if current > end {
				current = start
			}
		}
	}()

	// Канал-источник задач для исполнителей.
	taskSourceChannel := make(chan Task, 10)

	// Создатель задач.
	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				close(taskSourceChannel)
				log.Printf("Creator exited")
				return
			default:
				task := NewTask()
				taskSourceChannel <- task
				log.Printf("Created task: %s", task)
			}
		}
	}(ctx)

	// Канал с задачами, которые были успешно выполнены.
	taskOkChannel := make(chan Task)
	// Канал с ошибками, возникшими при выполнении задач.
	taskErrChannel := make(chan error)

	// Исполнитель задач.
	wg.Add(1)
	go func() {
		defer wg.Done()
		for task := range taskSourceChannel {
			err := task.Execute()
			if err != nil {
				taskErrChannel <- err
				log.Printf("%s errored out", task)
				continue
			}
			taskOkChannel <- task
			log.Printf("%s was executed successfully", task)
		}
		close(taskOkChannel)
		close(taskErrChannel)
		log.Printf("Executor exited")
	}()

	okResults := make(map[TaskID]Task)
	errors := []error{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for task := range taskOkChannel {
			log.Printf("Processing %s result for presentation", task)
			okResults[task.id] = task
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for err := range taskErrChannel {
			log.Printf("Processing error for presentation: %v", err)
			errors = append(errors, err)
		}
	}()

	wg.Wait()

	fmt.Print("\nSUMMARY:\n\n")

	fmt.Print("\n\tERRORS\n\n")
	for _, v := range errors {
		fmt.Printf("\t%v\n", v)
	}
	fmt.Print("\n")

	fmt.Print("\tDONE TASKS\n\n")
	for _, v := range okResults {
		fmt.Printf("\t%s\n", v)
	}
	fmt.Print("\n")
}
