package task

import (
	"fmt"
	"strings"
	"time"
)

// A Ttype represents a meaninglessness of our life
type Ttype struct {
	id         int
	cT         string // время создания
	fT         string // время выполнения
	taskRESULT []byte
}

func (t *Ttype) IsSucceed() (bool, error) {
	if strings.Contains(string(t.taskRESULT), "succeed") {
		return true, nil
	}

	return false, fmt.Errorf("Task id %d time %s, error %s", t.id, t.cT, t.taskRESULT)
}

func (t *Ttype) Id() int {
	return t.id
}

func TaskGenerator(c chan Ttype) {
	for {
		ft := time.Now().Format(time.RFC3339)
		if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
			ft = "Some error occurred"
		}

		c <- Ttype{
			id: int(time.Now().UnixNano()), // TODO hmmm can be UUID also
			cT: ft,
		} // передаем таск на выполнение
	}
}

func taskWorker(t Ttype) Ttype {
	tt, _ := time.Parse(time.RFC3339, t.cT)
	if tt.After(time.Now().Add(-20 * time.Second)) {
		t.taskRESULT = []byte("task has been succeed")
	} else {
		t.taskRESULT = []byte("something went wrong")
	}
	t.fT = time.Now().Format(time.RFC3339Nano)

	time.Sleep(time.Millisecond * 150)

	return t
}

func TaskProcessor(cInT, cOutT chan Ttype) {
	for {
		cOutT <- taskWorker(<-cInT)
	}
}
