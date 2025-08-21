package processor

import (
	"context"
	"log"
	"time"
)

type ExecutionService struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func NewExecutionService(parent context.Context) *ExecutionService {
	ctx, cancel := context.WithCancel(parent)
	return &ExecutionService{
		ctx:    ctx,
		cancel: cancel,
	}
}

func (e *ExecutionService) Start() {
	log.Println("ExecutionService started...")

	go func() {
		for {
			select {
			case <-e.ctx.Done():
				log.Println("ExecutionService stopped.")
				return
			default:
				e.executeTask()
				time.Sleep(2 * time.Second)
			}
		}
	}()
}

func (e *ExecutionService) Stop() {
	e.cancel()
}

func (e *ExecutionService) executeTask() {
	log.Println("ExecutionService is executing a task...")
}
