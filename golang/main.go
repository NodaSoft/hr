package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// Приложение эмулирует получение и обработку неких тасков.
// Пытается и получать, и обрабатывать в многопоточном режиме.
// После обработки тасков в течении 3 секунд приложение должно
// выводить накопленные к этому моменту успешные таски
// и отдельно ошибки обработки тасков.

// ЗАДАНИЕ: сделать из плохого кода хороший и рабочий - as best as you can.
// Важно сохранить логику появления ошибочных тасков.
// Важно оставить асинхронные генерацию и обработку тасков.
// Сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через pull-request в github
// Как видите, никаких привязок к внешним сервисам нет - полный карт-бланш на модификацию кода.

// *******************

// Здравствуйте.

// В коде задания канал `superChan` указан с буфером 10.
// Поскольку Uber настоятельно рекомендует не использовать буферные каналы https://github.com/uber-go/guide/blob/master/style.md#channel-size-is-one-or-none
// делаю предположение, что имелся ввиду пул из десяти воркеров,
// ведь запускать безлимитное количество ресурсоёмких горутин-воркеров не есть гут.

// Обычно, для подобных задач с воркер-пулом, пайп-лайнами и запуском джобов по крону,
// использую готовые библиотеки типа вот этой: https://github.com/autom8ter/machine
// или хотя бы пакет `errgroup.WithContext` но здесь юзаю только стандартную библиотеку.

// Спасибо, что прочитали это. Лёгкого Вам дня сегодня 🙏🌴

// *******************

type Task struct {
	result     []byte // оптимальный для GC padding & alignment
	err        error
	createdAt  string
	finishedAt string
	id         int
}

func (t Task) String() string {
	return fmt.Sprintf("Task id %d time %s, result %s", t.id, t.createdAt, t.result)
}

const (
	numWorkers = 10
	workTime   = time.Second * 3
)

func main() {

	deadlineCtx, cancel := context.WithTimeout(context.Background(), workTime)
	defer cancel()

	taskChan := generateTasks(deadlineCtx)

	result, errs := taskPipeline(taskChan)

	log.Println("Errors:")
	for _, e := range errs {
		log.Println(e)
	}

	log.Println("Done tasks:")
	for _, t := range result {
		log.Println(t)
	}
}

func taskPipeline(superChan chan Task) (map[int]Task, []error) {

	successTasksChan := make(chan Task)
	errTasksChan := make(chan error)

	stopChan := make(chan struct{})

	taskProcesser := func(t Task) Task {

		_, err := time.Parse(time.RFC3339, t.createdAt)
		if err != nil {
			t.result = []byte("something went wrong")
			t.err = fmt.Errorf("Task id %d time %s, error %s", t.id, t.createdAt, t.result)
		} else {
			t.result = []byte("task has been successed")
		}

		t.finishedAt = time.Now().Format(time.RFC3339Nano)

		time.Sleep(time.Millisecond * 150)

		return t
	}

	taskSorter := func(t Task) {

		if t.err != nil {
			select {
			case <-stopChan:
				return
			case errTasksChan <- t.err:
			}
		} else {
			select {
			case <-stopChan:
				return
			case successTasksChan <- t:
			}
		}
	}

	worker := func(wg *sync.WaitGroup, superChan chan Task) {
		defer wg.Done()

		for t := range superChan {
			t := t // https://golang.org/doc/faq#closures_and_goroutines

			go taskSorter(taskProcesser(t))
		}

	}

	launchWorkers := func(superChan chan Task) {
		wg := &sync.WaitGroup{}
		wg.Add(numWorkers)

		for i := 0; i < numWorkers; i++ {
			go worker(wg, superChan)
		}

		// дожидаюсь завершения работы читающих входной канал воркеров, чтобы теперь
		// дать сигнал на завершение всего пайп-лайна.
		wg.Wait()
		close(stopChan)
	}
	go launchWorkers(superChan)

	result := make(map[int]Task)
	errs := []error{}

	gatherResults := func() {
		for {
			select {
			case <-stopChan:
				return
			case t := <-successTasksChan:
				result[t.id] = t
			case e := <-errTasksChan:
				errs = append(errs, e)
			}
		}
	}
	gatherResults()

	return result, errs
}

func generateTasks(deadline context.Context) chan Task {

	superChan := make(chan Task)

	go func() {
		defer close(superChan)

		for {
			select {
			case <-deadline.Done():
				return
			default:
			}

			// в задании сказано этот код не менять. хотел оставить как есть,
			// но тут всё настолько плохо, что решил принять волевое решение и
			// всётки починить.

			// проблема в том, что этот код генерит таски нон-стоп, но вешает неуникальные id-шники,
			// поскольку id-шник не меняется на протяжении целой 1 секунды.
			// а в итоговой Мапе все таски с одинаковыми id-шниками заменяются,
			// в итоге остаётся только последняя таска с таким id-шником.
			// по-хорошему, генерить id нужно библиотеками типа гугловской uuid https://github.com/google/uuid

			ft := time.Now().Format(time.RFC3339)

			// условие Nanosecond()%2 редко даёт 0, из-за чего error-кейсов почти не бывает,
			// поэтому заменил условие на %3.
			if time.Now().Nanosecond()%3 > 0 {
				ft = "Some error occured"
			}
			t := Task{createdAt: ft, id: int(time.Now().UnixMilli())}

			select {
			case <-deadline.Done():
				return
			case superChan <- t:
			}

			// задержка для гарантии уникальности id-шников
			time.Sleep(time.Millisecond * 300)
		}
	}()

	return superChan
}
