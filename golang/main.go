package main

import (
	"fmt"
	"time"
)

// Приложение эмулирует получение и обработку неких тасков. Пытается и получать, и обрабатывать в многопоточном режиме.
// После обработки тасков в течении 3 секунд приложение должно выводить накопленные к этому моменту успешные таски и отдельно ошибки обработки тасков.

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

//создание тасков
func taskCreturer(a *chan Ttype, b *chan bool) {
	for {
		ft := time.Now().Format(time.RFC3339Nano)
		if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
			ft = "Some error occured"
		}
		select {
			case stop := <- *b:
				if stop == true {
					close(*a)
					close(*b)
					return
				}
			case *a <- Ttype{cT: ft, id: int(time.Now().UnixMilli())}: // передаем таск на выполнение
				time.Sleep(time.Millisecond)
		}
	}
}

//получение тасков
func taskWorker(a, done *chan Ttype, undone *chan error, b *chan bool) {
	for {
		select {
			case stop:= <- *b: 
				if stop == true {
					close(*b)
					return
				}
			case task := <- *a:
				sorter(worker(task), &(*done), &(*undone))
		}	
	}
}

//обработчик тасков
func worker(a Ttype) Ttype {
	_, err := time.Parse(time.RFC3339, a.cT)
	if err != nil {
		a.taskRESULT = []byte("something went wrong")
	} else {
		a.taskRESULT = []byte("task has been successed")
	}
		
	a.fT = time.Now().Format(time.RFC3339Nano)

	time.Sleep(time.Millisecond * 150)
	return a
}

//сортировщик тасков
func sorter(t Ttype, done *chan Ttype, undone *chan error) {
	if string(t.taskRESULT[14:]) == "successed" {
		*done <- t
	} else {
		*undone <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT, t.taskRESULT)
	}
}

//сборщик тасков
func taskResulter(result map[int]Ttype, err *[]error, done *chan Ttype, undone *chan error, b *chan bool) {
	for {
		select {
			case stop := <- *b:
				if stop == true {
					close(*done)
					close(*undone)
					close(*b)
					return
				}
			case r := <- *done:
				result[r.id] = r
			case r := <- *undone:
				*err = append(*err, r)
		}
	}
}

func main() {	
	//управляющие каналы
	stopCreturer := make(chan bool)
	stopWorker := make(chan bool)
	stopResulter := make(chan bool)
	
	//таски
	superChan := make(chan Ttype, 10)
	//создание тасков
	go taskCreturer(&superChan, &stopCreturer)
	
	//каналы сортировки
	doneTask := make(chan Ttype)
	undoneTask := make(chan error)
	//получение тасков
	go taskWorker(&superChan, &doneTask, &undoneTask, &stopWorker)
	
	//результаты обработки
	result := map[int]Ttype{}
	err := []error{}
	//сборщик тасков
	go taskResulter(result, &err, &doneTask, &undoneTask, &stopResulter)	
	
	//жизненный цикл
	time.Sleep(time.Second * 3)
	stopCreturer <- true
	stopWorker <- true
	stopResulter <- true

	//вывод результатов
	println("Errors:")
	for i, r := range err {
		fmt.Println(i+1, r)
	}
	println("Done tasks:")
	i := 0
	for r := range result {
		i++
		fmt.Println(i, r)
	}
}
