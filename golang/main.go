package main

import (
  "fmt"
  "github.com/google/uuid"
  "log"
  "sync"
  "time"
)

// task represents a meaninglessness of our life
type task struct {
  id         uuid.UUID  // id: time.Now().UnixNano() / 1000 в момент создания таски не гарантирует уникальных id
  createdAt  *time.Time // createdAt: время создания
  finishedAt *time.Time // finishedAt: время выполнения
}

// processTask: какая-то логика с задачей которая выполняется какое-то время
func (t *task) processTask() {
  time.Sleep(time.Millisecond * 150)
  now := time.Now()
  if t.createdAt != nil {
    t.finishedAt = &now
  }
}

// taskSort: решает успешно или нет выполнена задача и отправляет в лог
func (t *task) taskSort(wg *sync.WaitGroup, logTask chan<- string) {
  defer wg.Done()

  t.processTask()
  if t.createdAt == nil {
    logTask <- fmt.Sprintf("task id %s failed", t.id.String())
  } else {
    logTask <- fmt.Sprintf("task id %s succeeded at %s time process task %s", t.id.String(), t.finishedAt.Format(time.RFC3339Nano), t.finishedAt.Sub(*t.createdAt).String())
  }
}

// taskCreator: создание задач
func taskCreator(wg *sync.WaitGroup, sleepTime time.Duration, tasks chan task) {
  defer wg.Done()
  startTime := time.Now()
  for ; startTime.After(time.Now().Add(-sleepTime)); {
    createdTime := time.Now()
    if createdTime.UnixNano()/1000%2 == 0 { // добавил /1000 чтобы убрать доли наносекунд иначе на конце всегда 0 и условие всегда положительное
      tasks <- task{id: uuid.New(), createdAt: &createdTime}
    } else {
      tasks <- task{id: uuid.New()} // при фейле поле createdAt: nil что позволяет сократить время вычисления в методе processTask
    }
  }

  return
}

// process: чтение из каналов и обработка
func process(wg *sync.WaitGroup, tasks chan task, logTask chan string) {
  for {
    select {
    case t := <-tasks:
      wg.Add(1)
      go t.taskSort(wg, logTask)
    case l := <-logTask:
      log.Println(l)
    }
  }
}

func main() {
  tasks := make(chan task, 10)
  logTask := make(chan string)

  var wg sync.WaitGroup

  wg.Add(1)
  go taskCreator(&wg, time.Second*3, tasks)
  go process(&wg, tasks, logTask)

  wg.Wait()

  close(logTask)
  close(tasks)
  log.Println("done!")
}
