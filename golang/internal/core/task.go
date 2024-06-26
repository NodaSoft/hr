package core

import "time"

type Task struct {
	Id       int64
	Created  time.Time // время создания
	Finished time.Time // время выполнения
	Result   []byte
}
