package taskEmitter

import (
	"github.com/sirupsen/logrus"
	"taskService/task"
	"time"
)

// Emitter
// структура, позволяющая генерировать таски
type Emitter struct {
	TaskChan chan *task.Task
	timeout  *time.Duration
	quitChan chan bool
}

func NewEmitter(taskChan chan *task.Task, timeout *time.Duration) *Emitter {
	return &Emitter{
		TaskChan: taskChan,
		timeout:  timeout,
		quitChan: make(chan bool),
	}
}

// Quit
// Посылаем сигнал на завершение эмиттеру тасок. Блокирующая операция
func (e *Emitter) Quit() {
	logrus.Warn("Emitter got shutdown signal")
	e.quitChan <- true
	logrus.Warn("Emitter was successfully shut down")
}

// EmitTasks
// Запускает рутину, генерирующую таски. Неблокирующая
func (e *Emitter) EmitTasks() {
	go e.emitTasksInternal()
}

func (e *Emitter) emitTasksInternal() {
	logrus.Info("Tasks are being emitted")
	for {
		select {
		case <-e.quitChan:
			return
		default:
		}

		currTask := task.NewTask()
		currTime := time.Now()
		// определение четности через битовый оператор (просто потому что могу и теоретически быстрее, хотя скорее всего
		// компилятор нормально оптимизирует точно так же и операцию через остаток от деления)
		if currTime.Nanosecond()&1 > 0 { // вот такое условие появления ошибочных тасков
			pastTime := currTime.Add(time.Second * -30)
			currTask.CreationTime = &pastTime
		} else {
			currTask.SetCreationTime(currTime)
		}

		if e.timeout != nil {
			time.Sleep(*e.timeout)
		}

		e.TaskChan <- currTask
	}
}
