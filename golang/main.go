package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// статусы выполнения задачи
const (
	// C - style константы используются чтобы визуально отделить константы от обычных переменных
	STATUS_CREATED uint = iota
	STATUS_CREATION_ERROR
	STATUS_EXECUTED
	STATUS_EXECUTION_ERROR
)

// A Ttype represents a meaninglessness of our life
type TaskContext struct {
	ID            int
	CreationTime  time.Time
	ExecutionTime time.Time
	Status        uint
}

func main() {
	// чтобы программа сама завершалась был использован контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	createdTasksChan := TaskCreator(ctx)
	executedTasksChan := TaskWorker(createdTasksChan)
	successfulTasks, errorTasks := TaskSorter(executedTasksChan)

	<-ctx.Done()

	PrintTasks(successfulTasks, errorTasks)
}

// TaskCreator создает задачи и передает дальше по пайплайну
// создание задач искусственно ограничено в 1 миллисекунду
func TaskCreator(ctx context.Context) chan TaskContext {
	out := make(chan TaskContext)

	go func() {
		defer close(out)

		tasksCounter := 0
		// искусственное замедление создание тасок по 1 в миллисекунду. Это сделано чтобы проще можно было тестировать
		// вручную
		ticker := time.NewTicker(time.Millisecond * 1)

		for range ticker.C {
			tasksCounter++
			now := time.Now()

			task := TaskContext{ID: tasksCounter, CreationTime: now, Status: STATUS_CREATED}

			// на ОС Ventura 13 (macOS) наносекунды округляются. Вот ишью https://github.com/golang/go/issues/22037
			// исходя из вышеописанной проблемы допущу вольность и поменяю наносекунды на милисекунды
			// само условие требуется сохранить по задаче
			if now.UnixMilli()%2 > 0 {
				task.Status = STATUS_CREATION_ERROR
			}

			select {
			case <-ctx.Done():
				return
			case out <- task:
			}
		}
	}()

	return out
}

// TaskWorker получает задачи по пайплайну и пере дает дальше
// исполнение задачи искусственно установлено 150 миллисекунд
func TaskWorker(in <-chan TaskContext) chan TaskContext {
	out := make(chan TaskContext)

	go func() {
		defer close(out)

		wg := &sync.WaitGroup{}

		for task := range in {
			if task.Status != STATUS_CREATED {
				task.Status = STATUS_EXECUTION_ERROR
				out <- task
				continue
			}

			wg.Add(1)
			go func(t TaskContext) {
				defer wg.Done()
				// иммитация выполнения процесса
				time.Sleep(time.Millisecond * 150)
				t.ExecutionTime = time.Now()
				t.Status = STATUS_EXECUTED
				out <- t
			}(task)
		}

		wg.Wait()
	}()

	return out
}

// TaskSorter получает выполненные задачи по пайплайну и сохраняет в 2 отдельные массивы.
// Первый для успешных, воторой для неуспешных.
// это сделано для сохранения логики вывода в оригинале“
func TaskSorter(in <-chan TaskContext) ([]TaskContext, []TaskContext) {
	successfulTasks := make([]TaskContext, 0)
	errorTasks := make([]TaskContext, 0)

	for task := range in {
		if task.Status == STATUS_EXECUTED {
			successfulTasks = append(successfulTasks, task)
			continue
		}

		errorTasks = append(errorTasks, task)
	}

	return successfulTasks, errorTasks
}

// PrintTasks выводит на экран успешные и неуспешные задачи
// время выводится в милисекундах чтобы было проще взглядом проверить
func PrintTasks(successfulTasks, errorTasks []TaskContext) {
	fmt.Println("Done tasks:")
	for _, task := range successfulTasks {
		fmt.Printf("ID: %d, created: %d, executed: %d, status: %d\n",
			task.ID, task.CreationTime.UnixMilli(), task.ExecutionTime.UnixMilli(), task.Status)
	}

	errorMsg := "something went wrong"
	for _, task := range errorTasks {
		fmt.Printf("Task id %d time %d, error %s\n", task.ID, task.CreationTime.UnixMilli(), errorMsg)
	}
}
