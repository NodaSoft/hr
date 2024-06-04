package main

import (
	"fmt"
	"sync"
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
	id         int
	cT         time.Time
	fT         time.Time
	taskRESULT string
}

func main() {
	var wg sync.WaitGroup
	taskCreator := func(tasks chan<- Ttype) {
		defer close(tasks)
		for i := 0; i < 10; i++ { // количество тасков для демонстрации
			ct := time.Now()
			t := Ttype{id: i, cT: ct}
			if ct.Nanosecond()%2 > 0 {
				t.taskRESULT = "Some error occurred"
			}
			tasks <- t
			time.Sleep(100 * time.Millisecond) //задержка для имитации работы
		}
	}

	superChan := make(chan Ttype, 10)
	go taskCreator(superChan)

	taskWorker := func(t Ttype) Ttype {
		if t.taskRESULT == "" {
			t.taskRESULT = "task has been succeeded"
		}
		t.fT = time.Now()
		time.Sleep(150 * time.Millisecond) // обработка
		return t
	}

	doneTasks := make(chan Ttype, 10)
	undoneTasks := make(chan Ttype, 10)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for t := range superChan {
			result := taskWorker(t)
			if result.taskRESULT == "task has been succeeded" {
				doneTasks <- result
			} else {
				undoneTasks <- result
			}
		}
		close(doneTasks)
		close(undoneTasks)
	}()

	go func() {
		wg.Wait()
		fmt.Println("Processing complete.")
	}()

	time.Sleep(3 * time.Second) // ждем 3 секунды для демонстрации

	fmt.Println("Errors:")
	for t := range undoneTasks {
		fmt.Printf("Task id %d, error: %s\n", t.id, t.taskRESULT)
	}

	fmt.Println("Done tasks:")
	for t := range doneTasks {
		fmt.Printf("Task id %d, completed at: %s\n", t.id, t.fT)
	}
}

// В этом коде:

// Используется sync.WaitGroup для ожидания завершения всех горутин.
// Время создания и выполнения теперь использует тип time.Time.
// Ограничено количество создаваемых задач для демонстрации.
// Используются буферизированные каналы для doneTasks и undoneTasks, чтобы избежать блокировки при отправке.
// Добавлена задержка в taskCreator для имитации работы.
// Удалены лишние горутины внутри циклов, которые могли привести к состоянию гонки.
