package main

import (
	"fmt"
	"time"
)

// Приложение эмулирует получение и обработку неких тасков. Пытается и получать, и обрабатывать в многопоточном режиме.
// Приложение должно генерировать таски 10 сек. Каждые 3 секунды должно выводить в консоль результат всех обработанных к этому моменту тасков (отдельно успешные и отдельно с ошибками).

// ЗАДАНИЕ: сделать из плохого кода хороший и рабочий - as best as you can.
// Важно сохранить логику появления ошибочных тасков.
// Важно оставить асинхронные генерацию и обработку тасков.
// Сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через pull-request в github
// Как видите, никаких привязок к внешним сервисам нет - полный карт-бланш на модификацию кода.

// Мы даем тестовое задание чтобы:
// * уменьшить время технического собеседования - лучше вы потратите пару часов в спокойной домашней обстановке, чем будете волноваться, решая задачи под взором наших ребят;
// * увеличить вероятность прохождения испытательного срока - видя сразу стиль и качество кода, мы можем быть больше уверены в выборе;
// * снизить число коротких собеседований, когда мы отказываем сразу же.

// Выполнение тестового задания не гарантирует приглашение на собеседование, т.к. кроме качества выполнения тестового задания, оцениваются и другие показатели вас как кандидата.

// Мы не даем комментариев по результатам тестового задания. Если в случае отказа вам нужен наш комментарий по результатам тестового задания, то просим об этом написать вместе с откликом.

// A Ttype represents a meaninglessness of our life
type Ttype struct {
	id         int
	cT         string // время создания
	fT         string // время выполнения
	taskRESULT []byte
}

func taskCreator(a chan<- Ttype) {
	tiker := time.NewTicker(10 * time.Second)
	defer tiker.Stop()
	defer close(a)
	for {

		select {
		case <-tiker.C:
			return
		default:
			ft := time.Now().Format(time.RFC3339)
			// На мом устройсте не получается генерировать ошибочные случаи из за промежутков генерации, для тестов заменял на Second()
			if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
				ft = "Some error occured"
			}
			a <- Ttype{cT: ft, id: int(time.Now().Unix())} // передаем таск на выполнение
		}
	}
}

func handle(a Ttype) Ttype {
	tt, err := time.Parse(time.RFC3339, a.cT)
	/*
		Добавлена обработка ошибки, как использование идеалогии Golang,
		а так же для уменьшения времени вполнения функции (избегаем остальной работы с функциями пакета time)
		Хотя решение при котором в случае ошибки устанавливается дефолтное время и таким образом устанавливается результат тоже красивое,
		однако требует дополнительного времени на сравнение дат.
	*/
	if err != nil {
		a.taskRESULT = []byte("something went wrong")
	} else {
		if tt.After(time.Now().Add(-20 * time.Second)) {
			a.taskRESULT = []byte("task has been successed")
		} else {
			a.taskRESULT = []byte("something went wrong")
		}
	}

	a.fT = time.Now().Format(time.RFC3339Nano)

	time.Sleep(time.Millisecond * 150)

	return a
}

func handleAndSort(done chan<- Ttype, undone chan<- error, tasks <-chan Ttype) {
	defer close(done)
	defer close(undone)
	for task := range tasks {
		// Обрабатываем задачу
		task := handle(task)
		// Распределяем результат
		if string(task.taskRESULT[14:]) == "successed" {
			done <- task
		} else {
			undone <- fmt.Errorf("Task id %d time %s, error %s", task.id, task.cT, task.taskRESULT)
		}
	}
}

func main() {
	superChan := make(chan Ttype, 10)
	doneTasks := make(chan Ttype, 10)
	undoneTasks := make(chan error, 10)

	go taskCreator(superChan)

	/*
		Совмещён функционал handle и sort для того чтобы не создавать ещё один промежуточный канал.
		Код менее читаемый, зато не использует лишнего места.
	*/
	go handleAndSort(doneTasks, undoneTasks, superChan)

	/*
		Раасмотрена идея сохранять успешно выполненные таски в мапе по UnixTimestamp,
		из зи того, что таски генерируются достаточно быстро и внутри результирующей мапы они будут накладываться друг на друга
		и чать выполненных тасков будет затираться.
		За этим мапа заменена на слайс тасков
	*/
	result := []Ttype{}
	err := []error{}
	tiker := time.NewTicker(3 * time.Second)
	defer tiker.Stop()

	/*
		Это можно вынести функцию для лучшей читаемости кода и для того чтобы продемонстрировать умение пользоваться sync.Waitgroup
		однако я оставил так для меньших затрат памяти и для того чтобы не создавать ещё один объект
	*/
	for {
		select {
		case t, ok := <-doneTasks:
			if !ok {
				doneTasks = nil
			} else {
				result = append(result, t)
			}
		case e, ok := <-undoneTasks:
			if !ok {
				undoneTasks = nil
			} else {
				err = append(err, e)
			}
		case <-tiker.C:
			// Изменено дабы видеть содержимое а не просто порядковые номера в слайсе
			println("Errors:")
			for _, r := range err {
				println(r.Error())
			}

			println("Done tasks:")
			for _, r := range result {
				println(r.id)
			}

			result = []Ttype{}
			err = []error{}

			if doneTasks == nil && undoneTasks == nil {
				return
			}
		}
	}
}
