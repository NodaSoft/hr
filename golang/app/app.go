package app

import (
	"context"
	"fmt"
	"golang/config"
	"golang/models"
	"golang/pipe"
	Pipe "golang/pipe/v1"
	"log"
	"sync"
	"time"
)

type App struct {
	wg         *sync.WaitGroup
	doneCh     chan struct{}
	tasksLimit int
	tasksCap   int
	success    pipe.Interface
	failed     pipe.Interface
}

func New(cfg *config.Config) *App {
	app := new(App)
	app.wg = &sync.WaitGroup{}
	app.tasksCap = cfg.TasksQueueLimit

	return app
}

func (app *App) Run(ctx context.Context, tasksLimit int) error {
	app.success = Pipe.New(ctx)
	app.failed = Pipe.New(ctx)
	// асинхронно генерим таски
	superCh := app.taskCreturer(ctx, app.tasksCap)
	// асинхронно обрабатываем таски
	app.wg.Add(1)
	go app.handleTasks(ctx, superCh)
	// ждем пока все отработает
	app.wg.Wait()
	// вычитываем успешные результаты
	// можно это делать также асинхронно, но не вижу смысла

	// вычитываем ошибки

	return nil
}

// асинхронно создает таски. По завершении работы закрывает канал передачи резульатов.
func (app *App) taskCreturer(ctx context.Context, limit int) chan models.Ttype {
	superCh := make(chan models.Ttype, limit)
	app.wg.Add(1)
	go func() {
		defer app.wg.Done()
		defer close(superCh) // по завершении цикла выключаем шарманку
		for i := 0; i < app.tasksLimit; i++ {
			select {
			case <-ctx.Done():
				return
			default:
				task := models.Ttype{ID: int(time.Now().Unix())}
				task.CreateTime = time.Now()
				if task.CreateTime.Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
					task.TaskRESULT = append(task.TaskRESULT, []byte("Some error occured")...)
					task.IsFailed = true
				}
				superCh <- task // передаем таск на выполнение
			}
		}
	}()

	return superCh
}

func (app *App) handleTasks(ctx context.Context, tasksCh chan models.Ttype) {
	// получение тасков
	for t := range tasksCh {
		if err := ctx.Err(); err != nil {
			return
		}
		t = app.taskWorker(t)

		app.wg.Add(1)
		go func() {
			defer app.wg.Done()
			app.taskSorter(ctx, t)
		}()
	}
}

// puts tasks into boxes depends on result.
func (app *App) taskSorter(ctx context.Context, t models.Ttype) {
	var err error

	if !t.IsFailed {
		err = app.success.Send(ctx, t)
	} else {
		terr := fmt.Errorf("task id %d time %s, error %s", t.ID, t.CreateTime, t.TaskRESULT)
		err = app.failed.Send(ctx, terr)
	}

	if err != nil {
		log.Println(err)
	}
}

// emulate complete time
func (app *App) taskWorker(task models.Ttype) models.Ttype {
	const past = -20 * time.Second
	task.CompleteTime = time.Now()
	if task.IsFailed {
		return task
	}

	if task.CreateTime.After(time.Now().Add(past)) {
		task.TaskRESULT = []byte("task has been successed")
	} else {
		err := "something went wrong"
		task.TaskRESULT = []byte(err)
		task.IsFailed = true
	}
	// хотелось бы удалить, но наверно это эмулятор времени работы таски
	time.Sleep(time.Millisecond * 150)

	return task
}
