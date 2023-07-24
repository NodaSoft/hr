package main

import (
	"context"
	"fmt"
	"sync"
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

// A Task represents a meaninglessness of our life
type Task struct {
	id           uint
	createdTime  time.Time // время создания
	finishedTime time.Time // время выполнения
	// результат тасков бинарный, поэтому сделал result как bool. Если надо будет байты - можно сделать байтами, но
	// для сохранения логики работы именно этого теста - байты не нужны
	result  bool
	success bool // true - значит умер, таска провалена
	mu      sync.Mutex
}

type Error struct {
	id          uint
	createdTime time.Time // время создания
	result      bool
	mu          sync.Mutex
}

func taskProducer(n int, ctx context.Context) <-chan Task {
	// предположим что мы чаще будем компилировать под 64-разрядную платформу
	// а также предположим, что у нас всегда один сервис, иначе нам надо использовать GUID в качестве id,
	// либо сторонний генертор id - например, базу данных
	var id uint64 = 0
	ch := make(chan Task, n)
	go func() {
		defer close(ch) // закрываем в функции, где открывали

		for {
			select {
			case <-ctx.Done(): // выйдем если прервут через контектс
				return
			default:
				t := time.Now()
				// логическое не - т.к. статус от лучше назвать от успеха, а не от ошибок
				success := !(t.Nanosecond()%2 == 0) // вот такое условие появления ошибочных тасков
				// передаем таск на выполнение
				// id - атомик, т.к. он инекрементится асиинхронно
				ch <- Task{createdTime: t, id: uint(atomic.AddUint64(&id, 1)), success: success}
			}
		}
	}()
	return ch
}

func (a *Task) taskWorker() {
	tt := a.createdTime
	a.result = tt.After(time.Now().Add(-20 * time.Second))
	time.Sleep(time.Millisecond * 150) // перенес сюда сон, как часть воркера, до вычисления времени окончания. Инача его надо удались совсем
	a.finishedTime = time.Now()
}

func main() {
	n := 10 // размер буфера в канал и в пуле воркеров
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	defer cancelFunc()

	superChan := taskProducer(n, ctx)                                            // 1. создаеем основной канал
	doneTasks, undoneTasks := processTasks(superChan, n, ctx)                    // 2. pipeline - канал на вход, каналы на выход, внутри пул воркеров
	resultOk, resultErr, waiter := collectFromChans(doneTasks, undoneTasks, ctx) // 3. pipeline - собираем из каналоов в структуры

	<-waiter // вместо time.Sleep(time.Second * 3)

	println("Errors:")
	for _, r := range resultErr {
		// форматирование - как можно ближе к потребителю, т.к. потребителей может быть несколько и у каждого
		// свое форматирование
		println(fmt.Sprintf("Task id %d time %s, error %s", r.id, r.createdTime, "something went wrong"))
	}

	println("Done tasks:")
	for r := range resultOk {
		println(r)
	}
}

func collectFromChans(doneTasks chan Task, undoneTasks chan Error, ctx context.Context) (map[uint]Task, []Error, chan bool) {
	resultOk := map[uint]Task{}
	resultErr := []Error{}

	// соберем из каналов к структуры. Для этого заблокирум этим каналом, пока не будет завршена работа в каналох
	waiter := make(chan bool, 1)
	go func() {
		for doneTasks != nil || undoneTasks != nil { // более-менее станндарный подход для выхода из чтения из N каналов
			select {
			case <-ctx.Done(): // устанавливаем условие для выхода, если контекст завершен
				doneTasks = nil
				undoneTasks = nil
				break
			case r, ok := <-doneTasks:
				if ok {
					r.mu.Lock() // все основные структуры в Го не thread-safe, надо блокировать в горутин
					resultOk[r.id] = r
					r.mu.Unlock()
				} else {
					doneTasks = nil // обнулим канал, и больше в этот case не попадем, а переменную будем использовать чтобы выйти из цикла
				}

			case r, ok := <-undoneTasks:
				if ok {
					r.mu.Lock() // все основные структуры в Го не thread-safe, надо блокировать в горутин
					resultErr = append(resultErr, r)
					r.mu.Unlock()
				} else {
					undoneTasks = nil
				}
			}
		}
		waiter <- true // собрали всё в структуры - можно печатать
	}()
	return resultOk, resultErr, waiter
}

func processTasks(superChan <-chan Task, workPoolSize int, ctx context.Context) (chan Task, chan Error) {
	doneTasks := make(chan Task)
	undoneTasks := make(chan Error)

	tasksorter := func(t Task) {
		if t.result {
			doneTasks <- t
		} else {
			// форматирование как можно ближе к потребителю, т.к. таких потребителей может быть несколько типов
			undoneTasks <- Error{id: t.id, createdTime: t.createdTime, result: t.result}
		}
	}

	go func() {
		defer close(doneTasks) // закрывам поближе к тому месту, где открывали
		defer close(undoneTasks)

		wg := &sync.WaitGroup{}

		for i := 0; i < workPoolSize; i++ { // пул воркеров
			wg.Add(1)
			go func() {
				defer wg.Done()
				// получение тасков
				for t := range superChan {
					select {
					case <-ctx.Done():
						return
					default:
						t.taskWorker()
						tasksorter(t) // go было удалено, вместо этого введен пул горутин - цикл по workPoolSize
					}
				}

			}()
		}
		wg.Wait() // ждем пока вс воркеры закончать получать задания и завершатся
		// мы тут никого не блокируем (т.к. горутина) , а просто ждем когда можно будет закрыть каналы
	}()

	return doneTasks, undoneTasks
}
