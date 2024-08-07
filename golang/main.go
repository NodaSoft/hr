package main

import (
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
func taskCreturer(a chan Ttype) {
	go func() {
		for {
			ft := time.Now().Format(time.RFC3339)
			if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
				ft = "Some error occured"
			}
			a <- Ttype{cT: ft, id: int(time.Now().Unix())} // передаем таск на выполнение
		}
	}()
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

	time.Sleep(time.Millisecond * 150)

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

	superChan := make(chan Ttype, 10)

	go taskCreturer(superChan)

	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)

	go func() {
		// получение тасков
		for t := range superChan {
			t = taskWorker(t)
			go taskSorter(t, doneTasks, undoneTasks)
		}
		close(superChan)
	}()

	result := map[int]Ttype{}
	err := []error{}
	go func() {
		for r := range doneTasks {
			go func() {
				muRes.Lock()
				result[r.id] = r
				muRes.Unlock()
			}()
		}
		for r := range undoneTasks {
			go func() {
				muErr.Lock()
				err = append(err, r)
				muErr.Unlock()
			}()
		}
		close(doneTasks)
		close(undoneTasks)
	}()

	time.Sleep(time.Second * 3)

	println("Errors:")
	muErr.RLock()
	for r := range err {
		println(r)
	}
	muErr.RUnlock()

	println("Done tasks:")
	muRes.RLock()
	for r := range result {
		println(r)
	}
	muRes.RUnlock()
}
