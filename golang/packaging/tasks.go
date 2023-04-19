package main

import (
	"fmt"
	"time"
)

// Сюда вынесено всё, что относится непосредственно к таскам. Это можно покрывать юнит тестами легко.
// Можно также выделить в отдельный пакет и поправить структуру тасок, дабы не фигачить ошибки в поле времени..

// A Task represents a meaninglessness of our life
type Task struct {
	id         int
	cT         string // время создания или текст ошибки создания
	fT         string // время выполнения
	taskRESULT []byte
}

// NewTask -- создаем таску в куче и сразу ее настраиваем
func NewTask() any {
	created := time.Now()
	task := new(Task)
	task.InitTask(created)
	time.Sleep(51 * time.Microsecond) // слишком шустро..

	// вот такое условие появления ошибочных тасков .. переползло сюда, из-за унификации эмулятора
	if time.Now().Nanosecond()%2 > 0 {
		task.cT = "Some error occured"
	}

	return task
}

// InitTask -- инициализация задачи. Предпочтительнее генератора NewTask() т.к. работает с "внешним" объектом
func (task *Task) InitTask(created time.Time) *Task {
	task.cT = created.Format(time.RFC3339)
	//task.id = int(created.Unix()) так было, идент не уникален!
	task.id = int(created.UnixNano())
	return task
}

// sortResult -- анализ итога обраблотки таски и возврат успешного результата или форматирование ошибки
// так было, нафиг не надо, это должен делать обработчик вообще-то. Оставлено как было..
func (task *Task) sortResult() (*Task, error) {
	// TODO не лучший способ определять успешность получения задачи, но мало ли..
	if string(task.taskRESULT[14:]) == "successed" {
		return task, nil
	} else {
		return nil, fmt.Errorf("Task id %d time %s, error %s", task.id, task.cT, task.taskRESULT)
	}
}

// worker -- обработчик тасок. Может получить не время, а ошибку генератора таски.
func (task *Task) worker() *Task {
	if tt, err := time.Parse(time.RFC3339, task.cT); err != nil {
		task.taskRESULT = []byte("something went wrong, ct=" + task.cT)
	} else {
		if tt.After(time.Now().Add(-406100 * time.Microsecond)) {
			task.taskRESULT = []byte("task has been successed")
		} else {
			task.taskRESULT = []byte("something went wrong")
		}
	}
	task.fT = time.Now().Format(time.RFC3339Nano)

	time.Sleep(time.Millisecond * 150)

	return task
}

// ProcessTask -- callback обработчика таски в канализатор
func ProcessTask(item any) (any, error) {
	task, ok := item.(*Task)
	if ok {
		return task.worker().sortResult()
	} else {
		return nil, fmt.Errorf("ProcessTask() got %T, must be *Task", item)
	}
}

// GuidTask -- callback получения идента таски для канализатора
func GuidTask(item any) (int, error) {
	task, ok := item.(*Task)
	if ok {
		return task.id, nil
	}
	return 0, fmt.Errorf("GuidTask() got %T, must be *Task", item)
}
