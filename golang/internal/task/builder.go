package task

import (
	"github.com/google/uuid"
	"time"
)

// Builder makes parallel creating new Task and sends them to the Worker
// BuilderPause making after each Task for the convenience of reading the log
type Builder struct {
	Tasks chan Task
}

// Start Builder for endless parallel Task creating with BuilderPause
func (b Builder) Start() {
	for {
		go b.build()
		time.Sleep(BuilderPause)
	}
}

func (b Builder) build() {
	task := Task{Id: uuid.New().ID(), Start: time.Now()}
	b.Tasks <- task
}
