package main

import (
	"context"
	"fmt"
	"math/rand"
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
type Task struct {
	id         int
	createTime string // время создания
	finishTime string // время выполнения
	taskResult []byte
}

const timeout = 1 // время таймаута пишем свое (за 1мс выполняется 50-100 итераций +-)

func main() {
	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*timeout)

	// функция принимает канал и пушит в него таск
	taskCreaturer := func(ctx context.Context, created chan Task) {
		for {
			// завершаем цикл при таймауте
			select {
			case <-ctx.Done():
				return
			default:
			}

			// пишем в канал таски
			go func() {
				formatedTime := time.Now().Format(time.RFC3339)
				if rand.Int()%2 == 0 { // https://stackoverflow.com/questions/57285292 | Запускал код на винде, поэтому немного сменил условие
					formatedTime = "Some error occured"
				}
				created <- Task{createTime: formatedTime, id: rand.Int()} // передаем таск на выполнение
			}()
		}
	}

	createdTaskChan := make(chan Task, 10)
	performedTaskChan := make(chan Task, 10)

	// функция принимает свеженькие таски и пушит выполненные в канал
	taskWorker := func(created, performed chan Task) {
		for {
			go func(t Task) {
				_, err := time.Parse(time.RFC3339, t.createTime)
				if err != nil {
					t.taskResult = []byte("something went wrong")
				} else {
					t.taskResult = []byte("Success: *success details*")
				}
				t.finishTime = time.Now().Format(time.RFC3339Nano)

				time.Sleep(time.Millisecond * 150)
				performed <- t
			}(<-created)
		}
	}

	doneTasks := make(chan Task, 100)
	undoneTasks := make(chan error, 100)

	// функция принимает канал с выполнеными тасками и сортирует по каналам с успешными и неудачными
	taskSorter := func(performed, successful chan Task, failed chan error) {
		for {
			go func(t Task) {
				if string(t.taskResult[:7]) == "Success" {
					successful <- t
				} else {
					failed <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.createTime, t.taskResult)
				}
			}(<-performed)
		}
	}

	// запускаем создание тасков
	go taskCreaturer(ctx, createdTaskChan)

	// выполяем таски и пушим в performed
	go taskWorker(createdTaskChan, performedTaskChan)

	// сортируем таски, которые приходят в performed канал
	go taskSorter(performedTaskChan, doneTasks, undoneTasks)

	result := map[int]Task{}
	errors := []error{}

	go func() {
		/* убрал горутины, потому что иначе возникает ошибка конкурентной записи
		(если сильно надо писать асинхронно, можно юзнуть мьютексы) */
		for r := range doneTasks {
			result[r.id] = r
		}
	}()

	go func() {
		for r := range undoneTasks {
			errors = append(errors, r)
		}
	}()

	time.Sleep(time.Second * 1)

	// выводим данные
	println("Errors:")
	for r := range errors {
		fmt.Println(r)
	}

	println("Done tasks:")
	for r := range result {
		fmt.Println(r)
	}
}
