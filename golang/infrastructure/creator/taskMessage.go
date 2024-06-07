package creator

type Stringable interface {
	String() string
}

type TaskMessage[TValue Stringable] interface {
	GetValue() TValue
	GetError() error
	IsError() bool
}
