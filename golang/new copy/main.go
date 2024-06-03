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
	taskResult string // !!! возможно сделать стринг
}

type ResultTasks struct {
	DoneTask  []Task
	ErrorTask []Task
}

var (
	generateTime        time.Duration = 3
	outputTime          time.Duration = 1
	generationfrequency time.Duration = 1000
)

func taskGenerate(superChan chan Task, wg *sync.WaitGroup) { // !!!функция только на отдавания в канал// !!! а переименновать
	defer wg.Done()
	start := time.Now()
	// var i int = 1
	for {
		if time.Since(start) >= generateTime*time.Second {
			close(superChan) // Закрываем канал после 10 секунд
			return           // Выходим из функции
		} // !!!
		createTime := time.Now().Format(time.RFC3339)
		if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков !!!
			createTime = "Some error occured"
			// fmt.Println("_____1")
		}
		wg.Add(1)
		// fmt.Println("1_1 open ", i)
		// i++
		superChan <- Task{createTime: createTime, id: int(time.Now().Unix())} // передаем таск на выполнение
		time.Sleep(generationfrequency * time.Millisecond)                    // Ограничение на 1 секунду

	}
	// !!! возможно здесь надо сделать слип 10 секунд, и счетчик горутин
}

func tickerResult(mu *sync.Mutex, resultTasks *ResultTasks, resultChan chan ResultTasks, stop chan bool) {

	// defer wg.Done()
	ticker := time.NewTicker(outputTime * time.Second)
	fmt.Println("tickerResult")

	go func() {
		_, ok := <-stop
		if ok {
			fmt.Println("_____TickerResult close")
			// wg.Done()
			mu.Lock()
			resultChan <- *resultTasks
			close(resultChan)
			return

		}
	}()

	i := 1
	for range ticker.C {

		mu.Lock()
		resultChan <- *resultTasks
		fmt.Println("Передаем на печать_ ", i)
		i++
		// fmt.Println(resultTasks.DoneTask)
		mu.Unlock()
		fmt.Println("new tik___")

		// fmt.Println("new tik__2_ ", ok)

		// _, ok := <-doneChan
		// if !ok {
		// 	fmt.Println("tickerResult close")
		// 	// wg.Done()
		// 	close(resultChan)
		// 	return
		// }
	}

}

func printResult(resultChan chan ResultTasks, wg *sync.WaitGroup) {
	defer wg.Done()
	defer fmt.Println("printResult close")

	for result := range resultChan {
		fmt.Println("____result.DoneTask")
		println("Done tasks:")
		for _, task_d := range result.DoneTask {
			fmt.Println(task_d)
			// fmt.Println("task_r")
		}
		println("Errors:")
		for _, task_e := range result.ErrorTask {
			fmt.Printf("Task id %d time %s, error %s\n", task_e.id, task_e.createTime, task_e.taskResult)
			fmt.Println("task_r")
		}
	}

}

func taskWorker(task Task, doneChan, errorChan chan Task, wg *sync.WaitGroup) { // !!! не нужно чтоб возвращал таск
	defer wg.Done()
	tt, _ := time.Parse(time.RFC3339, task.createTime)
	// time.Sleep(3 * time.Second)
	if tt.After(time.Now().Add(-20 * time.Second)) {

		task.taskResult = "task has been successed"
		task.finishTime = time.Now().Format(time.RFC3339Nano)
		doneChan <- task
		// fmt.Println("k2__1")
	} else {
		task.taskResult = "something went wrong"
		task.finishTime = time.Now().Format(time.RFC3339Nano)
		errorChan <- task
		// fmt.Println("k2__2")
	}

	// fmt.Println("k2__")
	time.Sleep(time.Millisecond * 150) //!!! не очень понятно зачем эта задержка, функциональной нагрузке не несет

}

func storage(current chan Task, storage *[]Task, wg *sync.WaitGroup, mu *sync.Mutex, stop chan bool) {
	defer func() { stop <- true }()
	defer wg.Done()
	for task := range current {
		mu.Lock()
		*storage = append(*storage, task) // !!! rпереименновать
		mu.Unlock()
		wg.Done()
	}

	// for _, tast := range *storage { //евременно
	// 	fmt.Println(tast.id)
	// }

	fmt.Println("storage close" /*, storage*/)

}

func main() {
	start := time.Now()
	fmt.Println("Start")
	// ticker := time.NewTicker(outputTime * time.Second)
	// _ = ticker
	var wg sync.WaitGroup
	superChan := make(chan Task, 10) // !!1 переименновать канал
	doneChan := make(chan Task, 10)
	errorChan := make(chan Task, 10)
	_ = errorChan
	resultChan := make(chan ResultTasks)
	stop := make(chan bool)
	result := ResultTasks{DoneTask: []Task{}, ErrorTask: []Task{}}
	mu := sync.Mutex{}
	// Генерируем таскu

	func() {
		// defer func() { stop <- true }()
		wg.Add(2)
		go storage(doneChan, &result.DoneTask, &wg, &mu, stop)
		go storage(errorChan, &result.ErrorTask, &wg, &mu, stop)
	}()

	go tickerResult(&mu, &result, resultChan, stop)
	go func() {
		wg.Add(1)
		go printResult(resultChan, &wg) // !!! попробовать без  go
	}()

	wg.Add(1)
	go taskGenerate(superChan, &wg)

	go func() {
		wg_t := sync.WaitGroup{}
		for task_c := range superChan {
			wg_t.Add(1)
			go taskWorker(task_c, doneChan, errorChan, &wg_t) //!!!каналы сделать однонаправленные // go
			// fmt.Println(result.DoneTask)
		}

		wg_t.Wait()
		close(doneChan)
		close(errorChan)
		// stop <- true
		// fmt.Println("stop close")
		fmt.Println("close superchan")
		// fmt.Println(result)
	}()

	wg.Wait()
	fmt.Println("wait close")
	// stop <- true
	fmt.Println("stop close")
	// close(doneChan)
	// time.Sleep(3 * time.Second)
	// !!! закрыть все не закрытые каналы
	finish := time.Now()
	fmt.Println("Finish - ", finish.Sub(start))

}
