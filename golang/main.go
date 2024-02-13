package main

import (
	"fmt"
	"github.com/danyducky/go-abcp/Services"
	"time"
)

const (
	WorkDurationSeconds = 5
	ChannelBandwidth    = 10
)

func main() {
	var worker = Services.NewTaskWorker(ChannelBandwidth)

	go worker.DoWork(WorkDurationSeconds * time.Second)

	for task := range worker.DoneTasks {
		fmt.Println("Success task: ", task.Id)
	}

	for task := range worker.FailedTasks {
		fmt.Println("Failed task: ", task.Id)
	}
}
