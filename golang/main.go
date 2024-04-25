package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// Приложение эмулирует получение и обработку неких тасков. Пытается и получать, и обрабатывать в многопоточном режиме.
// Должно выводить успешные таски и ошибки обработки тасков по мере выполнения.
// ЗАДАНИЕ: сделать из плохого кода хороший и рабочий - as best as you can.
// Важно сохранить логику появления ошибочных тасков.
// Сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через pull-request в github
// Как видите, никаких привязок к внешним сервисам нет - полный карт-бланш на модификацию кода.

type Task struct {
	id        int
	createdAt string
}

type HandledTask struct {
	task       Task
	finishedAt string
	err        error
}

func createTasks(done <-chan struct{}) <-chan Task {
	taskCh := make(chan Task, 10)
	go func() {
		defer close(taskCh)
		for {
			select {
			case <-done:
				return
			default:
				ft := time.Now().Format(time.RFC3339)
				if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
					ft = "Some error occured"
				}
				taskCh <- Task{createdAt: ft, id: int(time.Now().Unix())*1000000000 + time.Now().Nanosecond()} // передаем таск на выполнение
			}
		}
	}()
	return taskCh
}

func handleTasks(inTasksCh <-chan Task, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		for t := range inTasksCh {
			time.Sleep(time.Millisecond * 150)

			hT := HandledTask{task: t}
			_, err := time.Parse(time.RFC3339, t.createdAt)
			if err != nil {
				hT.err = err
				fmt.Printf("Error in task %d, error %s\n", hT.task.id, hT.err)
			} else {
				hT.finishedAt = time.Now().Format(time.RFC3339Nano)
				fmt.Printf("Done task %d at %s\n", hT.task.id, hT.finishedAt)
			}
		}
		wg.Done()
	}()
}

func main() {
	done := make(chan struct{})

	producedTasksCh := createTasks(done)

	var wg sync.WaitGroup
	CPUCount := runtime.NumCPU()
	for i := 0; i < CPUCount; i++ {
		handleTasks(producedTasksCh, &wg)
	}

	time.Sleep(time.Second * 3)
	close(done)
	wg.Wait()
}
