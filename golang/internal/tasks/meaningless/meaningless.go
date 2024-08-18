package meaninglessTask

import (
	"fmt"
	"test_task/pkg/workerPool"
	"time"
)

const (
	TASK_SUCCESS_MESSAGE = "task has been successed"
	TASK_ERROR_MESSAGE   = "something went wrong"
)

var _ workerPool.Task = (*MeaninglessTask)(nil)

// A MeaninglessTask represents a meaninglessness of our life
type MeaninglessTask struct {
	ID           int
	CreationTime string
	FinishTime   string
	result       string
}

func New() *MeaninglessTask {
	currentTime := time.Now().Format(time.RFC3339)
	if time.Now().Nanosecond()%2 > 0 {
		currentTime = "Some error occured"
	}

	return &MeaninglessTask{
		ID:           int(time.Now().Unix()),
		CreationTime: currentTime,
	}
}

func (mt *MeaninglessTask) Process() {
	tt, _ := time.Parse(time.RFC3339, mt.CreationTime)
	if tt.After(time.Now().Add(-20 * time.Second)) {
		mt.result = TASK_SUCCESS_MESSAGE
	} else {
		mt.result = TASK_ERROR_MESSAGE
	}
	mt.FinishTime = time.Now().Format(time.RFC3339)
}

func (mt *MeaninglessTask) Result() string {
	return mt.result
}

func (mt *MeaninglessTask) IsSuccess() bool {
	return mt.result == TASK_SUCCESS_MESSAGE
}

func (mt *MeaninglessTask) Error() error {
	return fmt.Errorf("ID: %d, Created at: %s, Error: %s", mt.ID, mt.CreationTime, mt.result)
}

func (mt *MeaninglessTask) Status() string {
	return fmt.Sprintf("ID: %d, Created at: %s, Finished at: %s", mt.ID, mt.CreationTime, mt.FinishTime)
}
