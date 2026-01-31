package logic

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
	"github.com/pterm/pterm"
)

// BFLAContext defines the parameters for Functional Level attacks
type BFLAContext struct {
	TargetURL string
}

var bflaMethods = []string{"POST", "PUT", "DELETE", "PATCH"}

func ExecuteMassBFLA(concurrency int) {
	pterm.DefaultSection.Println("Phase 9.9: Industrialized BFLA Engine")

	GlobalDiscovery.mu.RLock()
	var targets []string
	for path := range GlobalDiscovery.Inventory {
		targets = append(targets, path)
	}
	GlobalDiscovery.mu.RUnlock()

	if len(targets) == 0 {
		utils.TacticalLog("Discovery pipeline is empty. Run 'map' first.")
		return
	}

	for _, path := range targets {
		ctx := &BFLAContext{TargetURL: CurrentSession.TargetURL + path}
		ctx.MassProbe(concurrency)
	}
}

func (b *BFLAContext) MassProbe(threads int) {
	methodChan := make(chan string, len(bflaMethods))
	var wg sync.WaitGroup

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for method := range methodChan {
				b.TamperAndProbeSilent(method)
			}
		}()
	}

	for _, m := range bflaMethods {
		methodChan <- m
	}
	close(methodChan)
	wg.Wait()
}

func (b *BFLAContext) TamperAndProbeSilent(method string) {
	req, _ := http.NewRequest(method, b.TargetURL, nil)
	req.Header.Set("User-Agent", "VaporTrace/2.1.0 (Phase 9.10 Industrialized)")
	
	activeToken := CurrentSession.AttackerToken
	if activeToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", activeToken))
	}

	resp, err := SafeDo(req, false, "BFLA-ENGINE")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		// PATCHED: Unified Logging with Phase 9.13 Tags
		utils.RecordFinding(db.Finding{
			Phase:    "PHASE III: EXPLOITATION",
			Target:   b.TargetURL,
			Details:  fmt.Sprintf("BFLA Method Matrix Success: %s returned %d", method, resp.StatusCode),
			Status:   "VULNERABLE",
			OWASP_ID: "API5:2023",
			MITRE_ID: "T1548.002", // Abuse Elevation Control Mechanism: Bypass User Account Control
			NIST_Tag: "DE.CM",     // Detection Processes
		})
	}
}