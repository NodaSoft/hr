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

type Task struct {
	ID         int
	CreateTime time.Time // время создания задачи
	FinishTime time.Time // время завершения выполнения задачи
	Result     string    // результат выполнения задачи
	Err        error     // ошибка, если таковая возникла при выполнении задачи
}

// taskCreator генерирует задачи и отправляет их в канал taskChan.
// когда приходит время таймаута завершаем создание тасок
func taskCreator(ctx context.Context, taskChan chan<- Task) {
	defer close(taskChan)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// Создаем новую задачу с уникальным ID и текущим временем создания.
			task := Task{
				ID:         int(time.Now().UnixNano()),
				CreateTime: time.Now(),
			}

			if task.CreateTime.Nanosecond()%2 != 0 {
				task.Err = fmt.Errorf("Some error occurred")
				task.Result = "something went wrong"
			} else {
				task.Result = "task has been successed"
			}

			taskChan <- task
		}
	}

}

// taskWorker обрабатывает задачи из канала taskChan и отправляет результаты в канал doneTaskChan.
func taskWorker(task Task, doneTaskChan chan<- Task, errChan chan<- error) {
	// Эмулируем выполнение задачи и устанавливаем время завершения.
	if task.CreateTime.After(time.Now().Add(-20 * time.Second)) {
		task.Result = "task has been successed"
	} else {
		task.Result = "something went wrong"
		task.Err = fmt.Errorf("something went wrong")
	}
	task.FinishTime = time.Now()

	// Эмулируем время обработки задачи.
	time.Sleep(time.Millisecond * 150)

	doneTaskChan <- task
	if task.Err != nil {
		errChan <- fmt.Errorf("Task id %d time %s, error %s", task.ID, task.CreateTime, task.Result)
	}
}

func main() {
	// Создаем контекст с таймаутом 10 секунд.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Создаем каналы для передачи задач, завершенных задач и ошибок.
	taskChan := make(chan Task, 10)
	doneTaskChan := make(chan Task, 10)
	errChan := make(chan error, 10)

	var wg sync.WaitGroup

	// Запускаем горутину для генерации задач.
	go taskCreator(ctx, taskChan)

	// Запускаем несколько горутин для обработки задач.
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for t := range taskChan {
				taskWorker(t, doneTaskChan, errChan)
			}
		}()
	}

	// Ожидаем завершения всех горутин обработки задач.
	go func() {
		wg.Wait()
		close(doneTaskChan)
		close(errChan)
	}()

	// Создаем WaitGroup для отслеживания завершения вывода результатов.
	var wgPrint sync.WaitGroup
	wgPrint.Add(1)
	go func() {
		defer wgPrint.Done()
		for {
			select {
			case r, ok := <-doneTaskChan:
				if !ok {
					return
				}
				fmt.Printf("ID: %d, Create Time: %s, Finish Time: %s, Result: %s\n", r.ID, r.CreateTime.Format("2006-01-02 15:04:05"), r.FinishTime.Format("2006-01-02 15:04:05"), r.Result)
			case err, ok := <-errChan:
				if !ok {
					return
				}
				fmt.Println("Error:", err)
			}
		}
	}()

	wgPrint.Wait()

}
