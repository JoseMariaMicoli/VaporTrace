package logic

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

// AggregatorConfig controls the background correlation service
type AggregatorConfig struct {
	Interval time.Duration
	Active   bool
}

// GlobalAggregator singleton
var GlobalAggregator = &AggregatorConfig{
	Interval: 10 * time.Second,
	Active:   true,
}

// StartContextAggregator launches the intelligence gathering engine in a goroutine
func StartContextAggregator() {
	utils.TacticalLog("[magenta]CONTEXT AGGREGATOR:[-] Service Started (Background Correlation Active)")
	go func() {
		for GlobalAggregator.Active {
			runCorrelationCycle()
			time.Sleep(GlobalAggregator.Interval)
		}
	}()
}

// runCorrelationCycle merges Loot and Map data into the Context Store
func runCorrelationCycle() {
	// 1. Process Loot Vault
	for _, item := range Vault {
		contextType := "Info"
		if item.Type == "JWT_TOKEN" || item.Type == "AWS_KEY" {
			contextType = "Credential"
		} else if item.Type == "EMAIL" {
			contextType = "PII"
		}

		// Correlate: Store the secret associated with its source URL
		db.StoreContext(db.ContextRow{
			Scope:    item.Source,
			DataType: contextType,
			Key:      item.Type,
			Value:    item.Value,
			Source:   "Loot-Vault",
		})
	}

	// 2. Process Discovery Inventory
	GlobalDiscovery.mu.RLock()
	for path, _ := range GlobalDiscovery.Inventory {
		if strings.Contains(strings.ToLower(path), "admin") || strings.Contains(strings.ToLower(path), "config") {
			db.StoreContext(db.ContextRow{
				Scope:    path,
				DataType: "HVT",
				Key:      "Endpoint-Keyword",
				Value:    path,
				Source:   "Recon-Map",
			})
		}
	}
	GlobalDiscovery.mu.RUnlock()
}

// EnrichCommandRequest is called by SafeDo/Transport to inject context-aware data
func EnrichCommandRequest(req *http.Request) {
	target := req.URL.String()

	// Retrieve relevant context from DB
	rows, err := db.GetContext(target)
	if err != nil {
		return
	}

	injected := false
	for _, row := range rows {
		if row.DataType == "Credential" {
			// Auto-Inject Bearer Token
			if row.Key == "JWT_TOKEN" && req.Header.Get("Authorization") == "" {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", row.Value))
				injected = true
			}
			// Auto-Inject AWS Keys
			if row.Key == "AWS_KEY" && req.Header.Get("X-Amz-Date") == "" {
				req.Header.Set("X-Amz-Date", time.Now().Format(time.RFC3339))
				injected = true
			}
		}
	}

	if injected {
		// Log to Context Tab (F5)
		utils.LogContext(fmt.Sprintf("[magenta]INJECTION >>[-] Applied credentials to [white]%s[-]", req.URL.Path))
	}
}
