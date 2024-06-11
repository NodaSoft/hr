package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

// Приложение эмулирует получение и обработку неких тасков. Пытается и получать, и обрабатывать в многопоточном режиме.
// Приложение должно генерировать таски 10 сек.
// Каждые 3 секунды должно выводить в консоль результат всех обработанных к этому моменту тасков (отдельно успешные и отдельно с ошибками).

// ЗАДАНИЕ: сделать из плохого кода хороший и рабочий - as best as you can.
// Важно сохранить логику появления ошибочных тасков.
// Важно оставить асинхронные генерацию и обработку тасков.
// Сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через pull-request в github
// Как видите, никаких привязок к внешним сервисам нет - полный карт-бланш на модификацию кода.

// A Ttype represents a meaninglessness of our life
type Task struct {
	Id         int64
	CreatedAt  time.Time
	ExecutedAt time.Time
	Result     []byte
	Error      error
}

func NewTask() Task {
	t := Task{
		Id:        time.Now().Unix(),
		CreatedAt: time.Now(),
		Result:    make([]byte, 0),
	}
	if time.Now().Nanosecond()%2 > 0 {
		t.Error = errors.New("Some error occured")
	}
	return t
}

type ErrTasks struct {
	mu       sync.Mutex
	BadTasks []error
}

func (et *ErrTasks) Add(err error) {
	et.mu.Lock()
	defer et.mu.Unlock()
	et.BadTasks = append(et.BadTasks, err)
}

func (et *ErrTasks) Results() string {
	et.mu.Lock()
	results := strings.Builder{}
	for _, er := range et.BadTasks {
		results.WriteString(er.Error())
		results.WriteByte('\n')
	}
	et.BadTasks = et.BadTasks[:0]
	et.mu.Unlock()
	return results.String()
}

func TaskGenerator(ctx context.Context, generatedTasks chan Task) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(generatedTasks)
				return
			default:
				generatedTasks <- NewTask()
			}
		}
	}()
}

func TaskWorker(tasks chan Task, result chan Task) {
	for t := range tasks {
		// elapsedTime := time.Since(t.CreatedAt)
		// elapsedTime > time.Duration(time.Second*20)
		if t.CreatedAt.After(time.Now().Add(-20 * time.Second)) {
			t.Result = []byte("task has been successed")
		} else {
			t.Error = errors.New("something went wrong")
		}
		t.ExecutedAt = time.Now()
		result <- t
	}
	close(result)
}

func TaskSorter(allTasks chan Task, completedTasks chan Task, badTasks chan error) {
	for t := range allTasks {
		if t.Error != nil {
			badTasks <- fmt.Errorf("Task id %d time %s, error %s", t.Id, t.CreatedAt, t.Error)
		} else {
			completedTasks <- t
		}
	}
}

func mergeResults(results []chan Task) chan Task {
	mergedResults := make(chan Task, 10)

	wg := &sync.WaitGroup{}
	for idx := range results {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for task := range results[i] {
				mergedResults <- task
			}
		}(idx)
	}

	return mergedResults
}

func main() {
	generatedTasks := make(chan Task, 10)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*10))
	defer cancel()
	go TaskGenerator(ctx, generatedTasks)

	results := make([]chan Task, 4)
	for i := 0; i < 4; i++ {
		workerResult := make(chan Task)
		results[i] = workerResult
		go TaskWorker(generatedTasks, workerResult)
	}

	mergedResults := mergeResults(results)
	completedTasks := make(chan Task)
	badTasks := make(chan error)
	go TaskSorter(mergedResults, completedTasks, badTasks)

	tasksResults := sync.Map{}
	go func() {
		for t := range completedTasks {
			tasksResults.Store(t.Id, t)
		}
	}()
	taskErrors := ErrTasks{}
	go func() {
		for err := range badTasks {
			taskErrors.Add(err)
		}
	}()

	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("Done tasks:")
			tasksResults.Range(func(key, value any) bool {
				fmt.Printf("%v %v\n", key, value)
				return true
			})
			fmt.Println("Errors:")
			fmt.Println(taskErrors.Results())
		}
	}
}
