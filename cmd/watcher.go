package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"leviathan-wake-cockpit/internal/config"
	keydb_database "leviathan-wake-cockpit/internal/database"
	service "leviathan-wake-cockpit/internal/services"
)

func main() {
	log.Println("Initializing Leviathan-Wake-Cockpit...")
	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatalf("FATAL: Could not load configuration: %v", err)
	}

	dbClient, err := keydb_database.NewKeyDBClient(cfg.KeyDBAddress)
	if err != nil {
		log.Fatalf("FATAL: Could not connect to KeyDB: %v", err)
	}
	log.Println("Successfully connected to KeyDB.")

	updaterService := service.NewWhaleUpdaterService(cfg, dbClient)
	streamService := service.NewStreamProcessorService(cfg, dbClient)
	log.Println("Internal services initialized.")

	go updaterService.Start()
	go streamService.Start()

	log.Println("Leviathan-Wake-Cockpit is now running. Press Ctrl+C to exit.")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")
	log.Println("Server gracefully stopped.")
}
