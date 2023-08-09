package steps

import (
	"context"
	st "main/internal/task/structs"
	"time"
)

func RunCreation(ctx context.Context, taskCh chan<- st.Task) {
	for {
		select {
		case <-ctx.Done():
			close(taskCh)
			return
		default:
			createTime := time.Now().Format(time.RFC3339)
			// on some machines, this always returns an even number
			if time.Now().Nanosecond()%2 > 0 {
				createTime = "Some error occurred"
			}
			// int(time.Now().Unix()) - not unique ids
			taskCh <- st.Task{CreateTime: createTime, Id: int(time.Now().Unix())}
		}
	}
}
