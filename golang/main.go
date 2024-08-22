package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
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

const (
	programTTL        = 10 * time.Second
	reportTime        = 3 * time.Second
	taskWorkTime      = 150 * time.Millisecond
	taskRelevanceTime = 20 * time.Second
)

var (
	ErrCreate  = errors.New("error while create")
	ErrTimeout = errors.New("error task timeout")
)

type Task struct {
	id         int
	createTime time.Time // время создания
	finishTime time.Time // время завершения
	err        error
}

func (t *Task) Work() {
	if !t.createTime.After(time.Now().Add(-taskRelevanceTime)) && t.err == nil {
		t.err = ErrTimeout
	}

	t.finishTime = time.Now()

	time.Sleep(taskWorkTime)
}

func (t Task) String() string {
	return fmt.Sprintf("Task id %d | create time: %s | work time: %s ", t.id, t.createTime.Format(time.RFC3339), t.finishTime.Format(time.RFC3339Nano))
}

type ErrorHandler struct {
	err []error
	mtx sync.Mutex
}

func (eh *ErrorHandler) Store(newErr error) {
	eh.mtx.Lock()
	defer eh.mtx.Unlock()
	eh.err = append(eh.err, newErr)
}

func (eh *ErrorHandler) LoadAllAndDelete() []error {
	eh.mtx.Lock()
	defer eh.mtx.Unlock()

	newErr := make([]error, len(eh.err))

	copy(newErr, eh.err)

	eh.err = nil

	return newErr
}

type TaskStorage struct {
	done        sync.Map
	undoneTasks ErrorHandler
}

func (ts *TaskStorage) PrintInfo() string {
	strB := strings.Builder{}

	strB.WriteString("Errors:\n")
	errs := ts.undoneTasks.LoadAllAndDelete()
	for _, err := range errs {
		strB.WriteString(fmt.Sprintf("%s\n", err.Error()))
	}

	strB.WriteString("\nDone tasks:\n")

	ts.done.Range(func(key, value any) bool {
		strB.WriteString(fmt.Sprintf("%s\n", value))
		return true
	})

	return strB.String()
}

func (ts *TaskStorage) Store(wg *sync.WaitGroup, t Task) {
	defer wg.Done()
	if t.err != nil {
		ts.undoneTasks.Store(fmt.Errorf("Task id %d time %s, error %w", t.id, t.createTime, t.err))
	} else {
		ts.done.Store(t.id, t)
	}
}

type TaskManager struct {
	ts TaskStorage
	wg sync.WaitGroup
}

func (tm *TaskManager) createTasks(ctx context.Context) chan Task {
	out := make(chan Task, 10)

	tm.wg.Add(1)
	go func(ctx context.Context, wg *sync.WaitGroup) {
		defer wg.Done()
	CRLP:
		for {
			select {
			case <-ctx.Done():
				break CRLP
			default:
				ft := time.Now()
				var err error
				if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
					err = ErrCreate
				}
				out <- Task{createTime: ft, id: int(time.Now().Unix()), err: err} // передаем таск на выполнение
			}

		}
		close(out)

	}(ctx, &tm.wg)

	return out
}

func (tm *TaskManager) processTasks(tasks <-chan Task) {

	defer tm.wg.Done()
	internal_wg := &sync.WaitGroup{}
	// получение тасков
	for t := range tasks {
		t.Work()
		internal_wg.Add(1)
		go tm.ts.Store(internal_wg, t)
	}
	internal_wg.Wait()

}

func (tm *TaskManager) Run(ctx context.Context) {

	taskChannel := tm.createTasks(ctx)
	go tm.processTasks(taskChannel)

	printTicker := time.NewTicker(reportTime)
EVLP:
	for {
		select {
		case <-ctx.Done():
			break EVLP
		case <-printTicker.C:
			fmt.Println(tm.ts.PrintInfo())
		}
	}

}

func (tm *TaskManager) Wait() {
	tm.wg.Wait()
}

func main() {
	ctx, canlce := context.WithCancel(context.Background())
	var tm TaskManager
	go tm.Run(ctx)
	time.Sleep(programTTL)
	canlce()
	tm.Wait()
}
