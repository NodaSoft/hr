package main

import (
	"context"
	"t/model"
	"time"
)

// Приложение эмулирует получение и обработку неких тасков. Пытается и получать, и обрабатывать в многопоточном режиме.
// Приложение должно генерировать таски 10 сек. Каждые 3 секунды должно выводить в консоль результат всех обработанных к этому моменту тасков (отдельно успешные и отдельно с ошибками).

// ЗАДАНИЕ: сделать из плохого кода хороший и рабочий - as best as you can.
// Важно сохранить логику появления ошибочных тасков.
// Важно оставить асинхронные генерацию и обработку тасков.
// Сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через pull-request в github
// Как видите, никаких привязок к внешним сервисам нет - полный карт-бланш на модификацию кода.

func main() {
	// Контекст, останавливающий создание задач через 10 секунд
	creatorCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Каналы будут закрыты внутри task.Create и task.Route
	superChan := make(chan model.Task, 10)
	doneTasks := make(chan model.Task)
	errorTasks := make(chan error)

	go model.Task.Create(model.Task{}, creatorCtx, superChan)
	go model.Task.Route(model.Task{}, superChan, doneTasks, errorTasks)

	results := model.Tasks{}
	errs := model.Errors{}

	// Выводим задачи каждые 3 секунды
	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			results.Print()
			errs.Print()
			if doneTasks == nil && errorTasks == nil {
				return
			}
		case t, ok := <-doneTasks:
			if ok {
				results[t.GetId()] = t
			} else {
				doneTasks = nil
			}
		case e, ok := <-errorTasks:
			if ok {
				errs = append(errs, e)
			} else {
				errorTasks = nil
			}
		}
	}
}
