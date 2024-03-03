package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"runtime"
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

// Job представляет собой задачу, которую нужно выполнить

//Доброго времени суток, исправлений виделось слишком много и учитывая что изначальный код
//на мой взгляд был попыткой реализовать паттерн pool workers, я решил это сделать.
//Код практически полностью изменен, но только так возможно выпонить п1 задания.

type Job struct {
	ID    string
	cT    time.Time // время создания
	fT    time.Time // время выполнения (с учетом очереди выполнения, если я правильно понял нас интересует именно это)
	Error error
}

// Worker представляет собой работника, который будет выполнять задачи
type Worker struct {
	ID          int
	TaskChannel chan Job
	QuitChannel <-chan struct{}
}

// NewWorker создает нового работника с указанным ID
func NewWorker(ctx context.Context, id int, jobChan chan Job) Worker {
	return Worker{
		ID:          id,
		TaskChannel: jobChan,
		QuitChannel: ctx.Done(),
	}
}

// Start запускает работника для выполнения задач
func (w *Worker) Start(wg *sync.WaitGroup) {
	go func() {

		for {
			func() {

				select {
				case job := <-w.TaskChannel:

					<-time.After(400 * time.Millisecond)

					job.fT = time.Now()

					if job.Error != nil {
						fmt.Println(job.Error.Error(), job.ID)
					} else {
						fmt.Printf("Job ID %s выполнен за %v\n", job.ID, job.fT.Sub(job.cT))
					}
					wg.Done()
				case <-w.QuitChannel:
					return
				}
			}()
		}
	}()
}

// Pool представляет собой пул работников
type Pool struct {
	Workers    []Worker
	JobChannel chan Job
	WG         sync.WaitGroup
}

// NewPool создает новый пул с указанным количеством работников
func NewPool(ctx context.Context, workerCount int) *Pool {
	pool := Pool{
		Workers:    make([]Worker, workerCount),
		JobChannel: make(chan Job),
	}

	for i := 0; i < workerCount; i++ {
		pool.Workers[i] = NewWorker(ctx, i, pool.JobChannel)

		pool.Workers[i].Start(&pool.WG)
	}

	return &pool
}

// SubmitJob отправляет задачу в пул
func (p *Pool) SubmitJob(job Job) {
	p.WG.Add(1)
	go func() {
		p.JobChannel <- job
	}()
}

// Shutdown завершает все работники пула после завершения задач
func (p *Pool) Shutdown() {
	p.WG.Wait()
	close(p.JobChannel)
}

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	pool := NewPool(ctx, runtime.NumCPU())
	defer cancel()
	tasks := make(chan Job)
	go taskCreator(ctx, tasks)

	// Отправка задач в пул
	for job := range tasks {
		pool.SubmitJob(job)
	}

	// Ожидание завершения задач и работников
	pool.WG.Wait()
	close(pool.JobChannel)
}

func taskCreator(ctx context.Context, task chan Job) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Creator stopping...")
				close(task)
				return

			case <-time.After(100 * time.Millisecond):
				cT := time.Now()

				var err error
				if time.Now().Second()%2 > 0 { // вот такое условие появления ошибочных тасков
					err = errors.New("Some error occured")
				}
				u := uuid.New()
				task <- Job{cT: cT, ID: u.String(), Error: err} // передаем таск на выполнение
			}
		}
	}()
}
