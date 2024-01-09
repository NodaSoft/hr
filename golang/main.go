package main

import (
	"fmt"
	"sync"
	"time"
)

type Ttype struct {
	id         int
	cT         time.Time 
	fT         time.Time 
	taskResult string
}

func main() {
	const maxWorkers = 10
	tasks := make(chan Ttype)
	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)
	var wg sync.WaitGroup

	go generateTasks(tasks)

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go worker(tasks, doneTasks, undoneTasks, &wg)
	}

	go closeChannels(doneTasks, undoneTasks, &wg)

	results, errors := collectResults(doneTasks, undoneTasks)

	printResults(results, errors)
}

func generateTasks(tasks chan<- Ttype) {
	for i := 0; ; i++ {
		tasks <- Ttype{id: i, cT: time.Now()}
		time.Sleep(100 * time.Millisecond) 
	}
}

func worker(tasks <-chan Ttype, doneTasks chan<- Ttype, undoneTasks chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	for t := range tasks {
		t, err := processTask(t)
		if err != nil {
			undoneTasks <- err
		} else {
			doneTasks <- t
		}
	}
}

func processTask(t Ttype) (Ttype, error) {
	time.Sleep(150 * time.Millisecond)
	t.fT = time.Now()

	if t.cT.Nanosecond()%2 > 0 {
		return t, fmt.Errorf("в задаче id %d произошла ошибка", t.id)
	}

	t.taskResult = "успешно"
	return t, nil
}

func closeChannels(doneTasks, undoneTasks chan Ttype, wg *sync.WaitGroup) {
	wg.Wait()
	close(doneTasks)
	close(undoneTasks)
}

func collectResults(doneTasks <-chan Ttype, undoneTasks <-chan error) (map[int]Ttype, []error) {
	results := make(map[int]Ttype)
	var errors []error
	for {
		select {
		case task, ok := <-doneTasks:
			if ok {
				results[task.id] = task
			}
		case err, ok := <-undoneTasks:
			if ok {
				errors = append(errors, err)
			}
		}
		if len(doneTasks) == 0 && len(undoneTasks) == 0 {
			break
		}
	}
	return results, errors
}

func printResults(results map[int]Ttype, errors []error) {
	fmt.Println("Ошибки:")
	for _, err := range errors {
		fmt.Println(err)
	}
	fmt.Println("Выполненные задачи:")
	for id, task := range results {
		fmt.Printf("ID: %d, Результат: %s\n", id, task.taskResult)
	}
}
