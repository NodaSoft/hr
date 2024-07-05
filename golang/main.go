package main

import (
	"fmt"
	"sync"
	"time"
)

type Ttype struct {
	id         int
	cT         string 
	fT         string 
	taskRESULT []byte
}

func main() {
	superChan := make(chan Ttype, 10)   
	doneTasks := make(chan Ttype, 10)  
	undoneTasks := make(chan error, 10) 

	var wg sync.WaitGroup  
	var mu sync.Mutex      

	
	taskCreator := func(a chan Ttype) {
		for {
			ft := time.Now().Format(time.RFC3339)
			if time.Now().Nanosecond()%2 > 0 { 
				ft = "Some error occured"
			}
			a <- Ttype{cT: ft, id: int(time.Now().Unix())} 
		}
	}

	
	taskWorker := func(a Ttype) Ttype {
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

	
	taskSorter := func(t Ttype) {
		if string(t.taskRESULT[14:]) == "successed" {
			doneTasks <- t
		} else {
			undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT, t.taskRESULT)
		}
	}

	
	go taskCreator(superChan)

	
	go func() {
		for t := range superChan {
			wg.Add(1) 
			go func(t Ttype) {
				defer wg.Done() 
				t = taskWorker(t)
				taskSorter(t)
			}(t)
		}
	}()

	
	go func() {
		wg.Wait()  
		close(doneTasks)
		close(undoneTasks)
	}()

	
	result := map[int]Ttype{}
	
	var err []error

	
	go func() {
		for r := range doneTasks {
			mu.Lock() 
			result[r.id] = r
			mu.Unlock() 
		}
	}()

	
	go func() {
		for r := range undoneTasks {
			mu.Lock() 
			err = append(err, r)
			mu.Unlock()
		}
	}()

	
	time.Sleep(time.Second * 3)

	
	fmt.Println("Errors:")
	mu.Lock() 
	for _, r := range err {
		fmt.Println(r)
	}
	mu.Unlock() 

	
	fmt.Println("Done tasks:")
	mu.Lock() 
	for _, r := range result {
		fmt.Println(r)
	}
	mu.Unlock() 
}
