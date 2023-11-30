package main

import (
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

// A Ttype represents a meaninglessness of our life
type Task struct {
	id         int
	createTime string // время создания
	finishTime string // время выполнения
	result     []byte
}

func (task Task) String() string {
	return fmt.Sprintf("id: %d\ncreateTime: %s\nfinishTime: %s\nresult: %s", task.id, task.createTime, task.finishTime, string(task.result))
}

func TaskCreater() <-chan Task {

	taskChannel := make(chan Task, 10)

	go func() {
		i := 0
		for {
			createTime := time.Now().Format(time.RFC3339)

			// вот такое условие появления ошибочных тасков
			if time.Now().Nanosecond()%2 > 0 {
				createTime = "Some error occured"
				close(taskChannel)
				break
			}

			// передаем таск на выполнение
			taskChannel <- Task{createTime: createTime, id: int(time.Now().Unix())}
			i++
		}
	}()

	return taskChannel
}

func TaskWorker(task Task) Task {

	createTime, _ := time.Parse(time.RFC3339, task.createTime)

	if createTime.After(time.Now().Add(-20 * time.Second)) {
		task.result = []byte("task has been successed")
	} else {
		task.result = []byte("something went wrong")
	}

	task.finishTime = time.Now().Format(time.RFC3339Nano)

	time.Sleep(time.Millisecond * 150)

	return task
}

func TaskSorter(taskChan <-chan Task) (<-chan Task, <-chan error) {

	done := make(chan Task)
	undone := make(chan error)

	go func() {
		for task := range taskChan {
			workedTask := TaskWorker(task)

			if string(workedTask.result[14:]) == "successed" {
				done <- workedTask
			} else {
				undone <- fmt.Errorf("Task id %d time %s, error %s", workedTask.id, workedTask.createTime, workedTask.result)
			}
		}
		close(done)
		close(undone)
	}()

	return done, undone
}

func main() {

	taskChan := TaskCreater()
	doneTasks, undoneTasks := TaskSorter(taskChan)

	result := map[int]Task{}
	err := []error{}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for r := range doneTasks {
			result[r.id] = r
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for r := range undoneTasks {
			err = append(err, r)
		}
	}()

	wg.Wait()

	time.Sleep(time.Second * 3)

	println("Errors:")
	for id, data := range err {
		println(fmt.Sprintf("%d : %s", id, data.Error()))
	}

	println("Done tasks:")
	for _, data := range result {
		fmt.Println(data)
		fmt.Println("-----")
	}
}
