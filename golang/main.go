package main

import (
	"context"
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

// taskCreturer отправляет созданную задачу Ttype в канал a
func taskCreturer(ctx context.Context, ch chan Ttype) {
	for {
		ft := time.Now().Format(time.RFC3339)
		if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
			ft = "Some error occured"
		}
		select {
		case ch <- Ttype{cT: ft, id: int(time.Now().Unix())}:
			continue
		case <-ctx.Done():
			close(ch)
			return
		}
	}
}

// taskWorker выставляет значение taskRESULT для экземпляра структуры Ttype, возвращает Ttype
func taskWorker(a Ttype) Ttype {
	tt, _ := time.Parse(time.RFC3339, a.cT)
	if tt.After(time.Now().Add(-20 * time.Second)) {
		a.taskRESULT = []byte("task has been successed")
	} else {
		a.taskRESULT = []byte("something went wrong")
	}
	a.fT = time.Now().Format(time.RFC3339Nano)

	return a
}

// taskSorter отправляет экземпляр Ttype в канал для выполненных задач, либо отправляет ошибку в канал для ошибок
func taskSorter(t Ttype, doneTasks chan Ttype, undoneTasks chan error) {
	if string(t.taskRESULT[14:]) == "successed" {
		doneTasks <- t
	} else {
		undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT, t.taskRESULT)
	}
}

func main() {

	var muRes sync.RWMutex // Mutex для мапы с результатами
	var muErr sync.RWMutex // Mutex для слайса с ошибками

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Каналы
	superChan := make(chan Ttype, 10)
	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)
	done := make(chan struct{}) // Канал для отправки сигнала о завершение работы горутин

	go taskCreturer(ctx, superChan)

	var wgSort sync.WaitGroup // Ожидание сортировки
	wgSort.Add(1)

	//
	go func() {
		defer wgSort.Done()
		for t := range superChan {
			t = taskWorker(t)
			taskSorter(t, doneTasks, undoneTasks)
		}
	}()

	go func() {
		// Ждём пока все таски будут отсортированы
		wgSort.Wait()
		close(doneTasks)
		close(undoneTasks)

	}()

	result := map[int]Ttype{}
	errRes := []error{}

	// Распределение задач после их выполнения
	processTasks := func() {
		resOk := true
		errOk := true
		for {
			select {
			case res, ok := <-doneTasks:
				if !resOk && !errOk {
					done <- struct{}{}
					return
				} else if !ok {
					resOk = ok
					break
				}

				muRes.Lock()
				result[res.id] = res
				muRes.Unlock()

			case err, ok := <-undoneTasks:
				if !resOk && !errOk {
					done <- struct{}{}
					return
				} else if !ok {
					errOk = ok
					break
				}
				muErr.Lock()
				errRes = append(errRes, err)
				muErr.Unlock()
			}
		}
	}
	go processTasks()

	// Вывод результатов каждые 3 секунды. Если горутины закончили работу, после вывода результата завершаем функцию.
	for {
		time.Sleep(time.Second * 3)

		println("Errors:")

		muErr.RLock()
		for i := range errRes {
			fmt.Println(errRes[i])
		}
		muErr.RUnlock()

		println("Done tasks:")
		muRes.RLock()
		for i := range result {
			fmt.Println(result[i])
		}
		muRes.RUnlock()

		select {
		case <-done:
			return
		default:
			continue
		}
	}
}
