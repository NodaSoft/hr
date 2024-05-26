package main

import (
    "fmt"
    "time"
    "sync"
)

type Ttype struct {
    id         int
    cT         string // время создания
    fT         string // время выполнения
    taskRESULT []byte
}

func main() {
    taskCreator := func(a chan<- Ttype, wg *sync.WaitGroup) {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for {
                ft := time.Now().Format(time.RFC3339)
                if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
                    ft = "Some error occured"
                }
                a <- Ttype{cT: ft, id: int(time.Now().Unix())} // передаем таск на выполнение
                time.Sleep(time.Millisecond * 100) // добавлена пауза, чтобы не перегружать систему
            }
        }()
    }

    superChan := make(chan Ttype, 10)
    var wg sync.WaitGroup

    taskCreator(superChan, &wg)

    taskWorker := func(a Ttype, wg *sync.WaitGroup) {
        wg.Add(1)
        go func() {
            defer wg.Done()
            tt, _ := time.Parse(time.RFC3339, a.cT)
            if tt.After(time.Now().Add(-20 * time.Second)) {
                a.taskRESULT = []byte("task has been successed")
            } else {
                a.taskRESULT = []byte("something went wrong")
            }
            a.fT = time.Now().Format(time.RFC3339Nano)

            time.Sleep(time.Millisecond * 150)

            superChan <- a // возвращаем обработанный таск в канал
        }()
    }

    doneTasks := make(chan Ttype)
    undoneTasks := make(chan error)

    taskSorter := func(t Ttype, wg *sync.WaitGroup) {
        wg.Add(1)
        go func() {
            defer wg.Done()
            if string(t.taskRESULT[14:]) == "successed" {
                doneTasks <- t
            } else {
                undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT, t.taskRESULT)
            }
        }()
    }

    go func() {
        // получение тасков
        for t := range superChan {
            taskWorker(t, &wg)
            time.Sleep(time.Millisecond * 50) // добавлена пауза, чтобы не перегружать систему
        }
        close(superChan)
    }()

    go func() {
        // сортировка тасков
        for {
            select {
            case t := <-doneTasks:
                taskSorter(t, &wg)
            case e := <-undoneTasks:
                taskSorter(Ttype{taskRESULT: []byte(e.Error())}, &wg)
            }
        }
    }()

    result := map[int]Ttype{}
    err := []error{}
    go func() {
        // сохранение тасков
        for {
            select {
            case t := <-doneTasks:
                result[t.id] = t
            case e := <-undoneTasks:
                err = append(err, e)
            }
        }
    }()

    time.Sleep(time.Second * 3)

    wg.Wait() // добавлена остановка, чтобы дождаться завершения всех горутин

    fmt.Println("Errors:")
    for _, r := range err {
        fmt.Println(r)
    }

    fmt.Println("Done tasks:")
    for _, r := range result {
        fmt.Println(r)
    }
}
