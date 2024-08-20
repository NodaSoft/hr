package main

import (
	"context"
	"time"

	"github.com/kostromin59/NodaSoft/tasks"
)

// Приложение эмулирует получение и обработку неких тасков. Пытается и получать, и обрабатывать в многопоточном режиме.
// Приложение должно генерировать таски 10 сек. Каждые 3 секунды должно выводить в консоль результат всех обработанных к этому моменту тасков (отдельно успешные и отдельно с ошибками).

// ЗАДАНИЕ: сделать из плохого кода хороший и рабочий - as best as you can.
// Важно сохранить логику появления ошибочных тасков.
// Важно оставить асинхронные генерацию и обработку тасков.
// Сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через pull-request в github
// Как видите, никаких привязок к внешним сервисам нет - полный карт-бланш на модификацию кода.

const printDuration = 3 * time.Second
const timeoutDuration = 10 * time.Second

func main() {
	manager := tasks.NewManager()

	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	created := tasks.RunCreator(ctx)
	handled := tasks.RunWorker(ctx, created)

	go func() {
		ticker := time.NewTicker(printDuration)
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
			case <-ticker.C:
				manager.Print()
			}
		}
	}()

	for task := range handled {
		if task.Err != nil {
			manager.AddFailed(task)
		} else {
			manager.AddSuccessed(task)
		}
	}
}
