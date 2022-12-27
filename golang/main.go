package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/NodaSoft/hr/golang/tasks"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнены остальных тасков

// В дальнейшем можно вынести константы в аргументы программы
const (
	WorkersCount = 10                     // Кол-во конкурентных обработчиков
	Interval     = time.Millisecond * 150 // Сколько ждать между созданием новой таски
	TTL          = 500 * time.Millisecond // Время жизни таски, иначе фейл
)

type undoneTaskError struct {
	id   int
	time string
	err  error
	t    string // type
}

func (e undoneTaskError) Error() string {
	return fmt.Sprintf("task id: %d, time: %s, error: %s: %s", e.id, e.time, e.t, e.err)
}

type doneTasks struct {
	mx sync.Mutex
	ts []tasks.Ttype
}

func (d *doneTasks) add(t tasks.Ttype) {
	d.mx.Lock()
	defer d.mx.Unlock()
	d.ts = append(d.ts, t)
}

type undoneTaskErrs struct {
	mx sync.Mutex
	es []error
}

func (d *undoneTaskErrs) add(t error) {
	d.mx.Lock()
	defer d.mx.Unlock()
	d.es = append(d.es, t)
}

func main() {
	exit := make(chan struct{})

	syschan := make(chan os.Signal, 1)
	signal.Notify(syschan, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-syschan
		close(exit)
	}()

	superChan := make(chan tasks.Ttype)

	// Даже если у нескольких тасок одинаковые id, выводим обе
	doneList := doneTasks{}
	undoneList := undoneTaskErrs{}

	var processor tasks.TtypeProcessor = tasks.TtypeProcessorMain{TTL: TTL}

	wg := &sync.WaitGroup{}
	wg.Add(WorkersCount)
	for i := 0; i < WorkersCount; i++ {
		go func() {
			defer wg.Done()
			// получение тасков
			for t := range superChan {
				err := processor.Process(&t)
				if err != nil {
					undoneList.add(undoneTaskError{
						id:   t.Id,
						time: t.CT,
						err:  err,
						t:    "incorrect task",
					})
				} else if t.ProcessingError != nil {
					undoneList.add(undoneTaskError{
						id:   t.Id,
						time: t.CT,
						err:  t.ProcessingError,
						t:    "processing error",
					})
				} else {
					doneList.add(t)
				}
			}
		}()
	}

	creatorStopped := tasks.TaskCreator(superChan, Interval, exit)

	<-creatorStopped
	wg.Wait()

	fmt.Println("Done tasks:")
	for _, r := range doneList.ts {
		fmt.Printf("%#v\n", r)
	}

	fmt.Println("Errors:")
	for _, e := range undoneList.es {
		_, _ = fmt.Fprintln(os.Stderr, e)
	}
}
