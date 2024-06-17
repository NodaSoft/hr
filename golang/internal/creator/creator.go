package creator

import (
	"context"
	"taskConcurrency/internal/domain/task"
	"time"
)

type Creator struct{}

func (c *Creator) Create(ctx context.Context, tasks chan<- task.Task) {
    go func(ctx context.Context) {
        defer close(tasks)
        for {
            select {
            case <-ctx.Done():
                return
            default:
                creationTime := time.Now().Format(time.RFC3339)
                if time.Now().UnixMilli()%2 > 0 {  // при .Nanosecond() получаем всегда 2 нуля ->
                    creationTime = "Some error occurred"
                }
                select {
                case tasks <- task.Task{CreationTime: creationTime, Id: int(time.Now().UnixNano())}: //// больше не генерируем таски с одинаковыми ID
                case <-ctx.Done():
                    return
                }
            }
        }
    }(ctx)
}
