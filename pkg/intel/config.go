package intel

import (
	"sync"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

// IntelConfig holds API keys and settings for external services
type IntelConfig struct {
	ShodanKey string
	mu        sync.RWMutex
}

var GlobalIntel = &IntelConfig{}

// ConfigureShodan sets the API key for Shodan operations
func ConfigureShodan(key string) {
	GlobalIntel.mu.Lock()
	defer GlobalIntel.mu.Unlock()
	GlobalIntel.ShodanKey = key
	utils.TacticalLog("[green]INTEL:[-] Shodan API Key updated.")
}

// GetShodanKey retrieves the key safely
func GetShodanKey() string {
	GlobalIntel.mu.RLock()
	defer GlobalIntel.mu.RUnlock()
	return GlobalIntel.ShodanKey
}
