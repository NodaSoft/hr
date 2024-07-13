package main

import (
	"context"
	"go.uber.org/zap"
	"sync"
	"time"
)

func main() {
	cfg := TaskProcessingConfig{
		GeneratingTasksDuration:     10,
		MaxProcessingWorkerDuration: 30,
		MaxHandleProcessedDuration:  30,

		UnprocessedTasksChannelBufferSize: 20,
		ProcessedTasksChannelBufferSize:   20,

		FillingTasksChannelWorkerCount:    10,
		ProcessingTasksChannelWorkerCount: 10,

		PrintingTasksPeriod: 3,
		//IsPrintTasksDetailed: false,
		IsPrintTasksDetailed: true,

		//LogLevel: int8(zap.DebugLevel),
		LogLevel: int8(zap.InfoLevel),
	}

	PanicOnError(ValidateTaskProcessingConfig(cfg))

	logger := NewLogger(cfg.LogLevel)

	ProcessingTasks(cfg, logger)
}

// ProcessingTasks is a function contains main logic of task processing.
func ProcessingTasks(cfg TaskProcessingConfig, logger *zap.Logger) {
	logger = logger.Named("ProcessingTasks")

	// changing parameters for emulation of creation and processing tasks
	//IncorrectTaskTimeDivision = 500
	//CreateTaskDuration = time.Millisecond * 500
	//ProcessTaskDuration = time.Millisecond * 500

	startTime := time.Now()
	requiredGeneratingEndTime := startTime.Add(time.Second * time.Duration(cfg.GeneratingTasksDuration))
	requiredProcessingEndTime := startTime.Add(time.Second * time.Duration(cfg.MaxProcessingWorkerDuration))
	requiredHandleProcessedEndTime := startTime.Add(time.Second * time.Duration(cfg.MaxHandleProcessedDuration))

	unprocessedTasksChannel := make(chan Task, cfg.UnprocessedTasksChannelBufferSize)
	processedTasksChannel := make(chan Task, cfg.ProcessedTasksChannelBufferSize)

	var fillTaskWg sync.WaitGroup
	var processingTasksWg sync.WaitGroup
	var handleProcessedTasksWg sync.WaitGroup

	logger.Info("Start", zap.Time("StartTime", startTime))

	for i := 0; i < cfg.FillingTasksChannelWorkerCount; i++ {
		ctx, cancel := context.WithDeadline(context.Background(), requiredGeneratingEndTime)
		defer cancel()

		fillTaskWg.Add(1)
		go FillTaskChannel(unprocessedTasksChannel, ctx, &fillTaskWg, logger)
	}

	for i := 0; i < cfg.ProcessingTasksChannelWorkerCount; i++ {
		ctx, cancel := context.WithDeadline(context.Background(), requiredProcessingEndTime)
		defer cancel()

		processingTasksWg.Add(1)
		go ProcessTaskChannel(unprocessedTasksChannel, processedTasksChannel, ctx, &processingTasksWg, logger)
	}

	ctx, cancel := context.WithDeadline(context.Background(), requiredHandleProcessedEndTime)
	defer cancel()

	handleProcessedTasksWg.Add(1)
	go HandleProcessedTasksChannel(
		processedTasksChannel,
		ctx,
		cfg.IsPrintTasksDetailed,
		time.Second*time.Duration(cfg.PrintingTasksPeriod),
		&handleProcessedTasksWg,
		logger,
	)

	fillTaskWg.Wait()
	close(unprocessedTasksChannel) // close unprocessedTasksChannel after all goroutines for filling are finished
	logger.Info("Filling tasks finished", zap.Duration("Duration", time.Since(startTime)))

	processingTasksWg.Wait()
	close(processedTasksChannel) // close processedTasksChannel after all goroutines for processing are finished
	logger.Info("Processing tasks finished", zap.Duration("Duration", time.Since(startTime)))

	handleProcessedTasksWg.Wait() // wait for all goroutines for processing processed tasks to finish
	logger.Info("Processing processed tasks finished", zap.Duration("Duration", time.Since(startTime)))

	logger.Info("Processing finished", zap.Duration("Duration", time.Since(startTime)))
}
