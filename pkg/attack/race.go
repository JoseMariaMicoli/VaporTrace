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

package attack

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/logic"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

// RaceConfig defines the parameters for a synchronization attack
type RaceConfig struct {
	TargetURL string
	Method    string
	Body      string // For POST requests
	Threads   int
}

// RunRace executes a synchronized parallel attack
func RunRace(config RaceConfig) {
	utils.TacticalLog(fmt.Sprintf("[magenta::b]RACE ENGINE:[-] Priming %d threads against %s...", config.Threads, config.TargetURL))

	// 1. Configure High-Concurrency Client
	// We need a custom transport to avoid connection pooling bottlenecks during the race
	raceClient := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        1000,
			MaxIdleConnsPerHost: 1000,
			DisableKeepAlives:   true, // Force new connections to prevent serialization
		},
	}

	// 2. The Synchronization Gate
	// All threads will block on this channel until we close it
	startGate := make(chan struct{})
	var wg sync.WaitGroup

	// Results channel
	type raceResult struct {
		Status int
		Length int64
		ID     int
	}
	results := make(chan raceResult, config.Threads)

	// 3. Spawn Workers
	for i := 0; i < config.Threads; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Pre-build request to minimize latency drift
			req, _ := http.NewRequest(config.Method, config.TargetURL, nil)
			if config.Body != "" {
				// Note: In real scenarios, you might need a strings.Reader here
				// For brevity, assuming GET or empty body for now
			}

			// Copy auth headers from session
			if logic.CurrentSession.AttackerToken != "" {
				req.Header.Set("Authorization", "Bearer "+logic.CurrentSession.AttackerToken)
			}
			req.Header.Set("X-Race-ID", fmt.Sprintf("%d", id))

			// BLOCK HERE: Wait for the starting gun
			<-startGate

			// FIRE!
			resp, err := raceClient.Do(req)
			if err != nil {
				return
			}
			defer resp.Body.Close()

			// Metrics
			body, _ := io.ReadAll(resp.Body)
			results <- raceResult{
				Status: resp.StatusCode,
				Length: int64(len(body)),
				ID:     id,
			}
		}(i)
	}

	// 4. The Countdown
	utils.TacticalLog("[blue]RACE:[-] Workers ready. Synchronizing...")
	time.Sleep(500 * time.Millisecond) // Let goroutines spin up and hit the gate
	utils.TacticalLog("[red::b]>>> EXECUTE <<<[-]")

	// 5. Release the Gate (Nanosecond precision start)
	close(startGate)

	// 6. Wait and Analyze
	wg.Wait()
	close(results)

	// 7. Differential Analysis
	statusCounts := make(map[int]int)
	lengthCounts := make(map[int64]int)
	total := 0

	for res := range results {
		statusCounts[res.Status]++
		lengthCounts[res.Length]++
		total++
	}

	// Report Results
	utils.TacticalLog(fmt.Sprintf("[green]RACE COMPLETE:[-] %d requests landed.", total))

	// Logic: If we successfully exploited a race, we often see mixed status codes
	// e.g., 5 requests got 200 OK (Coupon applied), 5 requests got 400 Bad Request
	if len(statusCounts) > 1 {
		utils.TacticalLog("[red]⚠ RACE ANOMALY DETECTED (Status Variance)[-]")
		for code, count := range statusCounts {
			utils.TacticalLog(fmt.Sprintf("  Code %d: %d responses", code, count))
		}

		// Log Finding
		utils.RecordFinding(db.Finding{
			Phase:        "PHASE III: RACE CONDITION",
			Command:      "race",
			Target:       config.TargetURL,
			Details:      fmt.Sprintf("Race Anomaly: Mixed Status Codes detected with %d threads", config.Threads),
			Status:       "CRITICAL",
			OWASP_ID:     "API6:2023",
			CVSS_Numeric: 8.5,
		})
	} else if len(lengthCounts) > 1 {
		utils.TacticalLog("[yellow]⚠ RACE WARNING (Length Variance)[-]")
		// Usually indicates one request succeeded differently than others
		utils.RecordFinding(db.Finding{
			Phase:   "PHASE III: RACE CONDITION",
			Command: "race",
			Target:  config.TargetURL,
			Details: "Race Anomaly: Response length variation detected",
			Status:  "HIGH",
		})
	} else {
		utils.TacticalLog("[gray]No variances detected (All responses identical).[-]")
	}
}
