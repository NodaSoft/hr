package main

import (
	"fmt"
	"time"
	"sync"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнены остальных тасков

// A Task represents a meaninglessness of our life
type Task struct {
	id         int
	cT         string // время создания
	fT         string // время выполнения
	result	   error
	taskRESULT []byte
}

type TaskFabric struct {
	cur_id int
}

func NewTaskFabric() TaskFabric {
	return TaskFabric{
		cur_id: 100,
	}
}

func (f *TaskFabric) CreateTask() Task {
	ft := time.Now().Format(time.RFC3339)
	if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков 
		ft = "Some error occured"
	}

	task := Task{
		id: f.cur_id,
		cT: ft,
	}
	f.cur_id++

	return task
}

type TaskWorker struct {
	goroutineChan chan interface{}
	taskSorter func(Task)
}

func NewTaskWorker(taskSorter func(Task)) *TaskWorker {
	const MAX_TASKS = 3

	return &TaskWorker{
		goroutineChan: make(chan interface{}, MAX_TASKS),
		taskSorter: taskSorter,
	}
}

func (w *TaskWorker) workTaskImpl(task Task) Task {
	tt, _ := time.Parse(time.RFC3339, task.cT)
	if tt.After(time.Now().Add(-20 * time.Second)) {
		task.result = nil
	} else {
		task.result = fmt.Errorf("something went wrong")
	}
	task.fT = time.Now().Format(time.RFC3339Nano)

	time.Sleep(time.Millisecond * 150)

	return task
}

func (w *TaskWorker) WorkTask(task Task) {
	w.goroutineChan <- struct{}{}

	go func() {
		task = w.workTaskImpl(task)

		w.taskSorter(task)

		<-w.goroutineChan
	}()
}

func SortTask(task Task, doneTasks chan Task, undoneTasks chan Task) {
	if task.result == nil {
		doneTasks <- task
	} else {
		undoneTasks <- task
	}
}

type TaskLogger struct {
	result map[int]Task
	err map[int]Task
	mutex sync.Mutex
}

func NewTaskLogger(doneTasks chan Task, undoneTasks chan Task) *TaskLogger {
	logger := &TaskLogger{
		result:	make(map[int]Task),
		err:	make(map[int]Task),
		mutex:	sync.Mutex{},
	}

	go func() {
		for {
			select {
			case task := <-doneTasks:
				logger.mutex.Lock()
				logger.result[task.id] = task
				logger.mutex.Unlock()
			case task := <-undoneTasks:
				logger.mutex.Lock()
				logger.err[task.id] = task
				logger.mutex.Unlock()
			}
		}
	}()

	return logger
}

func (t *TaskLogger) PrintTasks() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	println("Done tasks:")
	for r := range t.result {
		fmt.Printf("%v\n", r)
	}

	println("Failed tasks:")
	for _, v := range t.err {
		fmt.Printf("Task id %d time %s, error %s\n", v.id, v.cT, v.result)
	}
}

func main() {
	taskCreator := func(taskChan chan Task) {
		taskFabric := NewTaskFabric()

		for {
			task := taskFabric.CreateTask()
			taskChan <- task
		}
	}

	taskChan := make(chan Task, 10)
	go taskCreator(taskChan)

	doneTasks := make(chan Task)
	undoneTasks := make(chan Task)

	go func() {
		taskSorter := func(task Task) {
			SortTask(task, doneTasks, undoneTasks)
		}

		worker := NewTaskWorker(taskSorter)
		for {
			select {
			case task := <-taskChan:
				worker.WorkTask(task)
			}
		}
	}()

	logger := NewTaskLogger(doneTasks, undoneTasks)

	time.Sleep(time.Second * 3)
	logger.PrintTasks();
}
