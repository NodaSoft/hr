package task

import (
	"fmt"
	"time"
)

// A Task represents a meaninglessness of our life
type Task struct {
	Id         int       // Идентификатор таска
	CreatedAt  time.Time // Время создания
	FinishedAt time.Time // Время выполнения
	Error      error     // Ошибка выполнения
	TaskResult string    // Результат выполнения
}

const (
	taskSucceed = "task has been succeed"
	taskFailed  = "something went wrong"
)

func (t Task) Do() Task {
	if t.CreatedAt.After(time.Now().Add(-20*time.Second)) && t.Error == nil {
		t.TaskResult = taskSucceed
	} else {
		t.TaskResult = taskFailed

		if t.Error == nil {
			t.Error = fmt.Errorf(t.TaskResult)
		}
	}

	t.FinishedAt = time.Now()

	time.Sleep(time.Millisecond * 150)

	return t
}

// String возвращает строковое представление таска.
func (t Task) String() string {
	result := fmt.Sprintf("Task id: %d, time: %s, result: %s", t.Id, t.CreatedAt.Format(time.RFC3339), t.TaskResult)

	if t.Error != nil {
		result = fmt.Sprintf("%s, error: %v", result, t.Error)
	}

	return result
}
