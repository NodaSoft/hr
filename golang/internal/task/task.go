package task

import (
	"fmt"
	"time"
)

type Task struct {
	Id             int
	CreationTime   time.Time
	ExpirationTime time.Time
	CompletionTime *time.Time
	Result         []byte
	Error          error
}

func (task Task) IsExpired() bool {
	return time.Now().After(task.ExpirationTime)
}

func (task Task) String() string {
	return fmt.Sprintf(`Task{id: %d, creationTime: %v, completionTime: %v, result: %s, error: %v}`,
		task.Id, task.CreationTime.Format(time.RFC3339),
		task.CompletionTime.Format(time.RFC3339Nano), string(task.Result),
		task.Error)
}
