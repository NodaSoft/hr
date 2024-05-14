package main

import (
	"context"
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
	cT         string // время создания
	fT         string // время выполнения
	taskRESULT []byte
}

const (
	_aggregateDuration = 3
)

type result struct {
	mu   sync.RWMutex
	data map[int]Ttype
}

func newResult() *result {
	return &result{
		data: make(map[int]Ttype),
	}
}

func generateTasks(a chan Ttype) {
	for { // i := 0; i < 10; i++ {
		ft := time.Now().Format(time.RFC3339)
		if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
			ft = "Some error occured"
		}
		a <- Ttype{cT: ft, id: int(time.Now().Unix())} // передаем таск на выполнение // Таски имюет одинаковые id и пишутся в мапу по одному ключу. Чтобы увидеть реальное кол-во успешных тасок, поможет .Nanosecond()
	}
	close(a)
}

func process(a Ttype) Ttype {
	tt, _ := time.Parse(time.RFC3339, a.cT)
	if tt.After(time.Now().Add(-20 * time.Second)) {
		a.taskRESULT = []byte("task has been successed")
	} else {
		a.taskRESULT = []byte("something went wrong")
	}
	a.fT = time.Now().Format(time.RFC3339Nano)

	time.Sleep(time.Millisecond * 150)

	return a
}

func sortResult(t Ttype, doneTasks chan<- Ttype, undoneTasks chan<- error) {
	if isSuccessful(t) {
		doneTasks <- t
	} else {
		undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT, t.taskRESULT)
	}
}

func isSuccessful(t Ttype) bool {
	return string(t.taskRESULT[14:]) == "successed"
}

func processTasks(tasks <-chan Ttype, done chan<- Ttype, undone chan<- error) {
	wg := &sync.WaitGroup{}

	for t := range tasks {
		wg.Add(1)
		go func(task Ttype) {
			sortResult(process(task), done, undone)
			wg.Done()
		}(t)
	}
	wg.Wait()

	close(done)
	close(undone)
}

func aggregateResults(ctx context.Context, done <-chan Ttype, undone <-chan error) (*result, []error) {
	res := newResult()
	var err []error

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func(ctx context.Context) {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case r, ok := <-done:
				if !ok {
					return
				}
				res.mu.Lock()
				res.data[r.id] = r
				res.mu.Unlock()
			}
		}
	}(ctx)

	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case r, ok := <-undone:
				if !ok {
					return
				}
				err = append(err, r)
			}
		}
	}(ctx)

	wg.Wait()
	return res, err
}

func outputResults(res *result, errors []error) {
	println("Errors:")
	for _, err := range errors {
		println(err.Error())
	}

	println("Done tasks:")
	res.mu.RLock()
	for r := range res.data {
		println(r)
	}
	res.mu.RUnlock()
}

func main() {
	inputChan := make(chan Ttype, 10)

	go generateTasks(inputChan)

	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)

	go processTasks(inputChan, doneTasks, undoneTasks)

	// "После обработки тасков в течении 3 секунд приложение должно выводить накопленные к этому моменту успешные таски и отдельно ошибки обработки тасков."
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*_aggregateDuration)
	defer cancel()

	res, errs := aggregateResults(ctx, doneTasks, undoneTasks)
	outputResults(res, errs)
}
