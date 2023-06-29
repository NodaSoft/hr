package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнены остальных тасков

// Task - название нашей таски
type Task struct {
	id           int64
	creationTime time.Time // время создания
	finalTime    time.Time // время выполнения
	taskRESULT   []byte
	err          error
}

func createTasks(a chan Task) {
	startTime := time.Now()

	var wgCreator sync.WaitGroup
	stopChan := make(chan bool)

	wgCreator.Add(1)
	go func(wg *sync.WaitGroup) {
		cnt := 0

		for {
			isBreak := false
			select {
			case stop, ok := <-stopChan:
				if ok && stop {
					isBreak = true
					break
				}
			default:
				finalTime := time.Now()

				nsecs := finalTime.Sub(startTime).Nanoseconds()

				newTask := Task{creationTime: finalTime, id: time.Now().UnixNano() + nsecs}

				// Костыль, так как мой компьютер (Macbook Pro 2019) округляет наносекунды до микросекунд при вызове UnixNano или time.Now().Nanosecond()
				if nsecs%2 == 1 { // вот такое условие появления ошибочных тасков
					newTask.err = errors.New("Some error occured")
				}
				a <- newTask // передаем таск на выполнение
				cnt += 1
			}

			if isBreak {
				break
			}
		}

		fmt.Printf("Total generated tasks %d\n", cnt)
		wg.Done()
	}(&wgCreator)

	time.Sleep(200 * time.Millisecond)
	// time.Sleep(3 * time.Second)

	stopChan <- true
	wgCreator.Wait()

}

func (t *Task) processTask() {
	tt := t.creationTime

	// if tt.After(time.Now().Add(-20 * time.Second)) {
	if time.Now().Sub(tt) < 20*time.Second && t.err == nil {
		t.taskRESULT = []byte("task has been successed")
	} else {
		t.taskRESULT = []byte("something went wrong")
	}
	t.finalTime = time.Now()

	// В исходном коде стоял sleep и при ошибке, и при успешной обработке. Так и надо было оставить, верно?
	time.Sleep(time.Millisecond * 150)
}

func main() {
	superChan := make(chan Task, 10) // Super Tengen Toppa Guren Lagann? :) Оставляю изначальное название, по смыслу подходит

	doneTasks := make(chan Task, 10)
	undoneTasks := make(chan Task, 10)

	tasksorter := func(t Task) {
		if string(t.taskRESULT[14:]) == "successed" {
			doneTasks <- t
		} else {
			undoneTasks <- t
		}
	}

	// Почти те же мапы, только добавили мьютексов
	var taskIDMutex sync.RWMutex
	taskIDs := map[int64]Task{}

	var errorsMutex sync.RWMutex
	errorTaskIDs := map[int64]Task{} // Хочется хранить таски, а не строки

	workers := 10
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)

		go func(wg *sync.WaitGroup, tasksChan chan Task, errorsChan chan Task) {
			ch1Broken := false
			ch2Broken := false

			for {
				select { // Одинаково обрабатываем оба канала, но пишем в разные мапы
				case task, ok := <-tasksChan:
					if !ok {
						// Уровень вложенности 6 табов, оставляю из вредности к исходному коду, который надо было исправить
						ch1Broken = true
						break
					}
					taskIDMutex.Lock()
					taskIDs[task.id] = task
					taskIDMutex.Unlock()

				case task, ok := <-errorsChan:
					if !ok {
						ch2Broken = true // Уровень вложенности по прежнему 6 табов
						break
					}
					errorsMutex.Lock()
					errorTaskIDs[task.id] = task
					errorsMutex.Unlock()
				}

				if ch1Broken && ch2Broken {
					break
				}

			}
			wg.Done()

		}(&wg, doneTasks, undoneTasks)
	}

	workersSuper := 10
	var wgSuper sync.WaitGroup
	for i := 0; i < workersSuper; i++ {
		wgSuper.Add(1)

		go func(wgSuper *sync.WaitGroup, ch chan Task) {
			cnt := 0

			// получение тасков
			for t := range ch {
				t.processTask()
				tasksorter(t)

				cnt += 1
			}

			wgSuper.Done()
		}(&wgSuper, superChan)
	}

	createTasks(superChan)
	close(superChan)

	wgSuper.Wait() // Значит, воркеры, заполняющие doneTasks и undoneTasks отработали и эти каналы можно закрывать

	close(doneTasks)
	close(undoneTasks)

	wg.Wait() // Ожидаем выполнение тасков (sleep на)

	println("Errors:")
	for _, t := range errorTaskIDs {
		fmt.Printf("Task id %d time %s, error %s\n", t.id, t.creationTime.Format(time.RFC3339Nano), string(t.taskRESULT))
	}

	println("Done tasks:")
	for _, t := range taskIDs {
		fmt.Printf("Task id %d time %s, SUCCESS %s\n", t.id, t.creationTime.Format(time.RFC3339Nano), string(t.taskRESULT))
	}

	// Проверка, что количество изначально сгенерированных и обработанных тасков совпадает
	fmt.Printf("Count of errors %d\nCount of success %d\nSum %d\n", len(errorTaskIDs), len(taskIDs), len(errorTaskIDs)+len(taskIDs))
}
