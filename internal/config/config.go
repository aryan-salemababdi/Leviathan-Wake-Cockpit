package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	KeyDBAddress        string `json:"keydb_address"`
	BinanceWebsocketURL string `json:"binance_websocket_url"`
	UpdateIntervalHours int    `json:"update_interval_hours"`
}

func Load(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = json.Unmarshal(file, &cfg)
	return &cfg, err
}
