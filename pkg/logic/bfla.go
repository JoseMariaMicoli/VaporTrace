package logic

import (
	"fmt"
	"net/http"
	"sync" // Required for Phase 9.9 Concurrency

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/pterm/pterm"
)

// BFLAContext defines the parameters for Functional Level attacks
type BFLAContext struct {
	TargetURL string
}

// Administrative verbs to test for Function Level escalation (Method Matrix)
var bflaMethods = []string{"POST", "PUT", "DELETE", "PATCH"}

// ExecuteMassBFLA (Phase 9.9) runs a Method Matrix attack across the pipeline
func ExecuteMassBFLA(concurrency int) {
	pterm.DefaultSection.Println("Phase 9.9: Industrialized BFLA Engine")

	var targets []string
	GlobalDiscovery.mu.Lock()
	for _, path := range GlobalDiscovery.Endpoints {
		targets = append(targets, path)
	}
	GlobalDiscovery.mu.Unlock()

	if len(targets) == 0 {
		pterm.Warning.Println("Discovery pipeline is empty. Run 'swagger' or 'map' first.")
		return
	}

	for _, url := range targets {
		pterm.Info.Printfln("Testing Method Matrix for: %s", url)
		
		ctx := &BFLAContext{TargetURL: url}
		ctx.MassProbe(concurrency)
	}
}

// MassProbe implements the Phase 9.3 Worker Pool for Verb Tampering
func (b *BFLAContext) MassProbe(threads int) {
	methodChan := make(chan string, len(bflaMethods))
	var wg sync.WaitGroup

	// 1. Spawn Worker Pool
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for method := range methodChan {
				b.TamperAndProbeSilent(method)
			}
		}()
	}

	// 2. Feed Method Matrix
	for _, m := range bflaMethods {
		methodChan <- m
	}
	close(methodChan)

	// 3. Wait for workers
	wg.Wait()
}

// TamperAndProbeSilent handles high-speed execution with SafeDo mirroring (Phase 9.6)
func (b *BFLAContext) TamperAndProbeSilent(method string) {
	req, _ := http.NewRequest(method, b.TargetURL, nil)
	
	activeToken := CurrentSession.AttackerToken
	if activeToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", activeToken))
	}

	// SafeDo (Phase 9.6) Mirroring
	resp, err := SafeDo(req, false, "BFLA-ENGINE")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Analysis: 2xx/3xx on unauthorized verbs indicates BFLA
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		pterm.Warning.Prefix = pterm.Prefix{Text: "HIT", Style: pterm.NewStyle(pterm.BgCyan, pterm.FgBlack)}
		pterm.Warning.Printfln(" BFLA HIT: %s accepted unauthorized %s", b.TargetURL, method)

		// PERSISTENCE
		db.LogQueue <- db.Finding{
			Phase:   "PHASE 9.9: INDUSTRIALIZED BFLA",
			Target:  b.TargetURL,
			Details: fmt.Sprintf("BFLA (Verb Tampering) confirmed via %s", method),
			Status:  "VULNERABLE",
		}

		// Explicit Hit-Mirroring for IR visibility
		go func() {
			hitReq, _ := http.NewRequest(method, b.TargetURL, nil)
			if activeToken != "" {
				hitReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", activeToken))
			}
			SafeDo(hitReq, true, "BFLA-ENGINE")
		}()
	}
}

// Probe remains for legacy surgical single-threaded testing
func (b *BFLAContext) Probe() {
	pterm.DefaultHeader.WithFullWidth(false).Println("BFLA Surgical Probe")
	b.MassProbe(1)
}