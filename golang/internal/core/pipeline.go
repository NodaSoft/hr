package core

import "context"

type PipilineElement func(context.Context, chan Task) chan Task

// Invoke all elements, return channel , that last element provided
func Pipeline(ctx context.Context, tch chan Task, elements ...PipilineElement) chan Task {
	for _, e := range elements {
		tch = e(ctx, tch)
	}
	return tch
}
