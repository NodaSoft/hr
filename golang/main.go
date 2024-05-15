package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type TaskType struct {
	id         int
	createTime string
	finishTime string
	taskResult []byte
}

type Result struct {
	mw     sync.Mutex
	result map[int]TaskType
	erorrs []error
}

func taskCreator(a chan<- TaskType, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			close(a)
			return
		default:
			createTime := time.Now().Format(time.RFC3339)
			if time.Now().Nanosecond()%2 > 0 {
				createTime = "Some error occured"
			}

			a <- TaskType{
				createTime: createTime,
				id:         int(time.Now().Unix()),
			}
		}
	}
}

func taskSorter(tt TaskType, doneTasks chan<- TaskType, undoneTasks chan<- error) {
	if string(tt.taskResult[14:]) == "successed" {
		doneTasks <- tt
	} else {
		undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", tt.id, tt.createTime, tt.taskResult)
	}
}

func taskWorker(tt TaskType) TaskType {
	tm, _ := time.Parse(time.RFC3339, tt.createTime)

	if tm.After(time.Now().Add(-20 * time.Second)) {
		tt.taskResult = []byte("task has been successed")
	} else {
		tt.taskResult = []byte("something went wrong")
	}
	tt.finishTime = time.Now().Format(time.RFC3339Nano)

	time.Sleep(time.Millisecond * 150)

	return tt
}

func resultGetter(res *Result, doneTasks <-chan TaskType, undoneTasks <-chan error) {
	for {
		select {
		case r, ok := <-doneTasks:
			if !ok {
				doneTasks = nil
			} else {
				res.mw.Lock()
				res.result[r.id] = r
				res.mw.Unlock()
			}
		case err, ok := <-undoneTasks:
			if !ok {
				doneTasks = nil
			} else {
				res.mw.Lock()
				res.erorrs = append(res.erorrs, err)
				res.mw.Unlock()
			}
		}
		if doneTasks == nil && undoneTasks == nil {
			break
		}
	}
}

func main() {
	superChan := make(chan TaskType, 10)
	doneTasks := make(chan TaskType, 10)
	undoneTasks := make(chan error, 10)
	ctx, _ := context.WithTimeout(context.Background(), time.Second*3)

	go taskCreator(superChan, ctx)

	var wg sync.WaitGroup

	for {
		tt, ok := <-superChan
		if !ok {
			break
		}
		wg.Add(1)
		go func(tt TaskType) {
			task := taskWorker(tt)
			taskSorter(task, doneTasks, undoneTasks)
			wg.Done()
		}(tt)
	}

	res := Result{result: make(map[int]TaskType)}
	go resultGetter(&res, doneTasks, undoneTasks)

	wg.Wait()
	close(doneTasks)
	close(undoneTasks)

	fmt.Println("Errors:")
	for _, e := range res.erorrs {
		fmt.Println(e)
	}

	fmt.Println("Done tasks:")
	for _, r := range res.result {
		fmt.Printf("Task id %d completed at %s\n", r.id, r.finishTime)
	}
}
