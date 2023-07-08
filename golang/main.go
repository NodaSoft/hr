package main

import (
	"errors"
	"fmt"
	"math/rand"
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

// todo: 
// - split package
// - errors
// - done chan/context

// A Ttype represents a meaninglessness of our life
type Task struct {
	Id         int
	CreateTime time.Time // время создания
	FinishTime time.Time // время выполнения
	Result []byte
	Err error
}

func (t *Task) Work() {	
	if t.Err == nil && t.CreateTime.After(time.Now().Add(-20 * time.Second)) {
		t.Result = []byte("task has been successed")
	} else {
		t.Result = []byte("something went wrong")

		err := errors.New("something went wrong")
		if t.Err != nil {
			t.Err = fmt.Errorf("%w; %w", err, t.Err)
		} else {
			t.Err = err
		}
	}
		
	t.FinishTime = time.Now()

	time.Sleep(time.Millisecond * 150)
}

func (t *Task) String() string {
	if t.Err == nil {
		return fmt.Sprintf("task id: %d, time: %s, result: %s",
			t.Id,
			t.CreateTime.Format(time.RFC3339),
			t.Result,
		)
	} else {
		return fmt.Sprintf("task id: %d, error: %s",
			t.Id,
			t.Err,
		)
	}
}

type TaskGenerator struct {
	Duration time.Duration
	BufferSize int
}

func (tg *TaskGenerator) Start() (doneTasks, undoneTasks <-chan *Task) {
	tasks := tg.create()
	worked := tg.process(tasks)
	return tg.sort(worked)
}

func (tg *TaskGenerator) create() <-chan *Task {
	ch := make(chan *Task, tg.BufferSize)

	// rand.Seed(time.Now().UnixNano())
	
	go func() {		
		defer close(ch)
		
		after := time.After(tg.Duration)
		
		for {
			var err error
			createTime := time.Now()
			rnd := rand.Int()

			// if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
			if rnd%2 > 0{
				err = errors.New("odd task")
			}

			select {
			case <-after:	
				// close(ch)			
				return
			case ch <- &Task{ // передаем таск на выполнение
				// Id:		int(createTime.Unix()),
				Id:         rnd,
				CreateTime: createTime,
				Err:        err,
			}:
			}
		}
	}()

	return ch
}

func (tg *TaskGenerator) process(tasks <-chan *Task) <-chan *Task {
	ch := make(chan *Task)

	go func() {
		defer close(ch)

		wg := &sync.WaitGroup{}
		// получение тасков
		for t := range tasks {
			wg.Add(1)
			go func(t *Task) {
				t.Work()
				ch <- t
				wg.Done()
			}(t)
		}
		wg.Wait()
	}()

	return ch
}

func (tg *TaskGenerator) sort(tasks <-chan *Task) (doneTasks, undoneTasks chan *Task) {
	doneTasks = make(chan *Task)
	undoneTasks = make(chan *Task)

	go func() {
		defer close(doneTasks)
		defer close(undoneTasks)

		for t := range tasks {
			if t.Err == nil {
				doneTasks <- t
			} else {
				undoneTasks <- t
			}
		}
	}()

	return
}

type TaskConsumer struct {
	taskGenerator *TaskGenerator
}

func (tu *TaskConsumer) Run() {
	doneTasks, undoneTasks := tu.taskGenerator.Start()

	type MxMap struct {
		sync.RWMutex
		data map[int]*Task
	}
	results := &MxMap{data: make(map[int]*Task)}
	errors := &MxMap{data: make(map[int]*Task)}
	
	wg := &sync.WaitGroup{}	
	wg.Add(2)
	go func() {
		defer wg.Done()

		for dt := range doneTasks {
			dt := dt
			
			wg.Add(1)
			go func() { // обертка go func в данном случае не нужна (показать работу sync.RWMutex)
				defer wg.Done()

				results.Lock()				
				results.data[dt.Id] = dt
				results.Unlock()
			}()
		}
	}()
	go func() {
		defer wg.Done()

		for ut := range undoneTasks {
			errors.data[ut.Id] = ut
		}
	}()
	wg.Wait()

	// печатать можно из каналов (показать работу sync.WaitGroup)
	fmt.Printf("Errors (%d):\n", len(errors.data))
	for _, e := range errors.data {
		fmt.Println(e)
	}

	fmt.Printf("Done tasks (%d):\n", len(results.data))
	for _, r := range results.data {
		fmt.Println(r)
	}
}

func main() {
	taskConsumer := &TaskConsumer{
		taskGenerator: &TaskGenerator{
			Duration: time.Second * 3,
			BufferSize: 10,
		},
	}

	taskConsumer.Run()
}
