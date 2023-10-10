package main

import (
    "fmt"
    "time"
    "sync"
    "strconv"
)

// Количество потоков и прочие константы для игры с потоками
const SenderThreadCount = 20
const ReceiverThreadCount = 10
const TotalMsgCount = 100
const SendInterval = 100
const ReadInterval = 100
const InitMsgNum = 10000

// Тип отправителя тасков
type TaskSender struct {
    OutTaskChan chan Task
    MsgNum int
    SendMutex sync.Mutex
    SendWG sync.WaitGroup
}

// Метод отправки тасков
func (sender *TaskSender) SendTask () {
    defer sender.SendWG.Done()
    needSend := true
    // Отправляем строго определенное число сообщений
    for needSend {
        var task Task
        task.SendResult = 0
        task.CreateTime = time.Now().Format(time.RFC3339)
        // дальнейшие операции должны быть защищены мьютексом, чтобы не было одинаковых ИД и отправить ровно N сообщений
        // ID генерируем с инкрементом, чтобы был уникален (хотя бы для сессии) - но в идеале тут GUID бы
        sender.SendMutex.Lock()
        sender.MsgNum++
        task.ID, _ = strconv.Atoi(fmt.Sprintf("%d%d", int(time.Now().Unix()), sender.MsgNum))
        needSend = sender.MsgNum <= InitMsgNum + TotalMsgCount
        sender.SendMutex.Unlock()
        // не стоит валить в кучу время и ошибку, делаем поле под сообщение об ошибке
        if time.Now().Nanosecond() % 2 > 0 {
            task.SendMessage = "Some error occured"
            task.SendResult = -1
        }
        // отправляем сообщение и ждем следующего цикла
        if needSend {
            //fmt.Println(len(sender.OutTaskChan))
            sender.OutTaskChan <- task
            time.Sleep(time.Millisecond * SendInterval)
        }
    }
}

// Стартуем отправителя тасков
func (sender *TaskSender) Start () {
    // буффер с запасом, от форс-мажоров чтения = SenderThreadCount * 4
    sender.OutTaskChan = make(chan Task, SenderThreadCount * 4)
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
    ReceiveWG sync.WaitGroup
    CommitSuccWG sync.WaitGroup
    CommitErrWG sync.WaitGroup
}

// запустить приемник тасков
func (receiver *TaskReceiver) Start (inputChan chan Task) chan int {
    // создаем каналы
    receiver.SuccesTasksChan = make(chan Task)
    receiver.ErrorTasksChan = make(chan error)
    receiver.InputTaskChan = inputChan
    receiver.SuccesMap = make(map[int]Task)
    // канал для завершения работы клиента
    manageChan := make(chan int)
    // открываем группы ожидания для потоков обработки
    receiver.ReceiveWG.Add(ReceiverThreadCount)
    receiver.CommitSuccWG.Add(ReceiverThreadCount)
    receiver.CommitErrWG.Add(ReceiverThreadCount)
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
    // ждем завершения все потоков обработки потдельно по группам
    receiver.ReceiveWG.Wait()
    // закрываем каналы сортировок
    close(receiver.SuccesTasksChan)
    close(receiver.ErrorTasksChan)
    // ожидаем завершения чтения сортировок
    receiver.CommitSuccWG.Wait()
    receiver.CommitErrWG.Wait()
    // отправим в управляющий канал сообщение, что все ОК
    go func () {manageChan <- 1} ();
    // счастливые выходим и завершаем основную программу
    return manageChan
}

// Получить таск
func (receiver *TaskReceiver) PerformTask (task *Task) {
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
func (receiver *TaskReceiver) SeparateTask (task Task, SuccesTasksChan chan Task, ErrorTasksChan chan error) {
    // нет ошибок отправки и получения
    if task.ReceiveResult == 0 && task.SendResult == 0 {
        SuccesTasksChan <- task
    } else {
        ErrorTasksChan <- fmt.Errorf("Task id %d time %s, send error: %s, receive error: %s", task.ID, task.CreateTime, task.SendMessage ,task.ReceiveMessage)
    }
}

// обработать результаты получения
func (receiver *TaskReceiver) Receive () {
    defer receiver.ReceiveWG.Done()
    for t := range receiver.InputTaskChan {
        receiver.PerformTask(&t)
        receiver.SeparateTask(t, receiver.SuccesTasksChan, receiver.ErrorTasksChan)
        time.Sleep(time.Millisecond * ReadInterval)
    }
}

// зафиксировать результаты получения
func (receiver *TaskReceiver) CommitSuccess () {
    defer receiver.CommitSuccWG.Done()
    for r := range receiver.SuccesTasksChan {
        receiver.CommitMutex.Lock()
        receiver.SuccesMap[r.ID] = r
        receiver.CommitMutex.Unlock()
    }
}

// зафиксировать ошибки получения
func (receiver *TaskReceiver) CommitErrors () {
    defer receiver.CommitErrWG.Done()
    for r := range receiver.ErrorTasksChan {
        receiver.ErrorList = append(receiver.ErrorList, r)
    }
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
    <-receiver.Start(sender.OutTaskChan)

    // выводим результат на экран
    receiver.PrintResult()
}
