package logic

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

type BFLAContext struct {
	TargetURL string
}

var bflaMethods = []string{"POST", "PUT", "DELETE", "PATCH"}

func ExecuteMassBFLA(concurrency int) {
	// FIX: Removed pterm
	utils.TacticalLog("[cyan::b]Phase 9.9: Industrialized BFLA Engine Started[-:-:-]")

	GlobalDiscovery.mu.RLock()
	var targets []string
	for path := range GlobalDiscovery.Inventory {
		targets = append(targets, path)
	}
	GlobalDiscovery.mu.RUnlock()

	if len(targets) == 0 {
		utils.TacticalLog("[yellow]Discovery pipeline is empty. Run 'map' first.[-]")
		return
	}

	for _, path := range targets {
		ctx := &BFLAContext{TargetURL: CurrentSession.TargetURL + path}
		ctx.MassProbe(concurrency)
	}
	utils.TacticalLog("[green::b]BFLA Engine Execution Completed.[-:-:-]")
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
		utils.RecordFinding(db.Finding{
			Phase:   "PHASE III: EXPLOITATION",
			Command: "bfla", // Zero-Touch Trigger
			Target:  b.TargetURL,
			Details: fmt.Sprintf("BFLA Method Matrix Success: %s returned %d", method, resp.StatusCode),
			Status:  "VULNERABLE",
		})
	}
}
