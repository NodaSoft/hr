package worker

import (
	"context"
	"testing"
	"time"

	"github.com/Quantum12k/hr/golang/internal/task"
)

func TestWorker_run(t *testing.T) {
	type fields struct {
		pendingTasksCh chan *task.Task
		DoneTasksCh    chan *task.Task
	}

	type args struct {
		timeout time.Duration
	}

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "normal_task_handling",
			fields: fields{
				pendingTasksCh: make(chan *task.Task),
				DoneTasksCh:    make(chan *task.Task),
			},
			args: args{
				timeout: time.Second,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), tt.args.timeout)
			defer cancel()

			w := &Worker{
				pendingTasksCh: tt.fields.pendingTasksCh,
				DoneTasksCh:    tt.fields.DoneTasksCh,
			}

			go w.run(ctx)

			go func() {
				for {
					select {
					case <-ctx.Done():
						return
					default:
						tt.fields.pendingTasksCh <- task.New()
					}
				}
			}()

			for range w.DoneTasksCh {
			}

			<-ctx.Done()
		})
	}
}
