package main

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнены остальных тасков

// Вариант1: исправление легаси, с минимальными изменениями кода и без вмешательства в структуры данных:
// @author: fvn-2023-04-18

const _TotalRunTime = time.Second * 3

// A Ttype represents a meaninglessness of our life
type Ttype struct {
	id         int
	cT         string // время создания
	fT         string // время выполнения
	taskRESULT []byte
}

func main() {
	superChan := make(chan Ttype, 10)

	go func(a chan Ttype, timeout time.Duration) {
		i := 0
		started := time.Now()
		for time.Since(started) < timeout {
			ft := time.Now().Format(time.RFC3339)
			if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
				// TODO mixed string use, must be refactor
				ft = "Some error occured"
			}
			a <- Ttype{cT: ft, id: int(time.Now().Unix())} // передаем таск на выполнение
			i++
		}
		log.Printf("emulator ended %v, %v, %d", started, time.Now(), i)
		close(a) // канал закрывает писатель!
	}(superChan, _TotalRunTime)

	task_worker := func(a Ttype) Ttype {
		tt, err := time.Parse(time.RFC3339, a.cT)
		if err != nil {
			a.taskRESULT = []byte(a.cT)
		} else {
			if tt.After(time.Now().Add(-20 * time.Second)) {
				a.taskRESULT = []byte("task has been successed")
			} else {
				a.taskRESULT = []byte("something went wrong")
			}
		}
		a.fT = time.Now().Format(time.RFC3339Nano)

		time.Sleep(time.Millisecond * 150)

		return a
	}

	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)

	tasksorter := func(t Ttype) {
		if string(t.taskRESULT[14:]) == "successed" {
			// fvn: запись во внешне закрываемый канал: потенциальная паника!
			doneTasks <- t
		} else {
			// fvn: аналогично..
			undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT, t.taskRESULT)
		}
	}

	go func() {
		i := 0
		var wg sync.WaitGroup
		// получение тасков отдельным потоком
		for t := range superChan {
			// fvn: многопоточная обработка
			go func(t Ttype) {
				wg.Add(1)
				tasksorter(task_worker(t))
				wg.Done()
			}(t)
			i++
		}
		wg.Wait()
		close(doneTasks)
		close(undoneTasks)
		log.Printf("workers ended %d", i)
	}()

	result := []int{}
	errs := []error{}
	var resSync sync.Mutex

	// fvn: можно и через atomic:
	var workers atomic.Int32
	go func() {
		i := 0
		// поток подсчета итогов
		workers.Add(1) // добавляем счетчик обработчиков
		for {
			var (
				ok1, ok2 bool
				task     Ttype
				err      error
			)

			select {
			case task, ok1 = <-doneTasks: // канал закрыт и нет данных?
				if ok1 {
					resSync.Lock()
					result = append(result, task.id)
					resSync.Unlock()
				}
			case err, ok2 = <-undoneTasks: // аналогично
				if ok2 {
					resSync.Lock()
					errs = append(errs, err)
					resSync.Unlock()
				}
			}
			if !ok1 && !ok2 {
				break // from for!
			}
			i++
		}
		log.Println("results ended %d", i)
		workers.Add(-1)
	}()

	time.Sleep(_TotalRunTime)
	for workers.Load() > int32(0) {
		time.Sleep(10 * time.Millisecond)
	}

	log.Printf("Errors: %d", len(errs))
	//for r := range errs {
	//	println(r)
	//}

	log.Printf("Done tasks: %d", len(result))
	//for r := range result {
	//	println(r)
	//}
}
