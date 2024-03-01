package main

import (
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
	taskRESULT []byte
}

// Tasker represents a tasker. It can generate tasks and process them.
type Tasker struct {
	tasks        chan Ttype
	doneTasks    chan Ttype
	undoneTasks  chan error
	generateTask bool
	wg           sync.WaitGroup
}

// NewTasker creates new tasker object.
func NewTasker() *Tasker {
	return &Tasker{
		tasks:       make(chan Ttype, 10),
		doneTasks:   make(chan Ttype),
		undoneTasks: make(chan error),
	}
}

// Run tasker.
func (t *Tasker) Run(taskCreator func() Ttype,
	taskWorker func(a Ttype) Ttype) (results map[int]Ttype, errs *[]error) {

	results = map[int]Ttype{}
	errs = &[]error{}

	// Result processing
	t.wg.Add(1)
	go func() {
		defer t.wg.Done()
		for r := range t.doneTasks {
			results[r.id] = r
		}
	}()

	// Error processing
	t.wg.Add(1)
	go func() {
		defer t.wg.Done()
		for r := range t.undoneTasks {
			*errs = append(*errs, r)
		}
	}()

	// Task processing
	go func() {
		wg := sync.WaitGroup{}
		for ta := range t.tasks {
			wg.Add(1)
			go func(ta Ttype) {
				defer wg.Done()
				ta = taskWorker(ta)
				if string(ta.taskRESULT[14:]) == "successed" {
					t.doneTasks <- ta
					return
				}
				t.undoneTasks <- fmt.Errorf("task id: %d, error: %s", ta.id,
					ta.taskRESULT)
			}(ta)
		}
		wg.Wait()
		close(t.doneTasks)
		close(t.undoneTasks)
	}()

	// Start task creating
	go func() {
		t.generateTask = true
		for t.generateTask {
			t.tasks <- taskCreator()
		}
		close(t.tasks)
	}()

	return
}

// Stop tasker.
func (t *Tasker) Stop() {
	t.generateTask = false
	t.wg.Wait()
}

func main() {

	// Create tasker object
	t := NewTasker()

	// Start task creation and processing
	results, errs := t.Run(

		// Task creator
		func() Ttype {
			time.Sleep(time.Millisecond * 150)

			ft := time.Now().Format(time.RFC3339)
			if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
				ft = "Some error occured"
			}

			return Ttype{cT: ft, id: int(time.Now().Unix())} // передаем таск на выполнение
		},

		// Task worker
		func(a Ttype) Ttype {
			tt, _ := time.Parse(time.RFC3339, a.cT)
			if tt.After(time.Now().Add(-20 * time.Second)) {
				a.taskRESULT = []byte("task has been successed")
			} else {
				a.taskRESULT = []byte("something went wrong")
			}
			a.fT = time.Now().Format(time.RFC3339Nano)
			time.Sleep(time.Millisecond * 150)

			return a
		},
	)

	// Sleep for 3 seconds and stop tasker
	time.Sleep(time.Second * 3)
	t.Stop()

	// Print error results
	println("Errors:")
	for _, e := range *errs {
		fmt.Println(e)
	}

	// Print success results
	println("Done tasks:")
	for _, r := range results {
		fmt.Println("task id:", r.id)
	}
}
