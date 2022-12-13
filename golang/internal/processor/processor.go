package processor

import (
	"fmt"
	"sync"
	"time"

	"github.com/fedoroko/nodasoft/test/internal/model"
	taskWorker "github.com/fedoroko/nodasoft/test/internal/worker"
)

type ITaskProcessor interface {
	Run()
	CloseAndPrint()
}

// taskProcessor создает таски, отправляет их воркерам,
// получает их обратно и записывает результаты.
// fan-in fan-out
type taskProcessor struct {
	workers          []taskWorker.Worker
	workersWG        *sync.WaitGroup
	unprocessedTasks chan model.Task
	processedTasks   chan model.Task
	succeedTasks     map[int]model.Task
	failedTasks      []error
	createQuit       chan struct{} // канал выхода горутины ответственной за создание тасков
	workersQuit      chan struct{} // канал выхода воркеров
	storeQuit        chan struct{} // канал выхода горутины ответственной за запись
	gracefulShutdown chan struct{} // канал подтверждающий, что все горутины завершены
}

// Run запускает воркеров, горутину,
// которая будет записывать результаты,
// и горутину, которая будет создавать таски
func (taskProcessor *taskProcessor) Run() {
	for _, worker := range taskProcessor.workers {
		go worker.ListenAndProcess()
	}

	go taskProcessor.listenAndStore()
	go taskProcessor.create()
}

// create создает таски, отправлет их в канал к воркерам
func (taskProcessor *taskProcessor) create() {
	for {
		select {
		case <-taskProcessor.createQuit:
			return
		default:
			task := createTask()
			taskProcessor.unprocessedTasks <- task
		}
	}
}

// listenAndStore получает обработанные таски от воркеров,
// сортирует их и записывает результаты
func (taskProcessor *taskProcessor) listenAndStore() {
	for {
		select {
		case <-taskProcessor.storeQuit:
			close(taskProcessor.gracefulShutdown)
			return
		case task := <-taskProcessor.processedTasks:
			taskProcessor.storeTask(task)
		}
	}
}

func createTask() model.Task {
	timestamp := time.Now()
	task := model.Task{
		ID:           int(timestamp.Unix()),
		CreatedAt:    timestamp.Format(time.RFC3339),
		Timestamp:    timestamp,
		IsSuccessful: true,
	}

	if timestamp.Nanosecond()%2 > 0 {
		task.CreatedAt = model.TaskCreationError
		task.IsSuccessful = false
	}

	return task
}

// storeTask при реализации fan-in fan-out итоговая запись происходит синхронно,
// не боимся конкуренции, мутекс не нужен
func (taskProcessor *taskProcessor) storeTask(task model.Task) {
	if task.IsSuccessful {
		taskProcessor.succeedTasks[task.ID] = task
		return
	}

	err := fmt.Errorf(
		"task id %d time %s, error %s",
		task.ID, task.CreatedAt, task.Result)
	taskProcessor.failedTasks = append(taskProcessor.failedTasks, err)
}

func (taskProcessor *taskProcessor) CloseAndPrint() {
	taskProcessor.close()

	fmt.Println("Errors:")
	for _, err := range taskProcessor.failedTasks {
		fmt.Println(err)
	}

	fmt.Println("Done tasks:")
	for taskID := range taskProcessor.succeedTasks {
		fmt.Println(taskID)
	}
}

// close останавливает горутину создающую таски,
// посылает сигнал о завершении воркерам и ждет их остановки,
// посылает сигнал о завершении записывающей горутине,
// и дожидается ее сигнала об успешном завершении.
func (taskProcessor *taskProcessor) close() {
	close(taskProcessor.createQuit)
	close(taskProcessor.workersQuit)
	taskProcessor.workersWG.Wait()
	close(taskProcessor.storeQuit)
	<-taskProcessor.gracefulShutdown
}

func NewTaskProcessor(numOfWorkers int) ITaskProcessor {
	unprocessedTasks := make(chan model.Task, numOfWorkers)
	processedTasks := make(chan model.Task, numOfWorkers)
	workerQuit := make(chan struct{})

	workers := make([]taskWorker.Worker, numOfWorkers)
	workersWG := &sync.WaitGroup{}
	for i := 0; i < numOfWorkers; i++ {
		workers[i] = taskWorker.NewWorker(
			taskWorker.Chs{
				UnprocessedTasks: unprocessedTasks,
				ProcessedTasks:   processedTasks,
				Quit:             workerQuit,
			},
			workersWG,
		)
	}
	return &taskProcessor{
		workers:          workers,
		workersWG:        workersWG,
		unprocessedTasks: unprocessedTasks,
		processedTasks:   processedTasks,
		succeedTasks:     make(map[int]model.Task),
		failedTasks:      make([]error, 0),
		createQuit:       make(chan struct{}),
		workersQuit:      workerQuit,
		storeQuit:        make(chan struct{}),
		gracefulShutdown: make(chan struct{}),
	}
}
