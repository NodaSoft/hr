package main

import (
	"context"
	"fmt"
	"runtime"
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

type someTask struct {
	id         int
	createTime time.Time
	finishTime time.Time
	result     []byte
	isSuccess  bool
}

func taskGenerator(taskChan chan<- someTask, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			close(taskChan)
			return
		default:
			createTime := time.Now()
			taskChan <- someTask{
				id:         int(createTime.Unix()),
				createTime: createTime,
			}
		}

	}
}

func taskWorker(taskChan <-chan someTask, resultChan chan<- someTask, ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range taskChan {
		if task.createTime.Nanosecond()%2 == 0 {
			task.result = []byte(fmt.Sprintf("Task id %d time %s, success", task.id, task.createTime.Format(time.RFC3339)))
			task.isSuccess = true
		} else {
			task.result = []byte(fmt.Sprintf("Task id %d time %s, error %s", task.id, task.createTime.Format(time.RFC3339), task.result))
			task.isSuccess = false
		}
		time.Sleep(500 * time.Millisecond)
		task.finishTime = time.Now()
		select {
		case <-ctx.Done():
			return
		case resultChan <- task:
		}
	}

}

func main() {
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	createdTaskChan := make(chan someTask, 20)
	resultChan := make(chan someTask, 20)

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go taskWorker(createdTaskChan, resultChan, ctx, &wg)
	}

	go taskGenerator(createdTaskChan, ctx)

	successTasks := make([]someTask, 0, 20)
	failedTasks := make([]someTask, 0, 20)

	printTicker := time.NewTicker(3 * time.Second)
	for {
		select {
		case <-printTicker.C:
			fmt.Println("Errors:")
			for _, task := range failedTasks {
				fmt.Println(string(task.result))
			}
			fmt.Println("Success:")
			for _, task := range successTasks {
				fmt.Println(string(task.result))
			}
		case <-ctx.Done():
			wg.Wait()
			return
		default:
			task := <-resultChan
			if task.isSuccess {
				successTasks = append(successTasks, task)
			} else if !task.isSuccess {
				failedTasks = append(failedTasks, task)
			}
		}
	}
}
