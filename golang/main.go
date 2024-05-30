package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Выполнил email: pircuser61@rambler.ru  tg: @pircuser61

// Приложение эмулирует получение и обработку неких тасков. Пытается и получать, и обрабатывать в многопоточном режиме.
// Приложение должно генерировать таски 10 сек. Каждые 3 секунды должно выводить в консоль результат всех обработанных к этому моменту тасков (отдельно успешные и отдельно с ошибками).

// ЗАДАНИЕ: сделать из плохого кода хороший и рабочий - as best as you can.
// Важно сохранить логику появления ошибочных тасков.
// Важно оставить асинхронные генерацию и обработку тасков.
// Сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через pull-request в github
// Как видите, никаких привязок к внешним сервисам нет - полный карт-бланш на модификацию кода.

// A Task represents a meaninglessness of our life
type Task struct {
	id       int
	created  string // время создания
	finished string // время выполнения
	result   string
}

type Worker struct {
	id       int
	mxDone   sync.Mutex
	mxUnDone sync.Mutex
	done     []Task // нет особого смысла писать всем воркерам общий слайс
	undone   []Task // что бы использовать буферизованный канал нужно заранее знать размер буфера
	// иначе воркерам придется ждать пока принтер не заберет таски
	// размер в общем случае неизвестен, хотя можно было взять 3000 / 150 ~~ 20 по каналу на воркер
}

const taskChanSize = 10
const workerCount = 2
const workDelay = time.Millisecond * 150
const printInterval = time.Second * 3
const expirationInterval = -20 * time.Second // -1 для воспроизведения ошибки
const errDivisor = 2                         // 200 для воспроизведения ошибки

func taskCreator(ctx context.Context, stop chan<- struct{}, a chan Task) {
	defer func() {
		slog.Info("creator stopped")
		stop <- struct{}{}
	}()
	var taskId int
	for ctx.Err() == nil {
		ft := time.Now().Format(time.RFC3339)
		//taskId = int(time.Now().Unix())  Id пересакаются
		taskId++
		if time.Now().Nanosecond()%errDivisor > 0 { // вот такое условие появления ошибочных тасков
			ft = "Some error occured"
		}

		slog.Debug("new task", slog.Int("id", taskId))
		a <- Task{created: ft, id: taskId} // передаем таск на выполнение
	}
}

func taskWorker(wg *sync.WaitGroup, tasks chan Task, w *Worker) {
	defer wg.Done()
	for {

		task, ok := <-tasks
		if !ok {
			slog.Info("worker stopped", slog.Int("Id", w.id))
			return
		}
		slog.Debug("task in work", slog.Int("worker", w.id), slog.Int("id", task.id))
		err := taskHandle(task)
		task.finished = time.Now().Format(time.RFC3339Nano)
		if err != nil {
			task.result = err.Error()
		} else {
			task.result = "task has been successed"
		}
		if err == nil {
			w.mxDone.Lock()
			w.done = append(w.done, task)
			w.mxDone.Unlock()
		} else {
			w.mxUnDone.Lock()
			w.undone = append(w.undone, task)
			w.mxUnDone.Unlock()
		}

	}
}

func taskHandle(a Task) error {
	tt, err := time.Parse(time.RFC3339, a.created)
	if err == nil {
		expiration := time.Now().Add(expirationInterval)
		if !tt.After(expiration) {
			err = errors.New("something went wrong (expiration)")
		}
	}
	time.Sleep(workDelay)
	return err
}

func taskPrinter(out io.Writer, stopChan <-chan struct{}, doneChan chan<- struct{},
	workers *[workerCount]Worker) {
	defer func() {
		slog.Info("printer stopped")
		doneChan <- struct{}{}
	}()

	ticker := time.Tick(printInterval)

	var hasData bool
	ok := true
_mainLoop:
	for ok {
		select {
		case <-stopChan:
			ok = false
		case <-ticker:
		}
		slog.Debug("==== Print", slog.Bool("OK", ok))
		//  нужно ли выводить пустые списки ...
		hasData = false
		for i := 0; i < workerCount && !hasData; i++ {
			hasData = len(workers[i].done) > 0 || len(workers[i].undone) > 0
		}
		if !hasData {
			continue _mainLoop
		}

		fmt.Fprintln(out, "Errors:")
		for i := 0; i < workerCount; i++ {
			workers[i].mxUnDone.Lock()
			for _, r := range workers[i].undone {
				fmt.Fprintf(out, "Task id %d time %s, error %s\n", r.id, r.created, r.result)
			}

			workers[i].undone = workers[i].undone[:0]
			workers[i].mxUnDone.Unlock()
		}

		fmt.Fprintln(out, "Done tasks:")
		for i := 0; i < workerCount; i++ {
			workers[i].mxDone.Lock()
			for _, r := range workers[i].done {
				fmt.Fprintln(out, r)
			}
			workers[i].done = workers[i].done[:0]
			workers[i].mxDone.Unlock()
		}
	}
}

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(log)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	superChan := make(chan Task, taskChanSize)

	wg := sync.WaitGroup{}
	wg.Add(workerCount)

	stopCreate := make(chan struct{})
	go taskCreator(ctx, stopCreate, superChan)

	var workers [workerCount]Worker
	for i := 0; i < workerCount; i++ {
		workers[i].id = i
		go taskWorker(&wg, superChan, &workers[i])
	}

	stopPrint := make(chan struct{})
	donePrint := make(chan struct{})

	go taskPrinter(os.Stdout, stopPrint, donePrint, &workers)

	<-stopCreate
	slog.Info("creater stopped, wait workers...")
	close(superChan)
	wg.Wait()

	slog.Info("workers stopped, wait for print...")
	stopPrint <- struct{}{}
	<-donePrint

	slog.Info("Done")
}
