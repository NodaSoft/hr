package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнены остальных тасков

// A Ttype represents a meaninglessness of our life
type Ttype struct {
	id         int
	cT         string // время создания
	fT         string // время выполнения
	taskRESULT []byte
}

func main() {

	// ID задачи формируется int(time.Now().Unix()) при таком подходе мы получаем дубли
	// решил использовать счетчик (конечно для промышленного решение не очень, для данной задачи подходит)
	var id int64
	id = 0

	taskCreturer := func(a chan Ttype) {
		go func() {
			for {
				ft := time.Now().Format(time.RFC3339)
				if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
					ft = "Some error occured"
				}
				// через atomic увеличиваем номер таска и передаём в канал для исполнения
				a <- Ttype{cT: ft, id: int(atomic.AddInt64(&id, 1))} // передаем таск на выполнение
			}
		}()
	}

	superChan := make(chan Ttype, 10)

	go taskCreturer(superChan)

	// воркер эмитирует обработку, его не трогал
	task_worker := func(a Ttype) Ttype {
		tt, _ := time.Parse(time.RFC3339, a.cT)
		if tt.After(time.Now().Add(-20 * time.Second)) {
			a.taskRESULT = []byte("task has been successed")
		} else {
			a.taskRESULT = []byte("something went wrong")
		}
		a.fT = time.Now().Format(time.RFC3339Nano)

		time.Sleep(time.Millisecond * 150)

		return a
	}

	// добавил буфер в каналы, так как ниже будем создавать несколько горутин с воркерами
	doneTasks := make(chan Ttype, 10)
	undoneTasks := make(chan error, 10)

	tasksorter := func(t Ttype) {
		if string(t.taskRESULT[14:]) == "successed" {
			doneTasks <- t
		} else {
			undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT, t.taskRESULT)
		}
	}

	// запускаем 10 горутин для обработки канала superChan
	// 10 штук согласно буферу канала superChan, тут есть поле для развития, нужно учитывать наполнение канала
	for i := 0; i < 10; i++ {
		go func() {
			// получение тасков
			for {
				// range поменял на select для работы с каналом в несколько горутин
				select {
				case <-superChan:
					t := task_worker(<-superChan)
					go tasksorter(t)
				default:
				}
			}
			// конечно этот код никогда не исполнится, но закрывать канал читателем не верно, может быть паника
			//close(superChan)
		}()
	}

	// map сразу делаем с выделением памяти, чтобы не алоцировать при вставке
	result := make(map[int]Ttype, 100)
	err := []error{}
	// осталась 1на горутина читатель, операции простые обрабатываются по поступлению (наполнятся длолжны медленне чем читаться из канала)
	go func() {
		for {
			// select  в данном случае удобен, так как мы можем обрабатывать несколько каналов
			select {
			case <-undoneTasks:
				r := <-undoneTasks
				err = append(err, r)
			case <-doneTasks:
				r := <-doneTasks
				result[r.id] = r
			default:
			}
		}
		//close(doneTasks)
		//close(undoneTasks)
	}()

	time.Sleep(time.Second * 3)

	println("Errors:")
	for r := range err {
		println(r)
	}

	println("Done tasks:")
	for r := range result {
		println(r)
	}
	// добавил общее количество обработанных тасков (просто удобно смотреть при запусках)
	println("=======")
	println(len(result))

	//P.S. Старался вносить в код минимум изменений, так конечно он требует рефакторинга. И есть огромное поле для развития с количеством воркеров и наполнением каналов и т.д.
}
