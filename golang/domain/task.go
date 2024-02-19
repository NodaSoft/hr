package domain

import "time"

// TaskStatus Represents status of a Task.
type state string

const (
	Created    state = "Created"     // [Created] Таск в этом статусе, когда только создан
	InProgress state = "In Progress" // [InProgress] Таск в этом статусе, когда взят в обработку
	Successful state = "Successful"  // [Successful] Таск в этом статусе, когда обработан успешно
	Errored    state = "Error"       // [Errored] Таск в этом статусе, когда обработан с ошибкой
)

type Task struct {
	ID         int
	CreatedAt  time.Time // Время создания таска
	FinishedAt time.Time // Время завершения обработки таска
	state      state     // Текущее состояние таска
	payload    string    // Сообщение
}

func NewTask() *Task {
	return &Task{
		ID:        getNextTaskId(),
		CreatedAt: time.Now(),
		state:     Created,
	}
}

func (t *Task) IsSuccessful() bool {
	return t.state == Successful
}

func (t *Task) MarkAsSuccessfullyCompleted() {
	t.FinishedAt = time.Now()
	t.state = Successful
}

func (t *Task) MarkAsErrored() {
	t.FinishedAt = time.Now()
	t.state = Errored
}

func (t *Task) MarkAsInProgress() {
	t.state = InProgress
}
