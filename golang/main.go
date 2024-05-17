package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

//При выполнении задания не менял исходный тип данных, и операции с его параметрами(задача была не в этом(всем добра))

const timeout = 3 * time.Second

var logger = log.Default()

type Ttype struct {
	id         int
	cT         string // время создания
	fT         string // время выполнения
	taskRESULT []byte
}

// "сохранить логику появления ошибочных тасков"
func generateTask(ctx context.Context, ch chan<- Ttype) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			ft := time.Now().Format(time.RFC3339)
			if time.Now().Nanosecond()%2 > 0 {
				ft = "Some error occured"
			}
			ch <- Ttype{
				cT: ft,
				id: int(time.Now().Unix()),
			}
		}
	}
}

func taskWorker(t Ttype) Ttype {
	tt, _ := time.Parse(time.RFC3339, t.cT)
	if tt.After(time.Now().Add(-20 * time.Second)) {
		t.taskRESULT = []byte("task has been successed")
	} else {
		t.taskRESULT = []byte("something went wrong")
	}
	t.fT = time.Now().Format(time.RFC3339Nano)

	time.Sleep(time.Millisecond * 150)

	return t
}

// Сделать правильную мультипоточность обработки заданий
func getResult(ctx context.Context, ch <-chan Ttype) (map[int]Ttype, map[int]error) {
	result := make(map[int]Ttype)
	err := make(map[int]error)
	mu := new(sync.Mutex)
	wg := sync.WaitGroup{}
	for {
		select {
		case t := <-ch:
			wg.Add(1)
			go func() {
				t = taskWorker(t)
				if string(t.taskRESULT[14:]) == "successed" {
					mu.Lock()
					result[t.id] = t
					mu.Unlock()
				} else {
					mu.Lock()
					err[t.id] = fmt.Errorf("id=%d; error=%s", t.id, string(t.taskRESULT))
					mu.Unlock()
				}
				wg.Done()
			}()

		case <-ctx.Done():
			wg.Wait()
			return result, err
		}
	}
}

func main() {
	ttypeChannel := make(chan Ttype, 10)
	defer close(ttypeChannel)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go generateTask(ctx, ttypeChannel) // Важно оставить асинхронные генерацию

	resultCtx, resultCancel := context.WithTimeout(context.Background(), timeout)
	defer resultCancel()

	result, err := getResult(resultCtx, ttypeChannel)
	if err != nil {
		logger.Println("Tasks's errors:")
		for _, val := range err {
			logger.Println(val.Error())
		}
	}
	logger.Println("Done tasks:")
	for _, r := range result {
		logger.Println(r.id)
	}
}
