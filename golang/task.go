package main

import (
	"fmt"
	"sync"
	"time"
)

// Task represents a work-unit to be done.
type Task struct {
	id           int64
	creationTime time.Time
	finishTime   time.Time
	finished     bool
	result       Result
	initError    error
	mu           sync.RWMutex
}

// NewTask returns a new Task
func NewTask(id int64, creationTime time.Time) *Task {
	t := &Task{
		id:           id,
		creationTime: creationTime,
		finished:     false,
		initError:    nil,
	}

	if t.IsCorrupted() {
		t.initError = fmt.Errorf("initialization error")
	}

	return t
}

// IsCorrupted reports whether the task is correct or not.
func (t *Task) IsCorrupted() bool {
	return t.creationTime.Nanosecond()%2 > 0
}

// Error returns error if task IsCorrupted or nil otherwise.
func (t *Task) Error() error {
	return t.initError
}

// ID returns Task's id.
func (t *Task) ID() int64 {
	return t.id
}

// CreationTime returns Task's creation time.
func (t *Task) CreationTime() time.Time {
	return t.creationTime
}

// FinishTime returns Task's finish time. This is default time.Time struct value if not Finished.
func (t *Task) FinishTime() time.Time {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.finishTime
}

// Finished returns whether the task has been finished or not.
func (t *Task) Finished() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.finished
}

// Result returns Task's Result. This is default Result struct value until the task is Finished.
func (t *Task) Result() Result {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.result
}

// Finish finishes the task with provided Result. It's concurrent-safe.
func (t *Task) Finish(result Result) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.result = result
	t.finished = true
	t.finishTime = time.Now()
}

func (t *Task) String() string {
	return fmt.Sprintf(
		"Task: {id: %d, creation time: %s, finish time: %s, finished: %t, result: %s}",
		t.id,
		t.creationTime.Format(time.RFC3339Nano),
		t.finishTime.Format(time.RFC3339Nano),
		t.finished,
		t.result.String())
}
