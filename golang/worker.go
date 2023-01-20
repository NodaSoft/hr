package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

//SomeError Ошибка из условия об ошибочных тасках
var SomeError = errors.New("Some error occured")

//WorkerFactory занимается производством Worker
//Так же включает в себе логику Observer что плохо.
//Observer нужно вынести в отдельный класс.
type WorkerFactory struct {
	//doneTask успешно выполненные задачи
	doneTasks chan *Task
	//undoneTasks Задачи которые не были вовремя завершены
	undoneTasks chan error
	//withErr Задачи с ошибкой логику которой требовалось сохранить
	withErr chan *Task

	//workers - список созданных воркеров, храним чтобы рассылать им уведомление о принудительном завершении работы
	//необходимо для graceful shutdown
	workers []*Worker
	//wg Счетчик живых воркеров
	wg sync.WaitGroup

	//destroy - Оповещение о том, что все воркеры завершили работу.
	destroy chan bool

	//Запускает обсервера за воркерами. только один раз
	once sync.Once
}

//observe завершает работу если все воркеры остановились
func (w *WorkerFactory) observe() {
	w.once.Do(func() {
		w.wg.Wait()
		w.Destroy()
	})
}

//Destroy Принудительно завершает работу Worker закрывает каналы
//тут можно добавить больше логики для graceful шатдауна
func (w *WorkerFactory) Destroy() {
	for _, worker := range w.workers {
		worker.destroy <- true
	}
	//ждем ответа воркеров
	w.wg.Wait()

	close(w.doneTasks)
	close(w.undoneTasks)
	close(w.withErr)


	w.destroy <- true
}

//Chans Возвращает список каналов
func (w *WorkerFactory) Chans() (done <-chan *Task, undone <-chan error, withErr <-chan *Task) {
	done = w.doneTasks
	undone = w.undoneTasks
	withErr = w.withErr
	return
}

//NewWorkerFactory создает новую WorkerFactory, а так же канал оповещение о завершении работы.
func NewWorkerFactory() (factory *WorkerFactory, destroy chan bool) {
	destroy = make(chan bool, 1)
	factory = &WorkerFactory{
		doneTasks:   make(chan *Task, 10),
		undoneTasks: make(chan error, 10),
		withErr:     make(chan *Task, 10),
		destroy:     destroy,
	}
	return
}

type Worker struct {
	//Тоже самое что и в WorkerFactory
	//только Send only
	doneTasks   chan<- *Task
	undoneTasks chan<- error
	taskWithErr chan<- *Task

	//канал ожидания приказа о завершении работы
	destroy chan bool

	//beging коллбэк, при стартре работы, что бы Обсервер подсчитывал работающих  воркеров
	begin func()
	//done коллбэк, при завершении работы см. выше
	done func()
	//once используется для того, чтобы worker мог подписаться на события только однажды
	once sync.Once
}

func (w *WorkerFactory) Worker() *Worker {
	return &Worker{
		doneTasks:   w.doneTasks,
		undoneTasks: w.undoneTasks,
		taskWithErr: w.withErr,

		destroy: make(chan bool, 1),

		done: func() {
			w.wg.Done()
		},
		begin: func() {
			w.wg.Add(1)
			go w.observe()
		},
	}
}

//Сортирует таски по каналам по условию успешные\неуспешные\слишком долго ждали обработ
//все таски которые слишком долго ждали обработки идут в канал undone в первую очередь
//(в т.ч c ошибкой которую требовалось сохранить)
// во вторую очередь все таски со внутренней ошбикой которую требовалось сохранить уходтя в канал withErr
// все остальные падают в канал done
func (w *Worker) order(t *Task) {
	if t.result == Wrong {
		w.undoneTasks <- fmt.Errorf("task id: %d, time: %s, result:%s, error: %s", t.id, t.since, t.result, t.error)
		return
	}

	if t.error != nil {
		w.taskWithErr <- t
		return
	}

	w.doneTasks <- t
}

// Subscribe - Worker подписывается на события. Один Worker може подписаться только на один канал
func (w *Worker) Subscribe(input <-chan *Task) {
	f := func() {
		w.subscribe(input)
	}
	w.once.Do(f)
}

func (w *Worker) subscribe(input <-chan *Task) {
	w.begin()
	go func() {
		defer w.done()

		for {
			select {
			case msg, ok := <-input:
				if !ok {
					return
				}

				w.do(msg)
				w.order(msg)
			case <-w.destroy:
				return
			}

		}
	}()
}

//do делает вид, что Worker работает. Ждет 500мс  вместо изначальных 150мс.
func (w *Worker) do(a *Task) {
	tt, _ := time.Parse(time.RFC3339, a.since)
	a.till = time.Now().Format(time.RFC3339Nano)
	a.result = Wrong

	if tt.After(time.Now().Add(-20 * time.Second)) {
		a.result = Success
	}

	time.Sleep(time.Millisecond * 500)
}
