package logic

import (
	"fmt"
	"net/http"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db" // Added Persistence
	"github.com/pterm/pterm"
)

type BFLAContext struct {
	TargetURL string
}

// Administrative verbs to test for Function Level escalation
var bflaMethods = []string{"POST", "PUT", "DELETE", "PATCH"}

func (b *BFLAContext) Probe() {
	pterm.DefaultHeader.WithFullWidth(false).Println("BFLA / Function Level Probe (API5:2023)")

	activeToken := CurrentSession.AttackerToken
	if activeToken == "" {
		pterm.Warning.Println("No Attacker Token found. Probing without Authorization header (Testing baseline)...")
	}

	// PATCH: Removed local client definition

	for _, method := range bflaMethods {
		pterm.Info.Printf("Testing Method Shuffling: [%s] -> %s\n", method, b.TargetURL)

		req, _ := http.NewRequest(method, b.TargetURL, nil)
		if activeToken != "" {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", activeToken))
		}

		// PATCH: Using GlobalClient
		resp, err := GlobalClient.Do(req)
		if err != nil {
			pterm.Error.Printf("Connection error for %s: %v\n", method, err)
			continue
		}
		defer resp.Body.Close()

		// Analysis: 2xx or 3xx on unauthorized verbs often indicates BFLA
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			pterm.Warning.Prefix = pterm.Prefix{Text: "VULN", Style: pterm.NewStyle(pterm.BgRed, pterm.FgWhite)}
			pterm.Warning.Printf("BFLA POTENTIAL: Server accepted %s (Status: %d)\n", method, resp.StatusCode)

			// PERSISTENCE HOOK
			db.LogQueue <- db.Finding{
				Phase:   "PHASE III: AUTH LOGIC",
				Target:  b.TargetURL,
				Details: fmt.Sprintf("BFLA Method Allowed: %s", method),
				Status:  "UNAUTHORIZED ACCESS",
			}
		} else {
			pterm.Info.Printf("Verb %s rejected (Status: %d)\n", method, resp.StatusCode)
		}
	}
}