package processorService

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"leviathan-wake-cockpit/internal/config"

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
	httpClient  *http.Client
	whitelist   map[string]bool
	// grpcClient
}

type InfuraSubscriptionMessage struct {
	Params struct {
		Subscription string `json:"subscription"`
		Result       string `json:"result"`
	} `json:"params"`
}

type JSONRPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      int           `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type TransactionDetailsResponse struct {
	Result *struct {
		From string `json:"from"`
		To   string `json:"to"`
	} `json:"result"`
}

func NewProcessorService(cfg *config.Config, db *redis.Client) *ProcessorService {
	return &ProcessorService{
		cfg:         cfg,
		keydbClient: db,
		httpClient:  &http.Client{Timeout: 5 * time.Second},
		whitelist:   make(map[string]bool),
	}
}

func (s *ProcessorService) Start(ctx context.Context) {
	log.Println("StreamProcessorService (Hot Path) has started...")

	for {
		if err := s.loadWhitelistFromKeyDB(ctx); err == nil {
			break
		}
		log.Printf("WARN: Could not load initial whitelist. Retrying in 10 seconds...")
		time.Sleep(10 * time.Second)
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			s.connectAndListen()
			log.Println("Connection lost. Reconnecting in 5 seconds...")
			time.Sleep(5 * time.Second)
		}
	}
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

	subscribeMsg := `{"jsonrpc":"2.0","id":1,"method":"eth_subscribe","params":["newPendingTransactions"]}`
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
	var subMsg InfuraSubscriptionMessage
	if err := json.Unmarshal(msg, &subMsg); err != nil || subMsg.Params.Result == "" {
		return
	}

	txHash := subMsg.Params.Result

	txDetails, err := s.fetchTransactionDetails(txHash)
	if err != nil || txDetails.Result == nil {
		return
	}

	from := strings.ToLower(txDetails.Result.From)
	to := strings.ToLower(txDetails.Result.To)

	if s.whitelist[from] || s.whitelist[to] {
		log.Printf("✅ !!! WHALE TRANSACTION DETECTED !!! Hash: %s, From: %s, To: %s", txHash, from, to)
	}
}

func (s *ProcessorService) fetchTransactionDetails(hash string) (*TransactionDetailsResponse, error) {
	httpURL := strings.Replace(s.cfg.ArbitrumWsUrl, "wss://", "https://", 1)
	httpURL = strings.Replace(httpURL, "/ws/", "/", 1)

	payload := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "eth_getTransactionByHash",
		Params:  []interface{}{hash},
	}
	payloadBytes, _ := json.Marshal(payload)

	req, _ := http.NewRequestWithContext(context.Background(), "POST", httpURL, bytes.NewBuffer(payloadBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var txResp TransactionDetailsResponse
	if err := json.NewDecoder(resp.Body).Decode(&txResp); err != nil {
		return nil, err
	}

	return &txResp, nil
}
