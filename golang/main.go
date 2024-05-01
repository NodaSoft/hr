package main

import (
	"context"
	"fmt"
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
type Task struct {
	id         int
	createTime         string // время создания
	finishTime         string // время выполнения
	taskResult []byte
}

func main() {
	ctx, finish := context.WithTimeout(context.Background(), 3 * time.Second)
	defer finish()

	var(
		doneTasks, createdTasks = make(chan Task), make(chan Task)
		undoneTasks = make(chan error)
		result []Task
		err []error
	)

	go taskCreator(ctx, createdTasks)

WORKERS:
	for {
		select {
		case <-ctx.Done():
			break WORKERS
		case task := <-createdTasks:
			go taskWorker(task, doneTasks, undoneTasks)
		}
	}

RESULT:
	for  {
		select {
		case dTask := <-doneTasks:
			result = append(result, dTask)
		case erTask := <-undoneTasks:
			err = append(err, erTask)
		default:
			break RESULT
		}
	}

	fmt.Println("Errors:")
	for e := range err {
		fmt.Println(e)
	}

	fmt.Println("Done tasks:")
	for _, r := range result {
		fmt.Printf("Task %d created at %s\n", r.id, r.createTime)
	}
}

func taskCreator(ctx context.Context, out chan<- Task)  {
	for {
		select {
		case <-ctx.Done():
			close(out)
			return
		default:
			formatTime := time.Now().Format(time.RFC3339)
			if time.Now().Nanosecond() % 2 > 0 { // вот такое условие появления ошибочных тасков
				formatTime = "Some error occured"
			}
			out <- Task{createTime: formatTime, id: int(time.Now().Unix())}
		}
	}
}

func taskWorker(t Task, d chan<- Task, u chan<- error) {
	t.finishTime = time.Now().Format(time.RFC3339Nano)
	_, err := time.Parse(time.RFC3339, t.createTime)
	if err == nil {
		t.taskResult = []byte("task has been successed")
		d <- t
	} else {
		t.taskResult = []byte("something went wrong")
		u <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.createTime, t.taskResult)
	}

	time.Sleep(time.Millisecond * 150)
}