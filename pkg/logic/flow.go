package logic

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/pterm/pterm"
	"github.com/tidwall/gjson"
)

// FlowContext stores state variables (e.g., {{user_id}}) captured from responses
var FlowContext = make(map[string]string)

// RunFlow executes the entire sequence in order (Phase 7.1)
func RunFlow() {
	if len(ActiveFlow) == 0 {
		pterm.Warning.Println("Tactical queue is empty.")
		return
	}

	pterm.DefaultSection.Println("Executing Sequential Logic Flow")

	for i, step := range ActiveFlow {
		executeStep(i, step, "FLOW_SEQUENCER")
	}
}

// RunStep executes a single isolated step for out-of-order testing (Phase 7.2)
func RunStep(index int) {
	if index < 0 || index >= len(ActiveFlow) {
		pterm.Error.Println("Invalid step index.")
		return
	}
	
	pterm.Warning.Printfln("PHASE 7.2: Probing Out-of-Order State Machine")
	executeStep(index, ActiveFlow[index], "STATE_PROBE")
}

// executeStep handles the shared logic of injection and execution
func executeStep(i int, step FlowStep, engine string) {
	// 1. Variable Injection
	finalURL := step.URL
	finalBody := step.Body
	for k, v := range FlowContext {
		placeholder := "{{" + k + "}}"
		finalURL = strings.ReplaceAll(finalURL, placeholder, v)
		finalBody = strings.ReplaceAll(finalBody, placeholder, v)
	}

	// 2. Build Request
	req, _ := http.NewRequest(step.Method, finalURL, bytes.NewBufferString(finalBody))
	
	// 3. Execution via Phase 6 Evasion Layer
	resp, err := SafeDo(req, true, engine)
	if err != nil {
		pterm.Error.Printfln("Step %d [%s] failed: %v", i+1, step.Name, err)
		return
	}
	defer resp.Body.Close()

	// 4. Data Extraction for State Mapping
	bodyBytes, _ := io.ReadAll(resp.Body)
	if step.ExtractPath != "" {
		captured := gjson.Get(string(bodyBytes), step.ExtractPath)
		if captured.Exists() {
			FlowContext[step.ExtractPath] = captured.String()
			pterm.Info.Printfln("Mapped State: {{%s}} = %s", step.ExtractPath, captured.String())
		}
	}

	pterm.Success.Printfln("Step %d [%s]: %d", i+1, step.Name, resp.StatusCode)
}