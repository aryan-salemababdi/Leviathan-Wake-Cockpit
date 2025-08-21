package main

import (
	"context"
	"leviathan-wake-cockpit/internal/config"
	keydb_database "leviathan-wake-cockpit/internal/database"
	"leviathan-wake-cockpit/internal/services/processor"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using enviroment variables.")
	}
	log.Println("Initializing StreamProcessorService...")

	cfg := config.Load()
	if cfg == nil {
		log.Fatalf("FATAL: Could not load configuration")
	}

	dbClient, err := keydb_database.NewKeyDBClient(cfg.KeydbAddress)
	if err != nil {
		log.Fatalf("FATAL: Could not connect to KeyDB: %v", err)
	}
	log.Println("Successfully connected to KeyDB.")

	processorSvc := processor.NewProcessorService(cfg, dbClient)
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
