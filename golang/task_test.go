package main

import (
	"errors"
	"testing"
	"time"
)

func TestNewTask(t *testing.T) {
	//проверяем что обычные таски  ок работают
	since := time.Date(2000, 1, 1, 0, 0, 0, 2, time.Local)
	task := NewTask(since)
	if task.error != nil {
		t.Fail()
	}

	// логика ошибочных тасков сохранена
	since = time.Date(2000, 1, 1, 0, 0, 0, 1, time.Local)
	task = NewTask(since)

	if !errors.Is(task.error, SomeError) {
		t.Fail()
	}
}
