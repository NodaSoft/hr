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

type task struct {
	taskID     int
	createdAt  time.Time
	executedAt time.Time
	taskResult []byte
}

func NewTask(id int, createdAt time.Time, executedAt time.Time, data []byte) *task {
	return &task{
		taskID:     id,
		createdAt:  createdAt,
		executedAt: executedAt,
		taskResult: data,
	}
}

type Runner interface {
	Run()
}

type taskFactory struct {
	size        int
	tasksToRun  chan *task
	doneTasks   chan *task
	failedTasks chan error

	timeout time.Duration
}

func NewTaskFactory(size int, wg *sync.WaitGroup, timeout time.Duration) Runner {
	factory := taskFactory{
		size:        size,
		tasksToRun:  make(chan *task, size),
		doneTasks:   make(chan *task),
		failedTasks: make(chan error),
		timeout:     timeout,
	}

	for i := 0; i < size; i++ {
		wg.Add(1)
		go factory.worker(wg)
	}

	go factory.handle(wg)

	return &factory
}

func (t *taskFactory) Run() {
	defer close(t.tasksToRun)
	for i := 0; i < t.size; i++ {
		t.tasksToRun <- NewTask(int(time.Now().Unix()), time.Now(), time.Now(), nil)
	}
}

func (t *taskFactory) checkTask(inTask *task) {
	if inTask.createdAt.After(time.Now().Add(-20 * time.Second)) {
		inTask.taskResult = []byte("task has been successed")
	} else {
		inTask.taskResult = []byte("something went wrong")
	}
	time.Sleep(time.Millisecond * 150)
}

func (t *taskFactory) dispatchTask(currentTask *task) {
	if string(currentTask.taskResult[14:]) == "successed" {
		t.doneTasks <- currentTask
	} else {
		t.failedTasks <- fmt.Errorf("task_id: %d, executedAt: %s, err: %s", currentTask.taskID,
			currentTask.executedAt, string(currentTask.taskResult))
	}
}

func (t *taskFactory) worker(wg *sync.WaitGroup) {
	defer wg.Done()
	select {
	case currentTask, ok := <-t.tasksToRun:
		if ok {
			t.checkTask(currentTask)
			go t.dispatchTask(currentTask)
		}
	}
}

func (t *taskFactory) handle(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	for {
		select {
		case done, ok := <-t.doneTasks:
			if ok {
				fmt.Printf("DONE: task_id: %d, executed at %s, result: %s\n", done.taskID,
					done.executedAt.Format(time.RFC3339Nano), string(done.taskResult))
			}
		case err, ok := <-t.failedTasks:
			if ok {
				fmt.Printf("FAILED: %s\n", err)
			}
		case <-time.After(time.Second * t.timeout):
			return
		}
	}
}

func main() {
	wg := sync.WaitGroup{}
	factory := NewTaskFactory(10, &wg, 1)
	go factory.Run()
	wg.Wait()
}
