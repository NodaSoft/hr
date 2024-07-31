package main

import "context"

// Pipe is a function that takes a channel, processes its values,
// and returns a new channel with the processed values from the old one.
type Pipe[T any] func(ctx context.Context, in <-chan T) <-chan T

// Source is a function that sends values in the output channel.
type Source[T any] func(ctx context.Context) <-chan T

// StartPipeline links all the pipelines into one big pipeline and starts Source.
func StartPipeline[T any](ctx context.Context, start Source[T], pipes ...Pipe[T]) <-chan T {
	out := start(ctx)
	for _, p := range pipes {
		out = p(ctx, out)
	}

	return out
}

// DiscardPipe is a pipeline for discarding everything that passes into its input.
// It returns nil output channel.
func DiscardPipe[T any](ctx context.Context, in <-chan T) <-chan T {
	go func() {
		for range in {
			if err := ctx.Err(); err != nil {
				break
			}
		}
	}()

	return nil
}
