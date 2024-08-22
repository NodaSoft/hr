package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// Приложение эмулирует получение и обработку неких тасков. Пытается и получать, и обрабатывать в многопоточном режиме.
// Приложение должно генерировать таски 10 сек.
// Каждые 3 секунды должно выводить в консоль результат всех обработанных к этому моменту тасков (отдельно успешные и отдельно с ошибками).

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
	taskRESULT string
}

func taskWorker(a Ttype) Ttype {
	_, err := time.Parse(time.RFC3339, a.cT)
	// если не менять в taskCreator логику создания ошибки,
	// то time.Parse не сможет парсить a.cT и вернет ошибку есть смысл по ней и выходить как мне кажется
	if err != nil {
		a.taskRESULT = "something went wrong"
		log.Println("Parsing Error")
		a.fT = time.Now().Format(time.RFC3339Nano)
		return a
	}
	a.taskRESULT = "task has been successed"
	a.fT = time.Now().Format(time.RFC3339Nano)
	time.Sleep(time.Millisecond * 151)
	return a
}

func main() {
	var wg sync.WaitGroup
	superChan := make(chan Ttype, 10)
	doneCh := make(chan struct{})

	// генерируем задачи 10 секунд
	ticker := time.NewTicker(10 * time.Second)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(superChan)
		// останавливаемся генерацию по тикеру
		for {
			select {
			case <-ticker.C:
				log.Println("Кончилось время генерации")
				ticker.Stop()
				return
			default:
				ft := time.Now().Format(time.RFC3339)
				if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
					ft = "Some error occured"
				}
				superChan <- Ttype{id: int(time.Now().Unix()), cT: ft} // передаем таск на выполнение
			}
		}
	}()

	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)

	taskSorter := func(t Ttype) {
		if t.taskRESULT == "task has been successed" {
			doneTasks <- t
		} else {
			undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT, t.taskRESULT)
		}
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		// получение тасков
		for t := range superChan {
			t = taskWorker(t)
			// если я правильно понял, то тут происходит обработка тасков и она должна быть асинхронной
			// оставил как есть, данные продолжат писаться в каналы вывода пока не кончатся в канале superChan
			go taskSorter(t)
		}
		log.Println("Кончились данные в superChan")
	}()

	result := map[int]Ttype{}
	err := []error{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		// такой подход позволит получить данные из канала superChan, обработать и сохранить, не должно остаться подвисших горутин
		for {
			select {
			case r := <-doneTasks:
				result[r.id] = r
			case r := <-undoneTasks:
				err = append(err, r)
			case <-doneCh:
				log.Println("Закрыт doneCh вышел из получения результатов")
				return
			}
		}
	}()

	// вывод данных каждые 3 секунды так же по тикеру
	tickerResult := time.NewTicker(3 * time.Second)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-tickerResult.C:
				log.Println("Errors:")
				for r := range err {
					log.Println("Error", r)
				}
				log.Println("Done tasks:")
				for r := range result {
					log.Println("Done", r)
				}
			case <-doneCh:
				log.Println("Закрыт doneCh вышел из печати результатов")
				tickerResult.Stop()
				return
			}
		}
	}()

	time.Sleep(12 * time.Second)

	close(doneCh)
	defer close(doneTasks)
	defer close(undoneTasks)
	wg.Wait()
}
