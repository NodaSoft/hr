package main

import (
	"context"
	"fmt"
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
	cT         string // время создания
	fT         string // время выполнения
	taskResult []byte
}

func NewTtype() *Ttype {
	t := time.Now()
	ct := t.Format(time.RFC3339)
	if t.Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков (странное условие но пусть будет)
		ct = "Some error occured"
	}
	return &Ttype{cT: ct, id: int(t.Unix())}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	superChan := make(chan *Ttype, 10)
	doneTasks := make(chan *Ttype)
	undoneTasks := make(chan error)
	go createTask(ctx, superChan)

	result := map[int]Ttype{}
	go func() {
		for r := range doneTasks {
			result[r.id] = *r
		}
	}()
	err := []error{}
	go func() {
		for r := range undoneTasks {
			err = append(err, r)
		}
	}()
	w := &sync.WaitGroup{}
	w.Add(1)
	go taskWorker(superChan, doneTasks, undoneTasks, w)
	w.Wait()
	close(undoneTasks)
	close(doneTasks)

	println("Errors:")
	for r := range err {
		println(r)
	}

	println("Done tasks:")
	for r := range result {
		fmt.Println(r)
	}
}

func taskWorker(ch <-chan *Ttype, doneTasks chan<- *Ttype, undoneTasks chan<- error, w *sync.WaitGroup) {
	defer w.Done()
	wg := sync.WaitGroup{}
	for t := range ch {
		wg.Add(1)
		go func(a *Ttype, doneTasks chan<- *Ttype, undoneTasks chan<- error) {
			defer wg.Done()
			if !IsTaskValid(a) {
				a.taskResult = []byte("something went wrong")
				undoneTasks <- fmt.Errorf("task id %d time %s, error %s", a.id, a.cT, a.taskResult)
				return
			}
			a.taskResult = []byte("task has been successed")
			a.fT = time.Now().Format(time.RFC3339Nano)
			time.Sleep(time.Millisecond * 1500)
			doneTasks <- a
		}(t, doneTasks, undoneTasks)
	}
	wg.Wait()
}

func createTask(ctx context.Context, a chan<- *Ttype) {
	//чтобы таски не дублировались с одинаковым ID будем их запускать через 1с
	a <- NewTtype()
	t := time.NewTicker(time.Second)
	for {
		select {
		case <-ctx.Done():
			{
				close(a)
				t.Stop()
				return
			}
		case <-t.C:
			{
				a <- NewTtype() // передаем таск на выполнение
			}
		}
	}
}

func IsTaskValid(a *Ttype) bool {
	t, err := time.Parse(time.RFC3339, a.cT)
	if err != nil || !t.After(time.Now().Add(-20*time.Second)) {
		return false
	}
	return true
}
