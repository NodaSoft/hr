package main

import (
	"context"
	"errors"
	"fmt"
	"hash/maphash"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	unDone = iota
	done
)

var (
	ErrCreateTask = errors.New("error while creating")
)

// если совсем упарываться можно выравнять поля структуры, чтобы она лучше вставала в память

// A TaskType represents happiness of being alive
type TaskType struct {
	id           int
	createTime   int64 // лучше хранить в юникс
	completeTime int64 // лучше хранить в юникс
	err          error
	resultCode int
	resultMsg []byte
}

func main() {
	taskChan := make(chan TaskType, 10)
	processedTask := make(chan TaskType, 10)
	ctx, cancel := context.WithCancel(context.Background())
	// по-хорошему нужно было бы разделить группы ожидания обработчика и воркера, но в нашем случае остави одну
	unionWg := new(sync.WaitGroup)
	go taskCreater(ctx, taskChan)
	go func() {
		for task := range taskChan {
			// если хотим многопоточную обработку, то результаты воркер должен в канал заносить
			unionWg.Add(1)
			go taskWorker(context.Background(), unionWg, processedTask, task)
		}
	}()

	// будем считать, что нам важен порядок: сначала вывести успешно выполненные, потом ошибки, тогда оставляем две мапы
	// для простоты будем использовать мапы из пакета синк, они медленней, чем использование обычного мьютекса рядом, но
	// так меньше писать :)
	var successTask, unsuccessTask sync.Map

	go func() {
		for task := range processedTask {
			unionWg.Add(1)
			// сортировку перенесем сюда, хоть работа всего в отдельных горутинах это один из плюсов Golang, но бесконечный перенос данных по каналам
			// может вызвать замедление, как минимум из-за того, что данные копируются при попадании в буфер, а не передаются напрямую
			// если логика сортировки будет намного больше, чем task.resultCode == unDone, то лучше, конечно, выделить новую горутину
			// с двумя каналами: успешные задачи, неуспешные
			if task.resultCode == unDone {
				// хоть с версии 1.21 или 1.20 Go может сам корректно передать текущую переменную цикла,
				// сделаем по старинке и укажем в аргментах анонимки таску
				go func (task TaskType)  {
					defer unionWg.Done()
					unsuccessTask.Store(task.id, task)
				}(task)
				continue
			}
			go func (task TaskType)  {
				defer unionWg.Done()
				successTask.Store(task.id, task)
			}(task)
		}
	}()
	// если работаем с многопоточной обработкой, то лучше сразу обзавестись безопасной остановкой
	stopChecker := make(chan os.Signal, 1)
	signal.Notify(stopChecker, syscall.SIGTERM, syscall.SIGINT)
	select {
	case <- stopChecker:
	case <- time.After(time.Second * 2):
	}
	log.Println("stopping program")
		cancel()
		close(taskChan)
		unionWg.Wait()
		close(processedTask)
	// как уже говорилось, будем считать, что вывод нужен именно такой, да и не думаю, что это важно в этом задании
	// и если все-таки порядок важен, да хранит господь, того, кто это запустит на мощном железе, и нужно будет вывести это все...
	println("Done tasks:")
	successTask.Range(func(key, _ any) bool {
		id := key.(int)
		println(id)
		return true
	})
	println("Errors:")
	unsuccessTask.Range(func(_, value any) bool {
		task := value.(TaskType)
		println(fmt.Errorf("task id %d, creation time %d, error %s", task.id, task.createTime, task.resultMsg))
		return true
	})
}

func taskCreater(ctx context.Context, taskCh chan TaskType) {
	// Возможно могут возникнуть ситуации, когда канал уже закроют, но мы не успеем обработать закрытие
	// и тогда лучше будет использовать какую-то такую конструкцию
	// и при записи в канал проверять на closed
	// var closed bool
	// go func () {
	// 	select {
	// 	case <- ctx.Done():
	// 		closed = true
	// 	}
	// }()

WORKER:
	for {
		select {
		case <- ctx.Done():
			break WORKER
		default:
			// такая генерация id самая дружелюбная для мультипотока
			outUint64 := new(maphash.Hash).Sum64()
			id := int(outUint64)
			if id < 0 {
				id = -id
			}
			var err error
			taskCT := time.Now()
			ns := taskCT.Nanosecond()
			if ns%2 > 0 {
				err = ErrCreateTask
			}
			taskCh <- TaskType{createTime: taskCT.Unix(), id: id, err: err} // передаем таск на выполнение
		}
	}
}

// лучше подготовить место для передачи контекста там, где он может понадобиться, чем потом мучаться и менять везде интерфейсы
func taskWorker(ctx context.Context, wg *sync.WaitGroup, outputTaskCh chan TaskType, task TaskType) {
	defer wg.Done()
	taskCT := time.Unix(task.createTime, 0)
	// можно обыграть лучше, но сделаем без дальнейшей проверки, если имеется ошибка
	if task.err != nil {
		task.resultMsg = []byte("something went wrong")
		task.resultCode = unDone
	} else if taskCT.After(time.Now().Add(-20 * time.Second)) {
		task.resultMsg = []byte("task has been successed")
		task.resultCode = done
	}

	task.completeTime = time.Now().Unix()

	time.Sleep(time.Millisecond * 150)
	outputTaskCh <- task
}
