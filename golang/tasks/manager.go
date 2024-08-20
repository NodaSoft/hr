package tasks

import (
	"errors"
	"fmt"
	"sync"
)

var (
	ErrTaskFailed = errors.New("task failed")
)

type Manager struct {
	mu        sync.RWMutex
	failed    []*Task
	successed []*Task
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) AddSuccessed(task *Task) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.successed = append(m.successed, task)
}

func (m *Manager) AddFailed(task *Task) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.failed = append(m.failed, task)
}

func (m *Manager) Successed() []*Task {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.successed
}

func (m *Manager) Failed() []*Task {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.failed
}

func (m *Manager) Print() {
	m.mu.RLock()
	defer m.mu.RUnlock()

	fmt.Println("Failed tasks:")
	for _, task := range m.failed {
		fmt.Printf("%d (%s)\n", task.ID, task.Err.Error())
	}

	fmt.Println("Done tasks:")
	for _, task := range m.successed {
		fmt.Printf("%d\n", task.ID)
	}
}
