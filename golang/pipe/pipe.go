package pipe

import (
	"context"
)

// конкурентно безопасная труба для обмена объектами.
type Interface interface {
	Send(ctx context.Context, input interface{}) error
	Get(ctx context.Context) (interface{}, error)
	Close()
}
