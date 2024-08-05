package main

import (
	"context"
	"fmt"

	//"math/rand"
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

type Task struct {
	id         int
	startTime  string // время создания
	finishTime string // время выполнения
	err        error
	result     []byte
}

// CreateTask создает и возвращает новое задание
func CreateTask() Task {
	t := Task{id: int(time.Now().Unix()),
		startTime: time.Now().Format(time.RFC3339),
	}

	if time.Now().Nanosecond()%2 > 0 { // Условие генерации задания с ошибкой
		t.err = fmt.Errorf("some error occured")
	}

	/*if n := rand.Intn(10); n%2 != 0 {
		t.err = fmt.Errorf("some error occured") // Алтернативное условие генерации задания с ошибкой
	}*/

	return t
}

// TaskCreator создает в горутине задания с помощью CreateTask и отправляет их в канал tasks
func TaskCreator(ctx context.Context) chan Task {
	tasks := make(chan Task)

	go func() {
		defer close(tasks)
		for {
			select {
			case tasks <- CreateTask():
			case <-ctx.Done():
				return
			}
			time.Sleep(150 * time.Millisecond) // каждые 150 миллисекунд создается новая задача
		}
	}()

	return tasks
}

// Worker принимает задания и "проводит работу", если в задании нет ошибки. Отправляет его в канал success.
// Если ошибка есть, то отправляет задание в канал withErr.
func Worker(ctx context.Context, in chan Task) (chan Task, chan Task) {
	success := make(chan Task)
	withErr := make(chan Task)

	go func() {
		defer close(success)
		defer close(withErr)
		for {
			select {
			case t, ok := <-in:
				if !ok {
					return
				}

				t.finishTime = time.Now().Format(time.RFC3339)

				if t.err != nil {
					t.result = []byte("something went wrong")
					select {
					case withErr <- t:
					case <-ctx.Done():
						return
					}
				} else {
					t.result = []byte("task has been successed")
					select {
					case success <- t:
					case <-ctx.Done():
						return
					}
				}
			case <-ctx.Done():
				return
			}

			time.Sleep(10 * time.Millisecond) // имитация работы
		}
	}()

	return success, withErr
}

// CollectAndPrintSuccess каждые 3 секунды выводит в консоль успешные задачи
func CollectAndPrintSuccess(ctx context.Context, success chan Task) {
	tasks := make(map[int]Task)
	ticker := time.NewTicker(3 * time.Second) // тиккер, чтобы каждые 3 секунды выводить успешные задачи
	go func() {
		defer ticker.Stop()
		for {
			select {
			case t := <-success:
				tasks[t.id] = t
			case <-ticker.C:
				fmt.Println("Success:", tasks)
			case <-ctx.Done():
				return
			}
		}
	}()
}

// CollectAndPrintErr каждые 3 секунды выводит в консоль задачи с ошибкой
func CollectAndPrintErr(ctx context.Context, withErr chan Task) {
	tasks := make(map[int]Task)
	ticker := time.NewTicker(3 * time.Second) // тиккер, чтобы каждые 3 секунды выводить задачи с ошибкой
	go func() {
		defer ticker.Stop()
		for {
			select {
			case t := <-withErr:
				tasks[t.id] = t
			case <-ticker.C:
				fmt.Println("Errors:", tasks)
			case <-ctx.Done():
				return
			}
		}
	}()
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // контекст, который эмулируют работу в течение 10 секунд
	defer cancel()

	newTasks := TaskCreator(ctx)
	success, withErr := Worker(ctx, newTasks)

	CollectAndPrintSuccess(ctx, success)
	CollectAndPrintErr(ctx, withErr)
	<-ctx.Done() // ждем закрытия канала, чтобы main не завершилась раньше времени
}
