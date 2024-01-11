package main

import (
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// Эта программа на языке Go реализует менеджер задач. Она создает задачи,
// обрабатывает их и выводит результаты.
//
// Основные компоненты программы:
// main(): Запускает менеджер задач, генерирует задачи, запускает пул обработчиков и
// выводит результаты.
//
// Task: Интерфейс для задачи, который определяет методы для обработки задачи, проверки
// на ошибки и вывода результатов.
//
// SimpleTask: Структура, реализующая интерфейс Task. Она имитирует выполнение задачи,
// проверяет наличие ошибок и выводит результаты.
//
// TaskManager: Структура для управления задачами. Она генерирует задачи, запускает пул
// обработчиков, ожидает завершения всех задач и выводит результаты.
//
// setupSlog(): Настраивает логирование.
//
// Все задачи генерируются асинхронно и обрабатываются в пуле горутин. Результаты каждой задачи
// сохраняются и выводятся после завершения всех задач.

const (
	timeOutputFormat   = "2006-01-02 15:04:05" // формат вывода времени
	generateMaxTime    = 10                    // максимальное время генерации задачи
	executeMaxTime     = 300                   // максимальное время выполнения задачи
	genereteTasksCount = 1000                  // количество генерируемых задач
	genereteTasksN     = 100                   // количество генерируемых задач
	genereteTasksTime  = 3 * time.Second       // время генерации задач
	workerPoolSize     = 10                    // количество горутин для обработки задач
	logLevel           = slog.LevelDebug       // уровень логирования
	// logLevel = slog.LevelError
)

var (
	appName    = "nodasofthr"
	appVersion = "debug"
)

// main запускает менеджер задач и ожидает завершения всех задач, после чего распечатывает результаты
func main() {
	setupSlog()
	slog.Info("Starting", "os", runtime.GOOS, "arch", runtime.GOARCH, "version", appVersion, "app", appName)
	taskManager := NewTaskManager()              // создаем менеджер задач
	taskManager.GenerateTasks(genereteTasksTime) // генерируем задачи в зависимости от типа параметра (genereteTasksTime или genereteTasksN)
	taskManager.StartWorkerPool(workerPoolSize)  // запускаем пул обработчиков из N горутин
	taskManager.PrintResults()                   // печатаем результаты
	taskManager.WaitForCompletion()              // ожидаем завершения всех задач
	slog.Info("Finished", "os", runtime.GOOS, "arch", runtime.GOARCH, "version", appVersion, "app", appName)
}

// setupSlog настраивает логирование
func setupSlog() {
	replace := func(_ []string, a slog.Attr) slog.Attr {
		if a.Key == "time" {
			t := a.Value.Time()
			return slog.String("time", t.Format(timeOutputFormat))
		}
		return a
	}
	h := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:       logLevel,
		AddSource:   appVersion == "debug",
		ReplaceAttr: replace,
	})
	slog.SetDefault(slog.New(h))
}

// Task определяет интерфейс для задачи
type Task interface {
	Process()
	IsError() error
	PrintResult()
}

// NewSimpleTask создает новую задачу
func NewSimpleTask() *SimpleTask {
	task := &SimpleTask{}
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(generateMaxTime))) // получаем рандомное время генерации задачи
	if runtime.GOOS == "windows" {
		task.id = int(time.Now().UnixMicro()) // получаем уникальный ID задачи
	} else {
		task.id = int(time.Now().UnixNano()) // получаем уникальный ID задачи
	}
	task.createTime = time.Now()
	return task
}

// SimpleTask реализует интерфейс Task
type SimpleTask struct {
	id         int
	createTime time.Time
	finishTime time.Time
	err        error
}

// Process симулирует выполнение задачи
func (t *SimpleTask) Process() {
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(executeMaxTime))) // получаем рандомное время выполнения задачи
	t.finishTime = time.Now()
	result := "ok"
	ruleErrorResult := false // правило для генерации ошибки
	if runtime.GOOS == "windows" {
		ruleErrorResult = t.finishTime.UnixMicro()%2 > 0 // в windows Nanosecond() всегда возвращает 0 (https://github.com/golang/go/issues/28084)
	} else {
		ruleErrorResult = t.finishTime.Nanosecond()%2 > 0
	}
	if ruleErrorResult {
		result = "process error"
		t.err = fmt.Errorf(result)
	}
	slog.Debug("Task processed", "id", t.id, "result", result, "started", t.createTime.Format(timeOutputFormat), "finished", t.finishTime.Format(timeOutputFormat))
}

// IsError возвращает ошибку
func (t *SimpleTask) IsError() error {
	return t.err
}

// PrintResult печатает результат
func (t *SimpleTask) PrintResult() {
	if t.err != nil {
		fmt.Printf("Task ID: %d, started at: %s, error: %s\n", t.id, t.createTime.Format(timeOutputFormat), t.err)
	} else {
		fmt.Printf("Task ID: %d, started at: %s, finished at: %s\n", t.id, t.createTime.Format(timeOutputFormat), t.finishTime.Format(timeOutputFormat))
	}
}

// TaskManager определяет структуру для управления задачами
type TaskManager struct {
	tasks          chan Task      // Очередь задач
	results        chan Task      // Очередь результатов
	tasksWg        sync.WaitGroup // Для ожидания завершения всех задач
	printWg        sync.WaitGroup // Для ожидания завершения печати результатов
	generatedTasks int32          // Количество cгенерированных задач
	processedTasks int32          // Количество обработанных задач
}

// NewTaskManager создает новый объект TaskManager
func NewTaskManager() *TaskManager {
	return &TaskManager{
		tasks:   make(chan Task, genereteTasksCount),
		results: make(chan Task),
	}
}

// GenerateTasks генерирует задачи в зависимости от типа параметра (принимает int или time.Duration) (generic функция)
func (m *TaskManager) GenerateTasks(numTasks interface{}) {
	go func(numTasks interface{}) {
		switch taskType := numTasks.(type) {
		case int:
			m.generateTasksN(taskType)
		case time.Duration:
			m.generateTasksSeconds(taskType)
		default:
			slog.Error("Unsupported type", "type", fmt.Sprintf("%T", numTasks))
			os.Exit(1)
		}
	}(numTasks)
}

// generateTasksN генерирует фиксированное количество задач
func (m *TaskManager) generateTasksN(numTasks int) {
	for i := 0; i < numTasks; i++ {
		task := NewSimpleTask()
		m.tasks <- task
		atomic.AddInt32(&m.generatedTasks, 1)
		slog.Debug("Task generated", "id", task.id, "task_num", i+1, "total", numTasks)
	}
	close(m.tasks)
}

// generateTasksSeconds генерирует задачи в течение заданного времени
func (m *TaskManager) generateTasksSeconds(seconds time.Duration) {
	i := 1
	deadline := time.Now().Add(seconds)
	for time.Now().Before(deadline) {
		m.tasks <- NewSimpleTask()
		atomic.AddInt32(&m.generatedTasks, 1)
		// логируем время генерации задачи
		slog.Debug("Generated", "id", i, "task num", i, "stop_after", time.Until(deadline).Seconds())
		i++
	}
	close(m.tasks)
}

// StartWorkerPool запускает пул работ
func (m *TaskManager) StartWorkerPool(numWorkers int) {
	for i := 0; i < numWorkers; i++ {
		m.tasksWg.Add(1)
		go m.worker()
	}
}

// worker запускает обработчик задач
func (m *TaskManager) worker() {
	defer m.tasksWg.Done()
	for task := range m.tasks {
		task.Process()
		atomic.AddInt32(&m.processedTasks, 1)
		m.results <- task
	}
}

// WaitForCompletion ожидает завершения всех задач
func (m *TaskManager) WaitForCompletion() {
	m.tasksWg.Wait() // Ожидаем завершения всех обработчиков задач
	close(m.results) // Теперь закрываем канал Results
	m.printWg.Wait() // Ожидаем завершения печати результатов
	slog.Info("All tasks completed", "generated", m.generatedTasks, "processed", m.processedTasks)
}

// PrintResults печатает результаты
func (m *TaskManager) PrintResults() {
	m.printWg.Add(1)
	go func() {
		defer m.printWg.Done()
		var Successes []Task
		var Errors []Task
		for result := range m.results {
			if result.IsError() == nil {
				Successes = append(Successes, result)
			} else {
				Errors = append(Errors, result)
			}
		}

		fmt.Printf("\n--- Completed Tasks (%d) --- \n", len(Successes))
		for _, r := range Successes {
			r.PrintResult()
		}

		fmt.Printf("\n--- Error Tasks (%d) --- \n", len(Errors))
		for _, r := range Errors {
			r.PrintResult()
		}
	}()
}
