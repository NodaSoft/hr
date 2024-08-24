package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/golang-cz/devslog"
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

// * creating tasks
func taskCreator(superChan chan Ttype) {
	go func() {
		timer := time.NewTimer(10 * time.Second)
		defer close(superChan)

		for {
			select {
			case <-timer.C:
				return
			default:
				ct := time.Now().Format(time.RFC3339)
				// ? keep it ?
				if (time.Now().Nanosecond()/1000)%2 > 0 { // вот такое условие появления ошибочных тасков
					ct = "Some error occured"
				}
				task := Ttype{cT: ct, id: int(time.Now().UnixMilli())}
				log.Debug("task created", slog.Int("task", task.id))
				superChan <- task // передаем таск на выполнение
				time.Sleep(time.Millisecond * 50)
			}
		}
	}()
}

// * worker
func task_worker(a Ttype) Ttype {
	tt, _ := time.Parse(time.RFC3339, a.cT)
	// ? keep it
	if tt.After(time.Now().Add(-20 * time.Second)) {
		a.taskRESULT = []byte("task has been successed")
	} else {
		a.taskRESULT = []byte("something went wrong")
	}
	a.fT = time.Now().Format(time.RFC3339Nano)
	time.Sleep(time.Millisecond * 150)
	return a
}

// * sorter
func tasksorter(t Ttype, doneTasks chan Ttype, undoneTasks chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	if string(t.taskRESULT) == "task has been successed" {
		//log.Debug("task sorted 'success'", slog.Int("task", t.id))
		doneTasks <- t
	} else {
		//log.Debug("task sorted 'error'", slog.Int("task", t.id))
		undoneTasks <- fmt.Errorf("task id %d time %s, error %s", t.id, t.cT, t.taskRESULT)
	}
}

var log *slog.Logger

func main() {
	// init logger
	log = InitLogger()

	superChan := make(chan Ttype, 10)
	doneTasks := make(chan Ttype)
	undoneTasks := make(chan error)

	// generating tasks
	log.Info("starting task generating")
	go taskCreator(superChan)

	// getting tasks & passing them for sorting
	go func() {
		defer close(doneTasks)
		defer close(undoneTasks)
		wg := &sync.WaitGroup{}
		for t := range superChan {
			wg.Add(1)
			complited := task_worker(t)
			go tasksorter(complited, doneTasks, undoneTasks, wg)
		}
		wg.Wait()
	}()

	// filling up tasks from doneTasks & undoneTasks
	result := map[int]Ttype{}
	err := []string{}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	mu := &sync.Mutex{}
	go func() {
		defer wg.Done()
		for r := range doneTasks {
			mu.Lock()
			result[r.id] = r
			mu.Unlock()
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		for r := range undoneTasks {
			err = append(err, r.Error())
		}
	}()

	// sending logs every 3s
	ticker := time.NewTicker(3 * time.Second)
	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				log.Info("Complited tasks", slog.Any("res", result))
				log.Error("Errors", err)
			}
		}
	}()
	wg.Wait()
	done <- struct{}{}
	log.Info("Complited tasks", slog.Any("res", result))
	log.Error("Errors", slog.Any("err", err))
}

// setting logger
func InitLogger() *slog.Logger {
	loggerPtr := flag.String("logger", "", "Specify the logger type")
	flag.Parse()

	var handler slog.Handler
	switch *loggerPtr {
	// debug lvl
	case "devslog_debug":
		slogOpts := &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		}
		opts := &devslog.Options{
			HandlerOptions:    slogOpts,
			MaxSlicePrintSize: 4,
			SortKeys:          true,
			TimeFormat:        "[04:05]",
			NewLineAfterLog:   true,
			DebugColor:        devslog.Magenta,
		}
		handler = devslog.NewHandler(os.Stdout, opts)
	// info lvl
	case "devslog":
		slogOpts := &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelInfo,
		}
		opts := &devslog.Options{
			HandlerOptions:    slogOpts,
			MaxSlicePrintSize: 4,
			SortKeys:          true,
			TimeFormat:        "[04:05]",
			NewLineAfterLog:   true,
		}
		handler = devslog.NewHandler(os.Stdout, opts)
	// info lvl
	default:
		handler = slog.NewTextHandler(os.Stdout, nil)
	}
	logger := slog.New(handler)
	return logger
}
