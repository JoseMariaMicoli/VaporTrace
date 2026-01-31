package logic

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
	"github.com/pterm/pterm"
)

type IntegrationContext struct {
	TargetURL       string
	IntegrationType string 
}

func (i *IntegrationContext) Probe() {
	pterm.DefaultHeader.WithFullWidth(false).Println("API10: Unsafe Consumption / Integration Probe")

	payloads := map[string]string{
		"GitHub-Spoof":      `{"repository": {"url": "http://169.254.169.254/latest/meta-data/"}, "sender": {"login": "vapor-trace"}}`,
		"Stripe-Spoof":      `{"id": "evt_test", "type": "customer.subscription.deleted", "data": {"object": {"metadata": {"internal_admin": "true"}}}}`,
		"Generic-Injection": `{"source": "third-party", "data": {"cmd": "sleep 10"}}`,
	}

	for name, body := range payloads {
		req, _ := http.NewRequest("POST", i.TargetURL, bytes.NewBuffer([]byte(body)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Vapor-Trace-Tactical", "API10-Probe")
		req.Header.Set("User-Agent", fmt.Sprintf("VaporTrace-%s-Scanner", i.IntegrationType))

		start := time.Now()
		resp, err := SafeDo(req, false, "API10-PROBE")
		duration := time.Since(start)

		if err != nil {
			continue
		}
		defer resp.Body.Close()

		pterm.Info.Printf("Payload: %-18s | Status: %d | Latency: %v\n", name, resp.StatusCode, duration)

		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			// PATCHED: Unified Logging with Phase 9.13 Tags
			utils.RecordFinding(db.Finding{
				Phase:    "PHASE IV: INJECTION",
				Target:   i.TargetURL,
				Details:  fmt.Sprintf("Unsafe Integration [%s]: Accepted unsigned %s payload", i.IntegrationType, name),
				Status:   "VULNERABLE",
				OWASP_ID: "API10:2023",
				MITRE_ID: "T1190", // Exploit Public-Facing Application
				NIST_Tag: "ID.RA", // Risk Assessment
			})
		}
	}
}