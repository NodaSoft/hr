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
	id     int
	cT     string // время создания
	fT     string // время выполнения
	valid  bool
	result []byte
}

func (t *Task) work() {
	createTime, _ := time.Parse(time.RFC3339, t.cT)
	if createTime.After(time.Now().Add(-20 * time.Second)) {
		t.valid = true
		t.result = []byte("task has been successed")
	} else {
		t.result = []byte("something went wrong")
	}
	t.fT = time.Now().Format(time.RFC3339Nano)

	time.Sleep(time.Millisecond * 150)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	taskChan := make(chan *Task, 10)

	// Горутина для генерации тасков
	go func() {
	loop:
		for {
			select {
			case <-ctx.Done():
				break loop
			default:
				now := time.Now()
				ct := now.Format(time.RFC3339)
				if now.Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
					ct = "Some error occurred"
				}
				taskChan <- &Task{id: int(now.Unix()), cT: ct, valid: false} // передаем таск на выполнение
			}
		}
		close(taskChan)
	}()

	doneTasks := make(chan *Task)
	taskErrors := make(chan error)

	// Горутина сортировки тасков
	go func() {
		// получение тасков
		for t := range taskChan {
			t.work()
			if t.valid {
				doneTasks <- t
			} else {
				taskErrors <- fmt.Errorf("task id %d time %s, error %s", t.id, t.cT, t.result)
			}
		}
	}()

	result := map[int]*Task{}
	// Горутина для сохранения валидных тасков
	go func() {
		for r := range doneTasks {
			result[r.id] = r
		}
	}()

	errors := []error{}
	// Горутина для сохранения ошибок
	go func() {
		for r := range taskErrors {
			errors = append(errors, r)
		}
	}()

	// Ждём
	time.Sleep(time.Second * 3)
	// Всё закрываем
	cancel()
	close(doneTasks)
	close(taskErrors)

	println("Errors:")
	for _, err := range errors {
		println(err)
	}

	println("Done tasks:")
	for _, r := range result {
		println(r)
	}
}
