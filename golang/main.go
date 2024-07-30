package main

import (
	"fmt"
	"slices"
	"time"
)

const (
	ErrTask = "Some error occured"

	SuccessTaskResult = "task has been successed"
	ErrTaskResult     = "something went wrong"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.
// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнены остальных тасков

// A Task represents a meaninglessness of our life
type Task struct {
	ID         int
	createTime string // время создания
	finishTime string // время выполнения
	taskResult []byte
}

func (t Task) ToError() error {
	return fmt.Errorf("Task id %d time %s, error %s", t.ID, t.createTime, t.taskResult)
}

func newTask() Task {
	finishTime := time.Now().Format(time.RFC3339)
	if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
		finishTime = ErrTask
	}

	return Task{
		ID:         int(time.Now().Unix()),
		createTime: finishTime,
	}
}

func createTasks(bufferSize int) (tasksCreatorChan chan Task) {
	tasksCreatorChan = make(chan Task, bufferSize)

	go func() {
		for _i := 0; _i < bufferSize; _i++ {
			tasksCreatorChan <- newTask()
		}
		close(tasksCreatorChan)
	}()

	return tasksCreatorChan
}

func Do(task Task) Task {
	tt, _ := time.Parse(time.RFC3339, task.createTime)
	if tt.After(time.Now().Add(-20 * time.Second)) {
		task.taskResult = []byte(SuccessTaskResult)
	} else {
		task.taskResult = []byte(ErrTaskResult)
	}

	task.finishTime = time.Now().Format(time.RFC3339Nano)
	time.Sleep(time.Millisecond * 150)

	return task
}

func Run(tasksChan chan Task) (resultChan chan Task) {
	resultChan = make(chan Task)
	go func() {
		for task := range tasksChan {
			resultChan <- Do(task)
		}
		close(resultChan)
	}()

	return resultChan
}

func Sort(resultChan chan Task) (doneTasks map[int]Task, errTasks []error) {
	// именнованные аргументы имеют nil, поэтому инициализируем
	doneTasks = make(map[int]Task)
	for task := range resultChan {
		// Некая обработка результата, но можно было флаг isError в task прописать
		// Вдруг API тасок такая будет
		if slices.Equal(task.taskResult, []byte(ErrTaskResult)) {
			errTasks = append(errTasks, task.ToError())
		} else {
			doneTasks[task.ID] = task
		}
	}

	time.Sleep(time.Second * 3)
	return
}

func main() {
	tasksSize := 10
	doneTasks, errTasks := Sort(Run(createTasks(tasksSize)))

	println("Errors:")
	for r := range errTasks {
		println(r)
	}

	println("Done tasks:")
	for r := range doneTasks {
		println(r)
	}
}
