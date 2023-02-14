package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"taskprocessor/internal/server"
	"time"
)

type config struct {
	timeout int
}

func main() {
	c := config{
		timeout: 10,
	}

	l := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	s := server.New(l, 10)
	s.Process()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	timeout := time.After(time.Duration(c.timeout) * time.Second)

	select {
	case x := <-interrupt:
		l.Println("Received a signal", x.String())
		break
	case <-timeout: // TODO can use context either
		l.Print("Timeout exceeded")
		s.Stop()
		break
	}

	s.ShowResult()
}
