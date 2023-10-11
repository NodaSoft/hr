package generator

import (
	"task_service/internal/domain"
	"time"

	"github.com/google/uuid"
)

type Creater struct {
	superCh chan domain.Task
}

func New(superCh chan domain.Task) Creater {
	return Creater{
		superCh: superCh,
	}
}

func (c Creater) Run() {
	for {
		ft := time.Now().Format(time.RFC3339)
		if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
			ft = "Some error occured"
		}
		uuid := uuid.New()
		c.superCh <- domain.Task{CreationTime: ft, ID: uuid.String()} // передаем таск на выполнение
	}

}
