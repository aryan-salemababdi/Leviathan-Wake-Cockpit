// cmd/processor/main.go
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"stream-processor-service/internal/config"
	// "leviathan/stream-processor/internal/database"
	// "leviathan/stream-processor/internal/service"
)

func main() {
	log.Println("Initializing StreamProcessorService...")

	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatalf("FATAL: Could not load configuration: %v", err)
	}
	log.Printf("Loaded config: %+v", cfg)
	log.Println("KeyDB connection placeholder is ready.")
	log.Println("gRPC client placeholder is ready.")

	// processorSvc := service.NewProcessorService(cfg, dbClient, execClient)
	log.Println("ProcessorService is assembled.")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// go processorSvc.Start(ctx)

	log.Println("StreamProcessorService is now RUNNING. Press Ctrl+C to exit.")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down StreamProcessorService...")
	cancel()
	// time.Sleep(2 * time.Second)
	log.Println("Service gracefully stopped.")
}
