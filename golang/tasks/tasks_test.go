package tasks

import (
	"errors"
	"testing"
	"time"
)

func TestTtypeProcessorMain_Process(t *testing.T) {
	p := TtypeProcessorMain{
		TTL: time.Microsecond * 0,
	}
	task := Ttype{
		Id: 132,
		CT: time.Now().Format(CTFormat),
	}

	time.Sleep(time.Millisecond * 10)

	err := p.Process(&task)
	if err != nil {
		t.Errorf("task: %#v, Process() error = %#v", task, err)
	}

	if !errors.Is(task.ProcessingError, TooOldTaskErr) {
		t.Errorf("task: %#v, Process() error = %#v, need = %#v", err, task, TooOldTaskErr)
	}

	p.TTL = time.Second
	task.CT = time.Now().Format(CTFormat)

	err = p.Process(&task)
	if err != nil {
		t.Errorf("task: %#v, Process() error = %#v", task, err)
	}
	if task.ProcessingError != nil {
		t.Errorf("task: %#v, Process() error = %#v", task, task.ProcessingError)
	}
	_, err = time.Parse(FTFormat, task.FT)
	if err != nil {
		t.Errorf("task: %#v, FT cannot be parsed using %s", task, FTFormat)
	}

	task.CT = "something except time"

	err = p.Process(&task)
	if !errors.As(err, &ErrCreationTime{}) {
		t.Errorf("task: %#v, Process() error type = %T, need = type %T", task, err, ErrCreationTime{})
	}
}

func TestTaskCreator(t *testing.T) {
	out := make(chan Ttype, 3)
	exit := make(chan struct{})

	stopped := TaskCreator(out, time.Millisecond*70, exit)

	time.Sleep(time.Millisecond * 300)

	select {
	case <-exit:
		t.Fatal("'exit' channel cannot be closed from inside TaskCreator")
	default:
	}

	select {
	case <-stopped:
		t.Fatal("'stopped' channel cannot be closed from inside TaskCreator")
	default:
	}

	task := <-out
	if task.Id == 0 {
		t.Error("id cannot be empty; maybe task is empty?")
	}

	outClosed := make(chan struct{})
	count := 0
	go func() {
		defer close(outClosed)
		for outTask := range out {
			if outTask.Id == 0 {
				t.Error("id cannot be empty; maybe task is empty?")
			}
			if _, err := time.Parse(CTFormat, outTask.CT); err != nil {
				t.Errorf("task: %#v, CT cannot be parsed using %s", outTask, CTFormat)
			}
			count++
		}
	}()

	close(exit)

	select {
	case <-outClosed:
	case <-time.After(time.Millisecond * 200):
		t.Fatal("output channel is not closed")
	}

	if count < 3 {
		t.Errorf("TaskCreator created too few tasks: %d", count)
	}

	select {
	case <-stopped:
	case <-time.After(time.Millisecond * 200):
		t.Error("'stopped' channel should be closed from inside TaskCreator when it exits")
	}
}
