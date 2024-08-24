package main

import (
	"context"
	"errors"
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

const (
	workingTime = 10
	tickTime    = 3
)

type Task struct {
	id int

	creationTime time.Time
	elapsedTime  time.Duration

	error error
}

type Hub struct {
	taskChan chan *Task

	successTasks map[int]struct{}
	failTasks    map[error]struct{}

	wg sync.WaitGroup
	mu sync.Mutex
}

func newTaskHub() *Hub {
	return &Hub{
		taskChan:     make(chan *Task),
		successTasks: make(map[int]struct{}),
		failTasks:    map[error]struct{}{},
	}
}

func (h *Hub) generateTasks(ctx context.Context) {
	defer h.wg.Done()
	defer close(h.taskChan)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			task := &Task{
				id: int(time.Now().Unix()),
			}

			task.creationTime = time.Now()

			if time.Now().Nanosecond()%2 > 0 {
				task.error = errors.New("some error occurred")
			}

			h.taskChan <- task
		}
	}
}

func (h *Hub) startLog(ctx context.Context) {
	ticker := time.NewTicker(tickTime * time.Second)

	defer h.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			h.mu.Lock()

			fmt.Println("=========== Show statistics ===========")

			fmt.Println("Errors:")
			for err, _ := range h.failTasks {
				fmt.Println(err)
			}

			fmt.Println("Done tasks ids:")
			for id, _ := range h.successTasks {
				fmt.Println(id)
			}

			h.failTasks = make(map[error]struct{})
			h.successTasks = make(map[int]struct{})

			h.mu.Unlock()
		}
	}
}

func processTask(task *Task) (bool, error) {
	if task.error != nil {
		return false, fmt.Errorf("processTask: %w", task.error)
	}

	if !task.creationTime.After(time.Now().Add(-20 * time.Second)) {
		return false, errors.New("something went wrong")

	}

	return true, nil
}

func (h *Hub) runWorker() {
	//Не добавлял ограничений по кол-ву воркеров, но если бы нужно было:
	//Или запустить n воркеров и передать в них контекст
	//Или сделать буферизированный канал и читать/писать в него, ограничивая кол-во возможных воркеров.
	for t := range h.taskChan {
		h.wg.Add(1)

		go func() {
			defer h.wg.Done()

			success, err := processTask(t)
			if err != nil {
				h.mu.Lock()
				defer h.mu.Unlock()

				h.failTasks[fmt.Errorf("task id: %d, time: %s, processTaskerror: %w", t.id, t.creationTime.Format(time.RFC3339), err)] = struct{}{}

				return
			}

			if success {
				h.mu.Lock()
				defer h.mu.Unlock()

				h.successTasks[t.id] = struct{}{}
			}
		}()

		time.Sleep(time.Millisecond * 150)
	}
}

func main() {
	fmt.Println("Starting work...")

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(workingTime*time.Second))
	defer cancel()

	hub := newTaskHub()

	hub.wg.Add(1)
	go hub.generateTasks(ctx)

	hub.wg.Add(1)
	go hub.startLog(ctx)

	go hub.runWorker()

	hub.wg.Wait()

	fmt.Println("Work Done!")
}
