package main

import (
	"time"

	"github.com/fedoroko/nodasoft/test/internal/processor"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

func main() {
	numOfWorkers := 10
	handler := processor.NewTaskProcessor(numOfWorkers)
	handler.Run()
	defer handler.CloseAndPrint()

	time.Sleep(time.Second * 3)
}
