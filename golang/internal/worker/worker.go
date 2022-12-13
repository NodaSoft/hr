package worker

import (
	"sync"
	"time"

	"github.com/fedoroko/nodasoft/test/internal/model"
)

type Chs struct {
	UnprocessedTasks <-chan model.Task
	ProcessedTasks   chan<- model.Task
	Quit             <-chan struct{}
}

type Worker struct {
	ch Chs
	wg *sync.WaitGroup
}

// ListenAndProcess ждет необработанные таски, обрабатывает их
// и отправляет записывающей горутине
func (worker *Worker) ListenAndProcess() {
	worker.wg.Add(1)
	defer worker.wg.Done()

	for {
		select {
		case <-worker.ch.Quit:
			return
		case task := <-worker.ch.UnprocessedTasks:
			processedTask := processTask(task)
			worker.ch.ProcessedTasks <- processedTask
		}
	}
}

// processTask обрабатывает таск.
// В оригинале время обработки проставляется раньше симуляции,
// вопринял это как баг, в противном случае 2 и 3 строки снизу нужно поменять местами.
func processTask(task model.Task) model.Task {
	if task.IsSuccessful && task.Timestamp.After(time.Now().Add(-20*time.Second)) {
		task.Result = []byte(model.TaskResultSuccess)
	} else {
		task.Result = []byte(model.TaskResultFail)
		task.IsSuccessful = false
	}

	time.Sleep(time.Millisecond * 150)
	task.ProcessedAt = time.Now().Format(time.RFC3339Nano)

	return task
}

func NewWorker(chs Chs, wg *sync.WaitGroup) Worker {
	return Worker{ch: chs, wg: wg}
}
