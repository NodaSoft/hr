package main

import (
	"fmt"
	"sync"
	"time"
)

// A Task represents a meaninglessness of our life
type Task struct {
	id          int
	cT          time.Time // время создания
	fT          time.Time // время выполнения
	isSucceeded bool      //статус выполнения
	taskRESULT  []byte    //результат выполнения
}

func taskCreator(mc chan<- Task, wt time.Duration) {

	timer := time.NewTimer(wt * time.Second) //таймер  работы создателя тасков

	go func() {
		defer timer.Stop()

		for i := 0; ; i++ {
			select {
			case <-timer.C:
				close(mc) //закрываем главный канал тк его должен закрыть сендер а не ресивер
				return
			default:
				ft := time.Now()
				if (ft.Nanosecond()/100)%2 > 0 { // вот такое условие появления ошибочных тасков //обавил/100 так как наносекуды возвращаются всегда
					// с двумя нулями в конце это условие иначе никогда не сработает
					ft = time.Time{} // улевое время будет маркером ошибки
				}
				mc <- Task{cT: ft, id: i} // передаем таск на выполнение

			}

		}
	}()

}

func taskWorker(t Task, st chan Task, wg *sync.WaitGroup) { //работаем над таском и отправляем его в канал на сортировку
	defer wg.Done()
	if t.cT.After(time.Now().Add(-20 * time.Second)) {
		t.taskRESULT = []byte("task has been successed")
		t.isSucceeded = true
	} else {
		t.taskRESULT = []byte("something went wrong")
		t.isSucceeded = false
	}
	t.fT = time.Now()

	time.Sleep(time.Millisecond * 150)

	st <- t

}

func taskSorter(t Task, wg *sync.WaitGroup, dt, ut *[]Task) { //сортировщик тасков
	defer wg.Done()

	dtm := sync.Mutex{}
	utm := sync.Mutex{}
	if t.isSucceeded { //читаем статус таска и пишем его в соответствующий массив не забыв при этом мютекс тк сортировщик асинхронный
		dtm.Lock()
		*dt = append(*dt, t)
		dtm.Unlock()
	} else {
		utm.Lock()
		*ut = append(*ut, t)
		utm.Unlock()
	}

}

func printer(wt time.Duration, dt, ut *[]Task) {

	ticker := time.NewTicker(wt * time.Second)
	defer ticker.Stop()
	mu := sync.Mutex{}
	lastDt := -0
	lastUt := -0

	for {
		select {
		case <-ticker.C:
			mu.Lock()
			fmt.Println("last update at %s", time.Now().String())
			println("done tasks")
			for ; lastDt < len(*dt); lastDt++ {

				println((*dt)[lastDt].id)

			}
			println("undone tasks")
			for ; lastUt < len(*ut); lastUt++ {

				println((*ut)[lastUt].id)

			}
			mu.Unlock()
		}
	}

}

func printTotal(dt, ut []Task) {
	println(len(dt))
	println("TOTAL RESULT :")
	println("Completed: ")
	for _, v := range dt {
		fmt.Printf("task id: %v Completed At %s \n", v.id, v.cT.Format(time.RFC3339Nano))
	}
	println("Not Completed: ")
	for _, v := range ut {
		fmt.Printf("task id: %v Failed At %s \n", v.id, v.cT.Format(time.RFC3339Nano))
	}

}

func main() {

	mainChan := make(chan Task, 10)
	sortChan := make(chan Task, 10)
	var wgMain sync.WaitGroup //делаем группу так как main завершит все горутины прежде чем они вообще начнутся
	var doneTasks []Task
	var undoneTasks []Task

	taskCreator(mainChan, 10) //запускаем таск креатор добавил в него время работы

	go func() {
		wgMain.Add(1)
		defer wgMain.Done()
		defer close(sortChan)

		var wgSort sync.WaitGroup //группа на случай асинхронных воркеров

		for t := range mainChan {

			wgSort.Add(1)

			/*go*/
			taskWorker(t, sortChan, &wgSort)
			//воркер все равно блокирует чтение из канала так что запускать сортер в горутине тут нет особо смысла
			//так что воркер шлет результат работы в сорт канал откуда мы и заберем таск на сортировку
			//так же можно запустить несколько воркеров в отдельной горутине каждый на этот случай ждем wgsort
			//что бы не закрыть его раньше чем все воркеры отработают
			//блокирующий воркер эт обутылочное горлышко в данном случае
		}
		wgSort.Wait()

	}()

	go func() {
		wgMain.Add(1)
		defer wgMain.Done()

		for t := range sortChan {
			wgMain.Add(1)
			go taskSorter(t, &wgMain, &doneTasks, &undoneTasks) //забираем таски из сорт канала на сортировку
		}

	}()
	go func() {
		printer(3, &doneTasks, &undoneTasks) //принтер промежуточных результатов
	}()

	time.Sleep(2 * time.Millisecond) //небольшой дилей на инициализацию горутин
	wgMain.Wait()                    //ждем main пока не завершатся все горутины

	fmt.Println("All tasks processed, channels closed.")
	printTotal(doneTasks, undoneTasks) //итоговый принтер результатов

}
