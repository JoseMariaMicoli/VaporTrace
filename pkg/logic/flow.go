package logic

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"sync"

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

// RunRace executes a high-concurrency synchronized attack (Phase 7.3)
func RunRace(index int, threads int) {
	if index < 0 || index >= len(ActiveFlow) {
		pterm.Error.Println("Invalid step index. Use 'flow list'.")
		return
	}

	step := ActiveFlow[index]
	var wg sync.WaitGroup
	
	// The Synchronized Starting Gate
	startGate := make(chan struct{})

	pterm.Warning.Printfln("PHASE 7.3: Priming %d concurrent threads against [%s]", threads, step.Name)

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func(threadID int) {
			defer wg.Done()
			
			// 1. Prepare Request (Inject current state)
			finalURL := step.URL
			finalBody := step.Body
			for k, v := range FlowContext {
				finalURL = strings.ReplaceAll(finalURL, "{{"+k+"}}", v)
				finalBody = strings.ReplaceAll(finalBody, "{{"+k+"}}", v)
			}
			
			req, _ := http.NewRequest(step.Method, finalURL, bytes.NewBufferString(finalBody))
			
			// 2. WAIT FOR GATE
			<-startGate 
			
			// 3. EXECUTE (Jitter is DISABLED for maximum collision probability)
			resp, err := SafeDo(req, false, "RACE_ENGINE") 
			
			if err == nil {
				// We only print successes or critical failures to avoid terminal flooding
				if resp.StatusCode < 400 {
					pterm.Success.Printfln("Thread %d | COLLISION SUCCESS: %d", threadID, resp.StatusCode)
				}
			}
		}(i)
	}

	pterm.Info.Println("All threads ready. Releasing synchronizer...")
	
	// BOOM: Close the channel to release all goroutines simultaneously
	close(startGate) 
	
	wg.Wait()
	pterm.Success.Println("Race Condition probe sequence complete.")
}