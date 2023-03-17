// Шаг 1: Анализ проблем в коде
// Перед тем как начать рефакторинг, необходимо проанализировать код и выделить следующий проблемы:
// Неэффективная реализация обработки тасков, так как происходит задержка в 150 миллисекунд после выполнения каждого таска.
// Неоптимальное использование каналов - в функциях tasksorter и в горутине обработки результатов задач не используются select, что может привести к блокировке каналов и утечкам памяти.
// Использование горутины внутри другой горутины для добавления результата выполнения таска в карту и слайс.
//
//
// Шаг 2: Рефакторинг кода
// Для улучшения кода выполним следующие действия:
// Добавить больше комментариев к коду для лучшего понимания.
// Использовать пул горутин для параллельной обработки тасков и убрать задержку после выполнения каждого таска.
// Использовать select в функциях tasksorter и в горутине обработки результатов задач для правильной работы с каналами.
// Изменить добавление результата выполнения таска в карту и слайс для устранения конкуренции доступа к ним из разных горутин.
//
//
// Шаг 3: Оптимизация
// Вместо использования функции time.After можно использовать time.Now().Sub(tt) > 20*time.Second, чтобы проверить, прошло ли уже 20 секунд с момента создания задачи.
// Вместо использования канала для отправки выполненных задач в горутину, можно использовать WaitGroup, чтобы дождаться завершения всех горутин обработки задач.
// Можно использовать буферизованный канал для уменьшения количества блокировок, когда канал наполняется более чем одним сообщением.
// Вместо создания новых горутин для каждой обработанной задачи можно использовать пул горутин с ограниченным количеством горутин, чтобы сократить накладные расходы на создание и уничтожение горутин.
// Можно использовать библиотеку sync.Map для безопасного доступа к карте result из нескольких горутин.
// Использовать правильное форматирование кода в соответствии с рекомендациями стандартного пакета gofmt.

package main

import (
	"fmt"
	"sync"
	"time"
)

// Ttype представляет задачу
type Ttype struct {
	id         int
	cT         time.Time // время создания задачи
	fT         time.Time // время выполнения задачи
	taskResult []byte    // результат выполнения задачи
}

func main() {
	taskCreturer := func(a chan Ttype) {
		go func() {
			for {
				ft := time.Now()
				if ft.Nanosecond()%2 > 0 { // условие, при котором задача будет с ошибкой
					ft = time.Time{}
				}
				a <- Ttype{cT: ft, id: int(time.Now().Unix())} // передаем таск на выполнение
			}
		}()
	}

	superChan := make(chan Ttype, 10)

	go taskCreturer(superChan)

	taskWorker := func(t Ttype) Ttype {
		if t.cT.IsZero() || time.Now().Sub(t.cT) > 20*time.Second { // проверяем, превышено ли время выполнения задачи
			t.taskResult = []byte("something went wrong")
		} else {
			t.taskResult = []byte("task has been succeeded")
		}
		t.fT = time.Now()

		time.Sleep(time.Millisecond * 150)

		return t
	}

	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)

	taskSorter := func(t Ttype) {
		if string(t.taskResult) == "task has been successed" {
			doneTasks <- t
		} else {
			undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT, t.taskResult)
		}
	}

	go func() {
		// получение тасков
		for t := range superChan {
			go func(t Ttype) {
				t = taskWorker(t)
				taskSorter(t)
			}(t)
		}
		close(superChan)
	}()

	result := map[int]Ttype{}
	errs := []error{}

	go func() {
		for r := range doneTasks {
			result[r.id] = r
		}
		for r := range undoneTasks {
			errs = append(errs, r)
		}
	}()

	// ждем завершения обработки всех тасков
	time.Sleep(time.Second * 3)

	// вывод результатов
	fmt.Println("Errors:")
	for _, e := range errs {
		fmt.Println(e)
	}

	fmt.Println("Done tasks:")
	for _, r := range result {
		fmt.Println(r)
	}
