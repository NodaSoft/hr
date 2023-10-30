package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

const SenderThreadCount = 20
const ReceiverThreadCount = 10
const TotalMsgCount = 100
const SendInterval = 100
const ReadInterval = 100
const InitMsgNum = 10000

type TaskSender struct {
	OutTaskChan chan Task
	MsgNum      int
	SendMutex   sync.Mutex
	SendWG      sync.WaitGroup
}

func (sender *TaskSender) SendTask() {
	defer sender.SendWG.Done()
	needSend := true

	for needSend {
		var task Task
		task.SendResult = 0
		task.CreateTime = time.Now().Format(time.RFC3339)

		sender.SendMutex.Lock()
		sender.MsgNum++
		task.ID, _ = strconv.Atoi(fmt.Sprintf("%d%d",
			int(time.Now().Unix()), sender.MsgNum),
		)
		needSend = sender.MsgNum <= InitMsgNum+TotalMsgCount
		sender.SendMutex.Unlock()

		if time.Now().Nanosecond()%2 > 0 {
			task.SendMessage = "Some error occured"
			task.SendResult = -1
		}
		if needSend {
			sender.OutTaskChan <- task
			time.Sleep(time.Millisecond * SendInterval)
		}
	}
}

func (sender *TaskSender) Start() {
	sender.OutTaskChan = make(chan Task, SenderThreadCount*4)
	sender.MsgNum = InitMsgNum
	go func() {
		sender.SendWG.Add(SenderThreadCount)
		for i := 0; i < SenderThreadCount; i++ {
			go sender.SendTask()
		}
		sender.SendWG.Wait()
		close(sender.OutTaskChan)
	}()
}

type Task struct {
	ID             int
	CreateTime     string
	ReceiveTime    string
	SendResult     int
	SendMessage    string
	ReceiveResult  int
	ReceiveMessage string
}

type TaskReceiver struct {
	InputTaskChan   chan Task
	SuccesTasksChan chan Task
	ErrorTasksChan  chan error
	SuccesMap       map[int]Task
	ErrorList       []error
	CommitMutex     sync.Mutex
	ReceiveWG       sync.WaitGroup
	CommitSuccWG    sync.WaitGroup
	CommitErrWG     sync.WaitGroup
}

func (receiver *TaskReceiver) Start(inputChan chan Task) {
	receiver.SuccesTasksChan = make(chan Task)
	receiver.ErrorTasksChan = make(chan error)
	receiver.InputTaskChan = inputChan
	receiver.SuccesMap = make(map[int]Task)

	receiver.ReceiveWG.Add(ReceiverThreadCount)
	receiver.CommitSuccWG.Add(ReceiverThreadCount)
	receiver.CommitErrWG.Add(ReceiverThreadCount)

	for i := 0; i < ReceiverThreadCount; i++ {
		go receiver.Receive()
	}
	for i := 0; i < ReceiverThreadCount; i++ {
		go receiver.CommitSuccess()
	}
	for i := 0; i < ReceiverThreadCount; i++ {
		go receiver.CommitErrors()
	}
	receiver.ReceiveWG.Wait()

	close(receiver.SuccesTasksChan)
	close(receiver.ErrorTasksChan)

	receiver.CommitSuccWG.Wait()
	receiver.CommitErrWG.Wait()
}

func (receiver *TaskReceiver) PerformTask(task *Task) {
	tt, _ := time.Parse(time.RFC3339, task.CreateTime)
	if tt.After(time.Now().Add(-20 * time.Second)) {
		task.ReceiveResult = 0
		task.ReceiveMessage = "task has been successed"
	} else {
		task.ReceiveResult = 10003
		task.ReceiveMessage = "message from the future"
	}
	task.ReceiveTime = time.Now().Format(time.RFC3339Nano)
}

func (receiver *TaskReceiver) SeparateTask(task Task, SuccesTasksChan chan Task, ErrorTasksChan chan error) {
	if task.ReceiveResult == 0 && task.SendResult == 0 {
		SuccesTasksChan <- task
	} else {
		ErrorTasksChan <- fmt.Errorf(
			"Task id %d time %s, send error: %s, receive error: %s",
			task.ID, task.CreateTime,
			task.SendMessage, task.ReceiveMessage,
		)
	}
}

func (receiver *TaskReceiver) Receive() {
	defer receiver.ReceiveWG.Done()
	for t := range receiver.InputTaskChan {
		receiver.PerformTask(&t)
		receiver.SeparateTask(t, receiver.SuccesTasksChan, receiver.ErrorTasksChan)
		time.Sleep(time.Millisecond * ReadInterval)
	}
}

func (receiver *TaskReceiver) CommitSuccess() {
	defer receiver.CommitSuccWG.Done()
	for r := range receiver.SuccesTasksChan {
		receiver.CommitMutex.Lock()
		receiver.SuccesMap[r.ID] = r
		receiver.CommitMutex.Unlock()
	}
}

func (receiver *TaskReceiver) CommitErrors() {
	defer receiver.CommitErrWG.Done()
	for r := range receiver.ErrorTasksChan {
		receiver.ErrorList = append(receiver.ErrorList, r)
	}
}

func (receiver TaskReceiver) PrintResult() {
	println("Errors:")
	for _, e := range receiver.ErrorList {
		println(e.Error())
	}

	println("Done tasks:")
	for r := range receiver.SuccesMap {
		println(r)
	}
}

func main() {
	sender := TaskSender{}
	sender.Start()

	receiver := TaskReceiver{}
	receiver.Start(sender.OutTaskChan)

	receiver.PrintResult()
}
