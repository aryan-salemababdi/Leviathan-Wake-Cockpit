package config

import (
	"os"
)

type Config struct {
	KeydbAddress           string
	ArbitrumWsUrl          string
	ExecutionServerAddress string
}

func Load() *Config {
	return &Config{
		KeydbAddress:           os.Getenv("KEYDB_ADDRESS"),
		ArbitrumWsUrl:          os.Getenv("ARBITRUM_WS_URL"),
		ExecutionServerAddress: os.Getenv("EXECUTION_SERVER_ADDRESS"),
	}
}
