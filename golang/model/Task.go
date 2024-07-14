package model

import (
	"fmt"
	"time"
)

const (
	TaskErrorMessage    = "Task processing failed"
	TaskSuccessMessage  = "Task has been successful"
	TaskProcessingDelay = 150 * time.Millisecond
)

// Task represents a meaninglessness of our life
type Task struct {
	Id         int
	CreateTime time.Time // время создания
	FinalTime  time.Time // время выполнения
	Result     []byte
	Error      error
}

func (t *Task) Check() {
	if t.CreateTime.After(time.Now().Add(-20 * time.Second)) {
		t.Result = []byte(TaskSuccessMessage)
	} else {
		t.Result = []byte(TaskErrorMessage)
	}
	tt, err := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	if err != nil {
		t.Error = err
	}
	t.FinalTime = tt

	time.Sleep(TaskProcessingDelay)
}

func (t *Task) Sort(successTasks chan<- Task, failTasks chan error) {
	if t.Error != nil {
		failTasks <- fmt.Errorf("Task Id %d time %s, error %s", t.Id, t.CreateTime, t.Result)
	} else {
		successTasks <- *t
	}
}
