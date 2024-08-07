package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Приложение эмулирует получение и обработку неких тасков. Пытается и получать, и обрабатывать в многопоточном режиме.
// Приложение должно генерировать таски 10 сек. Каждые 3 секунды должно выводить в консоль результат всех обработанных к этому моменту тасков (отдельно успешные и отдельно с ошибками).

// ЗАДАНИЕ: сделать из плохого кода хороший и рабочий - as best as you can.
// Важно сохранить логику появления ошибочных тасков.
// Важно оставить асинхронные генерацию и обработку тасков.
// Сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через pull-request в github
// Как видите, никаких привязок к внешним сервисам нет - полный карт-бланш на модификацию кода.

// A Ttype represents a meaninglessness of our life
type Ttype struct {
	id         int
	cT         string // время создания
	fT         string // время выполнения
	taskRESULT []byte
}

// taskCreturer отправляет созданную задачу Ttype в канал a
func taskCreturer(ctx context.Context, ch chan Ttype) {
	for {
		ft := time.Now().Format(time.RFC3339)
		if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
			ft = "Some error occured"
		}
		select {
		case ch <- Ttype{cT: ft, id: int(time.Now().Unix())}:
			continue
		case <-ctx.Done():
			close(ch)
			return
		}
	}
}

// taskWorker выставляет значение taskRESULT для экземпляра структуры Ttype, возвращает Ttype
func taskWorker(a Ttype) Ttype {
	tt, _ := time.Parse(time.RFC3339, a.cT)
	if tt.After(time.Now().Add(-20 * time.Second)) {
		a.taskRESULT = []byte("task has been successed")
	} else {
		a.taskRESULT = []byte("something went wrong")
	}
	a.fT = time.Now().Format(time.RFC3339Nano)

	return a
}

// taskSorter отправляет экземпляр Ttype в канал для выполненных задач, либо отправляет ошибку в канал для ошибок
func taskSorter(t Ttype, doneTasks chan Ttype, undoneTasks chan error) {
	if string(t.taskRESULT[14:]) == "successed" {
		doneTasks <- t
	} else {
		undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT, t.taskRESULT)
	}
}

func main() {

	var muRes sync.RWMutex
	var muErr sync.RWMutex

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	superChan := make(chan Ttype, 10)

	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)
	done := make(chan struct{})

	go taskCreturer(ctx, superChan)

	var wgSort sync.WaitGroup // Ожидание сортировки
	wgSort.Add(1)

	go func() {
		defer wgSort.Done()
		for t := range superChan {
			t = taskWorker(t)
			taskSorter(t, doneTasks, undoneTasks)
		}
	}()

	result := map[int]Ttype{}
	errRes := []error{}

	go func() {
		resOk := true
		errOk := true
		for {
			select {
			case res, ok := <-doneTasks:
				if !resOk && !errOk {
					done <- struct{}{}
					return
				} else if !ok {
					resOk = ok
					break
				}

				muRes.Lock()
				result[res.id] = res
				muRes.Unlock()

			case err, ok := <-undoneTasks:
				if !resOk && !errOk {
					done <- struct{}{}
					return
				} else if !ok {
					errOk = ok
					break
				}
				muErr.Lock()
				errRes = append(errRes, err)
				muErr.Unlock()
			}
		}
	}()
	go func() {
		// Ждём пока все таски будут отсортированы
		wgSort.Wait()
		close(doneTasks)
		close(undoneTasks)

	}()

	for {
		time.Sleep(time.Second * 3)

		println("Errors:")

		muErr.RLock()
		for i, _ := range errRes {
			fmt.Println(errRes[i])
		}
		muErr.RUnlock()

		println("Done tasks:")
		muRes.RLock()
		for i := range result {
			fmt.Println(result[i])
		}
		muRes.RUnlock()

		select {
		case <-done:
			return
		default:
			continue
		}
	}
}
