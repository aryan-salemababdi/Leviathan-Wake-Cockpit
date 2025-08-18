package service

import (
	"leviathan-wake-cockpit/internal/config"
	"log"

	"github.com/go-redis/redis/v8"
)

type WhaleUpdaterService struct {
	cfg      *config.Config
	dbClient *redis.Client
}

func NewWhaleUpdaterService(cfg *config.Config, db *redis.Client) *WhaleUpdaterService {
	return &WhaleUpdaterService{cfg: cfg, dbClient: db}
}

func (s *WhaleUpdaterService) Start() {
	log.Println("Starting WhaleUpdaterService (Cold Path)...")
}
