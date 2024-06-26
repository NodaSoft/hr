package main

import (
	"context"
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
type Ttype struct {
	id         int
	cT         string // время создания
	fT         string // время выполнения
	taskRESULT []byte
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	some := make(chan struct{})

	// taskCreturer creates tasks
	taskCreturer := func(a chan Ttype, ch chan struct{}) {
		go func() {
			for {
				select {
				case <-ch:
					fmt.Println("stop task generator")
					return
				default:
					ft := time.Now().Format(time.RFC3339)
					if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
						ft = "Some error occured"
					}
					a <- Ttype{cT: ft, id: int(time.Now().Unix())} // передаем таск на выполнение

				}
			}
		}()
	}

	superChan := make(chan Ttype, 10)

	go taskCreturer(superChan, some)

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

	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)

	tasksorter := func(t Ttype) {
		if string(t.taskRESULT[14:]) == "successed" {
			doneTasks <- t
		} else {
			undoneTasks <- fmt.Errorf("task id %d time %s, error %s", t.id, t.cT, t.taskRESULT)
		}
	}

	go func() {
		// получение тасков
		for t := range superChan {
			t = task_worker(t)
			go tasksorter(t)
		}
		close(superChan)
	}()

	result := map[int]Ttype{}
	err := []error{}
	ticker := time.NewTicker(time.Second * 3)

	// read tasks
	for {
		select {
		case t := <-doneTasks:
			result[t.id] = t
		case e := <-undoneTasks:
			err = append(err, e)
		case <-ctx.Done():
			some <- struct{}{}
			ticker.Stop()
			fmt.Println("context done")
			return
		case <-ticker.C:
			// вывод всех обработанных к этому моменту тасков (накопительный результат)
			fmt.Println("Result:")
			for _, v := range result {
				println(string(v.taskRESULT))
			}

			fmt.Println("Errors:")
			for _, v := range err {
				println(v.Error())
			}
		}
	}

}
