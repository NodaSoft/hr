package workers

import (
	"fmt"
	"task_service/internal/domain"
	"time"
)

type SimpleWorker struct {
}

var ResultSuccess = domain.WokrkResult("task has been successed")
var ResultError = domain.WokrkResult("something wrong")

func NewSimpleWorker() SimpleWorker {
	return SimpleWorker{}
}

func (worker SimpleWorker) Handle(task domain.Task) (domain.Task, error) {
	var err error
	tt, _ := time.Parse(time.RFC3339, task.CreationTime)
	if tt.After(time.Now().Add(-20 * time.Second)) {
		task.Result = ResultSuccess
	} else {
		task.Result = ResultError
		err = fmt.Errorf("task id %s time %s, error %s", task.ID, task.CreationTime, task.Result)
	}

	task.ExecutionTime = time.Now().Format(time.RFC3339Nano)

	time.Sleep(time.Millisecond * 150)

	return task, err
}
