package main

import (
	"context"
	"os"
	"strconv"
	"strings"
	"taskhandler/internal/config"
	"taskhandler/internal/core"
	log "taskhandler/internal/logger"
	"time"

	"github.com/alecthomas/units"
)

// Приложение эмулирует получение и обработку неких тасков. Пытается и получать, и обрабатывать в многопоточном режиме.
// Приложение должно генерировать таски 10 сек. Каждые 3 секунды должно выводить в консоль результат всех обработанных к этому моменту тасков (отдельно успешные и отдельно с ошибками).

// ЗАДАНИЕ: сделать из плохого кода хороший и рабочий - as best as you can.
// Важно сохранить логику появления ошибочных тасков.
// Важно оставить асинхронные генерацию и обработку тасков.
// Сделать правильную мультипоточность обработки заданий.
// Обновленный код отправить через pull-request в github
// Как видите, никаких привязок к внешним сервисам нет - полный карт-бланш на модификацию кода.

func main() {
	// Startup, may panic, which is ok
	// The app shouldn't work without logger or with wrong config
	config.InitConfig()
	log.InitGlobalLogger()
	log.Info("Service preparations successed")

	// Starting global timer
	ctx, cancel := context.WithTimeout(context.Background(), config.C.Service.Timeout)
	defer cancel()

	// Magic .Time to generate incorrect tasks
	serviceStarted := time.Now()

	// Return result core.Task with result field
	resultWorker := func(t core.Task) core.Task {
		if t.Created == serviceStarted {
			t.Result = []byte("something went wrong")
		} else {
			t.Result = []byte("task has been successed")
		}
		t.Finished = time.Now()

		// Not sure about this, but i'll leave this sleep untouched
		// I guess its emulating real work
		time.Sleep(time.Millisecond * 150)

		return t
	}

	// Can only handle input specific for resultWorker
	separator := func(t core.Task) int8 {
		if (string)(t.Result) == "something went wrong" {
			return 1
		}
		if (string)(t.Result) == "task has been successed" {
			return 0
		}
		return -1
	}

	errors := make(chan error, 10)
	// Initial channel, will be filled and closed by FillChannel
	tasks := make(chan core.Task, 10)

	chain := []core.PipilineElement{
		core.FactoryToPipeElem(core.FillChannel, &core.BrokenFactory{ServiceStarted: serviceStarted}),
		core.HandlerToPipeElem(core.HandleTasks, resultWorker),
		core.SeparatorToPipeElem(core.SeparateBrokenTasks, separator, errors),
	}

	// Preparation for flushing
	done := core.Pipeline(ctx, tasks, chain...)
	flusher := time.NewTicker(config.C.Service.FlushRate)

	// It is insane that we need to hold all done tasks in memory for full 3 seconds
	// Ofc we can flush it to disk every 50 tasks or so and then read it every .FlushRate
	// After we read it again and write to stdout and stderr
	errorBuilder := strings.Builder{}
	errorBuilder.Grow(int(units.Mebibyte) * 8)
	doneBuilder := strings.Builder{}
	doneBuilder.Grow(int(units.Mebibyte) * 8)

	// allFlushed is for safe exit
	// after every last
	allFlushed := make(chan struct{})
	// Single Threaded read from both channels
	// Do not need to synchronize .Reset()
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info("Last flush. Reason: context.Done()")
				flush(&errorBuilder, &doneBuilder)
				allFlushed <- struct{}{}
				return

			case <-flusher.C:
				flush(&errorBuilder, &doneBuilder)
				errorBuilder.Reset()
				doneBuilder.Reset()
				errorBuilder.Grow(int(units.Mebibyte) * 8)
				doneBuilder.Grow(int(units.Mebibyte) * 8)

			case t := <-done:
				// Sprintf is slow and we have A LOT of tasks
				writeTask(&doneBuilder, t)

			case err := <-errors:
				errorBuilder.WriteString(err.Error())
				errorBuilder.WriteByte('\n')
			}
		}
	}()
	<-allFlushed
	// Just for everything to clean up
	time.Sleep(time.Millisecond * 100)
}

// utilitary functions
func writeTask(builder *strings.Builder, t core.Task) {
	builder.WriteString("Done core.Task. Id: ")
	builder.WriteString(strconv.Itoa((int)(t.Id)))
	builder.WriteString("Result: ")
	builder.WriteString((string)(t.Result))
	builder.WriteByte('\n')

}
func flush(errorBuilder *strings.Builder, doneBuilder *strings.Builder) {
	log.Info("Flushing.")
	log.Info("Used Memory in errorBuilder: ", errorBuilder.Len(), "While cap is:", errorBuilder.Cap())
	log.Info("Used Memory in doneBuilder: ", doneBuilder.Len(), "While cap is:", doneBuilder.Cap())
	os.Stderr.WriteString("\n\tError Tasks:\n")
	os.Stderr.WriteString(errorBuilder.String())
	os.Stdout.WriteString("\n\tDone Tasks:\n")
	os.Stdout.WriteString(doneBuilder.String())

}
