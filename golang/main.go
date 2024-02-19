package main

import (
	"context"
	"fmt"
	"nodasoft-golang/services"
	"time"
)

const (
	ContextTimeoutSec = 5
	RealWaitSec       = 3
	DOWorkTimeSec     = 10
	Bandwidth         = 10
)

func main() {
	start := time.Now()
	// мы можем контролить извне, например таким таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), ContextTimeoutSec*time.Second)
	go func() {
		// или мы хотим остановиться по какой-то причине,
		// не дожидаясь наступления таймаута от контекста
		time.Sleep(RealWaitSec * time.Second)
		cancel()

	}()

	orc := services.NewOrchestratorWithContext(ctx, Bandwidth, DOWorkTimeSec*time.Second)

	orc.Do()

	end := time.Now().Sub(start)
	time.Sleep(10 * time.Millisecond)
	fmt.Printf("time elapsed: %s\n", end)

}
