package main

import (
	"abcp-golang/pkg/result"
	"abcp-golang/pkg/task"
	"sync"
)

// Приложение эмулирует получение и обработку неких тасков. Пытается и получать, и обрабатывать в многопоточном режиме.
// Приложение должно генерировать таски 10 сек. Каждые 3 секунды должно выводить в консоль результат всех обработанных к этому моменту тасков (отдельно успешные и отдельно с ошибками).

// ЗАДАНИЕ: сделать из плохого кода хороший и рабочий - as best as you can.
// Важно сохранить логику появления ошибочных тасков.
// Важно оставить асинхронные генерацию и обработку тасков.
// Сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через pull-request в github
// Как видите, никаких привязок к внешним сервисам нет - полный карт-бланш на модификацию кода.

// обрабатывает канал готовых тасков и канал ошибок
func parseResults(doneTasks chan task.Tasks, undoneTasks chan error) (*result.Results, *sync.WaitGroup) {
	// res := result.NewResults()
	res := result.New()
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		for task := range doneTasks {
			res.AddResult(task)
		}
		wg.Done()
	}()

	go func() {
		for err := range undoneTasks {
			res.AddError(err)
		}
		wg.Done()
	}()

	return res, wg
}

func main() {
	tasksChan := task.Generator()

	doneTasks, undoneTasks := task.Work(tasksChan)
	results, wg := parseResults(doneTasks, undoneTasks)

	// печатаем результаты
	done := make(chan struct{})
	result.Print(done, results)
	wg.Wait()
	close(done)
}
