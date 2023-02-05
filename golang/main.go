package main

import (
	"fmt"
	"time"
	"runtime"
)

// A Ttype represents a meaninglessness of our life
type Ttype struct {
	id         int
	cT         string // время создания
	fT         string // время выполнения
	taskRESULT []byte
}

func taskCreturer (a chan Ttype) {
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

func task_worker (a Ttype) Ttype {
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

func tasksorter (t Ttype, doneTasks chan Ttype, undoneTasks chan error) {
    if string(t.taskRESULT[14:]) == "successed" {
        doneTasks <- t
    } else {
        undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT, t.taskRESULT)
    }
}

func getresulttask (doneTasks chan Ttype, undoneTasks chan error, result map[int]Ttype, err []error) {
    for r := range doneTasks {
        go func() {
            result[r.id] = r
        }()
    }
    for r := range undoneTasks {
        go func() {
            err = append(err, r)
        }()
    }
}

func main() {
    runtime.GOMAXPROCS(4) //Колличество ядер для параллельного выполнения,а не конкурентного

	result := map[int]Ttype{}
	err := []error{}

	superChan := make(chan Ttype, 10)
	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)

	go taskCreturer(superChan)

	go func() {
		// получение и обработка тасков
		for t := range superChan {
			t = task_worker(t)
			go tasksorter(t, doneTasks, undoneTasks)
			go getresulttask(doneTasks, undoneTasks,result,err)
		}
		close(superChan)
        close(doneTasks)
        close(undoneTasks)
	}()

    //Собственно ниже вывод результатов
    //но можно дописать параллельные корутины по обработке этих
    //тасков с условием в несколько потоков

	time.Sleep(time.Second * 4)

	println("Errors:")
	for r := range err {
		println(r)
	}

	println("Done tasks:")
	for r := range result {
		println(r)
	}
}

