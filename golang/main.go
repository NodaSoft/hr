package main

import (
	"context"
	"fmt"
	"log"
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

// Тип, представляющий таски. В зависимости от назначения тасков может понадобиться добавить поле TaskReturn, 
// для хранения результата успешного выполнения таска. Если результаты слишком разнообразны, можно сделать структуру дженериком
type Task struct {
	Id         int
	StartTime  string 
	EndTime    string 
	TaskFailed bool
	TaskError  error
}

// Создаёт таски. Прекращает когда ctx отменён
func taskCreator(ctx context.Context) (createdTasks <- chan Task) {
	tc := make(chan Task, 10)
	done := ctx.Done()
	go func() {
		// Можно заместо автоинкремента использовать UUID(есть пакет от гугла)
		id := 0 
	loop:
		for {
			select {
			case <- done:
				log.Println("Creation Ended")
				break loop
			default:
				now := time.Now()
				ct := now.Format(time.RFC3339)
				if now.Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
					ct = "Some error occured"
				}
				id++
				tc <- Task{StartTime: ct, Id: id} // передаем таск на выполнение
			}
		}
		close(tc)
	}()
	return tc
}

// Выполняет один таск
func runTask(t Task) (completed Task) {
	tt, err := time.Parse(time.RFC3339, t.StartTime)
	if err != nil {
		t.TaskFailed = true
		t.TaskError = err
		return t
	}
	if tt.After(time.Now().Add(-20 * time.Second)) {
		t.TaskFailed = false
	} else {
		t.TaskFailed = true
		t.TaskError = fmt.Errorf("Timed out")
	}
	t.EndTime = time.Now().Format(time.RFC3339Nano)

	time.Sleep(time.Millisecond * 150)

	return t
}

// Запускает все таски из канала в многопоточном режиме.
// Автоматически закрывает канал, который возращает
func taskRunner(tc <-chan Task) (completed chan Task) {
	completed = make(chan Task, 10)
	go func() {
		wg := sync.WaitGroup{}
		for task := range tc {
			wg.Add(1)
			go func(task Task) {
				completed <- runTask(task)
				wg.Done()
			}(task)
		}
		wg.Wait()
		close(completed)
	}()
	return
}

// Функция выводит сперва таски с ошибками исполнения, а потом успешно завершённые. 
// Может занимать значительное время(на моей машине за 3 секундами исполнения следуют 6 секунд вывода без tmux и 12-13 с ним)
func taskPrinter(tc chan Task) {
	result := map[int]Task{}
	errored := map[int]Task{}
	for t := range tc {
		if !t.TaskFailed {
			result[t.Id] = t
		} else {
			errored[t.Id] = t
		}
	}
	log.Println("Errors:")
	for id, task := range errored {
		log.Printf("\t%d: %s", id, task.TaskError.Error())
	}
	log.Println("Successful:")
	for id := range result {
		log.Printf("\t%d", id)
	}
}

func main() {
	log.SetFlags(0)

	creatorCtx, cancelCreator := context.WithTimeout(context.Background(), time.Second * 3)
	defer cancelCreator()

	taskChan := taskCreator(creatorCtx)

	doneChan := taskRunner(taskChan)

	taskPrinter(doneChan) // Это не нужно выполнять асинхронно т.к. вывод в консоль в любом случае является синхронной операцией
}
