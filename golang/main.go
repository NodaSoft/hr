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

// A tType represents a meaninglessness of our life
type tType struct {
	id         int
	cT         string // время создания
	fT         string // время выполнения
	taskRESULT []byte
}

func main() {
	var wg sync.WaitGroup
	superChan := make(chan tType, 10)
	doneTasks := make(chan tType, 10)
	undoneTasks := make(chan error, 10)
	stopChan := make(chan struct{})

	taskCreator := func(a chan tType) {
		defer wg.Done()
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-stopChan:
				return
			case <-ticker.C:
				return
			default:
				ft := time.Now().Format(time.RFC3339)
				if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
					ft = "an error has occurred"
				}
				a <- tType{cT: ft, id: int(time.Now().Unix())} // передаем таск на выполнение
			}
		}
	}

	wg.Add(1)
	go taskCreator(superChan)

	taskWorker := func(a tType) tType {
		tt, _ := time.Parse(time.RFC3339, a.cT)
		if tt.After(time.Now().Add(-20 * time.Second)) {
			a.taskRESULT = []byte("task has been completed successfully")
		} else {
			a.taskRESULT = []byte("something went wrong")
		}
		a.fT = time.Now().Format(time.RFC3339Nano)

		time.Sleep(time.Millisecond * 150)

		return a
	}

	taskSorter := func(t tType) {
		defer wg.Done()
		if string(t.taskRESULT) == "task has been completed successfully" {
			doneTasks <- t
		} else {
			undoneTasks <- fmt.Errorf("task id %d time %s, error %s", t.id, t.cT, t.taskRESULT)
		}
	}

	go func() {
		for t := range superChan {
			wg.Add(1)
			go taskSorter(taskWorker(t))
		}
	}()

	result := map[int]tType{}
	err := []error{}
	var mu sync.Mutex

	go func() {
		for r := range doneTasks {
			mu.Lock()
			result[r.id] = r
			mu.Unlock()
		}
	}()

	go func() {
		for e := range undoneTasks {
			mu.Lock()
			err = append(err, e)
			mu.Unlock()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				mu.Lock()
				fmt.Println("Errors:")
				for _, e := range err {
					fmt.Println(e)
				}
				fmt.Println("Done tasks:")
				for id, r := range result {
					fmt.Printf("ID: %d, Result: %s\n", id, r.taskRESULT)
				}
				mu.Unlock()
			case <-stopChan:
				return
			}
		}
	}()

	time.Sleep(10 * time.Second)
	close(stopChan)
	wg.Wait()

	close(doneTasks)
	close(undoneTasks)
	close(superChan)

	// Последние принты
	mu.Lock()
	fmt.Println("Final Errors:")
	for _, e := range err {
		fmt.Println(e)
	}
	fmt.Println("Final Done tasks:")
	for id, r := range result {
		fmt.Printf("ID: %d, Result: %s\n", id, r.taskRESULT)
	}
	mu.Unlock()
}
