package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

const debug = false

// Простите за глобал, но логгер чисто для дебага
var log = func() *zap.Logger {
	lc := zap.NewDevelopmentConfig()
	lc.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	if debug {
		lc.Level.SetLevel(zap.DebugLevel)
	} else {
		lc.Level.SetLevel(zap.InfoLevel)
	}

	log, err := lc.Build()
	if err != nil {
		panic(err.Error())
	}
	return log
}()

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

// Приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме.
// В конце должно выводить успешные таски и ошибки при выполнении остальных тасков.

// A TimeTask represents a meaninglessness of our life
type TimeTask struct {
	id           int
	creationTime string // время создания
	finishTime   string // время выполнения
	taskResult   []byte
}

// spamTasks яростно спамит тасками в Processor.
func spamTasks(ctx context.Context, p *Processor[TimeTask]) {
	for {
		now := time.Now()
		ft := now.Format(time.RFC3339)
		if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
			ft = "Some error occured"
		}
		log.Debug("SPAM")
		// передаем таск на выполнение
		if err := p.AddTask(ctx, TimeTask{creationTime: ft, id: int(now.Unix())}); err != nil {
			// Ну эту ошибку я считай сам придумал, поэтому особо умно обрабатывать её не нужно.
			if errors.Is(err, ErrStopped) {
				return
			}
			if ctx.Err() != nil {
				return
			}
			// all other unexpected err
			log.Panic("unexpected error in spammer", zap.Error(err))
		}
	}
}

// StupidWorker пинает SmartWorker чтобы он работал, а потом тупит 150мс.
func StupidWorker(data TimeTask) (TimeTask, error) {
	log.Debug("WORK")
	res, err := smartWorker(data)

	// В любом случае это костыль чтобы проц не сгорел, поэтому без разницы куда его пихать.
	time.Sleep(time.Millisecond * 150)

	return res, err
}

// smartWorker описывает какую-то умную логику обработки задачи.
func smartWorker(t TimeTask) (_ TimeTask, err error) {
	// Запись времени окончания обработки даже в случае ошибки
	defer func() {
		t.finishTime = time.Now().Format(time.RFC3339Nano)
		// Я перенес форматирование ошибки в defer только ради того, чтобы эту ошибку можно было записать в taskResult.
		if err != nil {
			t.taskResult = []byte(err.Error())
			// не люблю кидать просто %s в строку, обычно использую %q, либо в формате прописываю кавычки.
			err = fmt.Errorf("Task id[%d], time[%s], error[%s]", t.id, t.creationTime, err)
		}
	}()

	tt, err := time.Parse(time.RFC3339, t.creationTime)
	if err != nil {
		return t, err
	}

	if tt.IsZero() || time.Since(tt) > 20*time.Second {
		return t, errors.New("something went wrong")
	}

	t.taskResult = []byte("task has been successed")
	return t, nil
}

func main() {
	var (
		p = NewProcessor[TimeTask](ProcessorParams[TimeTask]{
			RunnerFunc: StupidWorker,
			NumWorkers: 10,
			InputCap:   10,
			ResCap:     10,
			ErrsCap:    10,
		})

		wg sync.WaitGroup

		taskResults = make(map[int]TimeTask)
		taskErrors  []error
	)

	spamCtx, spamCancel := context.WithCancel(context.Background())
	defer spamCancel()

	// Запуск обрабатывателя.
	results, errs, done := p.Start()

	wg.Add(1)
	go func() {
		defer wg.Done()
		spamTasks(spamCtx, p)
	}()

	// Обработка результата.
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			case res := <-results:
				log.Debug("READED RESULT", zap.Int("id", res.id))
				taskResults[res.id] = res
			}
		}
	}()

	// Обработка ошибок
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			case err := <-errs:
				log.Debug("READED ERROR ", zap.Error(err))
				taskErrors = append(taskErrors, err)
			}
		}
	}()

	sleep := time.Now()
	log.Info("start sleep")
	time.Sleep(time.Second * 3) // Ждун
	log.Info("stop sleep", zap.Duration("time", time.Since(sleep)))
	spamCancel() // отмена спаммера

	start := time.Now()
	log.Info("start wait")
	// Максимум возможно 10 успешныз задач, каждая из которых блокируются на 150мс,
	// я решил подождать 11 задач чтоб наверняка.
	stopCtx, stopCancel := context.WithTimeout(context.Background(), 11*150*time.Millisecond)
	defer stopCancel()
	if err := p.Stop(stopCtx); err != nil {
		log.Fatal("unable to stop processor properly", zap.Error(err))
	}

	// Дожидаемся пока обработается результат и ошибки
	wg.Wait()
	log.Info("stop wait", zap.Duration("time", time.Since(start)))

	// Я подумал что важно выводить информацию именно после окончания работы, а не в процессе.

	log.Info("Errors:", zap.Int("total", len(taskErrors)))

	taskIDs := maps.Keys(taskResults)
	slices.Sort(taskIDs)

	log.Info("Done tasks:", zap.Ints("taskIDs", taskIDs))
}
