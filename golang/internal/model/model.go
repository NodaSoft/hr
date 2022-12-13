package model

import "time"

type Task struct {
	ID           int
	CreatedAt    string
	ProcessedAt  string
	Result       []byte
	IsSuccessful bool
	Timestamp    time.Time
}

const (
	TaskCreationError = "some error"
	TaskResultSuccess = "task has been succeed"
	TaskResultFail    = "something went wrong"
)
