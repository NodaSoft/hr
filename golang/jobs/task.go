package jobs

import "time"

// A Task represents a meaninglessness of our life
type Task struct {
	ID        int
	CreatedAt time.Time // время создания
}

type TaskResult struct {
	ID         int
	CreatedAt  time.Time // время создания
	FinishedAt time.Time // время выполнения
	Payload    string
	Error      error
}

func NewTask(id int) *Task {
	return &Task{ID: id, CreatedAt: time.Now()}
}
