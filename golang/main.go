package main

import (
	"fmt"
	"sync"
	"time"
)

// Приложение эмулирует получение и обработку неких тасков. Пытается и получать, и обрабатывать в многопоточном режиме.
// После обработки тасков в течении 3 секунд приложение должно выводить накопленные к этому моменту успешные таски и отдельно ошибки обработки тасков.

// ЗАДАНИЕ: сделать из плохого кода хороший и рабочий - as best as you can.
// Важно сохранить логику появления ошибочных тасков.
// Важно оставить асинхронные генерацию и обработку тасков.
// Сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через pull-request в github
// Как видите, никаких привязок к внешним сервисам нет - полный карт-бланш на модификацию кода.

type Task struct {
	id             int
	creationTime   string
	completionTime string
	result         []byte
}

func main() {
	scheduleTasks := func(ch chan Task) {
		go func() {
			for {
				nowStr := time.Now().Format(time.RFC3339)
				if time.Now().Nanosecond()%2 > 0 {
					nowStr = "Some error occured"
				}
				// Используем наносекунды, чтобы избежать коллизий, но лучше вообще
				// генерировать случайный ID (например, GUID):
				taskID := int(time.Now().UnixNano())
				ch <- Task{creationTime: nowStr, id: taskID}
			}
		}()
	}

	superChan := make(chan Task, 10)

	go scheduleTasks(superChan)

	processTask := func(task Task) Task {
		// Тут надо удостоверяться, что строка успешно парсится в дату:
		parsedCreationTime, err := time.Parse(time.RFC3339, task.creationTime)
		// Здесь учитываем результат парсинга, помимо изначально заданного условия:
		if err == nil && parsedCreationTime.After(time.Now().Add(-20*time.Second)) {
			task.result = []byte("task has been successed")
		} else {
			task.result = []byte("something went wrong")
		}
		task.completionTime = time.Now().Format(time.RFC3339Nano)

		time.Sleep(time.Millisecond * 150)

		return task
	}

	doneTasksChan := make(chan Task)
	taskErrorsChan := make(chan error)

	sortTask := func(task Task) {
		if string(task.result)[14:] == "successed" { // Тут задумана индексация по строке, а не по массиву байтов
			doneTasksChan <- task
		} else {
			taskErrorsChan <- fmt.Errorf("Task id %d time %s, error %s", task.id, task.creationTime, task.result)
		}
	}

	go func() {
		for task := range superChan {
			task = processTask(task)
			go sortTask(task)
		}
		// Тут не нужно закрывать канал, потому что, если процесс и доберется до конца этой горутины, канал уже будет закрыт,
		// ведь итерация по каналу прекращается, только когда канал закрывается. А закрытие уже закрытого канала — ПАНИКА.
		//
		// Кстати, нет смысла закрывать во всем скрипте какой-либо канал, потому что все они обрабатывают информацию из superChan,
		// а его никогда не будет смысла закрывать, потому что scheduleTasks работает бесконечно.
		// Но если очень хочется закрыть каналы -- это можно сделать после time.Sleep(time.Second * 3)
	}()

	doneTasks := map[int]Task{}
	taskErrors := []error{}
	// Защищаем параллельно обновляемые данные мьютексами:
	doneTasksMutex := sync.Mutex{}
	taskErrorsMutex := sync.Mutex{}

	// Параллелизируем обработку успешных и ошибочных случаев. Не дожидаемся закрытия каналов:
	go func() {
		for doneTask := range doneTasksChan {
			doneTask := doneTask // Учитываем то, как в Голанге переменные захватываются анонимными функциями
			go func() {
				doneTasksMutex.Lock() // Получаем эксклюзивный доступ
				defer doneTasksMutex.Unlock()
				doneTasks[doneTask.id] = doneTask
			}()
		}
	}()
	go func() {
		for taskErr := range taskErrorsChan {
			taskErr := taskErr // Учитываем то, как в Голанге переменные захватываются анонимными функциями
			go func() {
				taskErrorsMutex.Lock() // Получаем эксклюзивный доступ
				defer taskErrorsMutex.Unlock()
				taskErrors = append(taskErrors, taskErr)
			}()
		}
	}()

	time.Sleep(time.Second * 3)

	// Ниже оборачиваем код в функции для удобного использования defer:

	println("Errors:")
	func() {
		taskErrorsMutex.Lock() // Получаем эксклюзивный доступ
		defer taskErrorsMutex.Unlock()
		for _, taskErr := range taskErrors {
			println(taskErr.Error())
		}
	}()

	println("Done tasks IDs:")
	func() {
		doneTasksMutex.Lock() // Получаем эксклюзивный доступ
		defer doneTasksMutex.Unlock()
		for taskID := range doneTasks {
			println(taskID)
		}
	}()
}
