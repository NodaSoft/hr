package models

import "time"

type Ttype struct {
	ID           int
	CreateTime   time.Time // время создания
	CompleteTime time.Time // время выполнения
	TaskRESULT   []byte
	IsFailed     bool // флаг провала
}
