package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Task:
// * Refactor to good coding practices
// * Keep failing tasks logic
// * Make proper multiprocessing for tasks handling
//
// Send all the fixes to GitHub via PR
//
// The app emulates tasks handling and attempts to receive and handle
// in multiple goroutines.
//
// Processed and errors of failed tasks could be printed upon exit

type config struct {
	LogLevel log.Level `envconfig:"LOG_LEVEL"`
}

type Task struct {
	id        uuid.UUID
	createdAt time.Time
	startedAt time.Time
	error     error
}

const (
	// TTL between task was created and executed
	taskExecutionMaxDelay = 20 * time.Second

	// just a delay for task execution emulation
	taskProcessingTime = 150 * time.Millisecond
)

func main() {
	cfg := config{
		LogLevel: log.WarnLevel,
	}

	if err := envconfig.Process("", &cfg); err != nil {
		panic(err)
	}

	log.SetLevel(cfg.LogLevel)

	taskGenerator := func(a chan Task) {
		go func() {
			for {
				now := timeNow()
				taskID := uuid.New()

				var err error
				if now.Nanosecond()%2 > 0 {
					err = errors.New("not even nanosecond: failed task")
				}

				log.WithFields(log.Fields{
					"component": "taskGenerator",
					"id":        taskID.String(),
					"error":     err,
				}).Trace("task generated")

				a <- Task{
					id:        taskID,
					createdAt: now,
					error:     err,
				}
			}
		}()
	}

	superChan := make(chan Task, 10)

	go taskGenerator(superChan)

	taskWorker := func(t Task) Task {
		log.WithFields(log.Fields{
			"component": "taskWorker",
			"id":        t.id.String(),
		}).Trace("task received")

		if t.createdAt.Before(timeNow().Add(-taskExecutionMaxDelay)) {
			t.error = errors.New("something went wrong: task is expired")
		}
		t.startedAt = timeNow()

		time.Sleep(taskProcessingTime)

		return t
	}

	doneTasks := make(chan Task)
	undoneTasks := make(chan error)

	taskSorter := func(t Task) {
		log.WithFields(log.Fields{
			"component": "taskSorter",
			"id":        t.id.String(),
			"error":     t.error,
		}).Trace("task passed for sorting")

		if t.error != nil {
			undoneTasks <- errors.Errorf(
				"Task id %s time %s, error %s",
				t.id.String(), t.createdAt, t.error.Error(),
			)
		} else {
			doneTasks <- t
		}
	}

	go func() {
		for t := range superChan {
			t = taskWorker(t)
			go taskSorter(t)
		}
		close(superChan)
	}()

	result := map[uuid.UUID]Task{}
	errs := []error{}
	resultMutex := &sync.Mutex{}
	errsMutex := &sync.Mutex{}

	go func() {
		for task := range doneTasks {
			go func(task Task) {
				resultMutex.Lock()
				defer resultMutex.Unlock()

				result[task.id] = task
			}(task)
		}
		for err := range undoneTasks {
			go func(err error) {
				errsMutex.Lock()
				defer errsMutex.Unlock()

				errs = append(errs, err)
			}(err)
		}
		close(doneTasks)
		close(undoneTasks)
	}()

	time.Sleep(time.Second * 3)

	// Func here just to isolate defers
	func() {
		resultMutex.Lock()
		errsMutex.Lock()
		defer resultMutex.Unlock()
		defer errsMutex.Unlock()

		fmt.Println("Errors:")
		for err := range errs {
			fmt.Println(err)
		}

		fmt.Println("Done tasks:")
		for r := range result {
			fmt.Println(r)
		}
	}()
}

func timeNow() time.Time {
	return time.Now().UTC()
}
