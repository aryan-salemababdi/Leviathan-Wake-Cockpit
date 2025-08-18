package service

import (
	"leviathan-wake-cockpit/internal/config"
	"log"

	"github.com/go-redis/redis/v8"
)

type StreamProcessorService struct {
	cfg      *config.Config
	dbClient *redis.Client
}

func NewStreamProcessorService(cfg *config.Config, db *redis.Client) *StreamProcessorService {
	return &StreamProcessorService{cfg: cfg, dbClient: db}
}

func (s *StreamProcessorService) Start() {
	log.Println("Starting StreamProcessorService (Hot Path)...")
}
