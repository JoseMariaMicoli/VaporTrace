package logic

import (
	"net/http"
	"sync"
	"time"
)

// EndpointEntry stores the path and the specific engines assigned by the analyzer.
// This is the core metadata structure that allows the Pipeline to route tasks.
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
// It replaces the previous simple map to allow for industrialized metadata tagging.
var GlobalDiscovery = DiscoveryStore{
	Inventory: make(map[string]*EndpointEntry),
}

// GlobalClient is the shared HTTP client used for tactical requests.
var GlobalClient = &http.Client{
	Timeout: 10 * time.Second,
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
			Engines: []string{}, // This will be populated by the AnalyzeDiscovery function in pipeline.go
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