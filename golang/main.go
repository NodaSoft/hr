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

// A Task represents a meaninglessness of our life
type Task struct {
	id           int
	creationTime string // time of creation
	execTime     string // time of execution
	sucsses      bool   // result of task evaluation; was changed to bool due to lack of purpuse
}

func main() {
	ch := make(chan Task)
	defer func() {
		close(ch)
		fmt.Println("Channel closed ch: ", ch)
	}()

	doneTasks := make(chan Task)
	defer func() {
		close(doneTasks)
		fmt.Println("Channel closed dt: ", doneTasks)
	}()

	undoneTasks := make(chan Task)
	defer func() {
		close(undoneTasks)
		fmt.Println("Channel closed undt: ", undoneTasks)
	}()

	go logic(ch)
	go taskSheduler(ch, doneTasks, undoneTasks)

	var flMut sync.RWMutex
	var scMut sync.RWMutex
	sucssesed := make(map[int]Task)
	failed := make(map[int]Task)

	go taskWrite(failed, sucssesed, doneTasks, undoneTasks, &flMut, &scMut)

	taskRead(failed, sucssesed, &flMut, &scMut)

	fmt.Println("Programm finished!")
}

// Unchanged core logic.
func logic(ch chan Task) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Print("logic panic prevented: ", err)
		}
	}()

	for {
		t := time.Now().Format(time.RFC3339)
		if time.Now().Nanosecond()%2 > 0 {
			t = " Some error occured (logic level)"
		}

		ch <- Task{
			creationTime: t,
			id:           int(time.Now().Unix()),
		}
		//time.Sleep(time.Millisecond * 100) // for debug purpese
	}
}

func taskEval(tsk Task) Task {
	_, err := time.Parse(time.RFC3339, tsk.creationTime)
	if err != nil {
		return Task{
			id:           tsk.id,
			creationTime: tsk.creationTime,
			execTime:     time.Now().Format(time.RFC3339Nano),
			sucsses:      false,
		}
	}

	return Task{
		id:           tsk.id,
		creationTime: tsk.creationTime,
		execTime:     time.Now().Format(time.RFC3339Nano),
		sucsses:      true,
	}
}

func taskSheduler(ch chan Task, dTsk chan Task, uTsk chan Task) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Print("taskSheduler panic prevented: ", err)
		}
	}()

	for {
		t, ok := <-ch
		if !ok {
			break
		}

		tEvl := taskEval(t)
		if tEvl.sucsses {
			dTsk <- tEvl
			continue
		}

		uTsk <- tEvl
	}
}

func taskWrite(
	fl map[int]Task,
	sc map[int]Task,
	dTsk chan Task,
	uTsk chan Task,
	flMut *sync.RWMutex,
	scMut *sync.RWMutex,
) {
	for {
		select {
		case t := <-uTsk:
			{
				flMut.Lock()
				fl[t.id] = t
				flMut.Unlock()
			}
		case t := <-dTsk:
			{
				scMut.Lock()
				sc[t.id] = t
				scMut.Unlock()
			}
		}
	}
}

func taskRead(
	fl map[int]Task,
	sc map[int]Task,
	flMut *sync.RWMutex,
	scMut *sync.RWMutex,
) {
	tk := time.NewTicker(time.Second * 3)
	done := time.NewTicker(time.Second * 10)

	for {
		select {
		case <-tk.C:
			{

				flMut.RLock()
				fmt.Println("Errors:")
				for id, t := range fl {
					fmt.Printf("Task id %v time %v, fail\n", id, t.creationTime)
					delete(fl, id)
				}
				flMut.RUnlock()

				scMut.RLock()
				fmt.Println("Done tasks:")
				for id, t := range sc {
					fmt.Printf("Task id %v time %v, success\n", id, t.creationTime)
					delete(sc, id)
				}
				scMut.RUnlock()
			}
		case <-done.C:
			fmt.Println("Done!")
			return
		}
	}
}
