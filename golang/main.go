package main

import (
	"fmt"
	"sync"
	"time"
)

// Приложение эмулирует получение и обработку неких тасков. Пытается и получать, и обрабатывать в многопоточном режиме.
// После обработки тасков в течении 3 секунд приложение должно выводить накопленные к этому моменту успешные таски и отдельно ошибки обработки тасков.

// ЗАДАНИЕ: сделать из плохого кода хороший и рабочий - as best as you can.
// Важно сохранить логику появления ошибочных тасков.
// Важно оставить асинхронные генерацию и обработку тасков.
// Сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через pull-request в github
// Как видите, никаких привязок к внешним сервисам нет - полный карт-бланш на модификацию кода.

func main() {
	taskCreturer := func(a chan Task) {
		for {
			now := time.Now()
			res := Task{id: now.Unix()}
			if now.Nanosecond()%2 == 0 { // вот такое условие появления рабочей таски
				res.start = now
			}
			a <- res // передаем таск на выполнение
		}
	}

	superChan := make(chan Task)
	go taskCreturer(superChan)

	doneTasks := make(chan Task)
	undoneTasks := make(chan Task)
	go func() {
		// получение тасков
		for t := range superChan {
			t.Work()
			go t.Sort(doneTasks, undoneTasks)
		}
	}()

	var (
		result = make(map[int64]Task, 0)
		resM   = new(sync.Mutex)
	)
	go func() {
		for r := range doneTasks {
			go func(r Task) {
				resM.Lock()
				defer resM.Unlock()
				result[r.id] = r
			}(r)
		}
	}()

	var (
		errs []error
		errM = new(sync.Mutex)
	)
	go func() {
		for t := range undoneTasks {
			go func(t Task) {
				errM.Lock()
				defer errM.Unlock()
				errs = append(errs, fmt.Errorf("task id: %d\tstart: %s\terror: %s", t.id, t.start.Format(time.RFC3339), t.cause))
			}(t)
		}
	}()

	time.Sleep(time.Second * 3)

	println("Errors:")
	for _, err := range errs {
		println(err.Error())
	}

	println("Done tasks:")
	for _, task := range result {
		println(task.String())
	}
}

// A Task represents a meaninglessness of our life
type Task struct {
	start     time.Time // время создания
	end       time.Time // время выполнения
	id        int64
	cause     string
	isSucceed bool
}

func (t *Task) Work() {
	defer func() {
		t.end = time.Now()
		time.Sleep(time.Millisecond * 150)
	}()
	if t.isSucceed = !t.start.IsZero(); t.isSucceed {
		t.cause = "task has been successed"
	} else {
		t.cause = "something went wrong"
	}
}

func (t Task) Sort(done, undone chan Task) {
	if t.isSucceed {
		done <- t
	} else {
		undone <- t
	}
}

func (t Task) String() string {
	return fmt.Sprintf("task id: %d | start: %s | end: %s | cause: %s", t.id, t.start.Format(time.RFC3339), t.start.Format(time.RFC3339Nano), t.cause)
}
