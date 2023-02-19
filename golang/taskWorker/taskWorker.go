package taskWorker

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
	"taskService/task"
	"time"
)

// Worker
// структура, реализующая обработку тасок
type Worker struct {
	outputChan chan *task.Task
	timeout    *time.Duration
}

func NewWorker(timeout *time.Duration) *Worker {
	return &Worker{
		timeout: timeout,
	}
}

func (w *Worker) SetOutput(outChan chan *task.Task) {
	w.outputChan = outChan
}

// Process
// главная функция-воркер для тасок
func (w *Worker) Process(currTask *task.Task, limiter chan bool, group *sync.WaitGroup) {
	defer func() {
		group.Done()
		logrus.Info("Finished task ", currTask.GetId())
		// вычитываем сообщение из канала, ограничивающего количество воркеров
		<-limiter
	}()
	// ставим промежуточный статус на случай, если логика обработки имеет шанс занять слишком много времени или зависнуть
	// в этом кейсе это бессмысленно, но пригодится для дебага после имплементации реальной логики
	currTask.Status = task.StatusInWork

	if currTask.CreationTime == nil {
		currTask.MakeFailed(fmt.Errorf("something went wrong"))
	} else if currTask.CreationTime.After(time.Now().Add(-20 * time.Second)) {
		currTask.MakeDone()
	} else {
		currTask.MakeFailed(fmt.Errorf("something went wrong"))
	}

	if w.timeout != nil {
		time.Sleep(*w.timeout)
	}

	currTask.SetFinishTime(time.Now())

	w.outputChan <- currTask
}
