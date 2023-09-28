package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"
)

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// A Task represents a meaninglessness of our life
//
type Task struct {
	id     int
	cT     string // время создания
	fT     string // время выполнения
	result []byte
}
// Метод (Task).String(), имплементирует интерфейс Stringer из модуля fmt
func (t Task) String() string {
	return fmt.Sprintf("Task { id:%d, cT:%s, fT:%s, result:%s }", t.id, t.cT, t.fT, t.result)
}
// Метод (Task).Successed(), возвращает булевское значение таска - успешность
func (t Task) Successed() bool {
	return string(t.result[14:]) == "successed"
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// TaskManager:
// The Solution demonstrates the skills in threads, pipes, slices, pointers, interfaces and
// object-oriented programming methods in Go, and also represents a meaninglessness of our life
//
type TaskManager struct {
	stop		bool			// стоп - сигнал 
	bufCap		int				// Ёмкость кольцевого буфера 
	iCount		int				// Максимальное число итераций. Если задать <=0 то выполняется бесконечный цикл
	iCounter	int				// Счетчик итераций. В случае бесконечного цикла, не используется
	superChan	chan Task		// Основной канал.
	doneTasks	chan Task		// Канал успешных тасков
	undoneTasks chan error		// Канал ошибок
	results 	map[int]Task	// Карта результатов
	errors		[]error			// Список ошибок
}

// Инициализация и Выполнение
func (m *TaskManager) Run () {
	if(m.bufCap == 0){ m.bufCap = 10 }								// Значения по умолчанию, если не заданы явно
	if(m.iCount == 0){ m.iCount = m.bufCap * 2 }					// 

	stdout, _ := exec.Command("sh", "-c", "ulimit -Sn").Output()	// Ограничения системы на софт, понадобится в следующей версии с реальной нагрузкой на воркеров
	go_limit_s, _ := strconv.Atoi(string(stdout[:len(stdout)-1]))  	// Удаляем перенос строки, приводим к типу int
	stdout, _  = exec.Command("sh", "-c", "ulimit -Hn").Output()	// Ограничения системы на хард, понадобится в следующей версии с реальной нагрузкой на БД
	go_limit_h, _ := strconv.Atoi(string(stdout[:len(stdout)-1]))  	// Удаляем перенос строки, приводим к типу int

	m.stop = false
	go m.Switch()											// Ожидает сигнал выключения и меняет значение m.stop = true

	m.superChan = make(chan Task, m.bufCap)					// Канал поступающих тасков
	m.doneTasks = make(chan Task, m.bufCap/2)				// Канал завершенных тасков
	m.undoneTasks = make(chan error, m.bufCap/2)			// Канал тасков с ошибками
	m.results = map[int]Task{}								// Карта результатов по id таска
	m.errors  = []error{}									// Список ошибок

	println("Буфер канала входящих таксов: ",m.bufCap)
	println("Буфер канала завершенных тасков: ",m.bufCap/2)
	println("Буфер канала ошибок: ",m.bufCap/2)
	println("Лимит вашей системы на активные процессы    (soft): ",go_limit_s)
	println("Лимит вашей системы на файловые дескрипторы (hard): ",go_limit_h)
	if m.iCount > 0 {
		println("Для выхода найжмите любую клавишу либо дождитесь завершения цикла (около 15 сек.): ",m.iCount)
	} else {
		println("Бесконечный цикл! Для выхода найжмите любую клавишу.")
	}

	var wg sync.WaitGroup									// счетчик горутин верхнего уровня
	wg.Add(4)												// Инкремент на эти 4 горутины
	go m.errorsCollector(&wg,m.errors,m.undoneTasks)		// Копит ошибки
	go m.resultCollector(&wg,m.results,m.doneTasks)			// Систематизирует результаты
	go m.Receiver(&wg,m.superChan,m.doneTasks,m.undoneTasks)// Принимает задачи, передает на обработку
	go m.Sender(&wg,m.superChan)							// Отправляет задачи
	wg.Wait()	 											// Ожидает завершения 4 горутин верхнего уровня
}

// Коллектор ошибок
func (m *TaskManager) errorsCollector (wg *sync.WaitGroup, err []error, c <-chan error) {
	defer wg.Done()
	for e := range c {
		m.errors = append(m.errors, e)
	}
}

// Коллектор успешных тасков
func (m *TaskManager) resultCollector (wg *sync.WaitGroup, res map[int]Task, c <-chan Task) {
	defer wg.Done()
	for task := range c {
		res[task.id] = task
	}
}

// Воркер с полезной нагрузкой
func (m *TaskManager) Worker (task Task) Task {
	t, _ := time.Parse(time.RFC3339, task.cT)
	if t.After(time.Now().Add(-20 * time.Second)) { // таски с ошибкой имеют другой формат cT и попадают в else
		task.result = []byte("task has been successed")
	} else {
		task.result = []byte("something went wrong")
	}
	time.Sleep(time.Millisecond * 150)				// Имитация длительной работы, а не вынужденный Sleep.
	task.fT = time.Now().Format(time.RFC3339Nano)	// Время завершения задачи
	return task
}

// Сортировщик тасков
func (m *TaskManager) Sorter (wg *sync.WaitGroup, task Task, done chan<- Task, undo chan<- error) {
	defer wg.Done()
	task = m.Worker(task)			// Ожидаем завершение синхроного вызова в рамках этой горутины.
	switch task.Successed() {		// Эквивалентно if else
	  case true:
		done <- task
	  default:
		undo <- fmt.Errorf(task.String()) // явный вызов .String(), т.к. fmt.Errorf не поддерживает интерфейс Stringer
	}
}

// Приемник тасков
func (m *TaskManager) Receiver (wg *sync.WaitGroup, c <-chan Task, done chan Task, undo chan error) {
	defer wg.Done()
	var wg2 sync.WaitGroup		// локальная группа дочерних горутин
	for task := range c {		// продолжается до закрытия и опустошения канала. При чтении из пустого и открытого, канал блокируется
		wg2.Add(1)				// Инкремент дочерних горутин
		go m.Sorter(&wg2,task,done,undo) // Обработать каждую полученную задачу в отдельном потоке
	}
	wg2.Wait()					// Ждать завершение дочерних горутин
	close(done)					// Закрыть каналы, чтобы НЕдочерние читающие горутины могли завершиться 
	close(undo)
}

// Выключатель. Ждет ненулевой сигнал и устанавливает переменную выключения
func (m *TaskManager) Switch () {
	os.Stdin.Read(make([]byte,1))	// ожидает ввода команды с клавиатуры
	m.stop = true					// приведет к завершению цикла Sender
	println("Pressed a Key, stop =", m.stop)
}
// Вариант без ввода с клавиатуры, сигнал или команда поступает из канала
func (m *TaskManager) Switch2 (c <-chan byte) {
	for range c {		// цикл заблокирован пока не поступит сигнал
		m.stop = true	// приведет к завершению цикла Sender
		break
	}
}

// Опциональный Счетчик цикла. При iCount <= 0 ничего не делает, только возвращает значение.
// при iCount <= 0 блок не исполняется, что приводит к бесконечному циклу в Sender
func (m *TaskManager) StopCounter() bool {
	if m.iCount > 0 {	   
		m.iCounter++
		if m.iCounter >= m.iCount {	m.stop = true }
	}
	return m.stop
}

// Передатчик тасков
func (m *TaskManager) Sender (wg *sync.WaitGroup, c chan<- Task) {
	defer wg.Done()
	// Т.к. размер буфера задан явно и нужно в конце вывести короткие списки, ограничиваем цикл либо емкостью буфера либо
	// числом итераций iCount. В последнем случае емкость буфера можно сделать меньше, а число итераций больше.
	// В реальной задаче цикл можно сделать бесконечным, а для выключения использовать stop-сигнал
	for {
		if m.StopCounter() { break }		// stop-сигнал и опциональный счетчик итераций совмещены
		t := time.Now()						// фиксируем время
		tf := t.Format(time.RFC3339)		// берем его так
		id := int(t.Nanosecond())			// и сяк, один раз, для соответствия проверяемого и сохраняемого значения.
		if id%2 > 0 {						// "вот такое искусственное условие появления ошибочных тасков"
			tf = "Some error occured"		// таски с нечетным временем получают признак ошибки
		}
		c <- Task{cT: tf, id: id}			// создаем экземпляр и передаем таск в канал.
	}
	close(c)	// Закрыть канал записи и выйти, это не мешает другим горутинам читать из этого канала
}

// Печатает ошибки и успешные таски
func (m *TaskManager) Log () {
	println("\x1b[31m" + "Errors: ",len(m.errors))				// Список ошибок красным цветом
	for i := range m.errors {
		if i > m.bufCap {
			break
		}
		println(m.errors[i].Error())
	}
	fmt.Println("\x1b[32m" + "Done tasks: ",len(m.results))		// Список успешных результатов зеленым цветом
	i:=0
	for key := range m.results{
		i++
		if i > m.bufCap {
			break
		}
		fmt.Println(m.results[key]) // fmt.Println умеет сам вызывать метод Task.String() т.к. имплементирован интерфейс Stringer
	}
	fmt.Println("\x1b[0m" + "The End")					// Конец программы, цвет по умолчанию
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Приложение эмулирует получение и обработку тасков в асинхронном многопоточном режиме
// В конце выводятся успешные таски и ошибки выполнения остальных тасков (при большом количестве выводится только часть и общее количество)
// Вся логика инкапсулирована в типе-классе TaskManager
func main() {
	m := TaskManager {						// Создаем Менеджера тасков с параметрами. Все параметры опциональны
		 bufCap:20,							// ёмкость кольцевого буфера, это же значение используется для ограничения вывода результатов
		 iCount:5000000,					// число итераций = тасков. Если <=0, то бесконечно!
	}
	m.Run()									// Запускаем. Синхронный вызов, ждем окончания работы или нажимаем любую клавишу
	m.Log()									// Выводим лог
}

