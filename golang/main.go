package main

import (
	"fmt"
	"sync"
	"time"
)

// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнения остальных тасков

// A Ttype represents a meaninglessness of our life
type Ttype struct {
	id         int
	cT         string // время создания
	fT         string // время выполнения
	taskRESULT []byte
}

func main() {
	taskCreturer := func(a chan Ttype) {
		go func() {
			for {
				ft := time.Now().Format(time.RFC3339)
				if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
					ft = "Some error occured"
				}
				a <- Ttype{cT: ft, id: int(time.Now().Unix())} // передаем таск на выполнение
			}
		}()
	}

	superChan := make(chan Ttype, 10)

	go taskCreturer(superChan)

	var wg sync.WaitGroup
	wg.Add(10)

	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)
	task_worker := func(a Ttype, wg *sync.WaitGroup, doneTasks chan Ttype, undoneTasks chan error) {
		defer wg.Done()
		tt, _ := time.Parse(time.RFC3339, a.cT)
		if tt.After(time.Now().Add(-20 * time.Second)) {
			a.taskRESULT = []byte("task has been successed")
			select {
			case doneTasks <- a:
			default:
			}
		} else {
			a.taskRESULT = []byte("something went wrong")
			select {
			case undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", a.id, a.cT, a.taskRESULT):
			default:
			}
		}
		a.fT = time.Now().Format(time.RFC3339Nano)
		time.Sleep(time.Millisecond * 150)
	}

	go func() {
		for t := range superChan {
			go task_worker(t, &wg, doneTasks, undoneTasks)
		}
		wg.Wait()
		close(doneTasks)
		close(undoneTasks)
	}()

	var mu sync.Mutex
	result := map[int]Ttype{}
	err := []error{}

	go func() {
		for r := range doneTasks {
			mu.Lock()
			result[r.id] = r
			mu.Unlock()
		}
	}()

	go func() {
		for r := range undoneTasks {
			mu.Lock()
			err = append(err, r)
			mu.Unlock()
		}
	}()

	time.Sleep(time.Second * 3)

	println("Errors:")
	for _, r := range err {
		println(r)
	}

	println("Done tasks:")
	for k := range result {
		println(k)
	}
}
