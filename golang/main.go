package main

import (
	"context"
	"fmt"
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
type Ttype struct {
	id           int64
	creationTime string // время создания
	finishTime   string // время выполнения
	fail         bool
	Result       []byte
}

func main() {
	taskCreator := func(ctx context.Context, a chan Ttype) {
		go func() {
			for {
				select {
				case <-ctx.Done():
					close(a)
					return
				default:
					creationTime := time.Now().Format(time.RFC3339)
					fail := false
					if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
						fail = true
					}
					a <- Ttype{fail: fail, creationTime: creationTime, id: time.Now().UnixNano()} // передаем таск на выполнение
				}
			}
		}()
	}
	ctx, stopCreator := context.WithCancel(context.Background())

	taskWorker := func(task Ttype) Ttype {
		creationTime, _ := time.Parse(time.RFC3339, task.creationTime)
		if task.fail || !creationTime.After(time.Now().Add(-20*time.Second)) {
			task.Result = []byte("something went wrong")
			task.fail = true
		} else {
			task.Result = []byte("task has been successed")
		}
		task.finishTime = time.Now().Format(time.RFC3339Nano)

		time.Sleep(time.Millisecond * 150)

		return task
	}

	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)
	defer close(doneTasks)
	defer close(undoneTasks)

	taskSorter := func(ctx context.Context, t Ttype) {
		if !t.fail {
			select {
			case <-ctx.Done():
				return
			default:
				doneTasks <- t
			}
		} else {
			select {
			case <-ctx.Done():
				return
			default:
				undoneTasks <- fmt.Errorf("task id %d time %s, error: %s", t.id, t.creationTime, t.Result)
			}
		}
	}

	taskChan := make(chan Ttype, 10)
	result := map[int64]Ttype{}
	err := []error{}
	
	go taskCreator(ctx, taskChan)

	go func() {
		// получение тасков
		for task := range taskChan {
			processedTask := taskWorker(task)
			go taskSorter(ctx, processedTask)
		}
	}()

	go func() {
		for task := range doneTasks {
			result[task.id] = task
		}
	}()

	go func() {
		for task := range undoneTasks {
			err = append(err, task)
		}
	}()

	time.Sleep(time.Second * 3)
	stopCreator()
	println("Errors:")
	for _, r := range err {
		println(r.Error())
	}

	println("Done tasks:")
	for r := range result {
		println(r)
	}
}
