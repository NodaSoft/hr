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

// A Task represents a meaninglessness of our life
type Task struct {
	id         int
	createdAt  string
	finishedAt string
	taskRESULT []byte
}

func createTasks() <-chan Task {
	ch := make(chan Task)

	go func() {
		defer close(ch)

		timer := time.After(10 * time.Second)
		for {
			select {
			case <-timer:
				return
			default:
				ft := time.Now().Format(time.RFC3339)
				if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
					ft = "Some error occured"
				}
				ch <- Task{createdAt: ft, id: int(time.Now().Unix())} // передаем таск на выполнение

				time.Sleep(time.Second) // для обеспечения уникальности id и читабельности вывода
			}
		}
	}()

	return ch
}

func handleTasks(tasks <-chan Task) (<-chan Task, <-chan error) {
	doneTasks := make(chan Task)
	undoneTasks := make(chan error)

	go func() {
		var wg sync.WaitGroup
		for t := range tasks {
			wg.Add(1)
			go func() {
				defer wg.Done()

				tt, _ := time.Parse(time.RFC3339, t.createdAt)
				if tt.After(time.Now().Add(-20 * time.Second)) {
					t.taskRESULT = []byte("task has been successed")
				} else {
					t.taskRESULT = []byte("something went wrong")
				}
				t.finishedAt = time.Now().Format(time.RFC3339Nano)

				time.Sleep(time.Millisecond * 150)

				if string(t.taskRESULT[14:]) == "successed" {
					doneTasks <- t
				} else {
					undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.createdAt, t.taskRESULT)
				}
			}()
		}
		wg.Wait()

		close(doneTasks)
		close(undoneTasks)
	}()

	return doneTasks, undoneTasks
}

func handleResults(doneTasks <-chan Task, undoneTasks <-chan error) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	errs := make([]error, 0)
	done := make([]int, 0)

	flushResult := func() {
		if len(errs) == 0 && len(done) == 0 {
			return
		}

		println("Errors:")
		for _, err := range errs {
			println(err.Error())
		}
		errs = make([]error, 0)

		println("Done tasks:")
		for _, d := range done {
			println(d)
		}
		done = make([]int, 0)

		println()
	}

	doneHandled, undoneHandled := false, false
	for !doneHandled && !undoneHandled {
		select {
		case <-ticker.C:
			flushResult()
		case err, ok := <-undoneTasks:
			if ok {
				errs = append(errs, err)
			} else {
				undoneHandled = true
			}
		case task, ok := <-doneTasks:
			if ok {
				done = append(done, task.id)
			} else {
				doneHandled = true
			}
		}
	}
	flushResult()
}

func main() {
	tasks := createTasks()
	doneTasks, undoneTasks := handleTasks(tasks)
	handleResults(doneTasks, undoneTasks)
}
