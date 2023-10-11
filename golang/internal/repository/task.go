package repository

import (
	"sync"
	"task_service/internal/domain"
)

type TaskRepository struct {
	tasks map[string]domain.Task
	mu    *sync.RWMutex
}

func NewTaskRepository() TaskRepository {
	return TaskRepository{
		tasks: make(map[string]domain.Task),
		mu:    &sync.RWMutex{},
	}
}

func (repo TaskRepository) List() map[string]domain.Task {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	return cloneMap[string, domain.Task](repo.tasks)
}

func (repo TaskRepository) Add(task domain.Task) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	repo.tasks[task.ID] = task
}
