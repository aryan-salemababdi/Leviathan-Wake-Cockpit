// cmd/processor/main.go
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"leviathan-wake-cockpit/internal/config"
	keydb_database "leviathan-wake-cockpit/internal/database"
	processorService "leviathan-wake-cockpit/internal/services"
)

func main() {
	log.Println("Initializing StreamProcessorService...")

	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatalf("FATAL: Could not load configuration: %v", err)
	}

	dbClient, err := keydb_database.NewKeyDBClient(cfg.KeydbAddress)
	if err != nil {
		log.Fatalf("FATAL: Could not connect to KeyDB: %v", err)
	}
	log.Println("Successfully connected to KeyDB.")

	processorSvc := processorService.NewProcessorService(cfg, dbClient)
	log.Println("ProcessorService is assembled.")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go processorSvc.Start(ctx)

	log.Println("StreamProcessorService is now RUNNING. Press Ctrl+C to exit.")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down StreamProcessorService...")
	cancel()
	time.Sleep(2 * time.Second)
	log.Println("Service gracefully stopped.")
}
