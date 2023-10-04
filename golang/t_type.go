package main

// A Ttype represents a meaninglessness of our life
type Ttype struct {
	id         int
	cT         string // время создания
	fT         string // время выполнения
	taskRESULT []byte
}

func (t *Ttype) Success() bool {
	return len(t.taskRESULT) > 14 && string(t.taskRESULT[14:]) == "successed"
}
