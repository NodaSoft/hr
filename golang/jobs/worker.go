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
				время в наносекундах всегда заканчивается на несколько нулей,
				и это условие никогда не срабатывает.

				Наверное это как-то связано с частотой процессора.
				Уменьшу до миллисекунд, чтобы оно выдавало ошибки.
			*/
			if task.CreatedAt.Nanosecond()/1e5%2 > 0 { // вот такое условие появления ошибочных тасков
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
