package taskRouter

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"sync"
	"taskService/task"
)

// Router
// структура, управляющая потоком тасок на разных этапах их жизни - получение, обработка, сортировка результатов
type Router struct {
	NewTasksChan      chan *task.Task
	FinishedTasksChan chan *task.Task
	FinishedTasks     map[uuid.UUID]*task.Task
	FailedTasks       map[uuid.UUID]*task.Task
	Worker            task.WorkerInterface
	MaxQueueSize      int
	queueLimiter      chan bool
	receiverQuitChan  chan bool
	sorterQuitChan    chan bool
	workersWG         *sync.WaitGroup
	sorterWG          *sync.WaitGroup
	receiverWG        *sync.WaitGroup
}

func NewRouter(work task.WorkerInterface, maxQueue int) *Router {
	finishedChan := make(chan *task.Task)
	work.SetOutput(finishedChan)
	return &Router{
		NewTasksChan:      make(chan *task.Task),
		FinishedTasksChan: finishedChan,
		FinishedTasks:     map[uuid.UUID]*task.Task{},
		FailedTasks:       map[uuid.UUID]*task.Task{},
		Worker:            work,
		MaxQueueSize:      maxQueue,
		queueLimiter:      make(chan bool, maxQueue-1),
		receiverQuitChan:  make(chan bool, 1),
		sorterQuitChan:    make(chan bool, 1),
		workersWG:         &sync.WaitGroup{},
		sorterWG:          &sync.WaitGroup{},
		receiverWG:        &sync.WaitGroup{},
	}
}

func (r *Router) GetInputChannel() chan *task.Task {
	return r.NewTasksChan
}

// Run
// Запускает рутину, проводящую таски по их жизненному циклу (получает, отправляет воркеру, сортирует после
// воркера). Неблокирующая
// TODO: добавить механизм на случай падения\зависания рисвера\сортера
func (r *Router) Run() {
	r.receiverWG.Add(1)
	go r.receiver()
	r.sorterWG.Add(1)
	go r.sorter()
}

// получаем таски, передаем воркеру. блокируется, если превышено количество сообщений в канале queueLimiter
// воркер по завершению вычитывает одно значение из канала, это позволяет не допустить переполнения при неограниченном
// потоке тасок
func (r *Router) receiver() {
	defer func() {
		r.receiverWG.Done()
	}()
	logrus.Info("router.receiver online")
	for currTask := range r.NewTasksChan {
		r.workersWG.Add(1)
		go r.Worker.Process(currTask, r.queueLimiter, r.workersWG)
		r.queueLimiter <- true // останавливаем по достижению лимита воркеров
	}
	logrus.Warn("Receiver was successfully shut down")
}

// сортируем обработанные таски по статусу (успешные\проваленные) и складываем в соответствующие карты
// к сортировщику идёт один канал вместо двух, потому что нафига
func (r *Router) sorter() {
	defer func() {
		r.sorterWG.Done()
	}()
	logrus.Info("router.sorter online")
	for currTask := range r.FinishedTasksChan {
		if currTask.Status == task.StatusDone {
			r.FinishedTasks[currTask.GetUUID()] = currTask
		} else {
			r.FailedTasks[currTask.GetUUID()] = currTask
		}
	}
	logrus.Warn("Sorter was successfully shut down")
}

// Quit
// ждём завершения всех воркеров и ресивера\сортера, завершаем работу роутера
func (r *Router) Quit() {
	logrus.Warn("Got router shutdown signal")

	close(r.NewTasksChan) // завершаем работу ресивера
	r.receiverWG.Wait()

	r.workersWG.Wait() // ждем окончания всех воркеров
	logrus.Warn("All jobs are finished")

	close(r.FinishedTasksChan) // завершаем работу сортера
	r.sorterWG.Wait()
	logrus.Warn("Router was successfully shut down")
}

func (r *Router) PrintErrors() {
	fmt.Printf("Errors (%d):\n", len(r.FailedTasks))
	for _, r := range r.FailedTasks {
		fmt.Println(r)
	}
}
func (r *Router) PrintFinished() {
	fmt.Printf("Done tasks (%d):\n", len(r.FinishedTasks))
	for _, r := range r.FinishedTasks {
		fmt.Println(r.String())
	}
}
