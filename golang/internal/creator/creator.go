package creator

import (
	"context"
	"taskConcurrency/internal/domain/task"
	"time"
)

type Creator struct{}

func (c *Creator) Create(ctx context.Context, tasks chan<- task.Task) {
	go func() {
		go func() {
			for {
				creationTime := time.Now().Format(time.RFC3339)
				if time.Now().UnixMilli()%2 > 0 { // при .Nanosecond() получаем всегда 2 нуля -> ошибочных всегда 0
					creationTime = "Some error occured"
				}
				tasks <- task.Task{CreationTime: creationTime,
					Id: int(time.Now().UnixNano())} // больше не генерируем таски с одинаковыми ID
			}
		}()
		<-ctx.Done()
		close(tasks)
	}()
}
