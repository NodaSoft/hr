package dispatcher

import (
	"context"
	"fmt"
	"os"
	"task_service/internal/domain"
)

type Dispatcher struct {
	doneTasks    chan domain.Task
	errorsTask   chan domain.Task
	shutdownChan chan os.Signal
}

func NewDispatcher(taskCh, errorTask chan domain.Task, shutdownChan chan os.Signal) *Dispatcher {
	return &Dispatcher{
		doneTasks:    taskCh,
		errorsTask:   errorTask,
		shutdownChan: shutdownChan,
	}
}

func (d Dispatcher) Dispatch(ctx context.Context) {
	fmt.Println("start dispatcher")
	for {
		select {
		case <-ctx.Done():
			fmt.Println("done")
			close(d.shutdownChan)
			close(d.doneTasks)
			close(d.errorsTask)
			return
		case <-d.shutdownChan:
			fmt.Println("start shutdown gracefully...")
			close(d.shutdownChan)
			close(d.errorsTask)
			close(d.doneTasks)
			return
		}

	}

}
