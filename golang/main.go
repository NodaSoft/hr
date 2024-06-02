package main

import (
	"fmt"
	"math/rand"
	"sync"
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

type Task struct {
	id         int       // уникальный идентификатор задач
	broken     bool      // поле обозначающее дефектность
	createdAt  time.Time // время создания
	finishedAt time.Time // время выполнения
	result     string    // результат выполнения задачи
}

// константы для настройки
const sheduleDuration = 10 * time.Second // длительность планирования
const sheduleFrequency = 1 * time.Second // частота планирования
const minProcessSeconds = 5              // минимальное время выполнения задачи
const maxProcessSeconds = 20             // максимальное время выполнения задачи
const displayFrequency = 3 * time.Second // частота вывода результатов
const dtFormat = time.DateTime           // так же форматы которые могут пригодится - RFC3339Nano

func main() {
	var wg sync.WaitGroup
	var mu = &sync.Mutex{}
	var doneTasks, errorTasks []Task

	// инициализируем канал задач и канал стоп-сигнал
	sheduleChan, sheduleStop := make(chan Task, 10), make(chan struct{})

	// запускам планирование задач
	go sheduleTasks(sheduleChan, sheduleStop, sheduleDuration, sheduleFrequency)

	// запускаем вывод результатов каждые 3 секунды
	go func() {
		ticker := time.NewTicker(displayFrequency)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				displayProgress(&doneTasks, &errorTasks)
			}
		}
	}()

	// читаем запланированные задачи пока не придет сигнал, далее ждем пока все задачи закончатся (с помощью wg)
	// используем wg созданный для отслеживания задач в том числе для ожидания стоп сигнала для избегания случаев где у нас обработка опережает создание задач
	// например когда первой генеируется задача с ошибкой (ее обработка моментальная)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case sheduledTask := <-sheduleChan:
				wg.Add(1)
				go processTask(sheduledTask, &doneTasks, &errorTasks, &wg, mu)
			case <-sheduleStop:
				fmt.Println("No more tasks to schedule.")
				return
			}
		}
	}()
	wg.Wait()
	displayProgress(&doneTasks, &errorTasks) // финальный вывод результатов
}

func (t *Task) toString() string {
	return fmt.Sprintf("id - %v, broken - %v, createdAt - %v, finishedAt - %v, result - %v", t.id, t.broken, t.createdAt.Format(dtFormat), t.finishedAt.Format(dtFormat), t.result)
}

// Функция создания задач
func createTask(lastId *int) Task {
	creationTime := time.Now()
	*lastId++
	newTask := Task{id: *lastId, createdAt: creationTime}
	if rand.Float64() < 0.44 {
		newTask.broken = true // Некий дефект в задаче который приведет к ошибке
	}
	return newTask
}

// Функция для заполнения канала shedule структурами Task в течении duration и с частотой frequency
func sheduleTasks(shedule chan<- Task, stop chan<- struct{}, duration, frequency time.Duration) {
	timeout := time.After(duration)
	id := -1
	for {
		select {
		case <-timeout:
			close(stop)
			fmt.Printf("Sheduled %v tasks.\n", id) // закрываем стоп канал вызывая срабатывание выхода из for select loop
			return
		default:
			newTask := createTask(&id)
			shedule <- newTask
			time.Sleep(frequency)
		}
	}
}

func processTask(task Task, doneTasks, errorTasks *[]Task, wg *sync.WaitGroup, mu *sync.Mutex) {
	defer wg.Done()

	// если сломано
	if task.broken {
		// записываем результат и время окончания обработки
		task.finishedAt = time.Now()
		task.result = "error"

		// безопасно записываем в слайс законченных с ошибкой задач
		mu.Lock()
		*errorTasks = append(*errorTasks, task)
		mu.Unlock()
		return
	}

	// имитируем процессинг задачи
	processDuration := time.Duration(randInt(minProcessSeconds, maxProcessSeconds)) * time.Second
	time.Sleep(processDuration)

	// записываем результат и время окончания обработки
	task.finishedAt = time.Now()
	task.result = "success"

	// безопасно записываем в слайс законченных задач
	mu.Lock()
	*doneTasks = append(*doneTasks, task)
	mu.Unlock()
}

func displayProgress(done, doneWithError *[]Task) {
	fmt.Println("DONE TASKS:")
	for _, task := range *done {
		fmt.Printf("%s\n", task.toString())
	}
	fmt.Println("DONE WITH ERROR TASKS:")
	for _, task := range *doneWithError {
		fmt.Printf("%s\n", task.toString())
	}
}

// Функция для генерация int в нужных диапазонах
func randInt(min, max int) int {
	return rand.Intn(max-min) + min
}
