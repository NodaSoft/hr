package main

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

var (
	waitingTime   = 1
	taskErrorCode = "ERROR"
)

func main() {
	var nTasks int = 10
	tm := NewTaskManager(nTasks)
	tm.GenerateTasks()
	tm.Start()
	tm.Print()
	tm.Observer()
}

type Task struct {
	id           int
	createdTime  time.Time
	finishedTime time.Time
	err          error
}

func (t *Task) Id() int {
	return t.id
}

func (t *Task) CreatedTime() time.Time {
	return t.createdTime
}

func (t *Task) FinishedTime() time.Time {
	return t.finishedTime
}

func (t *Task) Err() error {
	return t.err
}

func (t *Task) setFinishedTime(v time.Time) {
	t.finishedTime = v
}

func (t *Task) setErr(v error) {
	t.err = v
}

// Processing our task and emulate some work
func (t *Task) Processing() {
	time.Sleep(time.Second * time.Duration(waitingTime))

	t.setFinishedTime(time.Now())
	var isError bool = t.finishedTime.UnixMicro()%2 > 0

	if isError {
		t.setErr(fmt.Errorf(taskErrorCode))
	}

	log.Printf("%d Task processed, finished at: %v", t.id, t.finishedTime)
}

func (t *Task) Print() {
	if t.Err() == nil {
		log.Printf("%d Task, created at: %v, finished at: %v", t.Id(), t.CreatedTime(), t.FinishedTime())
	} else {
		log.Printf("%d Task, created at: %v, finished with error: %v", t.Id(), t.CreatedTime(), t.Err())
	}
}

func NewTask() Task {
	return Task{
		id:          int(time.Now().UnixMicro()),
		createdTime: time.Now(),
	}
}

type TaskManager struct {
	queueTasks   chan Task
	queueResults chan Task
	wgTasks      sync.WaitGroup
	wgResults    sync.WaitGroup
	nTask        int
	generated    int32
	proccessed   int32
}

func NewTaskManager(nTask int) TaskManager {
	return TaskManager{
		queueTasks:   make(chan Task, nTask),
		queueResults: make(chan Task),
		nTask:        nTask,
	}
}

func (m *TaskManager) QueueTasks() chan Task {
	return m.queueTasks
}

func (m *TaskManager) QueueResults() chan Task {
	return m.queueResults
}

func (m *TaskManager) NumberOfTasks() int {
	return m.nTask
}

func (m *TaskManager) addTaskInQueue(v Task) {
	m.queueTasks <- v
}

func (m *TaskManager) addResultInQueue(v Task) {
	m.queueResults <- v
}

// Generated N tasks
func (m *TaskManager) GenerateTasks() {
	go func() {
		for i := 0; i < m.NumberOfTasks(); i++ {
			task := NewTask()
			m.addTaskInQueue(task)
			atomic.AddInt32(&m.generated, 1)
			log.Printf("%d Task generated at: %v", task.Id(), task.CreatedTime())
		}
		close(m.QueueTasks())
	}()
}

// Start workers for N tasks
func (m *TaskManager) Start() {
	for i := 0; i < m.NumberOfTasks(); i++ {
		m.wgTasks.Add(1)
		go m.startWorker()
	}
}

func (m *TaskManager) startWorker() {
	defer m.wgTasks.Done()
	for task := range m.QueueTasks() {
		task.Processing()
		atomic.AddInt32(&m.proccessed, 1)
		m.addResultInQueue(task)
	}
}

// Printing successes and errors in our tasks
func (m *TaskManager) Print() {
	m.wgResults.Add(1)
	go func() {
		defer m.wgResults.Done()
		success := make([]Task, 0)
		errors := make([]Task, 0)

		for result := range m.QueueResults() {
			if result.Err() == nil {
				success = append(success, result)
			} else {
				errors = append(errors, result)
			}
		}
		log.Println("Completed tasks")
		for _, r := range success {
			r.Print()
		}
		log.Println("Count of completed tasks: ", len(success))
		log.Println("Finished with error tasks")
		for _, r := range errors {
			r.Print()
		}
		log.Println("Count of finished with error tasks: ", len(errors))
	}()
}

// Observer will wait until workers finish and we will print results
func (m *TaskManager) Observer() {
	m.wgTasks.Wait()
	close(m.QueueResults())
	m.wgResults.Wait()
	log.Printf("Total generated tasks: %d \n", m.generated)
	log.Printf("Total processed tasks: %d \n", m.proccessed)
}
