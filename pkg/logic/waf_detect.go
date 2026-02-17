/*
Copyright (c) 2026 José María Micoli
Licensed under the Business Source License 1.1
Change Date: 2033-02-17
Change License: Apache-2.0

You may:
✔ Study
✔ Modify
✔ Use for internal security testing

You may NOT:
✘ Offer as a commercial service
✘ Sell derived competing products
*/

package logic

import (
	"fmt"
	"sync"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

// WAFDetectionStats holds comprehensive WAF/IDS detection statistics
type WAFDetectionStats struct {
	RateLimitBlocks int
	WAFBlocks       int
	Redirects       int
	ServerErrors    int
	Detected        bool
	DetectionTypes  []string
}

// WAFDetectionManager manages WAF detection state
type WAFDetectionManager struct {
	mu              sync.RWMutex
	Stats           WAFDetectionStats
	LastUpdate      int64
	MonitoringState bool
}

var globalWAFDetection = &WAFDetectionManager{
	MonitoringState: true,
	Stats: WAFDetectionStats{
		RateLimitBlocks: 0,
		WAFBlocks:       0,
		Redirects:       0,
		ServerErrors:    0,
		Detected:        false,
		DetectionTypes:  []string{},
	},
}

// GetWAFDetectionStats returns current WAF detection statistics as a map
func GetWAFDetectionStats() map[string]interface{} {
	globalWAFDetection.mu.RLock()
	defer globalWAFDetection.mu.RUnlock()

	stats := map[string]interface{}{
		"rate_limit_blocks": globalWAFDetection.Stats.RateLimitBlocks,
		"waf_blocks":        globalWAFDetection.Stats.WAFBlocks,
		"redirects":         globalWAFDetection.Stats.Redirects,
		"server_errors":     globalWAFDetection.Stats.ServerErrors,
		"detected":          globalWAFDetection.Stats.Detected,
	}

	// Get rate limit backoff status (tracks 429 responses)
	runsStats := GetRateLimitStatus()
	if runsStats != nil {
		if count, ok := runsStats["count"].(int); ok {
			stats["rate_limit_blocks"] = count
		}
	}

	// If we have significant rate limiting detected, flag as potential WAF
	rateLimit := 0
	if rlValue, ok := stats["rate_limit_blocks"].(int); ok {
		rateLimit = rlValue
	}
	if rateLimit > 3 {
		stats["detected"] = true
	}

	return stats
}

// ReportWAFDetection outputs comprehensive WAF detection status
func ReportWAFDetection() {
	globalWAFDetection.mu.RLock()
	defer globalWAFDetection.mu.RUnlock()

	utils.TacticalLog("[magenta::b]PHASE 6.4: WAF DETECTION ENGINE[-:-:-]")
	utils.TacticalLog("[yellow]Status:[-] [cyan]ACTIVE & MONITORING")
	utils.TacticalLog("")

	utils.TacticalLog("[yellow]Monitored Patterns:[-]")
	utils.TacticalLog("  [blue]•[-] Rate Limits (429)")
	utils.TacticalLog("  [blue]•[-] WAF Blocks (403)")
	utils.TacticalLog("  [blue]•[-] Honeypots (Custom Redirects)")
	utils.TacticalLog("  [blue]•[-] Signature-based Injection (500 errors)")

	utils.TacticalLog("")
	utils.TacticalLog("[yellow]Current Statistics:[-]")
	utils.TacticalLog(fmt.Sprintf("  Rate Limit (429): [cyan]%d[-] blocks", globalWAFDetection.Stats.RateLimitBlocks))
	utils.TacticalLog(fmt.Sprintf("  WAF Blocks (403): [cyan]%d[-] blocks", globalWAFDetection.Stats.WAFBlocks))
	utils.TacticalLog(fmt.Sprintf("  Redirects (30x): [cyan]%d[-]", globalWAFDetection.Stats.Redirects))
	utils.TacticalLog(fmt.Sprintf("  Server Errors (50x): [cyan]%d[-]", globalWAFDetection.Stats.ServerErrors))

	utils.TacticalLog("")
	if globalWAFDetection.Stats.Detected {
		utils.TacticalLog("[red]⚠ WAF/IDS DETECTED[-] - Recommend switching to '[green]stealth silent[-]' mode")
		if len(globalWAFDetection.Stats.DetectionTypes) > 0 {
			utils.TacticalLog("[yellow]Detection Types:[-]")
			for _, dt := range globalWAFDetection.Stats.DetectionTypes {
				utils.TacticalLog(fmt.Sprintf("  [red]•[-] %s", dt))
			}
		}
	} else {
		utils.TacticalLog("[green]✓ No active WAF detection patterns observed")
	}

	utils.TacticalLog("")
	utils.TacticalLog("[yellow]Evasion Recommendation:[-] Enable 'stealth silent' mode for WAF-protected targets")
	utils.TacticalLog("[cyan]Tip:[-] Use 'loot list' to see captured WAF responses")
}

// UpdateWAFStats updates WAF detection statistics
func UpdateWAFStats(rateLimitDelta, wafBlocksDelta, redirectsDelta, serverErrorsDelta int) {
	globalWAFDetection.mu.Lock()
	defer globalWAFDetection.mu.Unlock()

	globalWAFDetection.Stats.RateLimitBlocks += rateLimitDelta
	globalWAFDetection.Stats.WAFBlocks += wafBlocksDelta
	globalWAFDetection.Stats.Redirects += redirectsDelta
	globalWAFDetection.Stats.ServerErrors += serverErrorsDelta

	// Determine if WAF is detected
	if globalWAFDetection.Stats.RateLimitBlocks > 3 ||
		globalWAFDetection.Stats.WAFBlocks > 0 ||
		globalWAFDetection.Stats.Redirects > 2 {
		globalWAFDetection.Stats.Detected = true
	}
}

// RecordWAFDetectionType adds a specific detection type to the list
func RecordWAFDetectionType(detectionType string) {
	globalWAFDetection.mu.Lock()
	defer globalWAFDetection.mu.Unlock()

	// Check if already in list
	for _, dt := range globalWAFDetection.Stats.DetectionTypes {
		if dt == detectionType {
			return
		}
	}

	globalWAFDetection.Stats.DetectionTypes = append(globalWAFDetection.Stats.DetectionTypes, detectionType)
	globalWAFDetection.Stats.Detected = true
}

// ResetWAFDetection resets all WAF detection statistics
func ResetWAFDetection() {
	globalWAFDetection.mu.Lock()
	defer globalWAFDetection.mu.Unlock()

	globalWAFDetection.Stats = WAFDetectionStats{
		RateLimitBlocks: 0,
		WAFBlocks:       0,
		Redirects:       0,
		ServerErrors:    0,
		Detected:        false,
		DetectionTypes:  []string{},
	}
}
