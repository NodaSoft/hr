package main

import (
	// "fmt"
	"fmt"
	// "sync"

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

// !!! Меняем на более читаемый нейминг
// !!!Вынесла функции из мейл

// A Task represents a meaninglessness of our life
type Task struct {
	id         int
	createTime string // время создания
	finishTime string // время выполнения
	taskResult []byte // !!! возможно сделать стринг
}

var tmpTime time.Duration = 4

func taskGenerate(a chan Task, wg *sync.WaitGroup) { // !!!функция только на отдавания в канал// !!! а переименновать

	// go func() {
	//сделать тайминг работы 10
	// defer wg.Done()
	start := time.Now()
	for {
		wg.Add(1)
		if time.Since(start) >= tmpTime*time.Second {
			close(a) // Закрываем канал после 10 секунд
			return   // Выходим из функции
		} // !!!
		// fmt.Println(1)

		generateTime := time.Now().Format(time.RFC3339)
		if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
			generateTime = "Some error occured"
		}
		a <- Task{createTime: generateTime, id: int(time.Now().Unix())} // передаем таск на выполнение
		// time.Sleep(1 * time.Second)                                    // Ограничение на 1 секунду

	}
	// time.Sleep(time.Second * 2)
	// }()
	// !!! возможно здесь надо сделать слип 10 секунд, и счетчик горутин
}

func taskWorker(a Task) Task {
	tt, _ := time.Parse(time.RFC3339, a.createTime)
	if tt.After(time.Now().Add(-2 * time.Second)) {
		a.taskResult = []byte("task has been successed")
	} else {
		a.taskResult = []byte("something went wrong")
	}
	a.finishTime = time.Now().Format(time.RFC3339Nano)

	time.Sleep(time.Millisecond * 150) //!!! не очень понятно зачем эта задержка, функциональной нагрузке не несет

	return a
}

func main() {
	var wg sync.WaitGroup
	superChan := make(chan Task, 10)

	// Генерируем таски
	// wg.Add(1)
	go taskGenerate(superChan, &wg)

	// taskWorker := func(a Task) Task {
	// 	tt, _ := time.Parse(time.RFC3339, a.createTime)
	// 	if tt.After(time.Now().Add(-20 * time.Second)) {
	// 		a.taskResult = []byte("task has been successed")
	// 	} else {
	// 		a.taskResult = []byte("something went wrong")
	// 	}
	// 	a.finishTime = time.Now().Format(time.RFC3339Nano)

	// 	time.Sleep(time.Millisecond * 150)

	// 	return a
	// }

	// workerTasks := make(chan Task)
	doneTasks := make(chan Task)
	// resultsChan := make(chan Task)
	// undoneTasks := make(chan error)

	// taskSorter := func(t Task) {
	// 	if string(t.taskResult[14:]) == "successed" {
	// 		doneTasks <- t
	// 	} else {
	// 		undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.createTime, t.taskResult)
	// 	}
	// }
	// result := map[int]Task{}
	// mu := sync.Mutex{}
	// err := []error{}
	// taskSorter := func(t Task) {
	// 	// fmt.Println(4)
	// 	if string(t.taskResult[14:]) == "successed" {
	// 		mu.Lock()
	// 		defer mu.Unlock()
	// 		result[t.id] = t
	// 		// fmt.Println("iiii", result)
	// 		// } else {
	// 		// 	s := fmt.Errorf("Task id %d time %s, error %s", t.id, t.createTime, t.taskResult)
	// 		// 	err = append(err, s)
	// 		// 	// fmt.Println(6)

	// 	}
	// }

	// taskSorter := func(t Task) {

	go func() {

		// получение тасков
		// fmt.Println(2)
		for t := range superChan {
			// wg.Add(1)
			t = taskWorker(t)
			// workerTasks <- t // !!! либо все таки функцие отдельн о сделать
			doneTasks <- t
			// resultsChan <- t
			// go taskSorter(t) //
		}
		close(doneTasks) // !!! Закрываем по счетчикам
	}()

	// go func() {
	// 	for t := range doneTasks {
	// 		mu.Lock()
	// 		defer mu.Unlock()
	// 		defer wg.Done()
	// 		// if string(t.taskResult) == "task has been successed" {
	// 		// fmt.Println(2)
	// 		doneTasks <- t
	// 		// } else {
	// 		// undoneTasks <- fmt.Sprintf("Task id %d time %s, error %s", t.id, t.createTime, t.taskResult)
	// 		// }
	// 	}
	// }()

	// fmt.Println(3)
	// fmt.Println("pppppp ", result)
	//
	// go func() {
	// 	for r := range doneTasks {
	// 		go func() {
	// 			result[r.id] = r
	// 		}()
	// 	}
	// 	for r := range undoneTasks {
	// 		go func() {
	// 			err = append(err, r)
	// 		}()
	// 	}
	// 	close(doneTasks)
	// 	close(undoneTasks)
	// }()

	// time.Sleep(time.Second * 5)

	// go func() {
	// 	for r := range workerTasks {
	// 		if string(r.taskResult[14:]) == "successed" {
	// 			// 	// 		result[t.id] = t
	// 			// 	// 		// fmt.Println("iiii", result)
	// 			// 	// 	} else {
	// 			// 	// 		s := fmt.Errorf("Task id %d time %s, error %s", t.id, t.createTime, t.taskResult)
	// 			// 	// 		err = append(err, s)
	// 			// 	// 		fmt.Println(6)

	// 		}
	// 	}
	// }()

	// println("Errors:")
	// for r := range err {
	// 	println(r)
	// }

	// println("Done tasks:")
	// for r := range result {
	// 	println(r)
	// }
	doneResults := make(chan Task)
	go func() {

		for task := range doneTasks {
			// mu.Lock()

			// println(task.id)
			doneResults <- task
			// mu.Unlock()
			// wg.Done()
		}
		close(doneResults)
	}()
	// for err := range undoneTasks {
	// 	errResults = append(errResults, err)
	// }

	for task_t := range doneResults {
		println("ooooo____________")
		// if task_t, ok := <-doneResults; ok {
		// _, ok := <-doneTasks // !!! возможно завершение надо сделать все таки через мьютексы
		// if !ok {
		// 	return
		// }

		println("Errors:")
		// for r := range err {
		// 	println(r)
		// 	println("error_t", r)
		// }

		println("Done tasks:")
		// for r := range doneResults {
		fmt.Println(task_t)
		// }
		fmt.Println("____")

		// go func() {
		wg.Wait()
		// close(doneTasks)
		// close(undoneTasks)
		// }()
		time.Sleep(3 * time.Second)
	}

	// for {
	// 	task, ok := <-superChan
	// 	if ok {
	// 		return
	// 	}
	// 	fmt.Println(task)
	// }
	// for time.Since(time.Now().Add(-4*time.Second)) > 0 {
	// 	fmt.Println(<-superChan)
	// }
	// wg.Wait()

}

// for task := range doneTasks {
// 		doneResults = append(doneResults, task)
// 	}

// 	for err := range undoneTasks {
// 		errResults = append(errResults, err)
// 	}
