package jobs

import (
	"fmt"
	"sync"
	"time"
)

func formatTime(t time.Time) string {
	return fmt.Sprintf(
		"[%02d:%02d:%02d:%04d]",
		t.Hour(),
		t.Minute(),
		t.Second(),
		t.Nanosecond()/1e5,
	)
}

// как я понял тайный замысел задания, оно пытается логать в параллели и асинхронно
// поэтому я сделал воркеров для логгера
func StartLogger(rc <-chan *TaskResult, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for result := range rc {
			var info string
			if result.Error != nil {
				info = fmt.Sprintf(
					"(ERROR)\u0009ID:%04d CreatedAt:%s FinishedAt:%s\n\u0009↳ Error: %s",
					result.ID,
					formatTime(result.CreatedAt.Local()),
					formatTime(result.FinishedAt.Local()),
					result.Error,
				)
			} else {
				info = fmt.Sprintf(
					"(OK)\u0009ID:%04d CreatedAt:%s FinishedAt:%s\n\u0009↳ Payload: %s",
					result.ID,
					formatTime(result.CreatedAt.Local()),
					formatTime(result.FinishedAt.Local()),
					result.Payload,
				)
			}
			fmt.Println(info)
		}
	}()
}
