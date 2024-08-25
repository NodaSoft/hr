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

// Мы даем тестовое задание чтобы:
// * уменьшить время технического собеседования - лучше вы потратите пару часов в спокойной домашней обстановке, чем будете волноваться, решая задачи под взором наших ребят;
// * увеличить вероятность прохождения испытательного срока - видя сразу стиль и качество кода, мы можем быть больше уверены в выборе;
// * снизить число коротких собеседований, когда мы отказываем сразу же.

// Выполнение тестового задания не гарантирует приглашение на собеседование, т.к. кроме качества выполнения тестового задания, оцениваются и другие показатели вас как кандидата.

// Мы не даем комментариев по результатам тестового задания. Если в случае отказа вам нужен наш комментарий по результатам тестового задания, то просим об этом написать вместе с откликом.

// A Ttype represents a meaninglessness of our life
type Ttype struct {
	id         int
	cT         string // время создания
	fT         string // время выполнения
	taskRESULT []byte
}

func taskCreator(a chan Ttype, timeout context.Context) {
out:
	for {
		cT_ := time.Now().Format(time.RFC3339)
		// Nanoseconds возвращает кол-во миллисекунд * 1000. Поменял на UnixMilli
		if time.Now().UnixMilli()%2 != 0 { // вот такое условие появления ошибочных тасков
			cT_ = "Some error occured"
		}
		select {
		case <-timeout.Done():
			close(a)
			break out
		default:
		}
		//много одинаковых id. Поменял Unix на UnixMilli
		a <- Ttype{id: int(time.Now().UnixMilli()), cT: cT_} // передаем таск на выполнение
		time.Sleep(100 * time.Millisecond)
	}
}

func taskWorker(a *Ttype) {
	_, err := time.Parse(time.RFC3339, a.cT)
	if err == nil {
		a.taskRESULT = []byte("task has been successed")
	} else {
		a.taskRESULT = []byte("something went wrong")
	}
	a.fT = time.Now().Format(time.RFC3339Nano)

	time.Sleep(150 * time.Millisecond)
}

func taskSorter(t Ttype, doneTasks chan<- Ttype, undoneTasks chan<- error) {
	if string(t.taskRESULT[14:]) == "successed" {
		doneTasks <- t
	} else {
		undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT, t.taskRESULT)
	}
}

func main() {
	superChan := make(chan Ttype, 10)
	timeout, _ := context.WithTimeout(context.Background(), 10*time.Second)

	go taskCreator(superChan, timeout)

	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)

	go func() {
		// получение тасков
		for t := range superChan {
			go func(t_ Ttype) {
				taskWorker(&t_)
				taskSorter(t_, doneTasks, undoneTasks)
			}(t)
		}
		//close(doneTasks)
		//close(undoneTasks)
	}()

	result := map[int]Ttype{}
	var err []error

	go func() {
		for r := range doneTasks {
			result[r.id] = r
		}
	}()
	go func() {
		for r := range undoneTasks {
			err = append(err, r)
		}
	}()

finish:
	for {
		time.Sleep(3 * time.Second)
		select {
		case <-timeout.Done():
			break finish
		default:
		}
		fmt.Println("Done")
		for k, v := range result {
			fmt.Println(k, v)
		}
		fmt.Println("Errors")
		for _, e := range err {
			fmt.Println(e)
		}
	}
}
