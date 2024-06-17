package monitor

import (
	"sync"
	"taskConcurrency/internal/domain/task"
	"time"
)

type Monitor struct{}

func (t *Monitor) PrintWithInterval(s time.Duration, doneTasks <-chan task.Task,
	undoneTasks <-chan error) {
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		result := make(map[int]task.Task)
		err := make(map[string]error)

		var resultLock sync.RWMutex
		var errLock sync.RWMutex

		go func() {
			for doneTask := range doneTasks {
				resultLock.Lock()
				result[doneTask.Id] = doneTask
				resultLock.Unlock()
			}
			wg.Done()
		}()

		go func() {
			for taskErr := range undoneTasks {
				errLock.Lock()
				err[taskErr.Error()] = taskErr
				errLock.Unlock()
			}
			wg.Done()
		}()

		go func() {
			for {
				println("Errors:")
				errLock.RLock()
				for k := range err {
					println(k)
				}
				errLock.RUnlock()
				time.Sleep(s)
			}
		}()
		go func() {
			for {
				println("Done tasks:")
				resultLock.RLock()
				for r := range result {
					println(r)
				}
				resultLock.RUnlock()

				time.Sleep(s)
			}
		}()
	}()
	wg.Wait()
}
