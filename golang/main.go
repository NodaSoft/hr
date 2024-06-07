package main

import (
	"TestTask/infrastructure/chanSplitter"
	"TestTask/infrastructure/creator"
	"bufio"
	"context"
	"os"
	"time"
)

func main() {
	println("Programm start")
	startTime := time.Now()
	const workerCount = 10
	ctx, closeFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer closeFn()

	superChan := make(chan creator.TaskMessage[SuperTask], workerCount)

	creatorCompleted := make(chan bool)
	var taskCreator creator.TaskCreator[SuperTask] = SuperTaskCreator{}
	go taskCreator.Start(ctx, superChan, creatorCompleted)

	taskWorker := GetSuperTaskWorker(10)
	workedTasksChan := make(chan creator.TaskMessage[SuperTask], workerCount)
	workerCompleted := make(chan bool)
	go taskWorker.Work(ctx, superChan, workedTasksChan, workerCompleted)

	doneTasks := make(chan creator.TaskMessage[SuperTask], workerCount)
	undoneTasks := make(chan creator.TaskMessage[SuperTask], workerCount)
	sorterCompleted := make(chan bool)
	taskSorter := chanSplitter.ChanSplitter[creator.TaskMessage[SuperTask]]{}
	taskSorter = taskSorter.WithCondition(doneTasks, func(m creator.TaskMessage[SuperTask]) bool { return !m.IsError() })
	taskSorter = taskSorter.WithCondition(undoneTasks, func(m creator.TaskMessage[SuperTask]) bool { return m.IsError() })
	go taskSorter.Split(ctx, workedTasksChan, sorterCompleted)

	result := []SuperTask{}
	err := []error{}
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case task := <-undoneTasks:
				err = append(err, task.GetError())
			case task := <-doneTasks:
				result = append(result, task.GetValue())
			}
		}
	}(ctx)

	go func(ctx context.Context) {
		consoleWriter := bufio.NewWriter(os.Stdout)
		defer consoleWriter.Reset(os.Stdout)
		for {
			timer := time.After(3 * time.Second)
			select {
			case <-ctx.Done():
				return
			case <-timer:
				write(consoleWriter, "Errors:\n")
				for _, errItem := range err {
					write(consoleWriter, errItem.Error())
					write(consoleWriter, "\n")
				}

				write(consoleWriter, "Done tasks:\n")
				for _, resultItem := range result {
					write(consoleWriter, string(resultItem.String()))
					write(consoleWriter, "\n")
				}
				write(consoleWriter, "\n")
				write(consoleWriter, "Elapsed from start: ")
				write(consoleWriter, time.Since(startTime).String())
				write(consoleWriter, "\n")
				write(consoleWriter, "-------------------\n")
				consoleWriter.Flush()
			}
		}
	}(ctx)

	<-ctx.Done()
	<-creatorCompleted
	<-workerCompleted
	<-sorterCompleted
	close(superChan)
	close(workedTasksChan)
	close(doneTasks)
	close(undoneTasks)
	println("Programm end")
}

func write(writer *bufio.Writer, s string) {
	for available := writer.Available(); len(s) > available; available, s = writer.Available(), s[available:] {
		writer.WriteString(s[:available])
		writer.Flush()
	}
	if len(s) > 0 {
		writer.WriteString(s)
	}
}
