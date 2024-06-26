package core

import "context"

// Simple and useful adapters to not create struct with method for every function
func FactoryToPipeElem(f func(context.Context, chan Task, TaskFactory) chan Task, factory TaskFactory) PipilineElement {
	return func(ctx context.Context, tch chan Task) chan Task {
		return f(ctx, tch, factory)
	}
}
func HandlerToPipeElem(f func(ctx context.Context, tasks chan Task, worker TaskWorker) chan Task, worker TaskWorker) PipilineElement {
	return func(ctx context.Context, tch chan Task) chan Task {
		return f(ctx, tch, worker)
	}
}
func SeparatorToPipeElem(f func(ctx context.Context, tch chan Task, separator func(Task) int8, separated chan error) chan Task, separator func(Task) int8, separated chan error) PipilineElement {
	return func(ctx context.Context, tch chan Task) chan Task {
		return f(ctx, tch, separator, separated)
	}
}
