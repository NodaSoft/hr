package main

import (
	"fmt"
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

//README:Код можно сильно упростить, выбросив избыточные горутины там, где они совершенно не нужны, например при подсчете результатов, но не стал, пусть будет с горутинами

// A TaskType represents a meaninglessness of our life
type TaskType struct {
	id         int
	cT         string // время создания
	fT         string // время выполнения
	taskRESULT []byte
}

type ChannelsType struct {
	taskCh         chan TaskType //Канал задач
	doneTasksCh    chan TaskType //Канал результатов
	undoneTasksCh  chan error    //Канал ошибок
	closeAllCh     chan bool     //Канал закрытия горутин
	closeAllChDone chan bool     //Канал признака, что все горутины закрыты
}

type allResultsType struct {
	results map[int]TaskType
	errs    []error
}

//И каналы инициализируем под капотом
func initChannels() (channels ChannelsType) {
	channels = ChannelsType{ //Кому как, а мне проще таскать данные в структурах
		taskCh:         make(chan TaskType, 10), //Ну, единственная задача зачем толстые каналы - чтобы мультизадачность не уперлась в производительность получателей, что приведет к завешиванию задатчика
		doneTasksCh:    make(chan TaskType, 10), //Ну раз уж начали толстые каналы открывать, пусть все такие будут
		undoneTasksCh:  make(chan error, 10),
		closeAllCh:     make(chan bool), //Тут один источник и один получатель, т.ч. вместимость 1 - то что нужно
		closeAllChDone: make(chan bool), //Тут один источник и один получатель, т.ч. вместимость 1 - то что нужно
	}
	return
}

//Вынес это в отдельные функции, иначе код не читаем. Можно было вынести и остальные, но не критично
func taskCreator(channels ChannelsType) {

	defer close(channels.taskCh) //Глушим в источнике, читающие горутины сами отвалятся
	defer fmt.Println("taskCreturer closed")

	for {
		var exit bool
		select {
		case <-channels.closeAllCh:
			exit = true
		default:
			var task TaskType
			task.cT = time.Now().Format(time.RFC3339)
			if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
				task.cT = "Some error occured"
			}
			task.id = int(time.Now().Unix())
			// передаем таск на выполнение
			channels.taskCh <- task
		}
		if exit {
			break
		}
		time.Sleep(100 * time.Millisecond) //А то очень быстро задачи плодятся. Обработчик тасков все пережует, у нас же мультипоточность
	}
}

func (task *TaskType) taskWorker() {
	tt, _ := time.Parse(time.RFC3339, task.cT)
	if tt.After(time.Now().Add(-20 * time.Second)) {
		task.taskRESULT = []byte("task has been successed") //Считаю что здесь может быть некое сообщение, которое мы не можем править, иначе ответы нужно слать по-человечески, а не вытаскивать его из строки
	} else {
		task.taskRESULT = []byte("something went wrong")
	}
	task.fT = time.Now().Format(time.RFC3339Nano)
	time.Sleep(500 * time.Millisecond)
	return
}

func (task *TaskType) taskSorter(channels ChannelsType) {
	if len(task.taskRESULT) >= 14 { //Защищаемся по длине посылки
		if string(task.taskRESULT[14:]) == "successed" {
			channels.doneTasksCh <- *task
			return
		}
	}
	channels.undoneTasksCh <- fmt.Errorf("Task id %d time %s, error %s", task.id, task.cT, task.taskRESULT)
}

//Здесь обрабатываются все задачи
func taskMixer(channels ChannelsType) {
	var wg sync.WaitGroup
	for task := range channels.taskCh { //При закрытии канала цикл закроется. Классная штука, никогда не пользовался.
		wg.Add(1)
		task := task //А то в горутину не залезает
		go func() { //Мы, наверное, хотим чтобы это выполнялось параллельно?
			task.taskWorker()
			task.taskSorter(channels)
			//Разбор ответов мог быть здесь
			wg.Done()
		}()
	}
	wg.Wait() //Ждем когда все горутины вернут ответ

	//Если код большой, то это можно писать в начале через defer, но и так нормально
	//Чистим за собой
	close(channels.doneTasksCh)
	close(channels.undoneTasksCh)
	channels.closeAllChDone <- true //Не важно что слать
	close(channels.closeAllChDone)
	fmt.Println("closeAllChDone closed") //Не успеет до закрытия программы, но пусть будет
}

//Здесь мы запрашиваем закрытие всех горутин и ждем их завершения
func closeAndWaitAllGorutines(channels ChannelsType) {
	channels.closeAllCh <- true //Не важно что отправлять, пусть будет true

	//Это завесит код до закрытия всех горутин
	//Можно было использовать wg или <-channels.closeAllChDone, но тогда не будет обратной связи по ожиданию
	//Обязательно ждем ответы от всех горутин, иначе программа закроется пока последние горутины не вернули ответ
	var exit bool
	for {
		select {
		case <-channels.closeAllChDone:
			exit = true
		default:
			fmt.Println("wait for all gorutines is done")
		}
		if exit {
			break
		}
		time.Sleep(time.Second)
	}
}

//А почему бы обработчик ответов тоже не поместить в горутину
//takeResults - обработчик ответов
//Т.к. здесь мы используем указать, то данные в нужную переменную в main,
func (allResults *allResultsType) takeResults(channels ChannelsType) {
	//Вот так должно быть посимпатичнее, хотя не вижу причин почему это не сделать сразу после taskSorter
	//Можно было это сделать и в одном цикле через select, но и так красиво
	go func() {
		for result := range channels.doneTasksCh { //При закрытии канала цикл закроется
			allResults.results[result.id] = result
		}
		defer fmt.Println("doneTasksCh closed")
	}()

	go func() {
		for err := range channels.undoneTasksCh { //При закрытии канала цикл закроется
			allResults.errs = append(allResults.errs, err)
		}
		defer fmt.Println("undoneTasksCh closed")
	}()
}

func main() {

	channels := initChannels() //В последний момент подумал, а почему бы и это не вынести из основного кода

	//задатчик
	go taskCreator(channels)

	//получение и обработка тасков
	go taskMixer(channels) //В последний момент подумал, а почему бы и это не вынести из основного кода

	allResults := allResultsType{
		results: make(map[int]TaskType), //Карту обязательно нужно инициализировать
	}

	go allResults.takeResults(channels) //Здесь go не обязателен, сделал для защиты от дурака, если внутри метода горутины случайно уберут и завесят все

	time.Sleep(3 * time.Second) //я так понял что мы обрываем здесь работу приложения, а не ждем завершения всех процессов (которые плодятся в бесконечном цикле)

	closeAndWaitAllGorutines(channels) //В последний момент подумал, а почему бы и это не вынести из основного кода

	fmt.Println("Errors:")
	for _, err := range allResults.errs {
		fmt.Println(err)
	}

	fmt.Println("Done tasks:")
	for taskId := range allResults.results {
		fmt.Println(taskId)
	}

	fmt.Println("Bye, bye")
}
