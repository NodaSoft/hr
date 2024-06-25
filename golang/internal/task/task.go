package task

import (
	"fmt"
	"time"
)

// Task structure
type Task struct {
	Id      int
	cTime   string // create time
	eTime   string // execution time
	tResult bool   // task result
}

// ITask is interface for Task structure
type ITask interface {
	Worker() (Task, error)
}

// Worker is divided tasks for created success or with error
func (t Task) Worker() (Task, error) {
	tt, _ := time.Parse(time.RFC3339, t.cTime)
	if tt.After(time.Now().Add(-20 * time.Second)) {
		t.tResult = true
	} else {
		t.tResult = false
	}
	t.eTime = time.Now().Format(time.RFC3339Nano)

	time.Sleep(time.Millisecond * 150)

	if t.tResult {
		return t, nil
	} else {
		return t, fmt.Errorf("task Id %d time %s, error something went wrong", t.Id, t.cTime)
	}
}

// Create creates tasks while main function not quited
func Create(c chan Task, quit <-chan time.Time) {
	for {
		select {
		case <-quit:
			return
		default:
			createdTime := time.Now().Format(time.RFC3339)
			if time.Now().Nanosecond()%2 > 0 {
				createdTime = "Some error occurred"
			}
			c <- Task{Id: int(time.Now().Unix()), cTime: createdTime}
		}
	}
}
