package main

import (
	"fmt"
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
	taskResult string // !!! возможно сделать стринг
}

var tmpTime time.Duration = 4

func taskGenerate(a chan Task, wg *sync.WaitGroup) { // !!!функция только на отдавания в канал// !!! а переименновать
	//сделать тайминг работы 10
	start := time.Now()
	for {

		if time.Since(start) >= tmpTime*time.Second {

			close(a) // Закрываем канал после 10 секунд
			fmt.Println("Generate - stop")
			return // Выходим из функции
		} // !!!
		fmt.Println("_____1")

		generateTime := time.Now().Format(time.RFC3339)
		if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
			generateTime = "Some error occured"
		}
		wg.Add(1)
		a <- Task{createTime: generateTime, id: int(time.Now().Unix())} // передаем таск на выполнение
		time.Sleep(1 * time.Second)                                     // Ограничение на 1 секунду

	}
	// !!! возможно здесь надо сделать слип 10 секунд, и счетчик горутин
}

func taskWorker(a Task, doneTasks chan Task) Task {
	tt, _ := time.Parse(time.RFC3339, a.createTime)
	if tt.After(time.Now().Add(-2 * time.Second)) {
		a.taskResult = "task has been successed"

	} else {
		a.taskResult = "something went wrong"
	}
	a.finishTime = time.Now().Format(time.RFC3339Nano)

	time.Sleep(time.Millisecond * 150) //!!! не очень понятно зачем эта задержка, функциональной нагрузке не несет

	doneTasks <- a
	fmt.Println("G")
	return a
}

func main() {
	var wg sync.WaitGroup
	superChan := make(chan Task, 10)
	fmt.Println("Start")
	// Генерируем таски
	go taskGenerate(superChan, &wg)

	// workerTasks := make(chan Task)
	doneTasks := make(chan Task)
	// resultsChan := make(chan Task)
	// undoneTasks := make(chan error)
	ticker := time.NewTicker(1 * time.Second)
	mu := sync.Mutex{}
	// println("ttttt____________")
	result := map[int]Task{}
	go func() {
		// println("sssss____________")
		for range ticker.C {
			println("Errors:")
			// for r := range err {
			// 	println(r)
			// 	println("error_t", r)
			// }

			println("Done tasks:")
			// for r := range doneResults {
			mu.Lock()
			for task_r := range result {
				// _ = task_r
				fmt.Println(task_r)
			}

			// fmt.Println("____1")
			mu.Unlock()
			// }
			// fmt.Println("____")
			// wg.Wait()
		}

	}()

	go func() {
		for task_c := range superChan {

			// _ = task

			// _ = task
			fmt.Println(5)
			task_c = taskWorker(task_c, doneTasks) //!!!каналы сделать однонаправленные

		}
	}()

	go func() {
		defer wg.Done()
		// fmt.Println("____2/1")
		for r := range doneTasks {
			fmt.Println("____2")
			// go func(r Task) {
			mu.Lock()
			result[r.id] = r
			mu.Unlock()

			// }(r)
		}
		// close(doneTasks)
	}()
	// fmt.Println("____5")
	go func() {
		for {
			// 	// fmt.Println("____7")
			// 	// _, ok := <-doneTasks
			// 	// if !ok {
			fmt.Println("____3")
			wg.Wait()
			return
			// fmt.Println("____6")
			// 	return
			// 	// канал закрыт
		}

		// 	// 	// close(doneTasks)
		// 	// 	// close(undoneTasks)
	}()

	// }

	// }()
	// time.Sleep(3 * time.Second)
	// wg.Wait()
	fmt.Println("Все горутины завершены")
}
