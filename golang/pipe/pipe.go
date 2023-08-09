package pipe

import (
	"context"
)

// конкурентно безопасная труба для обмена объектами.
type Interface interface {
	// отправить
	Send(ctx context.Context, input interface{}) error
	// получить
	Get(ctx context.Context) (interface{}, error)
	Close()
}
