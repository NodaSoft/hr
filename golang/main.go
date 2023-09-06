package main

import (
	"fmt"
	"log"
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

type TaskExecutor struct {
	Name                string
	concurrrencyLimiter chan struct{}
	waitGroup           sync.WaitGroup
}

func NewTaskExecutor(name string, concurrrencyLimit int) *TaskExecutor {
	if name == "" {
		log.Panic("set name for simplify understanding of logs")
	}
	if concurrrencyLimit < 1 {
		log.Panic("concurrrencyLimit must be greater than 0")
	}
	return &TaskExecutor{
		Name:                name,
		concurrrencyLimiter: make(chan struct{}, concurrrencyLimit),
	}
}

// Execute executes a task in a worker pool.
// If there is no free worker, the function is suspended until a free worker is available.
func (te *TaskExecutor) Execute(id int, task func()) {
	if te.concurrrencyLimiter == nil {
		log.Panic("TaskExecutor was not initialised, use NewTaskExecutor")
	}
	te.concurrrencyLimiter <- struct{}{}
	te.waitGroup.Add(1)
	go func(id int, task func()) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("TaskExecutor %s: job #%d panic: %v", te.Name, id, err)
			}
			<-te.concurrrencyLimiter
			te.waitGroup.Done()
		}()
		task()
	}(id, task)
}

// Wait waits for all tasks to complete
func (te *TaskExecutor) Wait() {
	te.waitGroup.Wait()
}

// Len returns the number of used workers
func (te *TaskExecutor) Len() int {
	return len(te.concurrrencyLimiter)
}

// Cap returns the number of workers in the poll
func (te *TaskExecutor) Cap() int {
	return cap(te.concurrrencyLimiter)
}

// A Ttype represents a meaninglessness of our life
// Этот комментарий, имя типа и имена полей оставляю как есть.
// Для совместимости, как будто этот пакет - не main ;-)
// Но это всё (кроме id) не соответствует требованиям к коду.
type Ttype struct {
	id         int
	cT         string // время создания
	fT         string // время выполнения
	taskRESULT []byte
}

// Do - Executes a task and sends the result to one of the channels
// (doneTasks or errors)
func (t *Ttype) Do(doneTasks chan Ttype, errors chan error) {
	tt, _ := time.Parse(time.RFC3339, t.cT)
	if tt.After(time.Now().Add(-20 * time.Second)) {
		t.taskRESULT = []byte("task has been successed")
	} else {
		t.taskRESULT = []byte("something went wrong")
	}
	t.fT = time.Now().Format(time.RFC3339Nano)

	time.Sleep(time.Millisecond * 150)

	if string(t.taskRESULT[14:]) == "successed" {
		doneTasks <- *t
	} else {
		errors <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT, t.taskRESULT)
	}
}

// String implements the Stringer interface.
func (t *Ttype) String() string {
	return fmt.Sprintf("ID:%d created:%s finished:%s Result: %s", t.id, t.cT, t.fT, t.taskRESULT)
}

// taskGenerator generates tasks and sends them to the taskQueue channel
// until it receives a message from the stop channel
func taskGenerator(taskQueue chan Ttype, stop <-chan time.Time) {
	defer close(taskQueue)
	for {
		ft := time.Now().Format(time.RFC3339)
		if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
			ft = "Some error occured"
		}
		select {
		case taskQueue <- Ttype{cT: ft, id: int(time.Now().Unix())}:
			// передаем таск в очередь
		case <-stop:
			return // завершаем генерацию задач
		}
	}
}

func main() {

	const TaskWorkersPoolSize = 5
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.LUTC | log.Lshortfile)

	var waitForResults sync.WaitGroup
	waitForResults.Add(2)

	// результаты работы оставляю без изменений (для совместимости)

	result := map[int]Ttype{}
	doneTasks := make(chan Ttype, TaskWorkersPoolSize)
	go func() {
		for task := range doneTasks {
			result[task.id] = task
		}
		waitForResults.Done()
	}()

	err := []error{}
	undoneTasks := make(chan error, TaskWorkersPoolSize)
	go func() {
		for fail := range undoneTasks {
			err = append(err, fail)
		}
		waitForResults.Done()
	}()

	taskQueue := make(chan Ttype, TaskWorkersPoolSize)
	taskExecutor := NewTaskExecutor("main", TaskWorkersPoolSize)
	var waitGroupForTaskStarter sync.WaitGroup
	waitGroupForTaskStarter.Add(1)
	go func() {
		defer waitGroupForTaskStarter.Done()
		// получение тасков
		for t := range taskQueue {
			t := t
			taskExecutor.Execute(t.id, func() {
				t.Do(doneTasks, undoneTasks)
			})
		}
	}()

	go taskGenerator(taskQueue, time.After(time.Second*3))

	waitGroupForTaskStarter.Wait()
	taskExecutor.Wait()
	close(doneTasks)
	close(undoneTasks)
	waitForResults.Wait()

	// вывод результатов работы изменил, т.к. исходный был практически бесполезным

	println("Errors:")
	for _, r := range err {
		fmt.Printf("%s\n", r)
	}

	println("Done tasks:")
	for _, r := range result {
		fmt.Printf("%s\n", r)
	}
}
