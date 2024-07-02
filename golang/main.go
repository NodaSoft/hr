package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sync"
	"time"
)

const (
	SERVICE_GENERATE_TASKS_FOR_X_SECONDS        = 10
	SERVICE_PRINT_HANDLED_TASKS_EVERY_X_SECONDS = 3

	TASKS_CHANNEL_BUFFER                    = 10
	TASK_HANDLING_TIMEOUT_THRESHOLD_SECONDS = 20
)

var (
	ai AutoIncrement
)

type AutoIncrement struct {
	sync.Mutex
	id int
}

func (ai *AutoIncrement) ID() int {
	ai.Lock()
	defer ai.Unlock()

	id := ai.id
	ai.id++
	return id
}

type Task struct {
	id         int
	createdAt  time.Time
	finishedAt time.Time

	failed  bool
	message string
}

type Service struct {
	output sync.Mutex
	wg     sync.WaitGroup

	context   context.Context
	tasks     chan *Task
	succesful []*Task
	failed    []*Task
}

func (s *Service) taskCreator() {
	for {
		select {
		case <-s.context.Done():
			return
		default:
			task := &Task{
				id:        ai.ID(),
				createdAt: time.Now(),
			}

			if time.Now().Nanosecond()%2 > 0 {
				task.failed = true
				task.message = "Some error occured"
			}

			s.tasks <- task
		}
	}
}

func (s *Service) tasksHandler() {
	for t := range s.tasks {
		go s.handleTask(t)
	}
}

func (s *Service) tasksOutput() {
	ticker := time.NewTicker(time.Second * SERVICE_PRINT_HANDLED_TASKS_EVERY_X_SECONDS)

	for {
		select {
		case <-s.context.Done():
			return
		case <-ticker.C:
			go func() {
				s.wg.Add(1)
				defer s.wg.Done()

				f := bufio.NewWriterSize(os.Stdout, 1<<16)
				defer f.Flush()

				s.output.Lock()
				failed := s.failed
				succesful := s.succesful
				s.failed = make([]*Task, 0)
				s.succesful = make([]*Task, 0)
				s.output.Unlock()

				for _, t := range failed {
					f.WriteString(fmt.Sprintf("Task id %d, time %s, error: %s\n", t.id, t.createdAt.Format(time.RFC3339), t.message))
				}

				for _, t := range succesful {
					f.WriteString(fmt.Sprintf("Task id %d, time %s, message: %s\n", t.id, t.createdAt.Format(time.RFC3339), t.message))
				}
			}()
		}
	}
}

func (s *Service) handleTask(t *Task) {
	if !t.failed {
		if t.createdAt.After(time.Now().Add(time.Second * -TASK_HANDLING_TIMEOUT_THRESHOLD_SECONDS)) {
			t.message = "Task has been successful"
		} else {
			t.failed = true
			t.message = "Something went wrong, timeout"
		}
	}

	t.finishedAt = time.Now()

	time.Sleep(time.Millisecond * 150)
	s.sortTask(t)
}

func (s *Service) sortTask(t *Task) {
	s.output.Lock()
	if t.failed {
		s.failed = append(s.failed, t)
	} else {
		s.succesful = append(s.succesful, t)
	}
	s.output.Unlock()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	service := &Service{
		context:   ctx,
		tasks:     make(chan *Task, TASKS_CHANNEL_BUFFER),
		succesful: make([]*Task, 0),
		failed:    make([]*Task, 0),
	}

	go service.taskCreator()
	go service.tasksHandler()
	go service.tasksOutput()

	// Task creation stopped
	time.Sleep(time.Second * SERVICE_GENERATE_TASKS_FOR_X_SECONDS)
	cancel()

	// Waiting for all the remaining prints
	service.wg.Wait()
}
