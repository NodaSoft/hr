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
	cT         string // время создания
	fT         string // время выполнения
	taskRESULT bool
	err        error
}

func main() {
	taskChan := make(chan *Ttype, 10)
	resultChan := make(chan *Ttype)
	var wg sync.WaitGroup

	// Создаем задачи
	wg.Add(1)
	go func() {
		defer wg.Done()
		createTask(taskChan)
		close(taskChan) // Закрываем канал задач после их создания
	}()

	// Выполняем задачи
	wg.Add(1)
	go func() {
		defer wg.Done()
		startTask(taskChan, resultChan)
	}()

	// Сортируем задачи
	wg.Add(1)
	go func() {
		defer wg.Done()
		sortTask(resultChan)
	}()

	wg.Wait()
	fmt.Println("Все задачи завершены")
}

func createTask(a chan<- *Ttype) {
	for i := 0; i < 10; i++ {
		ft := time.Now().Format(time.RFC3339)
		if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
			ft = "Some error occurred"
		}
		a <- &Ttype{cT: ft, id: i} // передаем таск на выполнение
		time.Sleep(time.Millisecond * 100)
	}
}

func startTask(taskChan <-chan *Ttype, resultChan chan<- *Ttype) {
	for task := range taskChan {
		tt, _ := time.Parse(time.RFC3339, task.cT)
		if tt.After(time.Now().Add(-20 * time.Second)) {
			task.taskRESULT = true
		}
		task.fT = time.Now().Format(time.RFC3339Nano)
		time.Sleep(time.Millisecond * 150)
		resultChan <- task
	}
	close(resultChan)
}

func sortTask(resultChan <-chan *Ttype) {
	results := make(map[int]*Ttype)
	errors := make([]error, 0)

	for task := range resultChan {
		if task.taskRESULT {
			results[task.id] = task
		} else {
			task.err = fmt.Errorf("task id %d time %s, error something went wrong", task.id, task.cT)
			errors = append(errors, task.err)
		}
	}

	// Вывод ошибок
	fmt.Println("Errors:")
	for _, err := range errors {
		fmt.Println(err)
	}

	// Вывод завершенных задач
	fmt.Println("Done tasks:")
	for id, task := range results {
		fmt.Printf("Task ID: %d, Creation Time: %s, Finish Time: %s\n", id, task.cT, task.fT)
	}
}
