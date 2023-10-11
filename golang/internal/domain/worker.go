package domain

type Worker interface {
	Handle(Task) (Task, error)
}
