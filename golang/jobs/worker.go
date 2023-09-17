package jobs

import (
	"fmt"
	"sync"
	"time"
)

func StartWorker(tc <-chan *Task, rc chan *TaskResult, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		for task := range tc {
			result := TaskResult{
				ID:        task.ID,
				CreatedAt: task.CreatedAt,
			}
			/*
				Не знаю в чем прикол,
				потому что у меня в любой Go проге
				время в наносекундах почти всегда заканчивается на несколько нулей,
				и это условие почти никогда не срабатывает.

				Чтобы увидеть эти ошибки нужно перебирать таски тысячами (go run . -q1000)

				Раскомментите /1е5 чтобы ошибки были чаще.
			*/
			if task.CreatedAt.Nanosecond() /* /1e5 */ %2 > 0 { // вот такое условие появления ошибочных тасков
				result.Error = fmt.Errorf("error occurred")
				result.Payload = ""
			} else {
				result.Payload = "task succeed"
			}

			time.Sleep(time.Millisecond * 150)
			result.FinishedAt = time.Now()

			rc <- &result
		}
	}()
}
