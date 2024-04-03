package main

import (
	"context"
	"fmt"
	"time"
)

type ResultCollector struct {
	taskProcessor *TaskProcessor

	doneTasks   []*Task
	undoneTasks []*Task
}

func NewResultCollector(puller Puller) *ResultCollector {
	collector := &ResultCollector{}
	collector.taskProcessor = &TaskProcessor{
		pull:    puller,
		process: collector.processTask,
		push:    func(*Task) error { return nil },
	}

	return collector
}

func (collector *ResultCollector) processTask(task *Task) {
	if task.err == nil {
		collector.doneTasks = append(collector.doneTasks, task)
	} else {
		collector.undoneTasks = append(collector.undoneTasks, task)
	}
}

func (collector *ResultCollector) printDone() {
	fmt.Println("Done tasks:")

	for _, task := range collector.doneTasks {
		fmt.Println(task.id)
	}
}

func (collector *ResultCollector) printErrors() {
	fmt.Println("Errors:")

	for _, task := range collector.undoneTasks {
		fmt.Printf("Task id %d time %s, error %s\n", task.id, task.createdAt.Format(time.RFC3339Nano), task.err)
	}
}

func (collector *ResultCollector) Run(ctx context.Context) {
	collector.taskProcessor.Run(ctx)
}
