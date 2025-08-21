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

	txDetails, err := s.fetchTransactionDetails(txHash)
	if err != nil || txDetails.Result == nil {
		return
	}

	from := strings.ToLower(txDetails.Result.From)
	to := strings.ToLower(txDetails.Result.To)

	if s.whitelist[from] || s.whitelist[to] {
		log.Printf("✅ Whale Transaction DETECTED! Hash: %s, From: %s, To: %s", txHash, from, to)
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
