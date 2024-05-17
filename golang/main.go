package main

import (
 "fmt"
 "time"
)

type Ttype struct {
 id         int
 cT         string
 fT         string
 taskRESULT []byte
}

func main() {
 taskCreturer := func(a chan Ttype) {
  go func() {
   for {
    ft := time.Now().Format(time.RFC3339)
    if time.Now().Nanosecond()%2 > 0 {
     ft = "Some error occured"
    }
    a <- Ttype{cT: ft, id: int(time.Now().Unix())}
    time.Sleep(200 * time.Millisecond) // добавлено небольшая пауза для имитации задержки
   }
  }()
 }

 superChan := make(chan Ttype, 10)

 go taskCreturer(superChan)

 task_worker := func(a Ttype, doneTasks chan Ttype, undoneTasks chan error) {
  tt, _ := time.Parse(time.RFC3339, a.cT)
  if tt.After(time.Now().Add(-20 * time.Second)) {
   a.taskRESULT = []byte("task has been successed")
   doneTasks <- a
  } else {
   a.taskRESULT = []byte("something went wrong")
   undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", a.id, a.cT, a.taskRESULT)
  }
 }

 doneTasks := make(chan Ttype)
 undoneTasks := make(chan error)

 go func() {
  for t := range superChan {
   go task_worker(t, doneTasks, undoneTasks)
  }
 }()

 time.Sleep(3 * time.Second)

 close(superChan)
 close(doneTasks)
 close(undoneTasks)

 fmt.Println("Errors:")
 for err := range undoneTasks {
  fmt.Println(err)
 }

 fmt.Println("Done tasks:")
 for task := range doneTasks {
  fmt.Println(task)
 }
}

