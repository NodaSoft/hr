package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
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



// ============================================
// Было принято решение поменять структуру Task (поменял название для простоты чтения), для более чистой логики было добавлено поле success, чтобы не определять успешность
// таска по куски строки taskResult, а также добавить поле err, опять же, на случай неуспешного выполнения таска.
// Весь функционал разделен в 4 функции, creator для создания, worker для имитации обработки, packager для выгрузки обработанных тасок и printer для отображения тех тасок,
// что packager успел выгрузить. Таким образом main() остается чистым и просто вызывает всю необходимую логику.
// Долго думал, стоит ли, но по итогу снес удачные и неудачные таски в один канал и в одном формате, таким образом (если представить другие сценарии), сохраняется гибкость
// в поведении printer'a, выводить можем что хотим и как хотим, так как для этого есть все данные, плюс это упрощает логику в packager.
// Сделал простой семафор через chan int, чтобы останавливать packager, в тот момент, как printer очищает прочитанные значения.
// Общий context имеет timeout в 10 секунд (время на генерация по заданию), printer срабатывает каждые 3 секунды, к этому моменту у него уже есть бэклог тасок для вывода.
// Некоторые куски хорошо было бы еще чуть-чуть разбить и структурировать, но в данном примере это бы только все усложнило.

// Ссылку на этот PR приложил в отклик на HH.
// ============================================



type Task struct {
	id          int     // Task ID
	createdAt   string  // Task creation time
	executedAt  string  // Task exection time
    success     bool    // Task creation status
    err         error   // Error in case of creation failure
}

// Creates new tasks, until the Context is canceled
func taskCreator(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup, newTasks chan<- Task) {
    for {
        select {
            case <-ctx.Done():
                cancel()
                close(newTasks)
                wg.Done()
                return 
            default:
                createdAt := time.Now().Format(time.RFC3339)
                success := true

                if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
                    success = false
                }

                id := rand.IntN(1000000 - 100000) + 100000

                newTasks <- Task{createdAt: createdAt, id: id, success: success}
        }
    }
}

// Simulates task execution Sorts executed tasks into corresponding channels
func taskWorker(ctx context.Context, wg *sync.WaitGroup, newTasks chan Task, processedTasks chan<- Task) {
    for t := range newTasks {
        select {
            case <- ctx.Done():
                close(processedTasks)
                wg.Done()
                return
            default:
                if !t.success {
                    t.err = errors.New("Something went wrong!")
                }

                t.executedAt = time.Now().Format(time.RFC3339Nano)

                time.Sleep(time.Millisecond * 150)

                select {
                    case processedTasks <- t:
                    default:
                }
        }
    }
}

// Handles sorted tasks and puts them into corresponding slices for printing
func taskPackager(ctx context.Context, wg *sync.WaitGroup, sem chan int, processedTasks <-chan Task, output *[]Task) {
    for {
        select {
            case <- ctx.Done():
                wg.Done()
                return
            default:
                sem <- 1
                processedTask := <- processedTasks
                *output = append(*output, processedTask)
                <- sem
        }
    }
}

// Prints packaged tasks
func taskPrinter(ctx context.Context, wg *sync.WaitGroup, sem chan int, output *[]Task) {
    for range time.Tick(3 * time.Second) {
        select {
            case <- ctx.Done():
                wg.Done()
                return
            default:

                var successfull, unsuccessfull string
                var successfullLen, unsuccessfullLen int

                for _, v := range *output {
                    if v.success {
                        successfullLen++
                        successfull = fmt.Sprintf("%s\nTask id:%d, createdAt:%s, executedAt:%s", successfull, v.id, v.createdAt, v.executedAt)
                    } else {
                        unsuccessfullLen++
                        unsuccessfull = fmt.Sprintf("%s\nTask id:%d, createdAt:%s, error:%v", unsuccessfull, v.id, v.createdAt, v.err)
                    }
                }

                fmt.Println(fmt.Sprintf("\n\nSuccessfull (%d):%s", successfullLen, successfull))
                fmt.Println(fmt.Sprintf("\nUnsuccessfull (%d):%s", unsuccessfullLen, unsuccessfull))
                sem <- 1
                *output = (*output)[:successfullLen+unsuccessfullLen]
                <- sem

        }
    }
}

func main() {
	newTasks := make(chan Task)
    processedTasks := make(chan Task)
	
    output := make([]Task, 0)
    sem := make(chan int, 1)
	
    wg := sync.WaitGroup{}
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

    // Start generating tasks
    wg.Add(1)
	go taskCreator(ctx, cancel, &wg, newTasks)

    // Start receiving generated tasks, simulate execution and sort them
    wg.Add(1)
    go taskWorker(ctx, &wg, newTasks, processedTasks)

    // Start packaging tasks into output
    wg.Add(1)
    go taskPackager(ctx, &wg, sem, processedTasks, &output)
   
    // Start printer with the ticker inside
    wg.Add(1)
    go taskPrinter(ctx, &wg, sem, &output)

    wg.Wait()
}
