package main

import (
	"TestTask/infrastructure/creator"
	"TestTask/infrastructure/worker"
	"time"
)

func GetSuperTaskWorker(workerCount int) worker.Worker[creator.TaskMessage[SuperTask], creator.TaskMessage[SuperTask]] {
	return worker.GetTaskWorker[SuperTask, creator.TaskMessage[SuperTask]](workerCount, SuperTaskExecutor)
}

func SuperTaskExecutor(mes creator.TaskMessage[SuperTask]) creator.TaskMessage[SuperTask] {
	mesValue := mes.GetValue()
	if mes.IsError() {
		mesValue.result = []byte("something went wrong")
	} else {
		mesValue.result = []byte("task has been successed")
	}
	mesValue.dateExecute = time.Now().Format(time.RFC3339Nano)

	time.Sleep(time.Millisecond * 150)

	return SuperTaskMessage{value: mesValue, err: mes.GetError()}
}
