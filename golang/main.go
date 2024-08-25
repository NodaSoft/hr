package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
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
	id         uint64
	cT         string // время создания
	fT         string // время выполнения
	taskRESULT []byte
}

var someError error = errors.New("Some error occurred")

func main() {
	taskCreator := func(a chan<- Ttype, ctxCreator context.Context) {
		var nextIdx uint64 = 0
		for {
			ft := time.Now().Format(time.RFC3339)

			// time.Now().Nanosecond() всегда заканчиается на много нулей, поэтому никогда не будет появляться ошибок
			// if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
			// 	ft = "Some error occured"
			// }

			// Условие появления ошибочных тасков с использованием рандома
			if rand.Intn(2) == 0 {
				ft = "Some error occurred"
			}

			nextIdx++
			select {
			case <-ctxCreator.Done():
				close(a)
				return
			case a <- Ttype{cT: ft, id: nextIdx}: // передаем таск на выполнение
				//
			}
		}
	}

	taskWorker := func(a Ttype) (*Ttype, error) {
		taskFinished, err := time.Parse(time.RFC3339, a.cT)
		if err == nil {
			if taskFinished.After(time.Now().Add(-20 * time.Second)) {
				a.taskRESULT = []byte("task has been successed")
			} else {
				a.taskRESULT = []byte("something went wrong")
				err = someError
			}
		} else {
			a.taskRESULT = []byte("parsing error")
			err = someError
		}
		a.fT = time.Now().Format(time.RFC3339Nano)
		time.Sleep(time.Millisecond * 150)
		return &a, err
	}

	taskSorter := func(t *Ttype, err error, doneTasks chan<- *Ttype, undoneTasks chan<- error) {
		if err == nil {
			doneTasks <- t
		} else {
			undoneTasks <- fmt.Errorf("task id %d time %s, error %s: %w", t.id, t.cT, string(t.taskRESULT), err)
		}
	}

	superChan := make(chan Ttype, 1000)
	// в течение 1 секунды будут создаваться таски
	ctxCreator, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	defer cancel()
	go taskCreator(superChan, ctxCreator)

	doneTasks := make(chan *Ttype)
	undoneTasks := make(chan error)

	go func() {
		// получение тасков
		var wg sync.WaitGroup
		for task := range superChan {
			wg.Add(1)
			go func(task Ttype) {
				defer wg.Done()
				t, err := taskWorker(task)
				taskSorter(t, err, doneTasks, undoneTasks)
			}(task)
		}
		wg.Wait()
		close(doneTasks)
		close(undoneTasks)
	}()

	result := make(map[uint64]*Ttype)
	errors := []error{}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for res := range doneTasks {
			result[res.id] = res
		}
	}()

	go func() {
		defer wg.Done()
		for err := range undoneTasks {
			errors = append(errors, err)
		}
	}()

	wg.Wait()

	fmt.Println("Errors:")
	for _, err := range errors {
		fmt.Printf("error: %v\n", err)
	}

	fmt.Println("Done tasks:")
	for id, res := range result {
		fmt.Printf("res-%d: %s\n", id, res.taskRESULT)
	}

	fmt.Printf("len err = %d, len res = %d\n", len(errors), len(result))
}
