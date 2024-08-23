package report

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

func SchedulePrint(ctx context.Context, wg *sync.WaitGroup, b *Builder, interval time.Duration) {
	defer func() {
		log.Println("scheduler print done")
		wg.Done()
	}()
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			log.Println("print last report")
			print(b.Report())
			return
		case <-ticker.C:
			print(b.Report())
		}
	}
}

func print(r Report) {
	sb := strings.Builder{}
	sb.WriteString("Succeeded:\n")
	for _, s := range r.succeeded {
		sb.WriteString(s)
	}

	sb.WriteString("Errors:\n")
	for _, s := range r.errors {
		sb.WriteString(s)
	}

	fmt.Println(sb.String())
	log.Println("err:", len(r.errors), "success:", len(r.succeeded))
}
