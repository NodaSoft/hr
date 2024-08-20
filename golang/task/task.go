package task

import (
	"time"
)

type Task struct {
	ID        int
	TimeStart time.Time
	TimeClose time.Time
	ErrorIs   bool
	Result    []byte
}

func New(i int) *Task {
	t := time.Now()

	return &Task{
		ID:        i,
		TimeStart: t,
		ErrorIs:   t.Nanosecond()%2 > 0,
	}
}

func (t *Task) Work() {
	if t.TimeStart.After(time.Now().Add(-20*time.Second)) && !t.ErrorIs {
		t.Result = []byte("success")
	} else {

		t.Result = []byte("wrong")
	}
	time.Sleep(time.Millisecond * 150)
}
