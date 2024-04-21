package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Приложение эмулирует получение и обработку неких тасков. Пытается и получать, и обрабатывать в многопоточном режиме.
// Должно выводить успешные таски и ошибки обработки тасков по мере выполнения.
// ЗАДАНИЕ: сделать из плохого кода хороший и рабочий - as best as you can.
// Важно сохранить логику появления ошибочных тасков.
// Сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через pull-request в github
// Как видите, никаких привязок к внешним сервисам нет - полный карт-бланш на модификацию кода.

// A Task represents the meaninglessness of our lives
type Task struct {
	id int
	// время создания
	createdAt time.Time
	// время выполнения
	finishedAt time.Time
	// результат операции
	result string
	// ошибка при выполнении
	err error
}

// Лучше вынести в отдельный класс для использования в будущем с интерфейсами и т.д.
type Producer struct{}

// Возвращает канал в который будут отправляться сгенерированные задачи
func (p Producer) Start(ctx context.Context, bufSize int) chan Task {
	// Закрыть канал следует по окончании контекста, т.к. он в отедьной горутине
	ch := make(chan Task, bufSize)

	// Отдельная горутина для удобства в использовании
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(ch)
				return
			default:
				ch <- p.createTask()
			}
		}
	}()

	return ch
}

func (Producer) createTask() Task {
	t := Task{
		createdAt: time.Now(),
		// использование наносекунд гарантирует уникальные id при условии что все они создаются в одном потоке
		id: int(time.Now().UnixNano()),
	}

	if t.createdAt.Nanosecond()%2 > 0 {
		t.err = fmt.Errorf("Some error occurred")
		t.result = "something went wrong"
	} else {
		t.result = "task has been successed"
	}

	return t
}

// Испольлзование объекта с каналами внутри позволит при необходимости в будущем легко добавить несколько получателей
type Worker struct {
	Source chan Task
	Done   chan Task
	Errs   chan error
}

// Начать произвождство задач и закрыть потом каналы
func (w *Worker) Generate() {
	defer close(w.Errs)
	defer close(w.Done)

	for t := range w.Source {
		finished := w.Process(t)
		if finished.err != nil {
			w.Errs <- finished.err
		} else {
			w.Done <- finished
		}

	}
}

// Обработку лучше вынести в отдельный метод для возможности тестирования, интерфейсов и т.д.
func (Worker) Process(task Task) Task {
	// Лучше оставить внутри функции, так как задержка в работе относится именно к реализации
	const waitTime = time.Millisecond * 150

	if task.createdAt.After(time.Now().Add(-20 * time.Second)) {
		task.result = "task has been successed"
	} else {
		// Не имеет смысла заполнять результат, так как он не должен использоваться если есть ошибка
		task.err = fmt.Errorf("something went wrong")
	}
	task.finishedAt = time.Now()

	time.Sleep(waitTime)

	return task
}

const (
	// Время работы генератора и обработчика
	processingTime = time.Second

	// Использование буферизованных каналов значительно повышает скорость работы
	taskChanSize = 10
	doneChanSize = 10
	errChanSize  = 10
)

func main() {
	// Контекст, задающий время работы / условие окончания
	processingCtx, cancel := context.WithTimeout(context.Background(), processingTime)
	defer cancel()

	producer := Producer{}
	// Сбор полчаемых задач в один канал
	newTasks := producer.Start(processingCtx, taskChanSize)

	worker := Worker{
		Source: newTasks,
		Done:   make(chan Task, doneChanSize),
		Errs:   make(chan error, errChanSize),
	}
	go worker.Generate()

	// Можно преаллоцировать место для обоих, но так как размер заранее неизвестен, а требования нет, можно оставить на нуле
	//
	// Так как используется только один поток для записи данных в хеш-таблицу, использование sync пакета не требуется:w
	result := make(map[int]Task)
	errs := make([]error, 0)

	aggregateWg := sync.WaitGroup{}

	aggregateWg.Add(1)
	go func() {
		defer aggregateWg.Done()
		for t := range worker.Done {
			result[t.id] = t
		}
	}()

	aggregateWg.Add(1)
	go func() {
		defer aggregateWg.Done()
		for err := range worker.Errs {
			errs = append(errs, err)
		}
	}()

	aggregateWg.Wait()

	fmt.Println("Errors:")
	for _, err := range errs {
		fmt.Println(err.Error())
	}

	fmt.Println("Done tasks:")
	for id := range result {
		fmt.Println(id)
	}
}
