package main

import (
	"fmt"
	"strings"
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

// A Ttype represents a meaninglessness of our life
type Ttype struct {
	id         int
	cT         string // время создания
	fT         string // время выполнения
	taskRESULT string // заменил на обычную строку, чтобы в дальнейшем просто смотреть на наличие суффикса
}

// TaskKiller - обертка для superChannel, нужна для исключения случая записи в закрытый канал
// Закрытие superChannel сигнализирует о необходимости в закрытии других каналов и о завершении работы воркеров
type TaskKiller struct {
	superChannel chan Ttype
	closed       bool
}

func main() {
	var wg sync.WaitGroup // добавил waitGroup для гарантии завершения работы воркеров
	wg.Add(5)
	taskKiller := TaskKiller{
		superChannel: make(chan Ttype, 10),
	}

	taskCreator := func(taskKiller *TaskKiller, wg *sync.WaitGroup) {
		for {
			if taskKiller.closed {
				break
			}
			ft := time.Now().Format(time.RFC3339)
			if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
				ft = "Some error occured"
			}
			taskKiller.superChannel <- Ttype{cT: ft, id: int(time.Now().Unix())} // передаем таск на выполнение
			time.Sleep(1 * time.Second)
		}
		close(taskKiller.superChannel)
		println("Creator is done")
		wg.Done()
	}

	go taskCreator(&taskKiller, &wg)
	go func(taskKiller *TaskKiller) {
		<-time.After(10 * time.Second)
		taskKiller.closed = true
	}(&taskKiller)

	taskWorker := func(task Ttype) Ttype {
		tt, _ := time.Parse(time.RFC3339, task.cT)
		if tt.After(time.Now().Add(-20 * time.Second)) {
			task.taskRESULT = "task has been successed"
		} else {
			task.taskRESULT = "something went wrong"
		}
		task.fT = time.Now().Format(time.RFC3339Nano)

		time.Sleep(time.Millisecond * 150)

		return task
	}

	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)

	taskSorter := func(task Ttype) {
		if strings.HasSuffix(task.taskRESULT, "successed") {
			doneTasks <- task
		} else {
			undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", task.id, task.cT, task.taskRESULT)
		}
		wg.Done()
	}

	go func() {
		// получение тасков
		for t := range taskKiller.superChannel {
			t = taskWorker(t)
			go taskSorter(t)
			wg.Add(1)
		}
		println("Unnamed worker 1 is done")
		wg.Done()
		close(doneTasks)
		close(undoneTasks)
	}()

	var (
		result = make(map[int]Ttype)
		err    []error

		resultMutex sync.Mutex // добавил мютексы для потокобезапасной записи в мапу и слайс. использовал обычный так как решил,
		errMutex    sync.Mutex // что имеет смысл чистить контейнеры после сообщения об их отработанных тасках
	)
	go func() {
		for r := range doneTasks {
			resultMutex.Lock()
			result[r.id] = r
			resultMutex.Unlock()
		}
		println("Unnamed worker 2 is done")
		wg.Done()
	}()

	go func() {
		for r := range undoneTasks {
			errMutex.Lock()
			err = append(err, r)
			errMutex.Unlock()
		}
		println("Unnamed worker 3 is done")
		wg.Done()
	}()

	// оберул в отдельную рутину
	go func() {
		for {
			time.Sleep(time.Second * 3)
			if taskKiller.closed {
				break
			}

			println("Errors:")
			errMutex.Lock()
			for r := range err {
				println(r)
			}
			err = err[:0]
			errMutex.Unlock()

			println("Done tasks:")
			resultMutex.Lock()
			for r := range result {
				println(r)
			}
			clear(result)
			resultMutex.Unlock()
		}
		println("Unnamed worker 4 is done")
		wg.Done()
	}()

	wg.Wait()
	println("main is done")
}
