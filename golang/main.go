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

// A Ttype represents a meaninglessness of our life
type Ttype struct {
	id         int
	cT         string // время создания
	fT         string // время выполнения
	taskRESULT []byte
}

const (
	workTimeSeconds  = 10
	printTimeSeconds = 3
)

func main() {
	// таймеры времени
	workTime := time.After(time.Second * workTimeSeconds)
	printTimer := time.NewTimer(time.Second * printTimeSeconds)

	// канал нужен чтобы корректно завершать программу после истечения времени и закрытия результирующего канала
	gracefulShutdownChan := make(chan struct{})

	// ассинхронное создание тасок
	taskCreturer := func(a chan Ttype) {
		go func() {
			for {
				select {
				case <-workTime:
					// по истечению времени закрываем канал и передаем сигнал о завершении работы программы
					close(a)
					gracefulShutdownChan <- struct{}{}
					return
				default:
					ft := time.Now().Format(time.RFC3339)
					if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
						ft = "Some error occured"
					}
					a <- Ttype{cT: ft, id: int(time.Now().Unix())} // передаем таск на выполнени
				}
			}
		}()
	}

	superChan := make(chan Ttype, 10)

	go taskCreturer(superChan)

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

	tasksorter := func(wg *sync.WaitGroup, t Ttype) {
		defer wg.Done()

		if string(t.taskRESULT[14:]) == "successed" {
			doneTasks <- t
		} else {
			undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT, t.taskRESULT)
		}
	}

	go func() {
		wg := &sync.WaitGroup{}

		// получение тасков
		for t := range superChan {
			wg.Add(1)
			go func(t Ttype) {
				t = task_worker(t)
				go tasksorter(wg, t)
			}(t)
		}

		// ждём закрытия канала, генерирующего таски
		// также ожидаем завершения работы горутин на запись в каналы результата и ошибок
		wg.Wait()
		close(doneTasks)
		close(undoneTasks)
	}()

	result := map[int]Ttype{}
	mux := &sync.Mutex{}

	err := []error{}

	// запись в буффер
	go func() {

		go func() {
			for r := range doneTasks {
				go func(r Ttype) {
					mux.Lock()
					result[r.id] = r
					mux.Unlock()
				}(r)
			}
		}()

		go func() {
			for r := range undoneTasks {
				go func(r error) {
					err = append(err, r)
				}(r)
			}
		}()

	}()

	for {
		select {
		case <-gracefulShutdownChan:
			return
		case <-printTimer.C:
			printTimer.Reset(time.Second * printTimeSeconds)
			fmt.Println("Errors:")
			for r := range err {
				fmt.Println(r)
			}

			fmt.Println("Done tasks:")
			for r := range result {
				fmt.Println(r)
			}
		}
	}
}
