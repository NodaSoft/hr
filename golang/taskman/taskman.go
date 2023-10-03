package taskman

////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// TaskManager: Package solution
// The Solution demonstrates the skills in goroutines, chans, slices, pointers, interfaces in Go,
// and also represents a meaninglessness of our life
////////////////////////////////////////////////////////////////////////////////////////////////////////////////

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// Public Constants

// ВНИМАНИЕ!
const DEBUG = true	// Препроцессинг в скрипте build_and_start удалит строки, содержащие "DEBUG", поэтому они однострочные: if DEBUG { ... }. Скорость увеличится.
					// Если просто собрать "go build .", то будет выводиться копмактный лог начала и завершения каждой функции или итерации. Медленнее, но наглядно.
const UintSize = 32 << (^uint(0) >> 32 & 1) // 32 или 64
const (
    MaxInt       = 1<<(UintSize-1) - 1 // 1<<31 - 1 или 1<<63 - 1
    MinInt       = -MaxInt - 1         // -1 << 31 или -1 << 63
    MaxUint uint = 1<<UintSize - 1     // 1<<32 - 1 или 1<<64 - 1
)

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// A Task represents a meaninglessness of our life
// Public:
type Task struct {
	id     int
	cT     string // время создания
	fT     string // время выполнения
	result []byte
	i      int    // индекс задачи по порядку, для отладки
}
// Public Метод (Task).String(), имплементирует интерфейс Stringer из модуля fmt
func (t Task) String() string {
	return fmt.Sprintf("Task { id:%d, cT:%s, fT:%s, result:%s }", t.id, t.cT, t.fT, t.result)
}
// Public Метод (Task).Successed(), возвращает булевское значение таска - успешность
func (t Task) Successed() bool {
	return string(t.result[14:23]) == "successed"
}

// Public vars:
var (
	BufCap		int				// Ёмкость кольцевого буфера 
	MaxCount	int				// Максимальное число итераций. Если <=0, то бесконечный цикл
	MaxResult	int				// Пороговое значение len для results и errors, при достижении которого данные будут выгружаться в хранилище
	Objects	= [...] string { "inChan", "doneTasks", "undoneTasks", "results", "errors", "result_cmd", "errors_cmd" } // Объекты мониторинга для микросервисов
)
// Private vars:
var (
	stop		bool			// стоп - сигнал 
	counter		int				// Счетчик итераций
	goCounter   int            	// Счетчик активных горутин
	goPeakCount int            	// Пиковое количество активных горутин

	inChan		chan Task		// Основной канал.
	inChanOpen	bool = false	// состояние для мониторинга

	doneTasks	chan Task		// Канал успешных тасков
	doneTOpen 	bool = false	// состояние для мониторинга

	undoneTasks chan error		// Канал ошибок
	undoneTOpen	bool = false	// состояние для мониторинга

	result_cmd	chan struct{}	// Канал команд для выгрузки результатов
	resCmdOpen 	bool = false	// состояние для мониторинга
	
	errors_cmd	chan struct{}	// Канал команд для выгрузки ошибок
	errCmdOpen 	bool = false	// состояние для мониторинга

	results 	map[int]Task	// Карта результатов. Ограничена MaxResult
	resUnlocked bool = true		// состояние для мониторинга
	muRes		sync.Mutex		// Мьютекс для блокировки во время изменений

	errors		[]error			// Список ошибок. Ограничен MaxResult
	errUnlocked bool = true		// состояние для мониторинга
	muErr		sync.Mutex		// Мьютекс для блокировки во время изменений
)


// Public functions

// Возврощает значение счетчика counter
func Counter() int {
	return counter
}

// Возврощает значение счетчика goCounter
func GoCounter() int {
	return goCounter
}

// Возврощает значение счетчика goPeakCount
func GoPeakCount() int {
	return goPeakCount
}

// Возвращает состояние объекта для мониторинга (канала, слайса или мапы) по строчному имени.
// сюда же будут добавляться другие объекты - память, дисковое пространство, сетевые интерфейсы и т.д.
// inChan, doneTasks, undoneTasks, results, errors, result_cmd, errors_cmd
func State(name string) (int, int, bool) {
	switch name {												// Used,	 All,	  Positiv | Negativ
	case "inChan":
		return len(inChan), cap(inChan), inChanOpen				// Len, Capacity, true(opened)|false(closed)
	case "doneTasks":
		return len(doneTasks), cap(doneTasks), doneTOpen		// Len, Capacity, true(opened)|false(closed)
	case "undoneTasks":
		return len(undoneTasks), cap(undoneTasks), undoneTOpen	// Len, Capacity, true(opened)|false(closed)
	case "results":
		l:=len(results); o:=0; if l%4 > 0 { o = 8 }				// Приближенная оценка, на основании допущения, что баскет заполняется на 1/2
		return l, l + l + o , resUnlocked						// Len, ~Capacity, true(Unlocked)|false(Locked)
	case "errors":
		return len(errors), cap(errors), errUnlocked			// Len, Capacity, true(Unlocked)|false(Locked)
	case "result_cmd":
		return len(result_cmd), cap(result_cmd), resCmdOpen		// Len, Capacity, true(opened)|false(closed)
	case "errors_cmd":
		return len(errors_cmd), cap(errors_cmd), errCmdOpen		// Len, Capacity, true(opened)|false(closed)
	default:
		return -1, -1, false									// для неизвестных объектов
	}
}


// Некая универсальная функция реализованная в пакете, которую можно использовать в других пакетах
func Factorial(v int64) int64 {
	if v > 1 {
		return v * Factorial( v - 1 )
	}
	return v
}

// Инициализация и Выполнение
func Run () {
	stdout, _  := exec.Command("sh", "-c", "ulimit -Sn").Output()	// Ограничения системы на софт, понадобится в следующей версии с реальной нагрузкой на воркеров
	go_limit_s, _ := strconv.Atoi(string(stdout[:len(stdout)-1]))  	// Удаляем перенос строки, приводим к типу int
	stdout, _  = exec.Command("sh", "-c", "ulimit -Hn").Output()	// Ограничения системы на хард, понадобится в следующей версии с реальной нагрузкой на БД
	go_limit_h, _ := strconv.Atoi(string(stdout[:len(stdout)-1]))  	// Удаляем перенос строки, приводим к типу int

	if(BufCap == 0){ BufCap = 10 }								// Значения по умолчанию, если не заданы явно
	if(MaxCount == 0){ MaxCount = BufCap * 2 }					// 
	if(MaxResult == 0){ MaxResult = 1000 }						// 

	stop = false
	go off()													// Ожидает сигнал выключения и меняет значение stop = true

	inChan = make(chan Task, BufCap)							// Канал поступающих тасков
	inChanOpen	= true											// состояние для мониторинга

	doneTasks = make(chan Task, BufCap/2)						// Канал завершенных тасков
	doneTOpen = true											// состояние для мониторинга

	undoneTasks = make(chan error, BufCap/2)					// Канал тасков с ошибками
	undoneTOpen	= true											// состояние для мониторинга

	results = map[int]Task{}									// Карта результатов по id таска
	errors  = []error{}											// Список ошибок

	result_cmd = make(chan struct{})							// Канал команд на выгрузку
	resCmdOpen = true											// состояние для мониторинга

	errors_cmd = make(chan struct{})							// Канал команд на выгрузку
	errCmdOpen = true											// состояние для мониторинга


	println("Буфер канала входящих таксов......................: ",cap(inChan))
	println("Буфер канала завершенных тасков...................: ",cap(doneTasks))
	println("Буфер канала ошибок...............................: ",cap(undoneTasks))
	println("Число логических процессоров..................nCPU: ",runtime.NumCPU())
	println("Количество потоков ОС для GO =....GOMAXPROCS(nCPU): ",runtime.GOMAXPROCS(8))
	println("Лимит системы на активные процессы............soft: ",go_limit_s)
	println("Лимит системы на файловые дескрипторы.........hard: ",go_limit_h)
	if MaxCount > 0 {
		println("Для выхода найжмите любую клавишу либо дождитесь завершения цикла: ",MaxCount)
	} else {
		println("Задан бесконечный цикл! Для выхода найжмите любую клавишу.")
	}

	var wg sync.WaitGroup								// счетчик горутин верхнего уровня
	wg.Add(6)											// Инкремент на эти 6 горутин
	go errorsUpload(&wg,errors_cmd)						// Ждет команду и Выгружает ошибки
	go resultUpload(&wg,result_cmd)						// Ждет команду и Выгружает результаты
	go errorsCollector(&wg,undoneTasks)					// Копит ошибки
	go resultCollector(&wg,doneTasks)					// Систематизирует результаты
	go taskReceiver(&wg,inChan,doneTasks,undoneTasks)	// Приемник - принимает задачи, передает на обработку
	go taskTransmitter(&wg,inChan)						// Передатчик - отправляет задачи
	wg.Wait()	 										// Ожидает завершения 6 горутин верхнего уровня
}

// Печатает ошибки и успешные таски
func Log () {
	println("\x1b[31m" + "Errors: ",len(errors))				// Список ошибок красным цветом
	for i := range errors {
		if i > BufCap { break	} 
		println(errors[i].Error())
	}
	fmt.Println("\x1b[32m" + "Done tasks: ",len(results))		// Список успешных результатов зеленым цветом
	i := 0
	for key := range results{
		i++
		if i > BufCap { break	}
		fmt.Println(results[key]) // fmt.Println умеет сам вызывать метод Task.String() т.к. имплементирован интерфейс Stringer
	}
	fmt.Println("\x1b[0m" + "The End")							// Конец программы, цвет по умолчанию
}

//////////////////////////////////////// Private functions /////////////////////////////////////////////////////////

// Коллектор ошибок. Два варианта блокировки.
func errorsCollector (wg *sync.WaitGroup, c <-chan error) {
	defer wg.Done()
	for e := range c {					
		if ( len(errors) >= MaxResult){
			select {
			case errors_cmd <- struct{}{}:			// Команда выгрузки
			default:
			}
			time.Sleep(5*time.Millisecond)			// Синхронный вариант
			for !resUnlocked {						// Такой способ намного оптимальнее по скорости, т.к. вызывается в 100000 раз реже, чем append
				time.Sleep(5*time.Millisecond)		// ждем завершение выгрузки
			}
		  } 
		//muErr.Lock();	errUnlocked=false			// Асинхронный вариант
		errors = append(errors, e)		
		//muErr.Unlock();	errUnlocked=true
		if DEBUG { print(" E") }
	}
}

// Коллектор успешных тасков. Два варианта блокировки.
func resultCollector (wg *sync.WaitGroup, c <-chan Task) {
	defer wg.Done()
	for task := range c { 	
		if (len(results) >= MaxResult){
			select {
			case result_cmd <- struct{}{}:			// Команда выгрузки
			default:
		   }
		   time.Sleep(5*time.Millisecond)			// Синхронный вариант
		   for !resUnlocked {						// Такой способ намного оптимальнее по скорости, т.к. вызывается в 100000 раз реже, чем добавление
		   		time.Sleep(5*time.Millisecond)		// ждем завершение выгрузки
		   }
		}
		//muRes.Lock();	resUnlocked=false			// Асинхронный вариант
		results[task.id] = task		
		//muRes.Unlock();	resUnlocked=true
		if DEBUG { print(" T") }
	}
}


// Выгрузка из коллектора ошибок.
func errorsUpload (wg *sync.WaitGroup, cmd <-chan struct{}){
	defer wg.Done()
	for range cmd { // Ожидаем очередную команду выгрузки. При закрытии канала цикл цикл и функция завершатся
		//storage <- errors	// сохранить данные в хранилище
		muErr.Lock();	errUnlocked=false		// Блокируем
		errors = errors[:0] 					// Очищаем коллектор
		muErr.Unlock();	errUnlocked=true		// Разблокируем
	}
}


// Выгрузка из коллектора результатов
func resultUpload (wg *sync.WaitGroup, cmd <-chan struct{}){
	defer wg.Done()
	for range cmd { // Ожидаем очередную команду выгрузки. При закрытии канала цикл цикл и функция завершатся
		//storage <- results // сохранить данные в хранилище
		muRes.Lock();	resUnlocked=false	// Блокируем
		clear(results)  					// очищаем коллектор
		muRes.Unlock();	resUnlocked=true  	// Разблокируем
	}
}


// Воркер с полезной нагрузкой
func taskWorker (task Task) Task {
	t, _ := time.Parse(time.RFC3339, task.cT)
	if t.After(time.Now().Add(-20 * time.Second)) { // таски с ошибкой имеют другой формат cT и попадают в else
		task.result = []byte("task has been successed")
	} else {
		task.result = []byte("something went wrong")
	}
	imax := (task.i % 10) + 1   					// Переменный факториал в пределах 1..10
	task.result = append(task.result, fmt.Sprintf(". Factorial(%d)=%d",imax,Factorial(int64(imax)))...)
	task.fT = time.Now().Format(time.RFC3339Nano)	// Время завершения задачи
	return task
}


// Сортировщик тасков
func taskSorter (wg *sync.WaitGroup, task Task, done chan<- Task, err chan<- error) {
	if DEBUG { print(" (W",task.i) } // Метка горутины на старте
	defer wg.Done()
	task = taskWorker(task)			 // Ожидаем завершение синхроного вызова в рамках этой горутины.
	switch task.Successed() {		 // Эквивалентно if else
	  case true:
		done <- task
		if DEBUG { print("\x1b[32m"," W",task.i,")}","\x1b[0m") } // Горутина завершилась
	  default:
		err <- fmt.Errorf(task.String()) // явный вызов .String(), т.к. fmt.Errorf не поддерживает интерфейс Stringer
		if DEBUG { print("\x1b[31m"," W",task.i,")}","\x1b[0m") } // Горутина завершилась
	}
	goCounter--   	// Cчетчик активных горутин для статистики.
}

// Приемник тасков
func taskReceiver (wg *sync.WaitGroup, c <-chan Task, done chan Task, undo chan error) {
	defer wg.Done()
	var wg2 sync.WaitGroup		// локальная группа дочерних горутин
	for task := range c {		// продолжается до закрытия и опустошения канала. При чтении из пустого и открытого, канал блокируется
		goCounter++
		if goPeakCount < goCounter { goPeakCount = goCounter }
		if DEBUG { print(" R",task.i) }		// Метка горутины перед стартом
		wg2.Add(1)							// Инкремент дочерних горутин
		go taskSorter(&wg2,task,done,undo) 	// Обработать каждую полученную задачу в отдельном потоке, счетчик зафиксировать.
	}
	if DEBUG { print("\n[\n Wait\n") }
	wg2.Wait()					// Ждать завершение дочерних пишущих горутин
	if DEBUG { print("\n]\n Close Collectors сhans\n") }
	close(done);		doneTOpen  = false	// Закрыть каналы, чтобы читающие горутины могли завершиться 
	close(undo);		undoneTOpen= false
	close(errors_cmd);  errCmdOpen = false
	close(result_cmd); 	resCmdOpen = false

}
	

// Выключатель. Ждет ненулевой сигнал и устанавливает переменную выключения
func off () {
	os.Stdin.Read(make([]byte,1))	// ожидает ввода команды с клавиатуры
	stop = true						// приведет к завершению цикла Sender
	if DEBUG { println("Pressed a Key, stop =", stop) }
}
// Вариант без ввода с клавиатуры, сигнал или команда поступает из канала
func off2 (c <-chan byte) {
	for range c {				// цикл заблокирован пока не поступит сигнал
		stop = true				// приведет к завершению цикла Sender
		break
	}
}


// Передатчик тасков
func taskTransmitter (wg *sync.WaitGroup, c chan<- Task) {
	defer wg.Done()
	// Т.к. размер буфера задан явно и нужно в конце вывести короткие списки, ограничиваем цикл либо емкостью буфера либо числом итераций MaxCount
	// Установив MaxCount=-1, цикл можно сделать бесконечным, а для выключения использовать stop-сигнал
	for {
		counter++
		if stop || MaxCount > 0 && MaxCount < counter { break } // условие окончания цикла
		if DEBUG { print(" {S", counter) }
		t := time.Now()											// фиксируем время 
		tf := t.Format(time.RFC3339)							// берем его так
		id := int(t.UnixNano())									// и сяк, один раз, для соответствия проверяемого и сохраняемого значения.
		if id%2 > 0 {											// "вот такое искусственное условие появления ошибочных тасков"
			tf = "Some error occured"							// таски с нечетным временем получают признак ошибки
		}                                           			
		c <- Task{cT: tf, id: id, i: counter}    				// создаем экземпляр и передаем таск в канал.
		if counter >= MaxInt { counter = 0 }					// в бесконечном цикле счетчик должен быть циклическим
		/*
		if counter % 30000 == 0 {								// Отладочный мониторинг состояния объектов
			for i := range Objects {							// В v.1.0 это будет в Web-интерфейсе через микросервис
				n := Objects[i]
				l, c, s := State(n)
				print(n,"[",l,",",c,",",s,"], ")
			}
			print("\n")
		}
		*/
		time.Sleep(time.Duration(goCounter) * time.Nanosecond)  // Пропорциональное замедление передатчика
	}
	close(c); inChanOpen=false;	// Закрыть канал записи и выйти, это не мешает другим горутинам читать из этого канала
}


