package main

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

// ЗАДАНИЕ:
// * сделать из плохого кода хороший;
// * важно сохранить логику появления ошибочных тасков;
// * сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через merge-request.

// приложение эмулирует получение и обработку тасков, пытается и получать и обрабатывать в многопоточном режиме
// В конце должно выводить успешные таски и ошибки выполнены остальных тасков

var ErrorResultBytes = []byte("Some error occured")

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	createTask := func() Ttype {
		ct := time.Now().Format(time.RFC3339)
		task := Ttype{cT: ct, id: int(time.Now().UnixNano())}
		if time.Now().Nanosecond()%2 > 0 { // вот такое условие появления ошибочных тасков
			task.taskRESULT = ErrorResultBytes
		}
		return task
	}

	cpu := runtime.NumCPU()
	processor := NewProcessor(cpu*2, cpu, createTask)
	result, err := processor.Loop(ctx)

	println("Errors:")
	for _, r := range err {
		fmt.Printf("%v\n", r)
	}

	println("Done tasks:")
	for _, r := range result {
		fmt.Printf("%d %s\n", r.id, string(r.taskRESULT))
	}
}
