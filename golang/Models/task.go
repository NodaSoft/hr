package Models

import "time"

// TaskStatus Represents status of a Task.
type TaskStatus string

const (
	InProgress TaskStatus = "In Progress" // [In Progress] task status.
	Successful TaskStatus = "Successful"  // [Successful] task status.
	Error      TaskStatus = "Error"       // [Error] task status.
)

// A Task represents a meaninglessness of our life.
type Task struct {
	Id          int64      // Task primary key.
	CreatedAt   time.Time  // Date of Task creation.
	CompletedAt time.Time  // Date of Task completion.
	Status      TaskStatus // Status of a Task.
	Message     string     // Task message.
}

// NewTask create an instance of a task.
func NewTask() *Task {
	var timeNow = time.Now()

	return &Task{
		Id:        timeNow.UnixNano(),
		CreatedAt: timeNow,
		Status:    InProgress,
	}
}

// Update Set new status of a Task with some message.
func (task *Task) Update(status TaskStatus, message string) {
	task.Status = status
	task.Message = message
}
