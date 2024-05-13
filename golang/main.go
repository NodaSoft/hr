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
// Совсем не понятно почму не синхронизировать? Вероятность одновременной записи очень большая
// Это может привести к потере данных или не правильной записи...
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

func taskSorter(t Ttype, doneTasks chan Ttype, undoneTasks chan error) {
	if string(t.taskRESULT[14:]) == "successed" {
		doneTasks <- t
	} else {
		undoneTasks <- fmt.Errorf("task id %d time %s, error %s", t.id, t.cT, t.taskRESULT)
	}
}

func taskReceiver(superChan chan Ttype, doneTasks chan Ttype, undoneTasks chan error) {
	for t := range superChan {
		t = taskWorker(t)
		go taskSorter(t, doneTasks, undoneTasks)
	}
	close(superChan)
}

func taskResultHandler(doneTasks chan Ttype, undoneTasks chan error, result map[int]Ttype, err *[]error) {
	for r := range doneTasks {
		result[r.id] = r
	}
	for r := range undoneTasks {
		*err = append(*err, r)
	}
	close(doneTasks)
	close(undoneTasks)
}

func main() {
	superChan := make(chan Ttype, 10)
	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)
	result := map[int]Ttype{}
	err := []error{}

	taskCreturer(superChan)
	go taskReceiver(superChan, doneTasks, undoneTasks)
	go taskResultHandler(doneTasks, undoneTasks, result, &err)

	time.Sleep(time.Second * 3)

	println("Errors:")
	for _, e := range err {
		println(e.Error())
	}

	println("Done tasks:")
	for _, r := range result {
		println(r.id)
	}
}
