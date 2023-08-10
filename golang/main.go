package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнены остальных тасков

// A Ttype represents a meaninglessness of our life
type Ttype struct {
	id         int
	cT         string // время создания
	fT         string // время выполнения
	taskRESULT []byte
}

func main() {
	// время работы программы
	ttl := time.Second * 3

	ctx, cancel := context.WithCancel(context.Background())
	// количество одновременно работающих воркеров
	numJobs := 10

	superChan := make(chan Ttype, numJobs)
	doneTasks := make(chan Ttype, numJobs)
	undoneTasks := make(chan error, numJobs)

	// запускаем генератор тасков
	go taskCreturer(superChan, ctx)
	// запускаем 10 воркеров обрабатывающих таски
	go newWorker(newTaskSorter(doneTasks, undoneTasks)).Run(numJobs, superChan, defaultTaskWorker)

	// создаем и запускаем объекты по хранению результатов и ошибок
	result := newResult()
	err := newErrors()
	for i := 0; i < numJobs; i++ {
		go receiveResult(doneTasks, result)
		go receiveErrors(undoneTasks, err)
	}

	time.Sleep(ttl)
	cancel()

	// вывод результатов и ошибок
	err.Print()
	result.Print()
}

// taskCreturer - принимает канал для отправки тасков и контекст, по которому прекращается отправка
func taskCreturer(task chan Ttype, ctx context.Context) {
loop:
	for {
		select {
		case <-ctx.Done():
			println("break taskCreturer")
			break loop

		default:
			task <- getTask() // передаем таск на выполнение
		}
	}
	close(task)
}

// getTask - генерирует и возвращает таску
func getTask() Ttype {
	ft := time.Now().Format(time.RFC3339)
	if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
		ft = "Some error occured"
	}
	return Ttype{cT: ft, id: int(time.Now().Unix())}
}

// defaultTaskWorker - примает таску, обрабатывает ее и возвращает обратно
func defaultTaskWorker(a Ttype) Ttype {
	tt, _ := time.Parse(time.RFC3339, a.cT)
	if tt.After(time.Now().Add(-20 * time.Second)) {
		a.taskRESULT = []byte("task has been successed")
	} else {
		a.taskRESULT = []byte("something went wrong")
	}
	a.fT = time.Now().Format(time.RFC3339Nano)

	time.Sleep(time.Millisecond * 150)

	return a
}

// taskSorter - сортировщик тасок
type taskSorter struct {
	doneTasks   chan Ttype
	undoneTasks chan error
}

func newTaskSorter(doneTasks chan Ttype, undoneTasks chan error) *taskSorter {
	return &taskSorter{doneTasks, undoneTasks}
}

func (sorter *taskSorter) Sort(t Ttype) {
	if string(t.taskRESULT[14:]) == "successed" {
		sorter.doneTasks <- t
	} else {
		sorter.undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT, t.taskRESULT)
	}
}

// ISorter - интерфейс сортировщика тасок - передается в воркер для сортировки после обработки
type ISorter interface {
	Sort(t Ttype)
}

// taskWorker - тип функции обработки тасок. принимает таску и возвращает ее в обработонном виде
type taskWorker func(a Ttype) Ttype

// taskWorker - воркер - примающий и обрабатывающий таски
type worker struct {
	sorter ISorter
}

func newWorker(sorter ISorter) *worker {
	return &worker{sorter}
}

func (w *worker) work(tasks chan Ttype, taskWorker taskWorker) {
	// получение тасков
	for t := range tasks {
		t = taskWorker(t)
		go w.sorter.Sort(t)
	}
}

// Run - запускает numJobs рутин, примающих и обрабатывающих таски
func (w *worker) Run(numJobs int, tasks chan Ttype, taskWorker taskWorker) {
	for i := 0; i < numJobs; i++ {
		go w.work(tasks, taskWorker)
	}
}

// errorsList - объект для хранения ошибок
type errorsList struct {
	mu         sync.RWMutex
	errorsList []error
}

func newErrors() *errorsList {
	return &errorsList{sync.RWMutex{}, make([]error, 0)}
}

func (e *errorsList) Put(err error) {
	e.mu.Lock()
	e.errorsList = append(e.errorsList, err)
	e.mu.Unlock()
}

func (e *errorsList) Print() {
	println("Errors:")
	e.mu.RLock()
	for r := range e.errorsList {
		println(r)
	}
	e.mu.RUnlock()
}

type IErrorsList interface {
	Put(err error)
}

// receiveErrors - запускает рутину принмающую ошибки по undoneTasks и добавлющую их в errorsList
func receiveErrors(undoneTasks chan error, errorsList IErrorsList) {
	for r := range undoneTasks {
		go errorsList.Put(r)

	}
}

// result - объект хранения результата
type result struct {
	mu     sync.RWMutex
	result map[int]Ttype
}

func newResult() *result {
	return &result{sync.RWMutex{}, make(map[int]Ttype)}
}

func (r *result) Put(t Ttype) {
	r.mu.Lock()
	r.result[t.id] = t
	r.mu.Unlock()
}

type IResult interface {
	Put(r Ttype)
}

func (r *result) Print() {
	println("Done tasks:")
	r.mu.RLock()
	for r := range r.result {
		println(r)
	}
	r.mu.RUnlock()
}

// receiveErrors - запускает рутину принмающую результаты по doneTasks и добавлющую их в result
func receiveResult(doneTasks chan Ttype, result IResult) {
	for r := range doneTasks {
		go result.Put(r)

	}
}
