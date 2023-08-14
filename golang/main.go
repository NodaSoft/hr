package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// ЗАДАНИЕ:
//
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнены остальных тасков

// A Ttype represents a meaninglessness of our life
type Ttype struct {
	isFailed   bool
	id         int
	cT         string // время создания
	fT         string // время выполнения
	taskRESULT []byte
}

type Tasker struct {
	doneTasks    []Ttype
	failedTasks  []Ttype
	incomeTaskCh chan Ttype
	resultTaskCh chan Ttype
	taskToMake   int
	ctx          context.Context
	wg           sync.WaitGroup
}

func (t *Tasker) makerGT() {

	t.incomeTaskCh = make(chan Ttype, t.taskToMake)
	t.resultTaskCh = make(chan Ttype, t.taskToMake)

	for i := 0; i < t.taskToMake; i++ { // create tasks in single thread, using channel as buffer

		taskTime := time.Now()
		newTask := Ttype{id: int(taskTime.Unix()), cT: taskTime.Format(time.RFC3339)}

		t.incomeTaskCh <- newTask // send all new tasks in buffered ch
	}
}

func (t *Tasker) runTaskWorkers() {

	for i := 0; i < t.taskToMake; i++ {

		t.wg.Add(1)

		go func() { //each task worker running in individual rutine, with respecting contxt

			task, ok := <-t.incomeTaskCh // get task to complete

			if !ok {
				fmt.Println("error occured while reading in ch ")
				//TODO: dosomething clever
				return
			}

			delay := time.NewTimer(time.Millisecond * 150) // to simulate working proccess using timer. 150 mls - from original task

			select { //two things may happend here: 1 OR work complete; 2 OR context is done

			case <-t.ctx.Done(): //context done before work finished, greacefully finishing work
				fmt.Println("context canncelled before work done")

				if !delay.Stop() { // avoid memory leaks
					<-delay.C
				}
				break

			case <-delay.C: //here work is done (or at least we get return value here)

				//---------------------------------- ORIGINAL LOGIC
				if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
					task.isFailed = true
					task.fT = "Some error occured"
					task.taskRESULT = []byte("something went wrong")

				} else {
					task.isFailed = false
					task.fT = time.Now().Format(time.RFC3339)
					task.taskRESULT = []byte("task has been successed")
				}
				//----------------------------------

				t.resultTaskCh <- task
				break
			}

			t.wg.Done()
		}()
	}
}

func (t *Tasker) listenResults() {

	for { //single thread listener
		select { //two things may happens: 1 ctx done (interrupted or all tasks complete); 2 we recieve result from worker
		case <-t.ctx.Done(): // context was cancelled while we are  waiting for responses (all complete, or error)
			//TODO:
			fmt.Println("CTX is done")

			return

		case result, ok := <-t.resultTaskCh:
			if !ok {
				//TODO:
				fmt.Println("ch closed")
				break
			}

			if !result.isFailed { //check status of task
				t.doneTasks = append(t.doneTasks, result)
			} else {
				t.failedTasks = append(t.failedTasks, result)
			}
			break
		}
	}

}

func (t *Tasker) resultPrinter() {

	fmt.Println("Succeed tasks:")
	for _, tr := range t.doneTasks {
		fmt.Println(tr)
	}

	fmt.Println("Failed tasks:")
	for _, tr := range t.failedTasks {
		fmt.Println(tr)
	}
}

func main() {

	ctx, cancelFn := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println("SIG I")
		fmt.Println(sig)
		cancelFn() // close context, and force all workers (and result listener) to break
	}()

	tasker := Tasker{taskToMake: 10, ctx: ctx} // manually set 10 tasks to create, in 10 threads

	tasker.makerGT()

	tasker.runTaskWorkers()

	go tasker.listenResults() //advanced thread to recieve results from task

	tasker.wg.Wait() // wait untill all workers complete their tasks
	cancelFn()       // close context, to tell result listener finish also

	tasker.resultPrinter()
}
