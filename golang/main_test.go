package main

import (
	"strings"
	"testing"
	"time"
)

func TestTaskWorker(t *testing.T) {
	tests := []struct {
		name string
		a    Ttype
		want string
	}{
		{"Test 1", Ttype{cT: time.Now().Format(time.RFC3339)}, "task has been successed"},
		{"Test 2", Ttype{cT: time.Now().Add(-30 * time.Second).Format(time.RFC3339)}, "something went wrong"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := taskWorker(tt.a)
			if string(got.taskRESULT) != tt.want {
				t.Errorf("taskWorker() = %v, want %v", string(got.taskRESULT), tt.want)
			}
		})
	}
}

func TestTaskSorter(t *testing.T) {
	tests := []struct {
		name string
		t    Ttype
		want string
	}{
		{"Test 1", Ttype{taskRESULT: []byte("task has been successed")}, "task has been successed"},
		{"Test 2", Ttype{taskRESULT: []byte("something went wrong")}, "something went wrong"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doneTasks := make(chan Ttype, 1)
			undoneTasks := make(chan error, 1)

			go taskSorter(tt.t, doneTasks, undoneTasks)

			select {
			case doneTask := <-doneTasks:
				if string(doneTask.taskRESULT) != tt.want {
					t.Errorf("taskSorter() = %v, want %v", string(doneTask.taskRESULT), tt.want)
				}
			case err := <-undoneTasks:
				if !strings.Contains(err.Error(), tt.want) {
					t.Errorf("taskSorter() = %v, want %v", err.Error(), tt.want)
				}
			}
		})
	}
}
