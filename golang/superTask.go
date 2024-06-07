package main

import (
	"TestTask/infrastructure/creator"
	"context"
	"fmt"
	"time"
)

type SuperTask struct {
	id          int
	dateCreate  string // время создания
	dateExecute string // время выполнения
	result      []byte
}

func (m SuperTask) String() string {
	return fmt.Sprintf("ID: %d | RESULT: %s | DATE_EXECUTE: %s | DATE_CREATE: %s", m.id, m.result, m.dateExecute, m.dateCreate)
}

type SuperTaskMessage struct {
	value SuperTask
	err   error
}

func (m SuperTaskMessage) GetValue() SuperTask {
	return m.value
}

func (m SuperTaskMessage) GetError() error {
	return m.err
}

func (m SuperTaskMessage) IsError() bool {
	return m.err != nil
}

type SuperTaskCreator struct{}

func (c SuperTaskCreator) Start(ctx context.Context, outCh chan<- creator.TaskMessage[SuperTask], complete chan<- bool) {
	defer func() { complete <- true }()
	for {
		task := SuperTask{
			id:         int(time.Now().UnixNano()), // С Unix() много дубликатов
			dateCreate: time.Now().Format(time.RFC3339),
		}

		var err error
		if time.Now().Nanosecond()%2 > 0 { // У меня возвращает числа в виде 1717763674048030300, всегда с 2 нулями на конце -> ошибок не бывает, я бы поставил UnixMilli()
			err = fmt.Errorf("ID: %d | MESSAGE: %s | DATE_CREATE: %s", task.id, "Some error occured", task.dateCreate)
		}

		select {
		case <-ctx.Done():
			return
		case outCh <- SuperTaskMessage{ // передаем таск на выполнение
			value: task,
			err:   err,
		}:
		}
	}
}
