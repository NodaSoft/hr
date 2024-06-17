package main

import (
	"context"
	"time"

	"taskConcurrency/internal/app"
)

func main() {
	a := app.App{}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	a.Do(ctx)
}
