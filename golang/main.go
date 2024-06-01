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

var (
	doneTasks = make(chan Ttype, 100)
 	undoneTasks = make(chan Ttype, 100)
	superChan = make(chan Ttype, 10)
	theEndChan = make(chan struct{})
	
)
const seccssPos = 14 

func taskCreaturer(a chan Ttype, theEndChan chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	defer close(a)
		for {
			select	{
			case <- theEndChan:
				return
			default:
			ft := time.Now().Format(time.RFC3339)
			if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
				ft = "Some error occured"
			}
			a <- Ttype{cT: ft, id: int(time.Now().Unix())} // передаем таск на выполнение
		}
	}
	
}

func task_worker(a Ttype) Ttype {
	tt, _ := time.Parse(time.RFC3339, a.cT) 
	if  tt.After(time.Now().Add(-20 * time.Second)) && a.id%2 == 0 {  // для генерации рандом а не только положительных тестов добавить && a.id%2 == 0
		a.taskRESULT = []byte("task has been successed")
	} else {
		a.taskRESULT = []byte("something went wrong")
	}
	a.fT = time.Now().Format(time.RFC3339Nano)

	time.Sleep(time.Millisecond * 150)

	return a
}

func tasksorter(t Ttype) {
	if string(t.taskRESULT[seccssPos:]) == "successed" {
		doneTasks <- t
	} else {
		undoneTasks <- t
	}
}

func ProcesTask(wg *sync.WaitGroup){
	defer wg.Done()
	for t := range superChan {
		procesTask := task_worker(t)
		tasksorter(procesTask)
	}
	
}

func iahochysleep(wg *sync.WaitGroup){
	defer wg.Done()
	for {
        select {
        case <-theEndChan:
            return
        default:
            time.Sleep(3 * time.Second) // Задержка в 3 секунды
            printRes(doneTasks, undoneTasks)
        }
	}
}
func printRes(doneTasks <-chan Ttype, undoneTasks <-chan Ttype){
	for {
        select {
        case t, ok := <-doneTasks:
			if !ok {
				doneTasks = nil
				continue
			}
            fmt.Printf("Done Task - ID: %d, Created: %s, Finished: %s, Result: %s\n", t.id, t.cT, t.fT, t.taskRESULT)
        case t, ok := <-undoneTasks:
			if !ok {
				undoneTasks = nil 
				continue
			}
            fmt.Printf("Undone Task - ID: %d, Created: %s, Finished: %s, Result: %s\n", t.id, t.cT, t.fT, t.taskRESULT)
		case <-theEndChan: 
            return
        }
    }
}
    

func main() {

	var wg sync.WaitGroup
	

	wg.Add(1)
	go taskCreaturer(superChan,theEndChan,&wg)

	wg.Add(1)
	go ProcesTask(&wg)
	
	wg.Add(1)
	go iahochysleep(&wg)
	
	
	time.Sleep(10 * time.Second)
	close(theEndChan)

	
	wg.Wait()

	close(doneTasks)
	close(undoneTasks)
}