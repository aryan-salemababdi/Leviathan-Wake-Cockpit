package processor

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func (s *ProcessorService) processTransactionMessage(msg []byte) {
	var subMsg InfuraSubscriptionMessage
	if err := json.Unmarshal(msg, &subMsg); err != nil || subMsg.Params.Result == "" {
		return
	}

	txHash := subMsg.Params.Result

	txDetails, err := s.FetchTransactionDetails(txHash)
	if err != nil || txDetails.Result == nil {
		return
	}

	from := strings.ToLower(txDetails.Result.From)
	to := strings.ToLower(txDetails.Result.To)
	value := txDetails.Result.Value

	if s.whitelist[from] || s.whitelist[to] {
		log.Printf("✅ Whale Tx DETECTED! Hash: %s, From: %s, To: %s, Value: %s wei",
			txHash, from, to, value)

		receipt, err := s.fetchTransactionReceipt(txHash)
		if err == nil && receipt.Result != nil {
			for _, logEntry := range receipt.Result.Logs {
				if len(logEntry.Topics) > 0 && logEntry.Topics[0] == "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef" {
					log.Printf("   📦 Token Transfer -> Contract: %s, Data(raw amount): %s", logEntry.Address, logEntry.Data)
				}
			}
		}
	}
}

func (s *ProcessorService) fetchTransactionReceipt(hash string) (*TransactionReceiptResponse, error) {
	httpURL := strings.Replace(s.cfg.ArbitrumWsUrl, "wss://", "https://", 1)
	httpURL = strings.Replace(httpURL, "/ws/", "/", 1)

	payload := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "eth_getTransactionReceipt",
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

	var receipt TransactionReceiptResponse
	if err := json.NewDecoder(resp.Body).Decode(&receipt); err != nil {
		return nil, err
	}

	return &receipt, nil
}

func (s *ProcessorService) FetchTransactionDetails(hash string) (*TransactionDetailsResponse, error) {
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
