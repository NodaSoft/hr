package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнены остальных тасков

// Задачи создаются и обрабатываются в непрерывном режиме, выполнения завершается срабатыванием сигнала Ctrl-c
// При этом основная горутина дожидается завершения задач которые остались в каналах и которые в данный момент выполняются
// Происходит мягкое завершения выполнения, после закрытия всех каналов выводится результат в консоль

// A Task represents a meaninglessness of our life
// Gaudeamus igitur
type Task struct {
	id              int
	createTime      string // время создания
	compilationTime string // время выполнения
	taskResult      string
	tError          bool
}

type TaskChannels struct {
	done   chan *Task
	undone chan error
	worker chan *Task
	sorter chan *Task
}

var taskChannels = TaskChannels{
	done:   make(chan *Task),
	undone: make(chan error),
	worker: make(chan *Task, 10),
	sorter: make(chan *Task, 10),
}

var exit chan struct{}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	wg := new(sync.WaitGroup)
	wg.Add(3) // taskCreturer, getTaskResultDone, getTaskResultUndone

	exit := make(chan struct{}, 1)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	result := map[int]Task{}
	err := []error{}

	go taskCreturer(ctx, wg)
	go getTaskOfChannelWorker()
	go taskSorter()
	go getTaskResultDone(&result, wg)
	go getTaskResultUndone(&err, wg)

	go func() { // запуск горутины которая будет отслеживать получение сигнала CTRL-C
		<-sig              // получение сигнала
		exit <- struct{}{} // отправка сигнала в цикл for горутины main
	}()
	for {
		<-exit    // получение сигнала на выход
		cancel()  // отправка сигнала на выход в контекст
		wg.Wait() // ожидает закрытия каналов и обработку оставшихся задач
		displayTaskResult(&result, &err)
		return
	}

	// time.Sleep(time.Second * 5)
}

func taskCreturer(ctx context.Context, wg *sync.WaitGroup) {
	createTask := func() *Task {
		ft := time.Now().Format(time.RFC3339Nano)
		if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
			ft = "Some error occured"
		}
		// time.Now().Unix() заменён на time.Now().UnixNano(), иначе задачи создаются очень быстро
		// и в результате тут (*result)[r.id] = *r переписывают друг друга
		return &Task{createTime: ft, id: int(time.Now().UnixNano())}
	}
	for {
		select {
		case <-ctx.Done(): // проверяем получение сигнала о завершении
			close(taskChannels.worker) // сигнал получен, закрываем канал
			wg.Done()
			return
		case <-time.After(time.Millisecond * 100): // уменьшения скорости создания задач, иначе вывод огромный
			taskChannels.worker <- createTask() // создаём задачи и отправляем
		}
	}
}

func getTaskOfChannelWorker() {
	wgWorker := new(sync.WaitGroup) // счётчик работающих задач
	for t := range taskChannels.worker {
		wgWorker.Add(1)           // увеличиваем счётчик на 1 задачу
		go t.taskWorker(wgWorker) // запуск каждой задачи в своей горутине(многозадачность)
	}
	// fmt.Println(&wgWorker)
	wgWorker.Wait()            // ожидаем когда завершатся все задачи
	close(taskChannels.sorter) // закрывыем канал сортировки
	return

}

func (t *Task) taskWorker(wgWorker *sync.WaitGroup) {
	defer wgWorker.Done() // задача завершена уменьшаем счётчик
	_, err := time.Parse(time.RFC3339, t.createTime)
	/*
		строка непонятна, зачем и для чего, от текущего времени вычитается 20 секунд и сверяется
		тогда нужно какое то условие, или time sleep, чтобы задачи менялись выполнено или нет
		 if tt.After(time.Now().Add(-20 * time.Second)) {
	*/
	if err != nil {
		t.tError = true
		t.taskResult = "something went wrong"
	} else {
		t.tError = false
		t.taskResult = "task has been successed"
	}
	t.compilationTime = time.Now().Format(time.RFC3339Nano)
	time.Sleep(time.Millisecond * 150)
	taskChannels.sorter <- t // задача обработана, отправка в сортировку
}

func taskSorter() {
	for t := range taskChannels.sorter {
		if t.tError == false {
			taskChannels.done <- t
		} else if t.tError == true {
			taskChannels.undone <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.compilationTime, t.taskResult)
		}
	}
	// срабатывает после того как все данные из канала taskChannel.sorter будут обработаны
	// так как канал taskChannel.sorter буферезированный, все задачи будут отсортированы
	close(taskChannels.done)
	close(taskChannels.undone)
}

func getTaskResultDone(result *map[int]Task, wg *sync.WaitGroup) {
	for r := range taskChannels.done {
		(*result)[r.id] = *r
	}
	defer wg.Done()
}

func getTaskResultUndone(err *[]error, wg *sync.WaitGroup) {
	for r := range taskChannels.undone {
		*err = append(*err, r)
	}
	defer wg.Done()
}

func displayTaskResult(result *map[int]Task, err *[]error) {
	fmt.Println("Done tasks: ", len(*result))
	for _, r := range *result {
		fmt.Printf("Task id %d create time %s, final time %s, result: %s \n", r.id, r.createTime, r.compilationTime, r.taskResult)
	}
	fmt.Println("Errors: ", len(*err))
	for _, r := range *err {
		fmt.Println(r)
	}
}
