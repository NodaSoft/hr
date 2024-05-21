package main

import (
	"context"
	"fmt"
	"time"
)

type Ttype struct {
	id     int
	cT     string
	fT     string
	result []byte
}

func taskCreate(a chan Ttype, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			close(a) // закрытие канала
			return
		default:
			ft := time.Now().Format(time.RFC3339)
			curTask := Ttype{cT: ft, id: int(time.Now().Unix())}
			//if time.Now().Second()%2 > 0 { // для теста, поставил проверку по секундам
			if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков. никогда не работает.
				curTask.result = []byte("Some error occured")
			}
			a <- curTask
			//time.Sleep(time.Millisecond * 100) // для теста, чтобы пропускать операции по каждой милисекунде
		}
	}
}

func taskWork(a Ttype) Ttype {
	tt, _ := time.Parse(time.RFC3339, a.cT)
	if tt.After(time.Now().Add(-20 * time.Second)) {
		if a.result != nil {
			a.result = []byte("something went wrong")
		} else {
			a.result = []byte("task has been successed")
		}
	}

	a.fT = time.Now().Format(time.RFC3339Nano)
	time.Sleep(time.Millisecond * 150)

	return a
}

func taskSort(curTask Ttype, doneTask map[int]Ttype, undoneTask map[int]error) {
	if string(curTask.result) == "task has been successed" {
		doneTask[curTask.id] = curTask
	} else {
		undoneTask[curTask.id] = fmt.Errorf("Task id %d time %s, error %s", curTask.id, curTask.cT, curTask.result)
	}
}

func main() {

	ctx, _ := context.WithTimeout(context.Background(), time.Second*3)
	superChan := make(chan Ttype, 10)

	result := map[int]Ttype{}
	err := map[int]error{}

	go taskCreate(superChan, ctx)

	go func() {
		for t := range superChan {
			t = taskWork(t)
			go taskSort(t, result, err)
		}
	}()

	time.Sleep(time.Second * 3)

	println("Errors:")
	for _, r := range err {
		fmt.Println(r)
	}

	println("Done tasks:")
	for r := range result {
		fmt.Println(r)
	}

}
