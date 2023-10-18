package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/NodaSoft/hr/golang/internal/tasks"
	"github.com/NodaSoft/hr/golang/internal/workers"
)

const poolSize = 10

func main() {
	generator := tasks.NewGenerator(poolSize)
	generator.Start()

	pool := workers.NewPool(poolSize, generator.Tasks())
	pool.Execute()

	time.Sleep(time.Second * 3)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	err := generator.Stop(ctx)
	err = errors.Join(err, pool.Shutdown(ctx))

	if err != nil {
		log.Fatal(err)
	}

	done, failed, err := pool.Result()

	if err != nil {
		log.Fatal(fmt.Errorf("failed to receive worker pool result: %w", err))
	}

	if len(failed) != 0 {
		log.Println("Failed tasks:")
	}

	for _, err := range failed {
		log.Println(err)
	}

	log.SetOutput(os.Stdout)
	log.Println("Done tasks:")

	for id := range done {
		log.Printf("id %d", id)
	}
}
