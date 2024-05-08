package main

import (
	"fmt"
	"sync"
	"time"
)

// Приложение эмулирует получение и обработку неких тасков. Пытается и получать, и обрабатывать в многопоточном режиме.
// После обработки тасков в течении 3 секунд приложение должно выводить накопленные к этому моменту успешные таски и отдельно ошибки обработки тасков.

// ЗАДАНИЕ: сделать из плохого кода хороший и рабочий - as best as you can.
// Важно сохранить логику появления ошибочных тасков.
// Важно оставить асинхронные генерацию и обработку тасков.
// Сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через pull-request в github
// Как видите, никаких привязок к внешним сервисам нет - полный карт-бланш на модификацию кода.

// A Ttype represents a meaninglessness of our life
type Ttype struct {
	id         int
	cT         time.Time // время создания
	fT         time.Time // время выполнения
	taskRESULT string
}

func main() {
	taskCreator := func(a chan Ttype) {
		for {
			ct := time.Now()
			id := ct.UnixNano()
			ft := ct
			var err string
			nano := ct.Nanosecond() / 1000
			if nano%2 > 0 { // вот такое условие появления ошибочных тасков
				err = "Some error occured"
			} else {
				err = "ok"
			}
			a <- Ttype{id: int(id), cT: ct, fT: ft, taskRESULT: err} // передаем таск на выполнение
			time.Sleep(100 * time.Millisecond)
		}
	}

	superChan := make(chan Ttype, 10)

	go taskCreator(superChan)
	done := make(chan Ttype)
	undone := make(chan error, 100)
	var wg sync.WaitGroup

	task_worker := func(a Ttype, wg *sync.WaitGroup, done chan<- Ttype, undone chan<- error) {
		defer wg.Done()
		if a.cT.After(time.Now().Add(-20*time.Second)) && a.taskRESULT == "ok" {
			a.taskRESULT = "task has been successed"
			a.fT = time.Now()
			done <- a
		} else {
			a.taskRESULT = "something went wrong"
			undone <- fmt.Errorf("Task id %d time %s, error %s", a.id, a.cT, a.taskRESULT)
			return
		}
	}

	go func() {
		// получение тасков
		for t := range superChan {
			wg.Add(1)
			go task_worker(t, &wg, done, undone)
		}
		close(superChan)
	}()

	err := []error{}
	var result sync.Map
	var mu sync.Mutex

	go func() {
		for {
			select {
			case s, ok := <-done:
				if !ok {
					done = nil
					continue
				}
				result.Store(s.id, s)
			case e, ok := <-undone:
				if !ok {
					undone = nil
					continue
				}
				mu.Lock()
				err = append(err, e)
				mu.Unlock()
			}
		}
	}()

	time.Sleep(time.Second * 3)
	wg.Wait()
	close(done)
	close(undone)

	fmt.Println("Errors:")
	for _, e := range err {
		fmt.Println(e)
	}

	fmt.Println("Done tasks:")
	result.Range(func(key, value interface{}) bool {
		s := value.(Ttype)
		fmt.Printf("Task id %d completed at %s\n", s.id, s.fT)
		return true
	})
}
