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

	s.connectAndListen(ctx)
}

func (s *ProcessorService) loadWhitelistFromKeyDB(ctx context.Context) error {
	jsonData, err := s.keydbClient.Get(ctx, WhitelistKey).Result()

	if err != nil {
		return err
	}

	var entries []WhitelistEntry

	if err := json.Unmarshal([]byte(jsonData), &entries); err != nil {
		return err
	}

	newWhitelist := make(map[string]bool)
	for _, entry := range entries {
		if entry.Blockchain == "Arbitrum" {
			newWhitelist[entry.Address] = true
		}
	}
	s.whitelist = newWhitelist
	log.Printf("Whitelist loaded into memory with %d Arbitrum addresses.", len(s.whitelist))
	return nil
}

func (s *ProcessorService) connectAndListen(ctx context.Context) {
	// اتصال به وب‌سوکت Alchemy
	conn, _, err := websocket.DefaultDialer.Dial(s.cfg.ArbitrumWsUrl, nil)
	if err != nil {
		log.Fatalf("FATAL: Could not connect to Arbitrum WebSocket: %v", err)
	}
	defer conn.Close()
	log.Println("Successfully connected to Arbitrum WebSocket.")

	// ارسال پیام Subscribe برای دریافت تراکنش‌های جدید
	subscribeMsg := `{"jsonrpc":"2.0","id":1,"method":"eth_subscribe","params":["alchemy_newPendingTransactions"]}`
	if err := conn.WriteMessage(websocket.TextMessage, []byte(subscribeMsg)); err != nil {
		log.Fatalf("FATAL: Could not subscribe to new transactions: %v", err)
	}
	log.Println("Subscribed to new pending transactions.")

	// حلقه اصلی برای خواندن پیام‌ها
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("ERROR: Error reading from WebSocket: %v. Reconnecting...", err)
			// TODO: پیاده‌سازی منطق اتصال مجدد
			return
		}

		// پردازش هر پیام در یک Goroutine جداگانه تا حلقه اصلی مسدود نشود
		go s.processTransactionMessage(message)
	}
}

// processTransactionMessage پیام را پردازش و فیلتر می‌کند
func (s *ProcessorService) processTransactionMessage(msg []byte) {
	// TODO: در اینجا باید ساختار دقیق پیام JSON از Alchemy را پیدا کرده و آن را Unmarshal کنیم
	// ساختار پیام معمولاً شامل from, to, value, data و ... است.
	// var tx Transaction
	// json.Unmarshal(msg, &tx)

	// فرآیند فیلترینگ (در اینجا به صورت مفهومی)
	// if s.whitelist[tx.From] || s.whitelist[tx.To] {
	//	 log.Printf("!!! WHALE TRANSACTION DETECTED: From %s To %s", tx.From, tx.To)

	//   TODO:
	//   ۱. رمزگشایی داده‌های تراکنش (tx.Data) برای فهمیدن نوع معامله (مثلاً Swap در Uniswap)
	//   ۲. ساخت پیام WhaleSignal با استاندارد Protobuf
	//   ۳. ارسال پیام از طریق کلاینت gRPC به سرور اجرا و ریسک
	// }
}
