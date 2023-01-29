package main

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
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

	workerNumber := 5
	g := errgroup.Group{}
	for i := 0; i < workerNumber; i++ {
		w := task_worker{
			ch:          superChan,
			doneTasks:   doneTasks,
			undoneTasks: undoneTasks,
		}
		g.Go(w.do)
	}
	g.Wait()
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

func createTask(ctx context.Context, a chan<- *Ttype) {
	//чтобы таски не дублировались с одинаковым ID будем их запускать через 1с
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
				t := time.Now()
				ct := t.Format(time.RFC3339)
				if t.Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков (странное условие но пусть будет)
					ct = "Some error occured"
				}
				a <- &Ttype{cT: ct, id: int(t.Unix())} // передаем таск на выполнение
			}
		}
	}

}

type task_worker struct {
	ch          <-chan *Ttype
	doneTasks   chan<- *Ttype
	undoneTasks chan<- error
}

func (t *task_worker) do() error {
	for a := range t.ch {
		isErr := false
		tt, err := time.Parse(time.RFC3339, a.cT)
		if err != nil || !tt.After(time.Now().Add(-20*time.Second)) {
			a.taskResult = []byte("something went wrong")
			isErr = true
			//return err
		} else {
			a.taskResult = []byte("task has been successed")
		}

		a.fT = time.Now().Format(time.RFC3339Nano)
		time.Sleep(time.Millisecond * 1500)
		if isErr {
			t.undoneTasks <- fmt.Errorf("task id %d time %s, error %s", a.id, a.cT, a.taskResult)
		} else {
			t.doneTasks <- a
		}
	}
	return nil
}
