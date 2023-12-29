package main

import (
    "fmt"
    "log"
    "sync"
    "time"
)

// EXERCISE:
// * make good code out of bad code;
// * it is important to preserve the logic of the appearance of erroneous tasks;
// * make proper multi-threading of task processing.
// Send updated code via merge-request.

// A Ttype represents a meaninglessness of our life
type Task struct {
    ID         int
	//Использование строк для представления времени не является идиоматическим в Go. Лучше использовать тип time.Time.
    CreatedAt  time.Time
    FinishedAt time.Time
    Result     string
    Error      error
}

func main() {
	//Функция Tasksorter проверяет, содержит ли результат задачи строку successed, что не является надежным способом проверки успешности задачи.
	//Функция Task_worker спит 150 миллисекунд, что может замедлить обработку задач.
	//Функция TaskCreturer создает задачи бесконечно. Лучше иметь способ остановить это.
    tasks := make(chan Task, 10)
    doneTasks := make(chan Task)
    undoneTasks := make(chan Task)

    var wg sync.WaitGroup

    // Create tasks
    wg.Add(1)
    go func() {
        defer wg.Done()
        for i := 0; i < 100; i++ {
            t := Task{
                ID:        i,
                CreatedAt: time.Now(),
            }
            if t.CreatedAt.Nanosecond()%2 > 0 {
                t.Error = fmt.Errorf("some error occurred")
            }
            tasks <- t
        }
        close(tasks)
    }()

    // Process tasks
    wg.Add(1)
    go func() {
        defer wg.Done()
        for t := range tasks {
            t.FinishedAt = time.Now()
            if t.Error == nil && t.FinishedAt.After(t.CreatedAt.Add(-20*time.Second)) {
				//Доступ к переменным result и err осуществляется из нескольких goroutines без синхронизации, что может привести к data races.
                t.Result = "task has been successful"
                doneTasks <- t
            } else {
                t.Result = "something went wrong"
                undoneTasks <- t
            }
        }
        close(doneTasks)
        close(undoneTasks)
    }()

    // Collect results
    wg.Add(1)
    go func() {
        defer wg.Done()
        for t := range doneTasks {
			//Функция println используется для печати ошибок и результатов, что не является идиоматическим в Go. Лучше использовать log.
            log.Printf("Done task: %v", t)
        }
        for t := range undoneTasks {
            log.Printf("Undone task: %v", t)
        }
    }()

    wg.Wait()
}
