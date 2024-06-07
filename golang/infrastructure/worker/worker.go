package worker

import (
	"context"
)

type Worker[TIn any, TOut any] interface {
	Work(context.Context, <-chan TIn, chan<- TOut, chan<- bool)
}
