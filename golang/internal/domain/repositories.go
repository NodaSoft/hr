package domain

type TaskRepository interface {
	List() map[string]Task
	Add(task Task)
}

type ErrorRepository interface {
	List() map[string]error
	Add(id string, err error)
}
