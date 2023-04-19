package main

import (
	"fmt"
	"log"
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
	taskCreturer := func(a chan Ttype) {
		// fvn: избыточный запуск отдельного потока: горутина в горутине
		go func() {
			i := 0
			// fvn: цикл эмуляции не имеет условий завершения, горутина не завершается. Только по концу работы main()
			for {
				ft := time.Now().Format(time.RFC3339)
				if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
					// fvn: ft -- строка в формате времени. Смешивание назначения строки - ошибки далее.
					ft = "Some error occured"
					log.Println("occured")
				}
				a <- Ttype{cT: ft, id: int(time.Now().Unix())} // передаем таск на выполнение
				// fvn: добавлено логирование работы
				i++
				log.Printf("cre=%d", i)
			}
		}()
	}

	// fvn: эмулятор работает шустрее, избыточность буферизации. Достаточно 1.
	superChan := make(chan Ttype, 10)

	// fvn: поток эмуляции тасок тут
	go taskCreturer(superChan)

	task_worker := func(a Ttype) Ttype {
		// fvn: не обработана ошибка, в поле может быть простой текст!
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

	// fvn: объявление канало должен делать писатель, он же должен их и закрывать
	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)

	tasksorter := func(t Ttype) {
		log.Println("sort")
		if string(t.taskRESULT[14:]) == "successed" {
			// fvn: запись во внешне закрываемый канал: потенциальная паника!
			doneTasks <- t
		} else {
			// fvn: аналогично..
			undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT, t.taskRESULT)
		}
	}

	go func() {
		// получение тасков
		for t := range superChan {
			log.Println("task")
			t = task_worker(t)
			// обработка тасков делается в этой горутине, а сортировка ответов - для каждой таски в своей. ТЗ - иное!
			go tasksorter(t)
		}
		// fvn: по сути сюда не попадаем, эмулятор создает таски быстрее, чем они обрабатываются
		// fvn: закрывать канал должен писатель..
		close(superChan)
		log.Println("closed super")
	}()

	result := map[int]Ttype{}
	err := []error{}
	go func() {
		// поток подсчета итогов
		log.Println("doned")
		// fvn: поскольку создание интенсивнее, этот канал никогда не пуст
		for r := range doneTasks {
			// fvn: бесполезная трата горутин.
			go func() {
				// fvn: map... -- не реентерабельна! Запись в множестве горутин..
				result[r.id] = r
			}()
		}
		log.Println("errored")
		// fvn: сюда в реальности никогда не попадем
		for r := range undoneTasks {
			// fvn: аналогичный избыток горутин
			go func() {
				// fvn: запись в общий слайс в горутине
				err = append(err, r)
			}()
		}
		close(doneTasks)
		close(undoneTasks)
		log.Println("closed doned,errored")
	}()

	time.Sleep(time.Second * 3)
	log.Println("timed 3s")

	println("Errors:")
	for r := range err {
		println(r)
	}

	println("Done tasks:")
	for r := range result {
		println(r)
	}
	// fvn: автоматическое закрытие горутин и каналов в состоянии "как есть", с необработанными данными.
}
