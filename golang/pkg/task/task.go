package task

import (
	"fmt"
	"time"
)

// A Tasks represents a meaninglessness of our life
type Tasks struct {
	Id           int
	creationTime string // время создания
	finishTime   string // время выполнения
	taskRESULT   []byte
}

// Генератора тасков.
// по условию 10 секунд генерирует таски и закрывает канал.
func Generator() chan Tasks {
	out := make(chan Tasks, 10)
	ticker := time.NewTicker(time.Second * 10)
	go func() {
		for {
			select {
			case <-ticker.C:
				close(out)
				return
			default:
				startTime := time.Now().Format(time.RFC3339)
				if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
					startTime = "Some error occured"
				}
				out <- Tasks{creationTime: startTime, Id: int(time.Now().Unix())} // передаем таск на выполнение
			}
		}
	}()

	return out
}

// Work запускает обработку тасков.
// Возвращает канал c выполненными тасками и
// канал ошибок.
func Work(in <-chan Tasks) (chan Tasks, chan error) {
	taskChan := make(chan Tasks)
	go func() {
		for task := range in {
			task = taskWorker(task)
			taskChan <- task
		}
		close(taskChan)
	}()
	doneTasks, undoneTasks := taskSorter(taskChan)

	return doneTasks, undoneTasks
}

// имитация работы
func taskWorker(task Tasks) Tasks {
	tt, _ := time.Parse(time.RFC3339, task.creationTime)
	if tt.After(time.Now().Add(-20 * time.Second)) {
		task.taskRESULT = []byte("task has been successed")
	} else {
		task.taskRESULT = []byte("something went wrong")
	}
	task.finishTime = time.Now().Format(time.RFC3339Nano)

	time.Sleep(time.Millisecond * 150)

	return task
}

// сортировка тасков на выполненные и нет
// информация по ошибочным таскам отправляются в канал ошибок
// возвращает канал таксков и канал ошибок
func taskSorter(tasks <-chan Tasks) (chan Tasks, chan error) {
	doneTasks := make(chan Tasks)
	undoneTasks := make(chan error)
	go func() {
		for t := range tasks {
			if string(t.taskRESULT[14:]) == "successed" {
				doneTasks <- t
			} else {
				undoneTasks <- fmt.Errorf("task id %d doneTasks, undoneTaskstime %s, error %s", t.Id, t.creationTime, t.taskRESULT)
			}
		}
		close(doneTasks)
		close(undoneTasks)
	}()

	return doneTasks, undoneTasks
}
