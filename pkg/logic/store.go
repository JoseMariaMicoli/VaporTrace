package logic

import (
	"net/http"
	"sync"
	"time"
)

// FlowStep defines a tactical action in a chain for Phase 7
type FlowStep struct {
	Name        string
	Method      string
	URL         string
	Body        string
	ExtractPath string // Logic for Phase 7.2 (JSON path to extract)
}

// ActiveFlow stores the sequence currently being recorded
var ActiveFlow []FlowStep

// EndpointEntry stores the path and the specific engines assigned by the analyzer.
type EndpointEntry struct {
	Path    string
	Engines []string
}

// DiscoveryStore manages the thread-safe inventory of discovered API endpoints.
type DiscoveryStore struct {
	Inventory map[string]*EndpointEntry
	mu        sync.RWMutex
}

// GlobalDiscovery is the centralized singleton for all discovery and engine operations.
var GlobalDiscovery = DiscoveryStore{
	Inventory: make(map[string]*EndpointEntry),
}

// GlobalClient (Unified Phase 8.4): Shared HTTP client for all tactical operations.
// Centralized here to avoid redeclaration collisions across the logic package.
var GlobalClient = &http.Client{
	Timeout: 30 * time.Second,
}

// SetGlobalClient allows the UI or main package to override the default client (e.g., for Proxy support).
func SetGlobalClient(c *http.Client) {
	GlobalClient = c
}

// AddEndpoint registers a new path in the inventory if it doesn't already exist.
func (ds *DiscoveryStore) AddEndpoint(path string) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	if _, exists := ds.Inventory[path]; !exists {
		ds.Inventory[path] = &EndpointEntry{
			Path:    path,
			Engines: []string{}, 
		}
	}
}

// GetEndpoints returns a list of all discovered paths for backward compatibility.
func (ds *DiscoveryStore) GetEndpoints() []string {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	keys := make([]string, 0, len(ds.Inventory))
	for k := range ds.Inventory {
		keys = append(keys, k)
	}
	return keys
}