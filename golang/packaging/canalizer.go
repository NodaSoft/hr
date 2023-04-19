package main

import (
	"log"
	"sync"
	"time"
)

// Canalizer -- можно выделить в отдельный, достаточно универсальный пакет поточной обработки через каналы Go
//
// @author Arhat109-20230419 keep this string only!
// @license: You may use as is, with any kinds on your opinion, without any demands or warranties.
// Выделен относительно универсальный набор функций, способных канализировать обмен чего угодно и как угодно. :)

type Canalizer struct {
	wgCreators sync.WaitGroup
	wgWorkers  sync.WaitGroup
	toChan     chan any
	resChan    chan any
	errChan    chan error
}

func NewCanalizer(toChan, resChan chan any, errChan chan error) *Canalizer {
	return &Canalizer{
		toChan:  toChan,
		resChan: resChan,
		errChan: errChan,
	}
}

// RunCreators -- генерация чего угодно в заданный канал общего назначения в несколько потоков, заданное время
// по завершению всех генераторов закрывает канал toChan самостоятельно
func (can *Canalizer) RunCreators(num int, timeout time.Duration, newItem func() any) {
	for ; num > 0; num-- {
		can.wgCreators.Add(1)
		go func(id int, timeout time.Duration) {
			started := time.Now()
			i := 0
			for time.Since(started) < timeout {
				can.toChan <- newItem()
				i++
			}
			log.Printf("Creator %d ended %d", id, i)
			can.wgCreators.Done()
		}(num, timeout)
	}
	can.WaitCreators()
	close(can.toChan)
}

// RunWorker -- обработчик чего угодно callback способом с возвратом результата или ошибки в каналы
// по завершению обработок закрывает каналы resChan и errChan самостоятельно
func (can *Canalizer) RunWorker(worker func(item any) (any, error)) {
	for t := range can.toChan {
		can.wgWorkers.Add(1)
		go func(item any) {
			result, err := worker(item)
			if err != nil {
				can.errChan <- err
			} else {
				can.resChan <- result
			}
			can.wgWorkers.Done()
		}(t)
	}
	can.WaitWorkers()
	close(can.resChan)
	close(can.errChan)
}

// RunResults -- сборка результатов любой структуры из потоков канальной обработки
// функция получает callback, который должен однозначно идентифицировать результат
func (can *Canalizer) RunResults(guid func(any) (int, error)) (map[int]any, []error) {
	var writeLocker sync.Mutex
	results := make(map[int]any)
	errs := make([]error, 0)
	ok1, ok2 := true, true
	for {
		var res any
		var er error
		select {
		case res, ok1 = <-can.resChan:
			if !ok1 {
				break // from select канал закрыт
			}

			id, err := guid(res)

			if err != nil {
				// нельзя выделять в функцию, т.к. не будет инлайниться из-за инлайн методов мьютекса!
				writeLocker.Lock()
				errs = append(errs, er)
				writeLocker.Unlock()
			} else {
				writeLocker.Lock()
				results[id] = res
				writeLocker.Unlock()
			}
		case er, ok2 = <-can.errChan:
			if !ok2 {
				break
			} // from select канал закрыт
			writeLocker.Lock()
			errs = append(errs, er)
			writeLocker.Unlock()
		}
		if !ok1 && !ok2 {
			break // from for закрыты оба канала
		}
	}
	return results, errs
}

func (can *Canalizer) WaitCreators() {
	can.wgCreators.Wait()
}

func (can *Canalizer) WaitWorkers() {
	can.wgWorkers.Wait()
}
