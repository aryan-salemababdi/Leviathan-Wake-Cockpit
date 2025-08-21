package processorService

import (
	"context"
	"encoding/json"
	"leviathan-wake-cockpit/internal/config"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

const WhitelistKey = "whale_whitelist"

type WhitelistEntry struct {
	Address    string `json:"address"`
	Blockchain string `json:"blockchain"`
}

type ProcessorService struct {
	cfg         *config.Config
	keydbClient *redis.Client
	whitelist   map[string]bool
	// grpcClient
}

func NewProcessorService(cfg *config.Config, db *redis.Client) *ProcessorService {
	return &ProcessorService{
		cfg:         cfg,
		keydbClient: db,
		whitelist:   make(map[string]bool),
	}
}

func (s *ProcessorService) Start(ctx context.Context) {
	log.Println("StreamProcessorService (Hot Path) has started...")

	// TODO
	err := s.loadWhitelistFromKeyDB(ctx)
	if err != nil {
		log.Printf("WARN: Could not load initial whitelist: %v. Will retry...", err)
	}

	s.connectAndListen()
}

func (s *ProcessorService) loadWhitelistFromKeyDB(ctx context.Context) error {
	jsonData, err := s.keydbClient.Get(ctx, WhitelistKey).Result()
	if err == redis.Nil {
		log.Println("No whitelist found in KeyDB, starting with empty list.")
		return nil
	} else if err != nil {
		return err
	}

	var inner string
	if err := json.Unmarshal([]byte(jsonData), &inner); err == nil {
		jsonData = inner
	}

	newWhitelist := make(map[string]bool)

	var entries []WhitelistEntry
	if err := json.Unmarshal([]byte(jsonData), &entries); err == nil && len(entries) > 0 {
		for _, entry := range entries {
			if entry.Blockchain == "Arbitrum" {
				newWhitelist[entry.Address] = true
			}
		}
		s.whitelist = newWhitelist
		log.Printf("Whitelist loaded into memory with %d Arbitrum addresses.", len(s.whitelist))
		return nil
	}

	var addresses []string
	if err := json.Unmarshal([]byte(jsonData), &addresses); err == nil && len(addresses) > 0 {
		for _, addr := range addresses {
			newWhitelist[addr] = true
		}
		s.whitelist = newWhitelist
		log.Printf("Whitelist loaded into memory with %d addresses.", len(s.whitelist))
		return nil
	}

	log.Printf("WARN: Whitelist data in KeyDB has unexpected format: %s", jsonData)
	return nil
}

func (s *ProcessorService) connectAndListen() {
	conn, _, err := websocket.DefaultDialer.Dial(s.cfg.ArbitrumWsUrl, nil)
	if err != nil {
		log.Fatalf("FATAL: Could not connect to Arbitrum WebSocket: %v", err)
	}
	defer conn.Close()
	log.Println("Successfully connected to Arbitrum WebSocket.")

	subscribeMsg := `{"jsonrpc":"2.0","id":1,"method":"eth_subscribe","params":["alchemy_newPendingTransactions"]}`
	if err := conn.WriteMessage(websocket.TextMessage, []byte(subscribeMsg)); err != nil {
		log.Fatalf("FATAL: Could not subscribe to new transactions: %v", err)
	}
	log.Println("Subscribed to new pending transactions.")

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("ERROR: Error reading from WebSocket: %v. Reconnecting...", err)
			return
		}

		go s.processTransactionMessage(message)
	}
}

func (s *ProcessorService) processTransactionMessage(msg []byte) {
	// TODO
}
