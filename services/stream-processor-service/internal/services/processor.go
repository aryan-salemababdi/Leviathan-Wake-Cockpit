package services

import (
	"context"
	"log"
)

type ProcessorService struct {
}

func NewProcessorService() *ProcessorService {
	return &ProcessorService{}
}

func (s *ProcessorService) Start(ctx context.Context) {
	log.Println("StreamProcessorService (Hot Path) has started...")

	// TODO

	<-ctx.Done()
	log.Println("Stopping StreamProcessorService...")
}
