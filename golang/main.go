package main

import (
	"context"
	"fmt"
	"strings"
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

// A Ttype represents a meaninglessness of our life
type Task struct {
	id         int
	createdAt  string
	finishedAt string
	status     string
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tasks := make(chan Task)
	handled := make(chan Task)

	results := make(map[int]Task)
	errs := make([]error, 0)

	duration := time.Duration(150) * time.Millisecond
	go taskGenerator(ctx, tasks, duration)

	go func() {
		for t := range handled {
			if strings.Contains(t.status, "successed") {
				results[t.id] = t
			} else {
				errs = append(errs, fmt.Errorf("task id %d time %s, error %s", t.id, t.createdAt, t.status))
			}
		}
	}()

	ticker := time.NewTicker(3 * time.Second)

	for {
		select {
		case <-ctx.Done():
			close(tasks)
			close(handled)
			fmt.Println("Time's up!")
			return
		case task := <-tasks:
			handled <- taskHadler(task)
		case <-ticker.C:
			fmt.Println("Errors:")
			for _, e := range errs {
				fmt.Println(e)
			}
			fmt.Println("Done tasks:")
			for _, r := range results {
				fmt.Println(r)
			}
		}
	}
}

func taskGenerator(ctx context.Context, superChan chan Task, duration time.Duration) {
	tick := time.NewTicker(duration)

	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			now := time.Now()

			if now.Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
				superChan <- Task{
					id:        int(time.Now().UnixMilli()),
					createdAt: now.Format(time.RFC3339),
					status:    "Some error occured",
				}
			} else {
				superChan <- Task{
					id:        int(time.Now().Unix()),
					createdAt: now.Format(time.RFC3339),
					status:    "Created",
				}
			}
		}
	}
}

func taskHadler(task Task) Task {
	created, err := time.Parse(time.RFC3339, task.createdAt)
	if err != nil {
		fmt.Println(err)
	}

	if created.After(time.Now().Add(-20 * time.Second)) {
		if strings.Contains(task.status, "error") {
			task.finishedAt = time.Now().Format(time.RFC3339)
		} else {
			task.finishedAt = time.Now().Format(time.RFC3339)
			task.status = "task has been successed"
		}
	} else {
		task.finishedAt = time.Now().Format(time.RFC3339)
		task.status = "something went wrong"
	}

	return task
}
