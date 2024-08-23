package main

import (
	"fmt"
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

// A Task represents a meaninglessness of our life
type Task struct {
	id       int
	created  string // время создания
	finished string // время выполнения
	result   []byte
}

func taskCreator(a chan Task) {
	startTime := time.Now()
	defer close(a)

	for time.Since(startTime) < 10*time.Second {
		now := time.Now()
		ft := now.Format(time.RFC3339)
		// здесь nano%200! потому что у меня все наносекунды *00, то есть четные
		if now.Nanosecond()%200 > 0 { // вот такое условие появления ошибочных тасков
			ft = "Some error occured"
		}
		a <- Task{created: ft, id: int(now.Unix())} // передаем таск на выполнение
		time.Sleep(100 * time.Millisecond)
	}
}

func taskWorker(a Task) Task {
	tt, _ := time.Parse(time.RFC3339, a.created)
	if tt.After(time.Now().Add(-20 * time.Second)) {
		a.result = []byte("task has been successed")
	} else {
		a.result = []byte("something went wrong")
	}
	a.finished = time.Now().Format(time.RFC3339Nano)
	time.Sleep(150 * time.Millisecond)
	return a
}

func main() {
	superChan := make(chan Task, 10)
	doneTasks := make(chan Task)
	undoneTasks := make(chan error)
	finished := make(chan struct{}) // после всех тасков

	var wg sync.WaitGroup
	var mu sync.Mutex

	go taskCreator(superChan)

	taskSorter := func(t Task) {
		t = taskWorker(t)
		defer wg.Done()

		if string(t.result[14:]) == "successed" {
			doneTasks <- t
		} else {
			undoneTasks <- fmt.Errorf("task id %d time %s, error %s", t.id, t.created, t.result)
		}
	}

	go func() {
		// получение тасков
		for t := range superChan {
			wg.Add(1)
			go taskSorter(t)
		}
	}()

	// слушаем выполненные
	result := map[int]Task{}
	go func() {
		for r := range doneTasks {
			mu.Lock()
			result[r.id] = r
			mu.Unlock()
		}
	}()

	// слушаем ошибки
	err := []error{}
	go func() {
		for r := range undoneTasks {
			mu.Lock()
			err = append(err, r)
			mu.Unlock()
		}
	}()

	printLog := func() {
		mu.Lock()
		fmt.Println("Errors:")
		for _, r := range err {
			fmt.Println(r)
		}
		fmt.Println("Done tasks:")
		for _, r := range result {
			fmt.Println(r)
		}
		result = map[int]Task{}
		err = []error{}
		mu.Unlock()
	}

	// периодический вывод результатов 3сек и по завершению
	ticker3 := time.NewTicker(3 * time.Second)
	defer ticker3.Stop()
	go func() {
		for {
			select {
			case <-ticker3.C:
				printLog()
			case <-finished:
				printLog()
				return
			}
		}
	}()

	time.Sleep(5 * time.Second) // старт генератора

	wg.Wait()
	close(doneTasks)
	close(undoneTasks)
	close(finished)
}
