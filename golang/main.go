package main

import (
	"fmt"
	"log"
	"time"
)

type Task struct {
	id              int
	createdTime     time.Time
	performanceTime time.Time
	result          []byte
	err             error
}

func generater(t chan<- *Task) {
	newtask := new(Task)
	newtask.id = int(time.Now().Unix())

	if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
		newtask.err = fmt.Errorf("Some error occurred")
	} else {
		newtask.createdTime = time.Now()
	}
	t <- newtask // передаем таск на выполнение
	log.Print("The task has been generated")
}

func handler(workingTasks <-chan *Task, doneTasks chan<- *Task, fellTasks chan<- *error) {
	for task := range workingTasks {
		if task.createdTime.After(time.Now().Add(-20 * time.Second)) {
			task.result = []byte("task has been successed")
			doneTasks <- task
		} else {
			taskerr := fmt.Errorf("task id %d time %s, result %s, error %s", task.id, task.createdTime, "something went wrong", task.err.Error())
			fellTasks <- &taskerr
		}
		task.performanceTime = time.Now()
		time.Sleep(time.Millisecond * 150)
		log.Print("The task has been processed")
	}
}

func tasks(workingTasks, doneTasks chan *Task, fellTasks chan *error) {
	go generater(workingTasks)
	go handler(workingTasks, doneTasks, fellTasks)
}

func logSuccsess(doneTasks <-chan *Task, result *[]*string) {
	for dt := range doneTasks {
		mes := fmt.Sprintf("Task has been done:%v", *dt)
		*result = append(*result, &mes)
	}
}

func logFault(fellTasks <-chan *error, result *[]*string) {
	for ft := range fellTasks {
		mes := fmt.Sprintf("Task has been fallen: %v", *ft)
		*result = append(*result, &mes)
	}
}

func logger(result *[]*string, done chan bool) {
	for dontstop := true; dontstop; {
		select {
		case <-done:
			log.Print("Received done signal, stopping logger")
			dontstop = false
		default:
			time.Sleep(time.Second * 3)
			log.Print(result)
			for _, res := range *result {
				log.Print(*res)
			}
			*result = (*result)[:0]
			log.Print("'result' slice has been reset")
		}
	}
}

func logTasks(doneTasks chan *Task, fellTasks chan *error, done chan bool) {
	var result []*string

	go logSuccsess(doneTasks, &result)
	go logFault(fellTasks, &result)
	logger(&result, done)
}

func main() {
	workingTasks := make(chan *Task)
	doneTasks := make(chan *Task)
	fellTasks := make(chan *error)
	done := make(chan bool)

	go func() {
		for {
			time.Sleep(time.Second * 1)
			tasks(workingTasks, doneTasks, fellTasks)
		}
	}()

	go func() {
		time.Sleep(10 * time.Second)
		close(done)
	}()

	logTasks(doneTasks, fellTasks, done)
}
