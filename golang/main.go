package main

import (
	"fmt"
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
				if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
					nowStr = "Some error occured"
				}
				taskID := int(time.Now().UnixNano())         // NOTE: но лучше вообще генерировать случайный ID (например, GUID)
				ch <- Task{creationTime: nowStr, id: taskID} // передаем таск на выполнение
			}
		}()
	}

	superChan := make(chan Task, 10)

	go scheduleTasks(superChan)

	processTask := func(task Task) Task {
		parsedCreationTime, err := time.Parse(time.RFC3339, task.creationTime)
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
		// получение тасков
		for task := range superChan {
			task = processTask(task)
			go sortTask(task)
		}
		// Тут не нужно закрывать канал, потому что, если процесс и доберется до конца этой горутины, канал уже будет закрыт,
		// ведь итерация по каналу прекращается, только когда канал закрывается. А закрытие уже закрытого канала — ПАНИКА.
	}()

	doneTasks := map[int]Task{}
	taskErrors := []error{}
	// Защищаем параллельно обновляемые данные мьютексами:
	doneTasksMutex := sync.Mutex{}
	taskErrorsMutex := sync.Mutex{}

	// Параллелизируем обработку успешных и ошибочных случаев. Не дожидаемся закрытия каналов:
	go func() {
		for doneTask := range doneTasksChan {
			doneTask := doneTask
			go func() {
				doneTasksMutex.Lock() // Получаем эксклюзивный доступ
				defer doneTasksMutex.Unlock()
				doneTasks[doneTask.id] = doneTask
			}()
		}
	}()
	go func() {
		for taskErr := range taskErrorsChan {
			taskErr := taskErr
			go func() {
				taskErrorsMutex.Lock() // Получаем эксклюзивный доступ
				defer taskErrorsMutex.Unlock()
				taskErrors = append(taskErrors, taskErr)
			}()
		}
	}()

	time.Sleep(time.Second * 3)

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
