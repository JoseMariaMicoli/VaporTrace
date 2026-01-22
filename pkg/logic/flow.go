package logic

import (
	"bytes"
	"net/http"

	"github.com/pterm/pterm"
)

// RunFlow executes the recorded sequence through the Evasion gatekeeper
func RunFlow() {
	if len(ActiveFlow) == 0 {
		pterm.Warning.Println("Tactical Flow is empty. Use 'flow add' to build a sequence.")
		return
	}

	pterm.DefaultSection.Println("Executing Business Logic Flow")

	for i, step := range ActiveFlow {
		req, err := http.NewRequest(step.Method, step.URL, bytes.NewBufferString(step.Body))
		if err != nil {
			pterm.Error.Printfln("Failed to build request for step %d: %v", i+1, err)
			continue
		}

		// Apply tactical session tokens
		if CurrentSession.AttackerToken != "" {
			req.Header.Set("Authorization", "Bearer "+CurrentSession.AttackerToken)
		}

		// Execute via Phase 6/9 Gatekeeper
		resp, err := SafeDo(req, true, "FLOW_ENGINE")
		if err != nil {
			pterm.Error.Printfln("Step %d [%s] failed: %v", i+1, step.Name, err)
			continue
		}
		resp.Body.Close()

		pterm.Success.Printfln("Step %d [%s]: %s %d OK", i+1, step.Name, step.Method, resp.StatusCode)
	}
}