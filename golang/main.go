package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

type Task struct {
	Id int
	CrtT,
	FinT time.Time
	Result  string
	Success bool
}

func (t *Task) WorkTask() {
	crtNano := t.CrtT.Nanosecond()
	if crtNano%2 == 0 {
		t.Result = "Task succeeded. Reason: " + strconv.Itoa(crtNano)
		t.Success = true
	} else {
		t.Result = "Task failed. Reason: " + strconv.Itoa(crtNano)
		t.Success = false
	}
	t.FinT = time.Now()
	time.Sleep(time.Millisecond * 150)
}

func (t *Task) String() string {
	return fmt.Sprintf("Task #%v took %v: %v", t.Id, t.FinT.Sub(t.CrtT), t.Result)
}

func genTasks(tCh chan Task, n int) {
	for i := 0; i < n; i++ {
		crT := time.Now()
		tCh <- Task{Id: i, CrtT: crT}
	}
	close(tCh)
}

func workTasks(tCh, rCh chan Task) {
	var wg sync.WaitGroup
	for task := range tCh {
		wg.Add(1)
		go func(task Task) {
			defer wg.Done()
			task.WorkTask()
			rCh <- task
		}(task)
	}
	wg.Wait()
	close(rCh)
}

const TasksCount = 20

func main() {

	tasksCh := make(chan Task, 10)
	taskResCh := make(chan Task, 10)

	go genTasks(tasksCh, TasksCount)
	go workTasks(tasksCh, taskResCh)

	succeededTasks := make([]*Task, 0, TasksCount)
	failedTasks := make([]*Task, 0, TasksCount)

	for task := range taskResCh {
		task := task
		if task.Success {
			succeededTasks = append(succeededTasks, &task)
		} else {
			failedTasks = append(failedTasks, &task)
		}
	}

	fmt.Printf("Succeeded tasks(%v): \n", len(succeededTasks))
	for _, sTask := range succeededTasks {
		fmt.Println(sTask)
	}
	fmt.Printf("Succeeded tasks(%v): \n", len(failedTasks))
	for _, fTask := range failedTasks {
		fmt.Println(fTask)
	}
}
