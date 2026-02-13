package logic

import (
	"fmt"
	"net/http"
	"net/http/httputil"
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

// === SPRINT 11.2: DataSilo for Autonomous Chaining ===
// DataSilo stores extracted loot (credentials, tokens, env vars) for use by subsequent actions
// Thread-safe storage for cross-action data dependencies
type DataSilo struct {
	mu   sync.RWMutex
	data map[string]interface{} // e.g., {"k8s_token": "eyJ0...", "admin_session": "xyz..."}
}

// GlobalDataSilo is the centralized singleton for autonomous action chaining
var GlobalDataSilo = &DataSilo{
	data: make(map[string]interface{}),
}

// Set stores a value in the DataSilo
func (ds *DataSilo) Set(key string, value interface{}) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.data[key] = value
}

// Get retrieves a value from the DataSilo
func (ds *DataSilo) Get(key string) (interface{}, bool) {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	val, ok := ds.data[key]
	return val, ok
}

// GetString retrieves a value as a string
func (ds *DataSilo) GetString(key string) string {
	if val, ok := ds.Get(key); ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// Exists checks if a key exists in the DataSilo
func (ds *DataSilo) Exists(key string) bool {
	_, ok := ds.Get(key)
	return ok
}

// List returns all keys in the DataSilo (for debugging)
func (ds *DataSilo) List() []string {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	keys := make([]string, 0, len(ds.data))
	for k := range ds.data {
		keys = append(keys, k)
	}
	return keys
}

// === TASK 4: TRAFFIC VAULT ===

// TrafficEntry represents a complete HTTP transaction for history
type TrafficEntry struct {
	ID           string
	Timestamp    time.Time
	Method       string
	URL          string
	StatusCode   int
	RequestDump  []byte
	ResponseDump []byte
	Duration     time.Duration
	Size         int64
}

// TrafficVault implements a thread-safe Ring Buffer
type TrafficVault struct {
	mu       sync.RWMutex
	Buffer   []*TrafficEntry
	Capacity int
	Head     int // Points to the next write position
	Count    int // Current number of items
	Lookup   map[string]*TrafficEntry
}

var GlobalVault = &TrafficVault{
	Capacity: 1000,
	Buffer:   make([]*TrafficEntry, 1000),
	Lookup:   make(map[string]*TrafficEntry),
}

// Add inserts a transaction into the ring buffer
func (tv *TrafficVault) Add(req *http.Request, resp *http.Response, duration time.Duration) {
	tv.mu.Lock()
	defer tv.mu.Unlock()

	reqDump, _ := httputil.DumpRequestOut(req, true)
	respDump, _ := httputil.DumpResponse(resp, true)

	id := fmt.Sprintf("%d", time.Now().UnixNano())

	entry := &TrafficEntry{
		ID:           id,
		Timestamp:    time.Now(),
		Method:       req.Method,
		URL:          req.URL.String(),
		StatusCode:   resp.StatusCode,
		RequestDump:  reqDump,
		ResponseDump: respDump,
		Duration:     duration,
		Size:         resp.ContentLength,
	}

	// Handle Ring Buffer Overwrite
	if tv.Count == tv.Capacity {
		// We are overwriting the oldest entry. Remove it from lookup map.
		oldest := tv.Buffer[tv.Head]
		if oldest != nil {
			delete(tv.Lookup, oldest.ID)
		}
	} else {
		tv.Count++
	}

	tv.Buffer[tv.Head] = entry
	tv.Lookup[id] = entry

	// Advance Head
	tv.Head = (tv.Head + 1) % tv.Capacity
}

// GetByID retrieves a specific entry
func (tv *TrafficVault) GetByID(id string) *TrafficEntry {
	tv.mu.RLock()
	defer tv.mu.RUnlock()
	return tv.Lookup[id]
}

// GetAll returns a slice of entries ordered from newest to oldest
func (tv *TrafficVault) GetAll() []*TrafficEntry {
	tv.mu.RLock()
	defer tv.mu.RUnlock()

	result := make([]*TrafficEntry, 0, tv.Count)

	// Start from the latest entry and move backwards
	idx := (tv.Head - 1 + tv.Capacity) % tv.Capacity

	for i := 0; i < tv.Count; i++ {
		if tv.Buffer[idx] != nil {
			result = append(result, tv.Buffer[idx])
		}
		idx = (idx - 1 + tv.Capacity) % tv.Capacity
	}
	return result
}
