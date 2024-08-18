package main

import (
	"fmt"
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

// A Ttype represents a meaninglessness of our life
type Task struct {
	id     int
	cT     time.Time
	fT     time.Time
	result []byte // Результат задачи
	error  error  // Ошибка задачи, удобнее хранить в отдельном поле
}

var (
	// Сделаем буфер на 10 тасок, чтобы  сгладить колебания между генерацией и обработкой задач
	superChan      = make(chan Task, 10)
	doneTasks      = make(chan Task, 10)
	undoneTasks    = make(chan Task, 10)
	waitGroup      sync.WaitGroup
	allDoneTasks   []Task // Срез для хранения всех завершенных задач
	allUndoneTasks []Task // Срез для хранения всех задач с ошибками
)

func taskCreator(superChan chan<- Task) {
	// Устанавливаем таймер на 10 секунд
	timer := time.NewTimer(10 * time.Second) // может быть не очень эффективно с тз CPU, но в данном случае это не критично
	defer close(superChan)                   // Безопасное закрытие канала по истечении таймера
	for {
		select {
		case <-timer.C:
			return
		default:
			currentTime := time.Now()
			nanoStr := strings.TrimRight(strconv.Itoa(currentTime.Nanosecond()), "0") // Удаление незначащих нулей справа
			nano, err := strconv.Atoi(nanoStr)
			if err != nil {
				fmt.Println("Error converting nanoseconds to int:", err)
				return
			}
			task := Task{
				id: int(currentTime.Unix()),
				cT: currentTime,
			}
			if nano%2 > 0 { // Случайное условие для симуляции ошибки
				task.error = fmt.Errorf("random error occurred")
			}
			superChan <- task
			time.Sleep(time.Millisecond * 1000) // Небольшая задержка для замедления генерации
		}
	}
}

// taskWorker обрабатывает задачи из superChan, добавляя результаты в doneTasks или undoneTasks.
func taskWorker(superChan <-chan Task) {
	defer waitGroup.Done()
	defer close(doneTasks)   
	defer close(undoneTasks) 

	for task := range superChan {
		task.fT = time.Now()
		if task.error == nil {
			task.result = []byte(fmt.Sprintf("Task %d completed successfully", task.id))
			doneTasks <- task
		} else {
			undoneTasks <- task
		}
		time.Sleep(time.Millisecond * 150) // Имитация обработки задачи
	}
}

// taskSorter выводит результаты каждые 3 секунды и завершает работу после закрытия канала результатов.
func taskSorter(quit <-chan struct{}) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-quit:
			return
		case <-ticker.C:
			appendTasks(doneTasks, &allDoneTasks)
			appendTasks(undoneTasks, &allUndoneTasks)

			fmt.Println("Errors:")
			for _, task := range allUndoneTasks {
				fmt.Printf("Task ID: %d, Created: %s, Error: %v\n", task.id, task.cT, task.error)
			}

			fmt.Println("Done tasks:")
			for _, task := range allDoneTasks {
				fmt.Printf("Task ID: %d, Created: %s, Finished: %s\n", task.id, task.cT, task.fT)
			}
		}
	}
}

func appendTasks(taskChan <-chan Task, allTasks *[]Task) {
	for {
		select {
		case task, ok := <-taskChan:
			if !ok {
				return // Выходим, если канал закрыт
			}
			*allTasks = append(*allTasks, task)
		default:
			return // Выходим, когда в канале временно нет задач
		}
	}
}

func main() {
	quit := make(chan struct{})

	go taskCreator(superChan)

	waitGroup.Add(1)
	go taskWorker(superChan)
	go taskSorter(quit)

	waitGroup.Wait() // Дожидаемся завершения taskWorker
	close(quit)      // Сигнализируем taskSorter о завершении работы
}

