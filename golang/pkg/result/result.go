package result

import (
	"abcp-golang/pkg/task"
	"fmt"
	"sync"
	"time"
)

type Results struct {
	results map[int]task.Tasks
	muRes   sync.RWMutex
	err     []error
	muErr   sync.RWMutex
}

func New() *Results {
	return &Results{
		results: make(map[int]task.Tasks),
		muRes:   sync.RWMutex{},
		err:     []error{},
		muErr:   sync.RWMutex{},
	}
}

// покобезопасная обработка готовых тасков
func (r *Results) AddResult(res task.Tasks) {
	r.muRes.Lock()
	defer r.muRes.Unlock()
	r.results[res.Id] = res
}

// потокобезопасная обработка ошибок
func (r *Results) AddError(err error) {
	r.muErr.Lock()
	defer r.muErr.Unlock()
	r.err = append(r.err, err)
}

// потокобезопасный вывод результатов обработки тасков
func (r *Results) PrintResults() {
	r.muRes.RLock()
	defer r.muRes.RUnlock()
	fmt.Println("Done tasks:")
	for _, res := range r.results {
		fmt.Println(res)
	}
}

// потокобезопасный вывод ошибок обработки тасков
func (r *Results) PrintErrors() {
	r.muErr.RLock()
	defer r.muErr.RUnlock()
	fmt.Println("Errors:")
	for _, err := range r.err {
		fmt.Println(err)
	}
}

// Каждые 3 секунды выводит в консоль результат всех обработанных к этому моменту тасков
// (отдельно успешные и отдельно с ошибками)
func Print(done chan struct{}, res *Results) {
	ticker := time.NewTicker(time.Second * 3)
	go func() {
		for {
			select {
			case <-ticker.C:
				res.PrintErrors()
				res.PrintResults()
			case <-done:
				return
			}
		}
	}()
}
