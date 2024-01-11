package main

import (
	"testing"
)

func TestNewTaskManager(t *testing.T) {
	tm := NewTaskManager()
	if tm == nil {
		t.Errorf("NewTaskManager() = %v, want non-nil", tm)
	}
}

func TestStartWorkerPool(t *testing.T) {
	tm := NewTaskManager()
	tm.StartWorkerPool(5)
}

func TestPrintResults(t *testing.T) {
	tm := NewTaskManager()
	tm.GenerateTasks(5)
	tm.StartWorkerPool(5)
	tm.PrintResults()
	tm.WaitForCompletion()
}

func TestGenerateTasksN(t *testing.T) {
	tm := NewTaskManager()
	tm.generateTasksN(5)
	if len(tm.tasks) != 5 {
		t.Errorf("GenerateTasksN(5) = %v, want 5", len(tm.tasks))
	}
}
