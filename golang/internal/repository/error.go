package repository

import (
	"sync"
)

type ErrorRepository struct {
	errors map[string]error
	mu     *sync.RWMutex
}

func NewErrorRepository() ErrorRepository {
	return ErrorRepository{
		errors: make(map[string]error),
		mu:     &sync.RWMutex{},
	}
}

func (repo ErrorRepository) List() map[string]error {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	return cloneMap[string, error](repo.errors)
}

func (repo ErrorRepository) Add(taskID string, err error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	repo.errors[taskID] = err
}
