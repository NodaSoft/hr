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

type service struct {
	tasks          chan Task
	completedTasks []Task
	failedTasks    []error
	wg             sync.WaitGroup
}

// A Task represents a meaninglessness of our life
// UPD: лучше сделать время через UNIX timestamp, чтобы был удобный поиск в дальнейшем
type Task struct {
	id         int
	createdAt  int64 // время создания
	executedAt int64 // время выполнения
	err        error
	result     []byte
}

func NewService() *service {
	return &service{
		tasks:          make(chan Task, 10),
		completedTasks: []Task{},
		failedTasks:    []error{},
		wg:             sync.WaitGroup{},
	}
}

func (s *service) taskCreator() {
	defer s.wg.Done()
	end := time.After(time.Second * 10)
	id := 1
	for {
		select {
		case <-end:
			close(s.tasks)
			return
		default:
			now := time.Now()
			if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
				err := fmt.Errorf("error on create")
				s.tasks <- Task{
					id:        id,
					err:       err,
					createdAt: now.UnixMilli(),
				}
			} else {
				s.tasks <- Task{
					id:        id,
					createdAt: now.UnixMilli(),
				}
			}
			id++
		}
	}
}

func (s *service) taskExecutor(task *Task) {
	now := time.Now()
	if task.err == nil {
		if time.UnixMilli(task.createdAt).After(now.Add(-20 * time.Second)) {
			task.result = []byte("Task has been completed.")
		} else {
			task.err = fmt.Errorf("something went wrong")
		}
		task.executedAt = now.UnixMilli()
	}

	time.Sleep(time.Millisecond * 150)
}

func (s *service) taskWorker() {
	defer s.wg.Done()
	for task := range s.tasks {
		t := &task
		s.taskExecutor(t)
	}
}

func main() {
	fmt.Println("Start")
	ticker := time.Tick(3 * time.Second)
	svc := NewService()

	svc.wg.Add(1)
	go svc.taskCreator()

	svc.wg.Add(1)
	go svc.taskWorker()

	done := make(chan struct{})
	go func() {
		svc.wg.Wait()
		close(done)
	}()

	for {
		select {
		case <-ticker:
			fmt.Println("Errors:")
			for _, f := range svc.failedTasks {
				fmt.Println(f)
			}

			fmt.Println("Completed:")
			for _, c := range svc.completedTasks {
				fmt.Println(c)
			}

		case <-done:
			fmt.Println("All tasks processed")
			return

		default:
			task := <-svc.tasks
			if task.err != nil {
				svc.failedTasks = append(svc.failedTasks, fmt.Errorf("ERROR | Id: %d | Time: %s | Reason: %s", task.id, time.UnixMilli(task.createdAt).Format(time.RFC822), task.err.Error()))
			} else {
				svc.completedTasks = append(svc.completedTasks, task)
			}
		}
	}
}
