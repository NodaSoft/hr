 package main

import (
	"fmt"
	"time"
	"sync"
	"strings"
)
// ***
// A.N.Chikilevsky  call8969081096@gmail.com(c) Enjoy!
// ****

// Приложение эмулирует получение и обработку неких тасков. Пытается и получать, и обрабатывать в многопоточном режиме.
// Приложение должно генерировать таски 10 сек. Каждые 3 секунды должно выводить в консоль результат всех обработанных к этому моменту тасков (отдельно успешные и отдельно с ошибками).

// ЗАДАНИЕ: сделать из плохого кода хороший и рабочий - as best as you can.
// Важно сохранить логику появления ошибочных тасков.
// Важно оставить асинхронные генерацию и обработку тасков.
// Сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через pull-request в github
// Как видите, никаких привязок к внешним сервисам нет - полный карт-бланш на модификацию кода.

type Ttype struct {
	id         int
	cT         string // время создания
	fT         string // время выполнения
	taskRESULT []byte
}

var wg sync.WaitGroup
var superChan = make(chan Ttype, 10)
var doneTasks = make(chan Ttype)
var undoneTasks = make(chan Ttype)

func createTask () {
	defer  wg.Done()
	            var result string
				dataCreate := time.Now().Format(time.RFC3339)
	            if time.Now().Nanosecond()%2 > 0 {
					result = "Some error occured here"
				} else {
					result = dataCreate
				}
	            superChan <- Ttype{cT: result, id: int(time.Now().UnixNano())}
	            time.Sleep(time.Second * 10)
}

func getTask() {
	defer  wg.Done()
			for t := range superChan {
				t1 := processingTask(t)
			    sortTask(t1)
		    }
}

func sortTask (t Ttype) {
	if (!(strings.TrimSpace(string(t.taskRESULT)) == "successed")) {
		// fmt.Printf("-- NON ++ SUCCES %d \n",t.id)
		undoneTasks <- t
	} else {
		// fmt.Printf("-- SUCCESS %d \n",t.id)
		doneTasks <- t
	}
}

func processingTask (a Ttype) Ttype {
	    tt, _ := time.Parse(time.RFC3339, a.cT)
 	
		if tt.After(time.Now().Add(-20 * time.Second)) {
		    a.taskRESULT = []byte("successed")
		} else {
			a.taskRESULT = []byte("something went wrong")
		}
	    a.fT = time.Now().Format(time.RFC3339Nano)
	    time.Sleep(time.Millisecond * 150)
		return a
}

func rangeTaskLeft() {
	time.Sleep(3 * time.Second)
	for {
		     for strUndonTask := range undoneTasks {
				 fmt.Printf("UNDONE: %d == \n",strUndonTask.id)
		     }
     }
}

func rangeTaskRight() {
	time.Sleep(3 * time.Second)
	for {
    		 for donеTask := range doneTasks {
			     fmt.Printf("DONE: %d == \n", donеTask.id)
		     }
    }
}

func contextTime () {
		  go createTask()
		  go getTask()
		  go rangeTaskRight()
		  go rangeTaskLeft()
}

func fun() {
	for timeout := time.After(10 * time.Second); ; {
        select {
        case <-timeout:
            return
        default:
		  wg.Add(1)
          contextTime()
        }
    }    
	wg.Wait()
}

func main() {
    fun()
    wg.Wait()
}


  
	 
