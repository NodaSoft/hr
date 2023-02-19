package task

import (
	"fmt"
	"github.com/google/uuid"
	"sync"
	"time"
)

const (
	StatusNoStatus int = iota
	StatusInWork
	StatusDone
	StatusFailed
)

var taskCounter int = 0

// Task
// Базовая структура для очень важных задач, которые решает наш сервис
type Task struct {
	uuid         uuid.UUID
	id           int
	CreationTime *time.Time // время создания
	FinishTime   *time.Time // время завершения
	Status       int
	Error        error
}

func (t *Task) GetId() int {
	return t.id
}

func (t *Task) GetUUID() uuid.UUID {
	return t.uuid
}

// WorkerInterface
// Интерфейс, описывающий базовый воркер для тасок
type WorkerInterface interface {
	Process(*Task, chan bool, *sync.WaitGroup)
	SetOutput(chan *Task)
}

// NewTask
// Автоматически создает UUID и id
func NewTask() *Task {
	uuidNew := uuid.New()
	taskCounter++
	return &Task{
		uuid:   uuidNew,
		id:     taskCounter,
		Status: StatusNoStatus,
		Error:  nil,
	}
}

// SetCreationTime
// для удобства присвоения времени создания таска (берёт значение, присваивает ссылку)
func (t *Task) SetCreationTime(newTime time.Time) {
	t.CreationTime = &newTime
}

// SetFinishTime
// Для удобства присвоения времени завершения таска (берёт значение, присваивает ссылку)
func (t *Task) SetFinishTime(newTime time.Time) {
	t.FinishTime = &newTime
}

func (t *Task) String() string {
	var execTime int64 = 0
	if t.CreationTime != nil && t.FinishTime != nil {
		execTime = t.FinishTime.Sub(*t.CreationTime).Milliseconds()
	}

	if t.Status == StatusFailed {
		return fmt.Sprintf("Task id: %d time: %s error: %s execTime: %d", t.id, t.CreationTime.Format(time.StampNano), t.Error.Error(), execTime)
	}
	return fmt.Sprintf("id: %d time: %s execTime: %d", t.id, t.CreationTime.Format(time.StampNano), execTime)
}

// MakeFailed
// На случай, если потом необходимо будет добавить логику на провал таски
func (t *Task) MakeFailed(err error) {
	t.Status = StatusFailed
	t.Error = err
}

// MakeDone
// На случай, если потом необходимо будет добавить логику на провал таски
func (t *Task) MakeDone() {
	t.Status = StatusDone
}
