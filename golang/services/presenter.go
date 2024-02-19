package services

import (
	"fmt"
	"nodasoft-golang/domain"
)

type presenter struct {
	tasks <-chan *domain.Task
}

func newPresenter(tasks <-chan *domain.Task) *presenter {
	return &presenter{tasks}
}

func (p *presenter) print() {
	for t := range p.tasks {
		if t.IsSuccessful() {
			p.printSuccessfulTask(t)
		} else {
			p.printErroredTask(t)
		}
	}
}

func (p *presenter) printSuccessfulTask(t *domain.Task) {
	fmt.Printf("task ID: %d, createdAt: %s, status: succesfully finished\n", t.ID, t.CreatedAt.String())
}

func (p *presenter) printErroredTask(t *domain.Task) {
	fmt.Printf("task ID: %d, createdAt: %s, status: failed\n", t.ID, t.CreatedAt.String())
}
