package services

import (
	"context"
	"time"
)

type Orchestrator struct {
	producer         *producer
	processor        *processor
	successPresenter *presenter
	errorPresenter   *presenter
	ctx              context.Context
	bandwidth        int
	processingTime   time.Duration
}

func NewOrchestrator(bandwidth int, processingTime time.Duration) *Orchestrator {
	return NewOrchestratorWithContext(context.Background(), bandwidth, processingTime)
}

func NewOrchestratorWithContext(ctx context.Context, bandwidth int, processingTime time.Duration) *Orchestrator {
	producer := newProducer(ctx, bandwidth)
	tasks := producer.getTasks()

	processor := newProcessor(ctx, tasks)
	success, errors := processor.resultChannels()

	successPres := newPresenter(success)
	errorPres := newPresenter(errors)

	return &Orchestrator{
		producer:         producer,
		processor:        processor,
		successPresenter: successPres,
		errorPresenter:   errorPres,
		ctx:              ctx,
		bandwidth:        bandwidth,
		processingTime:   processingTime,
	}
}

func (o *Orchestrator) Do() {
	go o.producer.generate()
	go o.successPresenter.print()
	go o.errorPresenter.print()
	go o.processor.doWork()

	select {
	case <-time.After(o.processingTime):
	case <-o.ctx.Done():
	}
}
