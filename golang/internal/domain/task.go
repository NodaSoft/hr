package domain

// A Task represents the full meaning of our life
type Task struct {
	ID            string
	CreationTime  string
	ExecutionTime string
	Result        WokrkResult
}

type WokrkResult []byte
