package main

import (
	"fmt"
	"runtime"
	"sync"
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
type Ttype struct {
	id         int
	cT         time.Time // время создания
	fT         time.Time // время выполнения
	taskRESULT []byte
}

func main() {
	taskCreateUsers := func(usersChan chan<- Ttype, err chan<- error, wg *sync.WaitGroup) {
		defer wg.Done()

		for i := 0; i < 10; i++ {
			ft := time.Now()

			if (time.Now().Nanosecond()/10000)%2 > 0 { // вот такое условие появления ошибочных тасков
				err <- fmt.Errorf("some error occurred")
				continue
			}
			usersChan <- Ttype{cT: ft, id: int(time.Now().Unix())} // передаем таск на выполнение
		}

		close(usersChan)
	}

	// All errors when users created or worker job
	errors := make(chan error)
	usersChan := make(chan Ttype, 10)

	var wgCreateUsers sync.WaitGroup

	wgCreateUsers.Add(1)

	go taskCreateUsers(usersChan, errors, &wgCreateUsers)

	taskWorker := func(usersChan <-chan Ttype, resChan chan<- Ttype, errChan chan error, wg *sync.WaitGroup) {
		defer wg.Done()

		for u := range usersChan {
			// More than 150 ms past since user created
			if time.Since(u.cT) > 150*time.Millisecond {
				errChan <- fmt.Errorf("something went wrong")
				continue
			}

			u.taskRESULT = []byte("task has been successed")
			u.fT = time.Now()

			// Time consuming task
			time.Sleep(time.Millisecond * 150)

			resChan <- u
		}
	}

	var wgWorkers sync.WaitGroup

	doneTasks := make(chan Ttype)

	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		wgWorkers.Add(1)

		go taskWorker(usersChan, doneTasks, errors, &wgWorkers)
	}

	// Close doneTasks channel when all workers has finished
	go func() {
		wgWorkers.Wait()
		close(doneTasks)
	}()

	// Close error channel when createUsers and workers finished
	go func() {
		wgCreateUsers.Wait()
		wgWorkers.Wait()
		close(errors)
	}()

	// Wait for done & error channels read out
	var wg sync.WaitGroup

	wg.Add(2)

	var result []Ttype

	go func() {
		defer wg.Done()

		for r := range doneTasks {
			result = append(result, r)
		}
	}()

	var err []error

	go func() {
		defer wg.Done()

		for r := range errors {
			err = append(err, r)
		}
	}()

	wg.Wait() // wait for all results and errors received

	println("Errors:")
	for _, e := range err {
		println(e.Error())
	}

	println("Done tasks:")
	for _, r := range result {
		println(r.id)
	}
}
