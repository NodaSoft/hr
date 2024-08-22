package main

import (
	"context"
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

// Мы даем тестовое задание чтобы:
// * уменьшить время технического собеседования - лучше вы потратите пару часов в спокойной домашней обстановке, чем будете волноваться, решая задачи под взором наших ребят;
// * увеличить вероятность прохождения испытательного срока - видя сразу стиль и качество кода, мы можем быть больше уверены в выборе;
// * снизить число коротких собеседований, когда мы отказываем сразу же.

// Выполнение тестового задания не гарантирует приглашение на собеседование, т.к. кроме качества выполнения тестового задания, оцениваются и другие показатели вас как кандидата.

// Мы не даем комментариев по результатам тестового задания. Если в случае отказа вам нужен наш комментарий по результатам тестового задания, то просим об этом написать вместе с откликом.

// A Sisyphus represents a meaninglessness of our life
type Sisyphus struct {
	ID                 int
	CreationTimeStamp  string // время создания
	ExecutionTimeStamp string // время выполнения
	TaskResult         []byte
}

func main() {
	taskQueue := make(chan Sisyphus, 10)
	doneTasks := make(chan Sisyphus, 10)
	undoneTasks := make(chan Sisyphus, 10)

	var wg sync.WaitGroup
	var mu sync.Mutex

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	taskCreator := func(a chan Sisyphus, ctx context.Context) {
		go func(a chan Sisyphus, ctx context.Context) {
			for {
				ft := time.Now().Format(time.RFC3339)
				if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
					ft = "Some error occured"
				}
				a <- Sisyphus{CreationTimeStamp: ft, ID: int(time.Now().Unix())} // передаем таск на выполнение

				select {
				case <-ctx.Done():
					close(a)
					return
				default:
					a <- Sisyphus{CreationTimeStamp: ft, ID: int(time.Now().Unix())} // передаем таск на выполнение
				}
			}
		}(a, ctx)
	}

	taskWorker := func(a Sisyphus) Sisyphus {
		tt, _ := time.Parse(time.RFC3339, a.CreationTimeStamp)
		if tt.After(time.Now().Add(-20 * time.Second)) {
			a.TaskResult = []byte("task has been successed")
		} else {
			a.TaskResult = []byte("something went wrong")
		}
		a.ExecutionTimeStamp = time.Now().Format(time.RFC3339Nano)

		return a
	}

	taskSorter := func(doneCh, undoneCh chan Sisyphus, t Sisyphus, mu *sync.Mutex) {
		mu.Lock()
		defer mu.Unlock()

		if strings.Contains(string(t.TaskResult), "success") {
			doneCh <- t
		} else {
			undoneCh <- t
		}
	}

	worker := func(ctx context.Context, mu *sync.Mutex) {
		for t := range taskQueue {
			processedTask := taskWorker(t)
			taskSorter(doneTasks, undoneTasks, processedTask, mu)
		}
	}

	taskPrinter := func(doneTasks, undoneTasks chan Sisyphus, ctx context.Context, wg *sync.WaitGroup) {
		defer wg.Done()

		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()

		doneComplete := false
		undoneComplete := false

		for {
			if doneComplete != false || undoneComplete != false {

			}
			select {
			case <-ctx.Done():

				close(doneTasks)
				close(undoneTasks)
				return
			case <-ticker.C:
				fmt.Println("Done Tasks:")
				for {
					select {
					case doneTask, ok := <-doneTasks:
						if !ok {
							doneComplete = true
							continue
						}
						fmt.Printf("Sisyphus ID: %d, Created: %s, Executed: %s, Result: %v\n", doneTask.ID, doneTask.CreationTimeStamp, doneTask.ExecutionTimeStamp, string(doneTask.TaskResult))
					default:
						doneComplete = true
					}
					if doneComplete == true {
						break
					}
				}

				fmt.Println("Undone Tasks:")
				for {
					select {
					case undoneTask, ok := <-undoneTasks:
						if !ok {
							undoneComplete = true
							continue
						}
						fmt.Printf("Sisyphus ID: %d, Created: %s, Executed: %s, Result: %v\n", undoneTask.ID, undoneTask.CreationTimeStamp, undoneTask.ExecutionTimeStamp, string(undoneTask.TaskResult))
					default:
						undoneComplete = true
					}
					if undoneComplete == true {
						break
					}
				}
			}
		}
	}

	wg.Add(3)
	go func() {
		taskCreator(taskQueue, ctx)
		wg.Done()
	}()

	go func() {
		worker(ctx, &mu)
		wg.Done()
	}()

	go func() {
		taskPrinter(doneTasks, undoneTasks, ctx, &wg)
		wg.Done()
	}()

	wg.Wait()
}
