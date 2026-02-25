package logic

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/pterm/pterm"
)

// BFLAContext defines the parameters for Functional Level attacks
type BFLAContext struct {
	TargetURL string
}

// Administrative verbs to test for Function Level escalation (API5:2023)
var bflaMethods = []string{"POST", "PUT", "DELETE", "PATCH"}

// ExecuteMassBFLA runs a Method Matrix attack across the discovery pipeline.
// It iterates through the GlobalDiscovery.Inventory metadata.
func ExecuteMassBFLA(concurrency int) {
	pterm.DefaultSection.Println("Phase 9.9: Industrialized BFLA Engine")

	GlobalDiscovery.mu.RLock()
	var targets []string
	// Every endpoint in the inventory is a target for BFLA by default
	for path := range GlobalDiscovery.Inventory {
		targets = append(targets, path)
	}
	GlobalDiscovery.mu.RUnlock()

	if len(targets) == 0 {
		pterm.Warning.Println("Discovery pipeline is empty. Run 'swagger' or 'map' first.")
		return
	}

	for _, path := range targets {
		pterm.Info.Printfln("Testing Method Matrix for: %s", path)
		
		// Ensure the full target URL is constructed
		ctx := &BFLAContext{TargetURL: CurrentSession.TargetURL + path}
		ctx.MassProbe(concurrency)
	}
}

// MassProbe implements the Worker Pool for Verb Tampering
func (b *BFLAContext) MassProbe(threads int) {
	methodChan := make(chan string, len(bflaMethods))
	var wg sync.WaitGroup

	// Initialize Worker Pool
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for method := range methodChan {
				b.TamperAndProbeSilent(method)
			}
		}()
	}

	// Feed verbs into the channel
	for _, m := range bflaMethods {
		methodChan <- m
	}
	close(methodChan)
	wg.Wait()
}

// TamperAndProbeSilent handles high-speed execution with SafeDo mirroring.
// It looks for unauthorized verbs that return success codes.
func (b *BFLAContext) TamperAndProbeSilent(method string) {
	req, _ := http.NewRequest(method, b.TargetURL, nil)
	req.Header.Set("User-Agent", "VaporTrace/2.1.0 (Phase 9.10 Industrialized)")
	
	activeToken := CurrentSession.AttackerToken
	if activeToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", activeToken))
	}

	// Execute via the networking gatekeeper (SafeDo)
	// We pass false for isHit initially as analysis happens after the call
	resp, err := SafeDo(req, false, "BFLA-ENGINE")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Analysis: 2xx/3xx on unauthorized verbs (like DELETE on a resource) indicates BFLA
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		pterm.Warning.Prefix = pterm.Prefix{Text: "HIT", Style: pterm.NewStyle(pterm.BgRed, pterm.FgWhite)}
		pterm.Warning.Printfln("BFLA Potential: Verb %s accepted at %s (Status: %d)", method, b.TargetURL, resp.StatusCode)

		// Persistence: Log the hit to the database
		db.LogQueue <- db.Finding{
			Phase:   "PHASE III: EXPLOITATION",
			Target:  b.TargetURL,
			Details: fmt.Sprintf("BFLA Method Matrix Success: %s returned %d", method, resp.StatusCode),
			Status:  "VULNERABLE",
		}
	}
}