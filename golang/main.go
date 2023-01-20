package main

import (
	"fmt"
	"time"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнены остальных тасков

func main() {
	c := NewPostman(50*time.Second, 11*time.Millisecond)
	superChan := c.Create()

	factory, destroyer := NewWorkerFactory()
	//создаем рабочих сколько нужно
	for i := 0; i < 2; i++ {
		w := factory.Worker()
		w.Subscribe(superChan)
	}

	doneTasks, undoneTasks, withErr := factory.Chans()

	console := Console{
		doneTasks:   doneTasks,
		undoneTasks: undoneTasks,
		withErr:     withErr,
		exitCode:    destroyer,
		print: func(a ...any) {
			fmt.Println(a...)
		},
	}

	console.Println()
}
