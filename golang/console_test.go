package main

import (
	"testing"
	"time"
)

func TestConsole_Println(t *testing.T) {
	doneTasks := make(chan *Task, 10)
	undoneTasks := make(chan error, 10)
	withErr := make(chan *Task, 10)
	exit := make(chan bool, 1)
	i := 0
	c := Console{
		doneTasks:   doneTasks,
		undoneTasks: undoneTasks,
		withErr:     withErr,
		exitCode:    exit,
		print: func(a ...any) {
			//просто подсчитывает сколько сообщений мы отправили на печать
			i++
		},
	}

	go c.Println()
	doneTasks <- &Task{}
	doneTasks <- &Task{}
	doneTasks <- &Task{}
	withErr <- &Task{}
	withErr <- &Task{}

	time.Sleep(1 * time.Second)
	if i != 5 {
		t.Fatalf("assertions don't eq %d !=%d", 5, i)
	}
}
