package main

import (
	"context"
	"fmt"
	"time"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнены остальных тасков

// A Ttype represents a meaninglessness of our life
type Ttype struct {
	id             int
	createTime     time.Time
	completionTime time.Time
	taskResult     []byte
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	taskCreator := func(ctx context.Context, tasksQueueChan chan Ttype) {
		id := 0

		for {
			select {
			case <-ctx.Done():
				close(tasksQueueChan)
				return
			default:
				var result []byte

				ct := time.Now()

				// На сколько я помню, точность наносекунд на современных ОС составляет сотни и тысячи. (у меня, по крайней мере это тысячи)
				// По этому с ошибками не будет генериться ни одной таски. По этому я подкорректировал условие
				ns := time.Now().Nanosecond() / 1000

				if ns%2 > 0 { // вот такое условие появления ошибочных тасков
					result = []byte("wrong task generated")
				}

				// Unix timestamp выражается в секундах. А, как минимум, первые 10 тасков будут генериться менее чем засекунду
				// По этому id таски, как минимум, у первых 10 штук будет одинаковый. По этому я заменил его на что-то более уникальное
				tasksQueueChan <- Ttype{id: id, createTime: ct, taskResult: result} // передаем таск на выполнение
				id++
			}
		}
	}

	tasksQueueChan := make(chan Ttype, 10)

	go taskCreator(ctx, tasksQueueChan)

	taskWorker := func(task Ttype) Ttype {
		switch {
		case len(task.taskResult) != 0:
			task.taskResult = []byte(fmt.Sprintf("something went wrong: %s", task.taskResult))
		case task.createTime.Before(time.Now().Add(-20 * time.Second)):
			task.taskResult = []byte("task is expired")
		default:
			task.taskResult = []byte("task has been succeed")
			task.completionTime = time.Now()
		}

		time.Sleep(time.Millisecond * 150)

		return task
	}

	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)

	taskSorter := func(task Ttype) {
		if string(task.taskResult[14:]) == "succeed" {
			doneTasks <- task
		} else {
			undoneTasks <- fmt.Errorf(
				"task id: %d, created at: %v, error: %s",
				task.id,
				task.createTime.Format(time.RFC3339),
				string(task.taskResult),
			)
		}
	}

	go func() {
		// получение тасков
		for t := range tasksQueueChan {
			t = taskWorker(t)
			go taskSorter(t)
		}

		close(tasksQueueChan)
		close(doneTasks)
		close(undoneTasks)
	}()

	result := make(map[int]Ttype, 0)

	go func() {
		for r := range doneTasks {
			result[r.id] = r
		}
	}()

	err := make([]error, 0)

	go func() {
		for r := range undoneTasks {
			err = append(err, r)
		}
	}()

	time.Sleep(time.Second * 3)

	println("Errors:")

	for i := range err {
		println(err[i].Error())
	}

	println("Done tasks:")

	for _, task := range result {
		fmt.Printf(
			"Task id: %d, created at: %s, result: %s\n",
			task.id,
			task.createTime.Format(time.RFC3339),
			task.taskResult,
		)
	}
}
