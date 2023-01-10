package main

import (
	"fmt"
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

// Review: переименовать структуру и поля
type Ttype struct {
	id int
	cT string // время создания // Review: время должно быть в юникс тайме
	fT string // время выполнения // Review: время должно быть в юникс тайме
	// Review: разбить на структуру с кодом выполнения и ошибкой
	taskRESULT []byte // Review: переделать в строку
}

func main() {
	// Review: грамматическая ошибка в названии переменной.
	// вынести в отдельную функцию
	// сделать конвейер
	taskCreturer := func(a chan Ttype) {
		go func() {
			for {
				ft := time.Now().Format(time.RFC3339)
				if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
					// Review: в одной переменной хранится и время и ошибка, при чем не информативная
					// исправить
					ft = "Some error occured"
				}
				// Review: ft - присваивается cT. Переделать
				a <- Ttype{cT: ft, id: int(time.Now().Unix())} // передаем таск на выполнение
			}
		}()
	}

	// Review: переименовать
	superChan := make(chan Ttype, 10)

	go taskCreturer(superChan)

	task_worker := func(a Ttype) Ttype {
		// Review: игнорирование ошибки. Плохо
		tt, _ := time.Parse(time.RFC3339, a.cT)
		if tt.After(time.Now().Add(-20 * time.Second)) {
			a.taskRESULT = []byte("task has been successed")
		} else {
			a.taskRESULT = []byte("something went wrong")
		}
		a.fT = time.Now().Format(time.RFC3339Nano)

		// Review: вероятно это имитация обработки. Надо сохранить
		time.Sleep(time.Millisecond * 150)

		return a
	}

	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)

	// Review: убрать замыкания
	tasksorter := func(t Ttype) {
		if string(t.taskRESULT[14:]) == "successed" {
			doneTasks <- t
		} else {
			undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT, t.taskRESULT)
		}
	}

	go func() {
		// получение тасков
		for t := range superChan {
			// Review: воркеры блочат друг друга. Распараллелить
			t = task_worker(t)
			go tasksorter(t)
		}
		close(superChan)
	}()

	// Review: избавиться от этого механизма, а параллельно выполнять вывод состояния задач
	result := map[int]Ttype{}
	err := []error{}
	go func() {
		for r := range doneTasks {
			go func() {
				result[r.id] = r
			}()
		}
		for r := range undoneTasks {
			go func() {
				// Review: гонка
				err = append(err, r)
			}()
		}
		close(doneTasks)
		close(undoneTasks)
	}()

	time.Sleep(time.Second * 3)

	println("Errors:")
	for r := range err {
		println(r)
	}

	println("Done tasks:")
	for r := range result {
		println(r)
	}
}
