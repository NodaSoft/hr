package main

import (
	"fmt"
	"slices"
	"sync"
	"time"
)

const (
	TASK_BUFFER_SIZE = 10
)

func getSuccessResult() []byte {
	return []byte("task has been successed")
}

func getFailureResult() []byte {
	return []byte("something went wrong")
}

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

// Во-первых, необходимо понять, допустимо ли модифицировать реализацию
// издателя. Из требования сохранения логики, заложенной в алгоритм генерации
// тасок, можно сделать вывод, что модификация издателя опасна, поскольку мы не
// понимаем, как он работает
func taskCreturer(a chan Ttype) {
	go func() {
		for {
			ft := time.Now().Format(time.RFC3339)
			if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
				ft = "Some error occured"
			}
			a <- Ttype{cT: ft, id: int(time.Now().Unix())} // передаем таск на выполнение
		}
	}()
}

// Что определяет, является ли таска ошибочной? Если код работал до нас и
// документация к нему отсутствует, то ошибочной таской является та таска,
// которую исходный алгоритм определяет как ошибочную.
//
// В данном случае ошибочными являются таски, у которых в поле ft записана
// строка "Some error occured", постольку и только постольку, что для них
// выражение tt.After(time.Now().Add(-20 * time.Second)) всегда будет ложным,
// ведь при ошибке парсинга функция time.Parse() возвращает пустой объект,
// интерпретируемый как 12 ночи 1-го января 1-го года от рождества Христова и,
// очевидно, данная метка времени не может быть позже какой-то другой метки
// времени. Также являются ошибочными таски, созданные более 20 секунд назад.
// Проверку на ошибку при парсинге, однако, следует добавить исключительно в
// качестве проверки предусловий
//
// Вместе с тем, то, что время выполения форматируется отлично от времени
// создания является очевидной ошибкой на этапе разработки, поскольку клиент не
// ожидает, что в одном и том же типе разные поля со временем будут использовать
// разное форматирование. Поскольку в существующем клиентском коде никто не
// использует поле ft и, соответственно, не успел завязаться на данное ошибочное
// поведение, то оно может и должно быть исправлено
func processTask(task Ttype) (Ttype, error) {
	if task.cT == "Some error occured" {
		task.taskRESULT = getFailureResult()
	} else {
		var (
			err         error
			taskCreated time.Time
		)

		if taskCreated, err = time.Parse(time.RFC3339, task.cT); err != nil {
			return task, fmt.Errorf("error parsing creation time: %s", err)
		}
		if time.Now().Sub(taskCreated) < (20 * time.Second) {
			task.taskRESULT = getSuccessResult()
		} else {
			task.taskRESULT = getFailureResult()
		}
	}
	task.fT = time.Now().Format(time.RFC3339)
	time.Sleep(time.Millisecond * 150)
	return task, nil
}

type ChanSubscription[T any] struct {
	SuccessChan chan T
	FailChan    chan error
	unsub       chan struct{}
}

func (s *ChanSubscription[T]) Unsubsribe() {
	s.unsub <- struct{}{}
}

type ChanEvent[T any] struct {
	Object T
	sub    *ChanSubscription[T]
}

func (e *ChanEvent[T]) Success(t T) {
	e.sub.SuccessChan <- t
}

func (e *ChanEvent[T]) Fail(err error) {
	e.sub.FailChan <- err
}

func SubscribeChan[T any](objects chan T, handler func(*ChanEvent[T])) *ChanSubscription[T] {
	sub := &ChanSubscription[T]{
		SuccessChan: make(chan T),
		FailChan:    make(chan error),
		unsub:       make(chan struct{}),
	}
	go func() {
		for {
			select {
			case obj := <-objects:
				event := &ChanEvent[T]{
					Object: obj,
					sub:    sub,
				}
				handler(event)
			case <-sub.unsub:
				close(sub.SuccessChan)
				close(sub.FailChan)
				return
			}
		}
	}()

	return sub
}

// Неизвестно, можно ли менять сигнатуру функции, поэтому при необходимости
// использования обработчика с другой сигнатурой разумно реализовать обёртку
//
// Откуда берется константа 14? Она берется из соображения, что в одном
// используемом издателем тексте ошибки присутствует слово "successed". Не
// существует никакой конвенции относительно того, что статус должен
// содержаться в оставшихся после 14-го байта тексте ошибки и точно так же не
// существует возможности, что сравнение полной строки приведет к какой-либо
// неоднозначности. Именно поэтому логично сравнивать полный текст результата
// с константой
func taskHandler(e *ChanEvent[Ttype]) {
	var (
		err           error
		processedTask Ttype
	)
	if processedTask, err = processTask(e.Object); err != nil {
		e.Fail(fmt.Errorf("error processing task: %s", err))
		return
	}

	if !slices.Equal(processedTask.taskRESULT, getSuccessResult()) {
		t := processedTask
		e.Fail(fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT, t.taskRESULT))
		return
	}

	e.Success(processedTask)
	return
}

func main() {
	taskChan := make(chan Ttype, TASK_BUFFER_SIZE)
	taskCreturer(taskChan)
	creatorSub := SubscribeChan(taskChan, taskHandler)

	result := make([]Ttype, 0)
	errors := make([]error, 0)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for r := range creatorSub.SuccessChan {
			result = append(result, r)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for err := range creatorSub.FailChan {
			errors = append(errors, err)
		}
	}()

	// Мы ничего не знаем об издателе, поэтому сами решаем в какой момент
	// перестать обрабатывать события
	time.Sleep(3 * time.Second)
	creatorSub.Unsubsribe()
	wg.Wait()

	fmt.Println("Errors:")
	for _, e := range errors {
		fmt.Println(e)
	}

	fmt.Println("Done tasks:")
	for _, r := range result {
		fmt.Println(r.id)
	}
}
