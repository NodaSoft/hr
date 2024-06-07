package creator

import "context"

type TaskCreator[TTask Stringable] interface {
	Start(context.Context, chan<- TaskMessage[TTask], chan<- bool)
}
