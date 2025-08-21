package processor

import (
	"context"
	"encoding/json"
	"log"

	"github.com/go-redis/redis/v8"
)

const WhitelistKey = "whale_whitelist"

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
		log.Printf("Whitelist loaded with %d Arbitrum addresses.", len(s.whitelist))
		return nil
	}

	var addresses []string
	if err := json.Unmarshal([]byte(jsonData), &addresses); err == nil && len(addresses) > 0 {
		for _, addr := range addresses {
			newWhitelist[addr] = true
		}
		s.whitelist = newWhitelist
		log.Printf("Whitelist loaded with %d addresses.", len(s.whitelist))
		return nil
	}

	log.Printf("WARN: Whitelist data has unexpected format: %s", jsonData)
	return nil
}
