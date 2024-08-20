package main

import (
	"errors"
	"fmt"
	"runtime"
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
// --------------------------------------------------------------------------------------------------

// Код проверил race detector'ом
// Все pull request в открытом доступе, я писал самостоятельно, не копировал
// В данном случае "разносить код по папочкам" не вижу необходимым
const (
	timeLimit     = 10 * time.Second
	printInterval = 3 * time.Second
	proceedTime   = 150 * time.Millisecond
	taskTimeout   = 2 * proceedTime
)

type MyTask struct {
	id         int
	createTime time.Time
	finishTime time.Time
	succeed    bool
	err        error
}

var mainWG sync.WaitGroup
var workerWG sync.WaitGroup

func main() {
	// Читаем число максимума гошных процессеров для правильной параллельности
	gomaxprocs := runtime.GOMAXPROCS(0)

	// Если установить размер буффера каналов tasks, processedTasks равным GOMAXPROCS,
	// воркеры не будут голодать и не будут ждать освобождения буффера канала
	tasks := make(chan MyTask, gomaxprocs)
	processedTasks := make(chan MyTask, gomaxprocs)
	quit := make(chan struct{})

	mainWG.Add(2)
	go taskCreator(tasks)
	workerWG.Add(gomaxprocs)
	// запускаем worker pool
	for i := 0; i < gomaxprocs; i++ {
		// Очевидно с имитацией тасков можно запустить намного больше горутин
		// Но для правильной параллельности использую GOMAXPROCS
		go taskWorker(tasks, processedTasks)
	}
	go taskSorter(processedTasks, quit)

	workerWG.Wait()
	quit <- struct{}{}
	close(processedTasks)
	mainWG.Wait()
}

func taskCreator(tasks chan<- MyTask) {
	defer mainWG.Done()
	defer close(tasks)

	//Таймер
	timer := time.NewTimer(timeLimit)
	for id := 0; ; id++ {
		select {
		case <-timer.C:
			return
		default:
			currTime := time.Now().UTC()
			task := MyTask{
				id:         id,
				createTime: currTime,
				finishTime: time.Time{},
				succeed:    false,
				err:        nil,
			}
			if time.Now().Nanosecond()>>2%2 > 0 {
				// Не совсем понял под "Важно сохранить логику появления ошибочных тасков", могу ли я менять условие?
				// поскольку time.Now() имеет разрешение 100 наносекунд, чтобы условие срабатывало
				// я сдвинул биты, при этом, если таски не задерживаются при записи в канал, то они будут иметь
				// одинаковое время создания, от чего ошибки нельзя навать случайными
				task.err = errors.New("something went wrong")
			}
			tasks <- task
		}
	}

}

func taskWorker(tasks <-chan MyTask, doneTasks chan<- MyTask) {
	defer workerWG.Done()

	for task := range tasks {
		// Довольно непонятная проверка таймаута. Почему таймаут такой большой?
		// Проверяю до обработки задачи
		// Предполагаю, что перед обработкой большая часть задач будет простаивать времени
		// чуть больше const proceedTime, но одна задача будет иметь простой 2*proceedTime
		// поскольку будет ожидать освобождения канала tasks.
		// Задача мнимая и непонятно что именно требуется сделать, и по какому условию считать таймаут
		// но если такая проверка должна срабатывать, то я поставлю таймаут равным 2*proceedTime
		// Если же нужно исключать задачи с аномально высоким простоем, я бы поставил 3*proceedTime,
		// но в этом примере условие не будет срабатывать
		timeout := !task.createTime.After(time.Now().UTC().Add(-taskTimeout))
		if timeout {
			task.err = errors.Join(task.err, errors.New("task timeout"))
		} else {
			// proceed imitation
			time.Sleep(proceedTime)
		}

		if task.err == nil {
			task.succeed = true
		}
		task.finishTime = time.Now().UTC()
		doneTasks <- task
	}
}

func taskSorter(doneTasks <-chan MyTask, quit <-chan struct{}) {
	defer mainWG.Done()
	succeedTasks := make([]MyTask, 0)
	notSucceedTasks := make([]MyTask, 0)

	ticker := time.NewTicker(printInterval)
	defer ticker.Stop()
	var task MyTask
	for {
		select {
		case task = <-doneTasks:
			if task.succeed {
				succeedTasks = append(succeedTasks, task)
			} else {
				notSucceedTasks = append(notSucceedTasks, task)
			}
		case <-ticker.C:
			mainWG.Add(1)
			go printTasks(succeedTasks, notSucceedTasks)
			succeedTasks = make([]MyTask, 0)
			notSucceedTasks = make([]MyTask, 0)
		case <-quit:
			mainWG.Add(1)
			printTasks(succeedTasks, notSucceedTasks)
			return
		}
		// передаем ресурсы процессора другим горутинам
		runtime.Gosched()
	}
}

func printTasks(succeed, notSucceed []MyTask) {
	defer mainWG.Done()
	builder := strings.Builder{} // поменьше syscall

	builder.WriteString("Succeeded tasks:\n")
	for i := 0; i < len(succeed); i++ {
		builder.WriteString(fmt.Sprintf("\tTask id: %d Time: %d nanoseconds\n",
			succeed[i].id, succeed[i].finishTime.Sub(succeed[i].createTime).Nanoseconds()))
	}

	builder.WriteString("Not succeeded tasks:\n")
	for i := 0; i < len(notSucceed); i++ {
		builder.WriteString(fmt.Sprintf("\tTask id: %d Time: %d nanoseconds Err: %s\n",
			notSucceed[i].id, notSucceed[i].finishTime.Sub(notSucceed[i].createTime).Nanoseconds(), notSucceed[i].err))
	}
	fmt.Print(builder.String())
}
