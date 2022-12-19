package task

import (
	"time"
)

const (
	MaxExecutionCount = 10
	MaxExecutionTime  = time.Duration(20) * time.Second
	BuilderPause      = time.Duration(2) * time.Second
)

// A Task represents a meaninglessness of our life
type Task struct {
	Id           uint32
	Start        time.Time
	Finish       time.Time
	ErrorMessage string
}
