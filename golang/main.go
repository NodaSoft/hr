package main

import (
	"context"
	"fmt"
	"log"
	"sync"
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
	success    bool   // стейт лучше в поле держать
	taskRESULT []byte
}

func taskCreator(ctx context.Context, dest chan Ttype) {
	for aid := 0; ; aid++ {
		ft := time.Now().Format(time.RFC3339)
		success := true
		if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
			// ft = "Some error occured"
			success = false
		}
		select {
		case <-ctx.Done(): // finish loop if context closed
			log.Printf("finish task creator")
			return
			// time.Now().Unix() в качестве id не подходит - меняется раз в секунду, поменял на автоинкремент (aid)
		case dest <- Ttype{cT: ft, id: aid, success: success}: // передаем таск на выполнение
		}
	}
}

func handleTask(task Ttype) Ttype {
	tt, err := time.Parse(time.RFC3339, task.cT)
	if task.success && err == nil && tt.After(time.Now().Add(-20*time.Second)) {
		task.taskRESULT = []byte("task has been successed")
	} else {
		task.taskRESULT = []byte("something went wrong")
	}
	task.fT = time.Now().Format(time.RFC3339Nano)
	time.Sleep(time.Millisecond * 150)
	return task
}

func main() {
	superChan := make(chan Ttype, 10)
	ctx, cancelTaskCreator := context.WithTimeout(context.TODO(), time.Second*10) // или таймаут, или cancel() для прерывания вечного цикла

	creatorMutex := sync.Mutex{} // для ожидания заверешния taskCreator

	go func() {
		creatorMutex.Lock()
		defer creatorMutex.Unlock()
		taskCreator(ctx, superChan)
	}()

	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)

	tasksorter := func(task Ttype) { // замыкание здесь оправдано
		if task.success { // string(t.taskRESULT[14:]) == "successed" {
			doneTasks <- task
		} else {
			undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", task.id, task.cT, task.taskRESULT)
		}
	}

	workersWaitGroup := sync.WaitGroup{}
	
	go func() { // цикл получения тасков, при этом исходно все таски выполняются по очереди, параллелится только роутинг
		for task := range superChan {
			workersWaitGroup.Add(1)
			go func(otask Ttype) { // а теперь и таски параллельно (почти)
				defer workersWaitGroup.Done()
				otask = handleTask(otask)
				tasksorter(otask)
			}(task) // передаем как параметр для избегания ошибки "захват замыканием переменной цикла"
		}
		log.Printf("finish tasks receiver")
		// close(superChan) // comment reason: закрывать канал из читателя, который, к тому же заполняется вечным циклом - некорректно
	}()

	time.Sleep(time.Second * 3) // даем генератору тасков поработать 3 сек
	
	cancelTaskCreator() // cancel context after timeout

	result := sync.Map{} // хотя и убрали рейс, но для примера: в случае, 
						 // если надо из множества горутин обновлять мапу - есть sync.Map
						 // или закрыть mutex'ом как в примере с массивом ошибок ниже

	errLock := sync.Mutex{} // хоть меняем и в одном потоке, но если из замыкания-корутины меняем переменную - то лучше защищать мьютексом
	err := make([]error, 0)

	// обрабатываем результаты

	resultGroup := sync.WaitGroup{}	
	go func() {
		for r := range doneTasks {
			resultGroup.Add(1)
			go func(task Ttype) {
				defer resultGroup.Done()
				result.Store(task.id, task) // для примера работы с sync.Map{}
			}(r)
		}
		log.Printf("finish success tasks collector")
	}()

	resultGroup.Add(1) // сбор ошибок тоже надо подождать

	go func() {
		defer resultGroup.Done()

		errLock.Lock()				
		defer errLock.Unlock()

		for r := range undoneTasks {
			err = append(err, r)
		}

		log.Printf("finish failed tasks collector")
	}()

	workersWaitGroup.Wait() // все горутины запущены, ждем пока воркеры закончат работу

	creatorMutex.Lock() // гарантируем что creatorMutex завершился и не попытается записать в закрытый канал
	defer creatorMutex.Unlock() 
	close(superChan) // close superChan (цикл получения тасков должен завершиться, если не отменили ещё)

	close(doneTasks) // запись в каналы doneTasks и undoneTasks закончена, можно закрывать
	close(undoneTasks)

	resultGroup.Wait() // ждем пока закончат работу сборщики результатов

	println("Errors:")

	for _, r := range err {
		println(r.Error())
	}

	println("Done tasks:")
	counter := 0 // https://github.com/golang/go/issues/20680 sync.Map{} has no Length() method

	result.Range(func(key any, value any) bool {
		println(fmt.Sprintf("%d : %v", key, value))
		counter++
		return true
	})

	log.Printf("total errors: %d, total success: %d", len(err), counter)
}
