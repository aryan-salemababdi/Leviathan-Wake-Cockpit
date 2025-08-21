package processor

import (
	"context"
	"log"
	"net/http"
	"time"

	"leviathan-wake-cockpit/internal/config"

	"github.com/go-redis/redis/v8"
)

type ProcessorService struct {
	cfg         *config.Config
	keydbClient *redis.Client
	whitelist   map[string]bool
	httpClient  *http.Client
}

func NewProcessorService(cfg *config.Config, db *redis.Client) *ProcessorService {
	return &ProcessorService{
		cfg:         cfg,
		keydbClient: db,
		whitelist:   make(map[string]bool),
		httpClient:  &http.Client{Timeout: 5 * time.Second},
	}
}

func (s *ProcessorService) Start(ctx context.Context) {
	log.Println("StreamProcessorService (Hot Path) has started...")

	for {
		if err := s.loadWhitelistFromKeyDB(ctx); err == nil {
			break
		}
		log.Println("WARN: Could not load initial whitelist. Retrying in 10 seconds...")
		time.Sleep(10 * time.Second)
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			s.connectAndListen(ctx)
			log.Println("Connection lost. Reconnecting in 5 seconds...")
			time.Sleep(5 * time.Second)
		}
	}
}
