package steps

import (
	"context"
	st "main/internal/task/structs"
	"time"
)

func RunProcessing(ctx context.Context, taskCh <-chan st.Task, doneCh chan<- st.Task, taskRelevanceTime time.Duration) {
	for {
		select {
		case <-ctx.Done():
			close(doneCh)
			return
		case task, ok := <-taskCh:
			if !ok {
				return
			}
			doneCh <- processTask(task, taskRelevanceTime)
		}
	}
}

func processTask(task st.Task, taskRelevanceTime time.Duration) st.Task {
	createTime, _ := time.Parse(time.RFC3339, task.CreateTime)
	if createTime.After(time.Now().Add(-taskRelevanceTime)) {
		task.Result = []byte("task has been successfully completed")
	} else {
		task.Result = []byte("something went wrong")
	}
	task.FinishTime = time.Now().Format(time.RFC3339Nano)

	time.Sleep(time.Millisecond * 150)

	return task
}
