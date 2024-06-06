package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
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

type Status int8

const (
	Processing Status = iota
	Success
	Error
)

// A Ttype represents a meaninglessness of our life
type Ttype struct {
	Id          int
	Status      Status    //статус
	Description string    //описание задачи
	CreatedAt   time.Time // время создания
	DoneAt      time.Time // время выполнения
}

func main() {

	superChan := make(chan Ttype, 10)
	//в Go стандартный вывод медленный, поэтому принято решение использовать буферизированный
	out := bufio.NewWriter(os.Stdout)
	//управление временем жизни программы через контекст
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	go taskProducer(ctx, superChan)

	doneTasks := make(chan Ttype)
	undoneTasks := make(chan Ttype)
	outChan := make(chan Ttype)

	go taskWorker(ctx, superChan, outChan)
	go taskSorter(ctx, outChan, doneTasks, undoneTasks)
	//можно было сделать обработку и на основе мапы, но на каналах более идиоматично

	go taskConsumer(doneTasks, out, "success-")
	go taskConsumer(undoneTasks, out, "error-")
	t := time.NewTicker(time.Second * 3)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("graceful shutdown")
			return
		case <-t.C:
			//каждые три секунды высвобождается поток вывода
			out.Flush()
		}
	}

}

func taskProducer(ctx context.Context, a chan Ttype) {
	defer close(a)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			createdAt := time.Now()
			s := Processing
			if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
				s = Error
			}
			a <- Ttype{Status: s, CreatedAt: createdAt, Id: int(time.Now().Unix())} // передаем таск на выполнение
			time.Sleep(time.Millisecond * 50)
		}

	}
}

func taskWorker(ctx context.Context, in <-chan Ttype, out chan<- Ttype) {
	defer close(out)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			for a := range in {

				if a.Status == Processing && a.CreatedAt.After(time.Now().Add(-20*time.Second)) {
					a.Description = "task has been successed"
					a.Status = Success
				} else {
					a.Status = Error
					a.Description = "something went wrong"
				}
				a.DoneAt = time.Now()
				out <- a

			}
		}
	}
}

func taskSorter(ctx context.Context, allTasks, doneTasks, undoneTasks chan Ttype) {
	defer close(doneTasks)
	defer close(undoneTasks)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			for t := range allTasks {
				if t.Status == Success {
					doneTasks <- t
				} else {
					undoneTasks <- t
				}
			}
		}
	}

}

func taskConsumer(tasks chan Ttype, out *bufio.Writer, taskDescription string) {
	for t := range tasks {
		fmt.Fprintln(out, taskDescription, t.Id, t.Description)
	}

}
