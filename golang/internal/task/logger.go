package task

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"
	"time"
)

// Logger writes done Tasks to the log
type Logger struct {
	Tasks chan Task
}

// Start Logger for wring done Tasks to the log
func (l Logger) Start() {
	w := tabwriter.NewWriter(os.Stdout, 15, 0, 4, ' ', 0)

	if _, err := w.Write([]byte("Id\tStart\tFinish\tError\n")); err != nil {
		log.Println(err)
	}

	for {
		task := <-l.Tasks
		l.log(w, task)
	}
}

func (l Logger) log(w *tabwriter.Writer, task Task) {
	output := fmt.Sprintf(
		"%d\t%s\t%s\t%s\n",
		task.Id,
		task.Start.Format(time.RFC3339),
		task.Finish.Format(time.RFC3339),
		task.ErrorMessage,
	)

	if _, err := w.Write([]byte(output)); err != nil {
		log.Println(err)
	}

	if err := w.Flush(); err != nil {
		log.Println(err)
	}
}
