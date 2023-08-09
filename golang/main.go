package main

import (
	"context"
	"flag"
	"golang/app"
	"golang/config"
	"log"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнены остальных тасков

// TODO: добавить контекст и обработку сигналов
func main() {
	var n = flag.Int("n", 10, "total tasks number")
	flag.Parse()

	cfg := &config.Config{TasksQueueLimit: 10}
	App := app.New(cfg)
	if err := App.Run(context.Background(), *n); err != nil {
		log.Println(err)
	}
}
