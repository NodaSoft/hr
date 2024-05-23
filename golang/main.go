package main

import (
	"fmt"
	"time"
)

// Приложение эмулирует получение и обработку неких тасков. Пытается их
// получать, и обрабатывать в многопоточном режиме.
// После обработки тасков в течении 3 секунд приложение должно выводить
// накопленные к этому моменту успешные таски и отдельно ошибки обработки
// тасков.

// ЗАДАНИЕ: сделать из плохого кода хороший и рабочий - as best as you
// can.
// Важно сохранить логику появления ошибочных тасков.
// Важно оставить асинхронные генерацию и обработку тасков.
// Сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через pull-request в github
// Как видите, никаких привязок к внешним сервисам нет - полный карт-
// бланш на модификацию кода.

// A Task is a processible task which keeps its processing status.
type Task struct {
	id           int
	createdTime  string
	finishedTime string
	result       []byte
}

func main() {
	sourceCh := make(chan Task, 10)

	go createTasks(sourceCh)

	doneCh := make(chan Task)
	errsCh := make(chan error)

	go processTasks(sourceCh, doneCh, errsCh)

	doneTasks := map[int]Task{}
	tasksErrs := []error{}

	go distributeDoneTasks(doneCh, doneTasks)
	// передаем вторым аргументом указатель на слайс, чтобы изменять
	// значение оригинального слайса внутри обрабатывающей функции
	go distributeTasksErrs(errsCh, &tasksErrs)

	time.Sleep(time.Second * 3)

	// заменяем встроенную функцию println, так как она может быть
	// удалена из языка в будущих релизах
	fmt.Println("Errors:")
	for _, err := range tasksErrs {
		fmt.Println(err)
	}

	fmt.Println("Done tasks:")
	for taskId := range doneTasks {
		fmt.Println(taskId)
	}
}

func createTasks(sourceCh chan Task) {
	for {
		createdTime := time.Now().Format(time.RFC3339)

		// вот такое условие появления ошибочных тасков
		if time.Now().Nanosecond()%2 > 0 {
			createdTime = "Some error occured"
		}

		// передаем таск на выполнение
		sourceCh <- Task{
			id:          int(time.Now().Unix()),
			createdTime: createdTime,
		}
	}
}

func processTasks(sourceCh, doneCh chan Task, errsCh chan error) {
	for task := range sourceCh {
		task = processTask(task)

		sortTask(task, doneCh, errsCh)
	}

	// выход из цикла произойдет автоматически, если канал будет закрыт

	// у нас нет необходимости закрывать каналы в данном приложении,
	// так как приложение работает 3 секунды, после этого выводит
	// результаты. После завершения память, выделенная под каналы,
	// освободится
}

func processTask(task Task) Task {
	parsedCreatedTime, _ := time.Parse(time.RFC3339, task.createdTime)

	if parsedCreatedTime.After(time.Now().Add(-20 * time.Second)) {
		task.result = []byte("task has been successed")
	} else {
		task.result = []byte("something went wrong")
	}

	task.finishedTime = time.Now().Format(time.RFC3339Nano)

	time.Sleep(time.Millisecond * 150)

	return task
}

func sortTask(task Task, doneCh chan Task, errsCh chan error) {
	if string(task.result[14:]) == "successed" {
		doneCh <- task

		return
	}

	errsCh <- fmt.Errorf("Task id %d time %s, error %s", task.id, task.createdTime, task.result)
}

func distributeDoneTasks(doneCh chan Task, doneTasks map[int]Task) {
	for task := range doneCh {
		doneTasks[task.id] = task
	}
}

func distributeTasksErrs(errsCh chan error, tasksErrs *[]error) {
	for task := range errsCh {
		// функция append возвращает новый слайс, когда в исходном не
		// хватает емкости для добавления новых элементов. Поэтому мы
		// меняем значение по адресу указателя на исходный слайс
		*tasksErrs = append(*tasksErrs, task)
	}
}
