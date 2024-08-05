package main

import (
	"fmt"
	"sync"
	"time"
)

type Ttype struct {
	id           int
	creationTime string // исходное время создания
	cT           string // время создания или ошибка
	fT           string // время выполнения
	taskRESULT   []byte
}

func main() {
	superChan := make(chan Ttype, 10)
	doneTasks := make(chan Ttype, 10)
	undoneTasks := make(chan Ttype, 10)
	var wg sync.WaitGroup
	var idCounter int
	var mu sync.Mutex

	// Функция для генерации тасков
	taskCreator := func(a chan Ttype) {
		go func() {
			defer close(a)
			start := time.Now()
			for time.Since(start) < 10*time.Second {
				mu.Lock()
				idCounter++
				mu.Unlock()
				creationTime := time.Now().Format(time.RFC3339)
				cT := creationTime
				if time.Now().Nanosecond()%2 > 0 { // условие появления ошибочных тасков
					cT = "Some error occurred"
				}
				a <- Ttype{id: idCounter, creationTime: creationTime, cT: cT} // передаем таск на выполнение
				time.Sleep(time.Millisecond * 100) // Предотвращаем переполнение канала слишком быстро
			}
		}()
	}

	// Функция для обработки тасков
	taskWorker := func(a Ttype) Ttype {
		if a.cT == "Some error occurred" {
			a.taskRESULT = []byte("Some error occurred")
		} else {
			tt, err := time.Parse(time.RFC3339, a.cT)
			if err == nil && tt.After(time.Now().Add(-20*time.Second)) {
				a.taskRESULT = []byte("task has been successful")
			} else {
				a.taskRESULT = []byte("something went wrong")
			}
		}
		a.fT = time.Now().Format(time.RFC3339Nano)
		time.Sleep(time.Millisecond * 150)
		return a
	}

	// Функция для сортировки тасков
	taskSorter := func(t Ttype) {
		if string(t.taskRESULT) == "task has been successful" {
			doneTasks <- t
		} else {
			undoneTasks <- t
		}
	}

	// Запускаем генерацию тасков
	taskCreator(superChan)

	// Запускаем обработку тасков
	go func() {
		for t := range superChan {
			wg.Add(1)
			go func(task Ttype) {
				defer wg.Done()
				processedTask := taskWorker(task)
				taskSorter(processedTask)
			}(t)
		}
		wg.Wait()
		close(doneTasks)
		close(undoneTasks)
	}()

	// Функция для вывода результатов
	printResults := func(doneTasks []Ttype, undoneTasks []Ttype) {
		fmt.Println("Errors:")
		for _, e := range undoneTasks {
			fmt.Printf("Task id %d created at %s error: %s\n", e.id, e.creationTime, string(e.taskRESULT))
		}

		fmt.Println("Done tasks:")
		for _, r := range doneTasks {
			fmt.Printf("Task id %d created at %s finished at %s result: %s\n", r.id, r.creationTime, r.fT, string(r.taskRESULT))
		}
	}

	// Каждые 3 секунды выводим результаты
	go func() {
		for {
			time.Sleep(3 * time.Second)
			var done []Ttype
			var undone []Ttype
			doneLoop:
			for {
				select {
				case t := <-doneTasks:
					done = append(done, t)
				default:
					break doneLoop
				}
			}
			undoneLoop:
			for {
				select {
				case t := <-undoneTasks:
					undone = append(undone, t)
				default:
					break undoneLoop
				}
			}
			printResults(done, undone)
		}
	}()

	// Ожидаем завершения всех горутин
	wg.Wait()

	// Вывод финальных результатов
	var finalDone []Ttype
	var finalUndone []Ttype
	for t := range doneTasks {
		finalDone = append(finalDone, t)
	}
	for t := range undoneTasks {
		finalUndone = append(finalUndone, t)
	}
	printResults(finalDone, finalUndone)
}
