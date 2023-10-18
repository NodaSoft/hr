package tasks

import (
	"sync"
	"time"
)

type Task struct {
	id          int // int for simplicity and readability (could be UUID)
	mux         *sync.RWMutex
	createdAt   string // time type is part of the task
	completedAt string // this type is for symmetry
	err         error
}

func New(id int, createdAt string) *Task {
	return &Task{
		id:        id,
		mux:       &sync.RWMutex{},
		createdAt: createdAt,
	}
}

// Sets task completion time and (optionally) error if err is not nil.
func (t *Task) MarkAsCompleted(err error) {
	t.mux.Lock()
	t.completedAt = time.Now().Format(time.RFC3339Nano)
	t.err = err
	t.mux.Unlock()
}

// Error is always nil if task is not completed.
func (t *Task) State() (isCompleted bool, withError error) {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.completedAt != "", t.err
}

func (t *Task) Id() int {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.id
}

func (t *Task) CreatedAt() string {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.createdAt
}

func (t *Task) CompletedAt() string {
	t.mux.RLock()
	defer t.mux.RUnlock()
	return t.completedAt
}
