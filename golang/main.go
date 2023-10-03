package main

import (
    "fmt"
    "time"
    "sync"
)

// Количество потоков
const SenderThreadCount = 5
const ReceiverThreadCount = 5

// Тип отправителя тасков
type TaskSender struct {
    OutTaskChan chan Task
    Seed int
    SendMutex sync.Mutex
}

// Метод отправки тасков
func (sender *TaskSender) SendTask () {
    for {
        var task Task
        task.SendResult = 0
        task.CreateTime = time.Now().Format(time.RFC3339)
        // ID генерируем с инкрементом, чтобы был уникален (хотя бы для сессии)
        task.ID = int(time.Now().Unix()) + sender.Seed
        sender.SendMutex.Lock()
        sender.Seed++
        sender.SendMutex.Unlock()
        // не стоит валить в кучу время и ошибку, делаем поле под сообщение об ошибке
        if time.Now().Nanosecond() % 2 > 0 {
            task.SendMessage = "Some error occured"
            task.SendResult = -1
        }
        sender.OutTaskChan <- task

        time.Sleep(time.Millisecond * 150)
    }
}

// Стартуем отправителя тасков
func (sender *TaskSender) Start () {
    sender.OutTaskChan = make(chan Task, 10)
    sender.Seed = 0
    for i := 0; i < SenderThreadCount; i++ {
        go sender.SendTask()
    }
}

// Тип обработчика тасков
type Task struct {
    ID int
    CreateTime string
    ReceiveTime string
    SendResult int
    SendMessage string
    ReceiveResult int
    ReceiveMessage string
}

// Тип приемника тасков
type TaskReceiver struct {
    InputTaskChan chan Task
    SuccesTasksChan chan Task
    ErrorTasksChan chan error
    SuccesMap map[int]Task
    ErrorList []error
    CommitMutex sync.Mutex
    ReceiveMutex sync.Mutex
}

// Получить таск
func (receiver TaskReceiver) PerformTask (task *Task) {
    tt, _ := time.Parse(time.RFC3339, task.CreateTime)
    // Проверим что таск не из будущего
    if tt.After(time.Now().Add(-20 * time.Second)) {
        task.ReceiveResult = 0
        task.ReceiveMessage = "task has been successed"
    } else {
        // протсто фейковый код ошибки, это можно потом в MAP вынести
        task.ReceiveResult = 10003
        task.ReceiveMessage = "message from the future"
    }
    task.ReceiveTime = time.Now().Format(time.RFC3339Nano)
}

// сортировка тасков - успех/провал
func (receiver TaskReceiver) SeparateTask (task Task, SuccesTasksChan chan Task, ErrorTasksChan chan error) {
    // нет ошибок отправки и получения
    if task.ReceiveResult == 0 && task.SendResult == 0 {
        SuccesTasksChan <- task
    } else {
        ErrorTasksChan <- fmt.Errorf("Task id %d time %s, send error: %s, receive error: %s", task.ID, task.CreateTime, task.SendMessage ,task.ReceiveMessage)
    }
}

// запустить приемник тасков
func (receiver *TaskReceiver) Start (inputChan chan Task) {
    // создаем каналы
    receiver.SuccesTasksChan = make(chan Task)
    receiver.ErrorTasksChan = make(chan error)
    receiver.InputTaskChan = inputChan
    receiver.SuccesMap = make(map[int]Task)
    // запускаем сканнирование входящих тасков
    for i := 0; i < ReceiverThreadCount; i++ {
        go receiver.Receive()
    }
    // запускаем сканнирование фиксации тасков после получения
    for i := 0; i < ReceiverThreadCount; i++ {
        go receiver.CommitSuccess()
    }
    for i := 0; i < ReceiverThreadCount; i++ {
        go receiver.CommitErrors()
    }
}

// обработать результаты получения
func (receiver TaskReceiver) Receive () {
    for t := range receiver.InputTaskChan {
        receiver.ReceiveMutex.Lock()
        receiver.PerformTask(&t)
        receiver.SeparateTask(t, receiver.SuccesTasksChan, receiver.ErrorTasksChan)
        receiver.ReceiveMutex.Unlock()
    }
    close(receiver.InputTaskChan)
}

// зафиксировать результаты получения
func (receiver *TaskReceiver) CommitSuccess () {
    for r := range receiver.SuccesTasksChan {
        receiver.CommitMutex.Lock()
        receiver.SuccesMap[r.ID] = r
        receiver.CommitMutex.Unlock()
    }
    close(receiver.SuccesTasksChan)
}

// зафиксировать ошибки получения
func (receiver *TaskReceiver) CommitErrors () {
    for r := range receiver.ErrorTasksChan {
        receiver.CommitMutex.Lock()
        receiver.ErrorList = append(receiver.ErrorList, r)
        receiver.CommitMutex.Unlock()
    }
    close(receiver.ErrorTasksChan)
}

// Вывести результат на экран
func (receiver TaskReceiver) PrintResult () {
    println("Errors:")
    for _, e := range receiver.ErrorList {
        println(e.Error())
    }

    println("Done tasks:")
    for r := range receiver.SuccesMap {
        println(r)
    }
}

// Основной метод приложения
func main() {
    // стартуем отправителя тасков
    var sender TaskSender
    sender.Start()
    // стартуем приемник тасков
    var receiver TaskReceiver
    receiver.Start(sender.OutTaskChan)

    // ждем
    time.Sleep(time.Second * 3)

    // выводим результат на экран
    receiver.PrintResult()
}
