package processor

import (
	"context"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func (s *ProcessorService) connectAndListen(ctx context.Context) {
	conn, _, err := websocket.DefaultDialer.Dial(s.cfg.ArbitrumWsUrl, nil)
	if err != nil {
		log.Printf("ERROR: Could not connect to Arbitrum WebSocket: %v", err)
		return
	}
	defer conn.Close()
	log.Println("Connected to Arbitrum WebSocket.")

	subscribeMsg := `{"jsonrpc":"2.0","id":1,"method":"eth_subscribe","params":["newPendingTransactions"]}`
	if err := conn.WriteMessage(websocket.TextMessage, []byte(subscribeMsg)); err != nil {
		log.Printf("ERROR: Could not subscribe: %v", err)
		return
	}
	log.Println("Subscribed to new pending transactions.")

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Disconnecting from WebSocket...")
			return
		case <-ticker.C:
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("ERROR: Could not send ping: %v", err)
				return
			}
		default:
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Printf("ERROR: Error reading from WebSocket: %v", err)
				return
			}
			go s.processTransactionMessage(message)
		}
	}
}
