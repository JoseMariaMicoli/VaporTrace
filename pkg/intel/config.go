/*
Copyright (c) 2026 José María Micoli
Licensed under {'license_type': 'BSL', 'change_date': '2033-02-17', 'convert_to': 'Apache-2.0'}

You may:
✔ Study
✔ Modify
✔ Use for internal security testing

You may NOT:
✘ Offer as a commercial service
✘ Sell derived competing products
*/

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
