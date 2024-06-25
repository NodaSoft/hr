package main

import (
	"context"
	"fmt"
	"taskhandler/internal/config"
	log "taskhandler/internal/logger"
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

// A Ttype represents a meaninglessness of our life
type Task struct {
	id       int64
	created  time.Time // время создания
	finished time.Time // время выполнения
	result   []byte
}
type brokenFactory struct {
	serviceStarted time.Time
}

const BROKEN_FACTORY_SLEEP_MS = 300

func (factory *brokenFactory) MakeTask() Task {
	time.Sleep(BROKEN_FACTORY_SLEEP_MS * time.Millisecond)
	now := time.Now()
	if now.Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
		//ft = "Some error occured"
		log.Debug("Broken factory sended broken task, again, someone fix it already")
		return Task{created: factory.serviceStarted, id: now.Unix()}
	}
	return Task{created: now, id: now.Unix()}

}

type TaskFactory interface {
	MakeTask() Task
}

func FillChannel(ctx context.Context, tch chan<- Task, factory TaskFactory) {
	for {
		log.Debug(".FillChannel iteration")
		select {
		case <-ctx.Done():
			close(tch)
			log.Info("FillChannel by context done, channel closed")
			return
		case tch <- factory.MakeTask():
		}
	}
}

func main() {
	// Startup, may panic, which is ok
	// The app shouldn't work without logger or with wrong config
	config.InitConfig()
	log.InitGlobalLogger()
	log.Info("Service preparations successed")
	ctx, _ := context.WithTimeout(context.Background(), config.C.Service.Timeout)

	// Magic .Time to generate incorrect tasks
	serviceStarted := time.Now()

	// Create tasks
	tasks := make(chan Task, 10)

	// Create workers
	go FillChannel(ctx, tasks, &brokenFactory{serviceStarted: serviceStarted})

	// Snailcase; fills the result field  of task
	task_worker := func(t Task) Task {
		if t.created.After(time.Now().Add(-20 * time.Second)) {
			t.result = []byte("task has been successed")
		} else {
			t.result = []byte("something went wrong")
		}

		// Not sure about this, but i'll leave this sleep untouched
		time.Sleep(time.Millisecond * 150)

		t.finished = time.Now()

		return t
	}

	doneTasks := make(chan Task)
	undoneTasks := make(chan error)

	tasksorter := func(t Task) {
		if t.created == serviceStarted {
			undoneTasks <- fmt.Errorf("Task id %d time %s, error %s", t.id, t.created, t.result)
			return
		}
		doneTasks <- t
	}

	go func() {
		// получение тасков
		for t := range tasks {
			t = task_worker(t)
			go tasksorter(t)
		}
		/* close(tasks) */
	}()

	result := map[int64]Task{}
	err := []error{}
	go func() {
		for r := range doneTasks {
			go func() {
				result[r.id] = r
			}()
		}
		for r := range undoneTasks {
			go func() {
				err = append(err, r)
			}()
		}
		close(doneTasks)
		close(undoneTasks)
	}()

	time.Sleep(time.Second * 3)

	println("Errors:")
	for r := range err {
		println(r)
	}

	println("Done tasks:")
	for r := range result {
		println(r)
	}
}
