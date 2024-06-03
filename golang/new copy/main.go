package main

import (
	"fmt"
	"sync"
	"time"
)

// A Task represents a meaninglessness of our life
type Task struct {
	id         int
	createTime string // время создания
	finishTime string // время выполнения
	taskResult string // результат выполнения таска
}

type ResultTasks struct {
	DoneTask  []Task // Канал для выполненных тасков
	ErrorTask []Task // Канал для тасков, которые отработали неккоректно
}

var (
	generateTime        time.Duration = 10 // Время, в течение которого генерируются таски, секунды
	outputTime          time.Duration = 3  // Частота появления  вывода результатов работы программы, секунды
	generationfrequency time.Duration = 1  // Время, требуемое для генерации одного таска, милисекунды
)

// tickerResult - функция отвечает за вывод в консоль результатов работы программы каждые 3 секунды
// в исходном коде функционал вывода каждые 3 секунды был реализован некорректно
func tickerResult(mu *sync.Mutex, resultTasks *ResultTasks, resultChan chan ResultTasks, stop chan bool) {
	ticker := time.NewTicker(outputTime * time.Second)
	go func() { // запускаем в горутине функцию, которая отлавливает завершения работы обработчика тасков
		_, ok := <-stop
		if ok {
			mu.Lock()
			resultChan <- *resultTasks
			close(resultChan)
			return
		}
	}()
	for range ticker.C { // срабатывает каждые 3 секунды
		mu.Lock()
		resultChan <- *resultTasks
		mu.Unlock()
	}
}

// taskGenerate - функция отвечает за генерацию тасков в течение заявленного времени, 10 секунд.
// В исходном коде код генерируеться бесконечно
func taskGenerate(generateChan chan Task, wg *sync.WaitGroup) {
	defer wg.Done()
	wg.Add(1)

	start := time.Now()
	for {
		if time.Since(start) >= generateTime*time.Second {
			close(generateChan) // Закрываем канал после 10 секунд
			return              // Выходим из функции
		}
		createTime := time.Now().Format(time.RFC3339)
		if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков !!! При Nanosecond негативные таски не генерируются, но если изменить на секунды  все окб оригинал также, по условию задания не меняем логи
			createTime = "Some error occured"
		}
		wg.Add(1)
		generateChan <- Task{createTime: createTime, id: int(time.Now().Unix())} // передаем таск на выполнение
		time.Sleep(generationfrequency * time.Millisecond)                       // Ограничениеб если хотим другую периодичность появления тасков
	}
}

// processingTasks  - функция отвечает за чтение из generateChan и запуска непосредственно функции выполнение тасков
// Вывела отделно для простотв чтения кода и посика ошибок
func processingTasks(generateChan, doneChan, errorChan chan Task) {
	wg_t := sync.WaitGroup{} // сделала внутри функции waitGroup, чтоб должаться окончания выполнения всех таска
	for task_c := range generateChan {
		wg_t.Add(1)
		go taskWorker(task_c, doneChan, errorChan, &wg_t)
	}
	wg_t.Wait()
	close(doneChan)
	close(errorChan)
}

// taskWorker - функция отвечает за выполнение таска
func taskWorker(task Task, doneChan, errorChan chan Task, wg *sync.WaitGroup) {
	defer wg.Done()
	tt, _ := time.Parse(time.RFC3339, task.createTime)
	if tt.After(time.Now().Add(-20 * time.Second)) {
		task.taskResult = "task has been successed"
		task.finishTime = time.Now().Format(time.RFC3339Nano)
		doneChan <- task
	} else {
		task.taskResult = "something went wrong"
		task.finishTime = time.Now().Format(time.RFC3339Nano)
		errorChan <- task
	}
	time.Sleep(time.Millisecond * 150)
}

// storage - функция отвечает за хранение отработанных тасков
func storage(current chan Task, storage *[]Task, wg *sync.WaitGroup, mu *sync.Mutex, stop chan bool) {
	defer func() { stop <- true }() // отлавливаем завершения работы функции
	defer wg.Done()

	wg.Add(1)
	for task := range current {
		mu.Lock()
		*storage = append(*storage, task)
		mu.Unlock()
		wg.Done()
	}
}

// printResult - функция отвечает за вывод результатов работы программы
func printResult(resultChan chan ResultTasks, wg *sync.WaitGroup) {
	defer wg.Done()
	wg.Add(1)
	for result := range resultChan {
		fmt.Println("Done tasks:")
		for _, task_d := range result.DoneTask {
			fmt.Println(task_d)
		}
		fmt.Println("Errors:")
		for _, task_e := range result.ErrorTask {
			fmt.Printf("Task id %d time %s, error %s\n", task_e.id, task_e.createTime, task_e.taskResult)
		}
	}
}

// main - функция запускает программу
func main() {
	var wg sync.WaitGroup
	generateChan := make(chan Task, 10) // !!1 переименновать канал
	doneChan := make(chan Task, 10)
	errorChan := make(chan Task, 10)
	resultChan := make(chan ResultTasks)
	stop := make(chan bool)

	result := ResultTasks{DoneTask: []Task{}, ErrorTask: []Task{}}
	mu := sync.Mutex{}
	// Отсеживаем время выводы в терминал результатов
	go tickerResult(&mu, &result, resultChan, stop)
	// Генерируем таски
	go taskGenerate(generateChan, &wg)
	// Запускаем таски
	go processingTasks(generateChan, doneChan, errorChan)
	// Храним все отработанные таски
	go storage(doneChan, &result.DoneTask, &wg, &mu, stop)
	go storage(errorChan, &result.ErrorTask, &wg, &mu, stop)
	// Распечатываем все отработанные таски н текущий момент
	printResult(resultChan, &wg)

	wg.Wait()
}
