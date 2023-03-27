package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнены остальных тасков

// --------------------------------------------------------------------------------------
// Правки:
// * Разделение структур. Структура задач(таски) и основная структура приложения(service)
// * Назначение функций сущностям. Единая точка входа для сервиса(конструктор)
// * Исправлен способ передачи даты задач. Зачем мы везде передаём дату строку, потом парсим, потом снова форматируем в строку.
//	 Передача времени, если необходимо тогда приведение к строке в отдельных функциях.
// * Исправлено определение ошибочной таски. Некорректное условие, ft и условие могли иметь разницу во времени
// * Исправлен id. Несколько тасков могли бы иметь одинаковый айдишник, если в течении одной секунды бы навалили тасков.
//	 Исправлено на Nanosecond(), но возможно использовать рандомную строку или uuid. Так же можем использовать hash256
// * Добавлен отдельно статус таски и результат(текст). Статус таски теперь является айдишник taskStatusError и taskStatusSuccess
// * Убрана задержка в main, приложение заканчивает работу по закрытыю контекста
// * Убраны лишние каналы и закрытия каналов
// * Произведён полный рефакторинг программы
// * и другое...
// --------------------------------------------------------------------------------------

// A Ttype represents a meaninglessness of our life
type Ttype struct {
	taskChannel     chan task
	successfulTasks []task
	failedTasks     []task
	ctx             context.Context
}

type task struct {
	id            int
	status        int
	result        []byte
	timeCreation  time.Time // время создания
	timeExecution time.Time // время выполнения
}

const (
	contextTimeout      = 5
	chanBuffer          = 10
	taskLifetimeSeconds = 300
)
const (
	taskStatusSuccess = iota
	taskStatusError
)

// NewService Конструктор. Инициализация сервиса и его компонентов
func NewService(timeout int) *Ttype {

	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(time.Duration(timeout)*time.Second))

	service := &Ttype{
		ctx:         ctx,
		taskChannel: make(chan task, chanBuffer),
	}
	go service.taskWorkerCreator()
	go service.taskListener()
	go service.lifetimeErrorTasks()
	go service.lifetimeSuccessTasks()
	return service
}

// Мониторим сколько уже лежит в нашем кеше успешная задача,
// чтобы не упасть по памяти если мы увеличим таймаут приложения
func (d *Ttype) lifetimeSuccessTasks() {
	for i, t := range d.successfulTasks {
		if time.Now().Sub(t.timeCreation).Seconds() > taskLifetimeSeconds {
			lastTask := d.successfulTasks[len(d.successfulTasks)-1]
			d.successfulTasks[i] = lastTask
			d.successfulTasks[len(d.successfulTasks)-1] = task{}
			d.successfulTasks = d.successfulTasks[:len(d.successfulTasks)-1]
		}
	}
	time.Sleep(time.Minute * 5)
}

// Мониторим сколько уже лежит в нашем кеше провальная задача,
// чтобы не упасть по памяти если мы увеличим таймаут приложения
func (d *Ttype) lifetimeErrorTasks() {
	for i, t := range d.failedTasks {
		if time.Now().Sub(t.timeCreation).Seconds() > taskLifetimeSeconds {
			lastTask := d.failedTasks[len(d.failedTasks)-1]
			d.failedTasks[i] = lastTask
			d.failedTasks[len(d.failedTasks)-1] = task{}
			d.failedTasks = d.failedTasks[:len(d.failedTasks)-1]
		}
	}
	time.Sleep(time.Minute * 5)
}

// Воркер для автоматического создания тасков
func (d *Ttype) taskWorkerCreator() {

	for {
		select {
		case <-d.ctx.Done():
			close(d.taskChannel)
			return
		default:
			timekeeper := time.Now()
			taskStatus := taskStatusSuccess
			if timekeeper.Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
				taskStatus = taskStatusError
			}

			// Отправляем таск на обработку
			d.taskChannel <- task{
				timeCreation: timekeeper,
				status:       taskStatus,
				id:           time.Now().Nanosecond(),
			}
		}
	}
}

// Слушатель созданных тасков из канала
func (d *Ttype) taskListener() {
	for {
		select {
		case t := <-d.taskChannel:
			t.taskHandler()
			go d.taskSorter(t)
		case <-d.ctx.Done():
			return
		default:
		}
	}
}

// Сортировщик заданий, успешные и неуспешные
func (d *Ttype) taskSorter(t task) {
	if t.status == taskStatusSuccess {
		d.successfulTasks = append(d.successfulTasks, t)
	} else {
		d.failedTasks = append(d.failedTasks, t)
	}
}

// Распечатаем текущие успешные задачи
func (d *Ttype) printSuccessTasks() {
	println("success tasks:")
	for _, task := range d.successfulTasks {
		printTaskInfo(task)
	}
}

// Распечатаем текущие проваленные задачи
func (d *Ttype) printFailedTasks() {
	println("failed tasks:")
	for _, task := range d.failedTasks {
		printTaskInfo(task)
	}
}

// Обработчик задач, вычисление статусов и других полей
func (t *task) taskHandler() {

	if t.status == taskStatusError {
		t.result = []byte("failed task")
		goto end
	}

	if t.timeCreation.After(time.Now().Add(-20 * time.Second)) {
		t.result = []byte("task has been successes")
	} else {
		t.result = []byte("something went wrong")
	}

end:
	t.timeExecution = time.Now()

	time.Sleep(time.Millisecond * 150)
}

// Ожидаем вылета по таймауту
func (d *Ttype) timer() {
	for {
		select {
		case <-d.ctx.Done():
			log.Println("context timeout...")
			return
		}
	}
}

// Вывод информации о задаче в определённом формате
func printTaskInfo(t task) {
	cTimestamp := t.timeCreation.Format(time.RFC3339)
	fmt.Printf("id: %d, status: %s, created_at: %s, result: %s\n", t.id, statusToCommon(t.status), cTimestamp, t.result)
}

// Получение статуса задания из айдишника в случае необходимости
func statusToCommon(status int) string {
	switch status {
	case taskStatusSuccess:
		return "success"
	case taskStatusError:
		return "error"
	default:
		return "unknown"
	}
}

func main() {

	log.Println("start application...")

	service := NewService(contextTimeout)
	service.timer()
	service.printFailedTasks()
	service.printSuccessTasks()
}
