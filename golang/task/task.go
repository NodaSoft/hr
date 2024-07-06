package task

import (
	"fmt"
	"time"
)

// A Task represents a meaninglessness of our life
type Task struct {
	id             int
	creationTime   string // Время создания
	completionTime string // Время выполнения
	result         []byte
	completed      bool
}

func Create() Task {
	creationTime := time.Now().Format(time.RFC3339)
	if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
		creationTime = "Some error occurred"
	}
	// Таски генерятся слишком быстро, чтобы брать просто юникс тайм по миллисекундам *И* делать это айдишником
	// Т.е. Желательно использовать что-то более подходящее для айди, но и так сойдет для этой задачи
	return Task{creationTime: creationTime, id: int(time.Now().UnixNano())}
}

func (t *Task) Work() {
	tt, _ := time.Parse(time.RFC3339, t.creationTime)
	if tt.After(time.Now().Add(-30 * time.Second)) {
		t.completed = true
		t.result = []byte("task has been successed")
	} else {
		t.completed = false
		t.result = []byte("something went wrong")
	}

	t.completionTime = time.Now().Format(time.RFC3339Nano)
	time.Sleep(time.Millisecond * 150)
}

func (t *Task) IsCompleted() bool {
	return t.completed
}

func (t *Task) String() string {
	if t.IsCompleted() {
		return fmt.Sprintf("Task id: %d time: %s, result: %s", t.id, t.creationTime, t.result)
	} else {
		return fmt.Sprintf("Task id: %d time: %s, error: %s", t.id, t.creationTime, t.result)
	}
}
