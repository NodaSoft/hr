package creator

import (
	"context"
	"testing"
	"time"

	"github.com/Quantum12k/hr/golang/internal/task"
)

func TestCreator_run(t *testing.T) {
	type fields struct {
		NewTasksCh chan *task.Task
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
			name: "normal_task_generation",
			fields: fields{
				NewTasksCh: make(chan *task.Task),
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

			c := &Creator{
				NewTasksCh: tt.fields.NewTasksCh,
			}

			go c.run(ctx)

			for range c.NewTasksCh {
			}

			<-ctx.Done()
		})
	}
}
