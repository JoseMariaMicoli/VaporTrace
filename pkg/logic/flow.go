package logic

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
	"github.com/tidwall/gjson"
)

// FlowContext stores state variables
var FlowContext = make(map[string]string)

func RunFlow() {
	if len(ActiveFlow) == 0 {
		utils.TacticalLog("[yellow]Tactical queue is empty.[-]")
		return
	}

	utils.TacticalLog("[cyan]Executing Sequential Logic Flow[-]")

	for i, step := range ActiveFlow {
		executeStep(i, step, "FLOW_SEQUENCER")
	}
}

func RunStep(index int) {
	if index < 0 || index >= len(ActiveFlow) {
		utils.TacticalLog("[red]Invalid step index.[-]")
		return
	}

	utils.TacticalLog("[yellow]PHASE 7.2: Probing Out-of-Order State Machine[-]")
	executeStep(index, ActiveFlow[index], "STATE_PROBE")
}

func executeStep(i int, step FlowStep, engine string) {
	finalURL := step.URL
	finalBody := step.Body
	for k, v := range FlowContext {
		placeholder := "{{" + k + "}}"
		finalURL = strings.ReplaceAll(finalURL, placeholder, v)
		finalBody = strings.ReplaceAll(finalBody, placeholder, v)
	}

	req, _ := http.NewRequest(step.Method, finalURL, bytes.NewBufferString(finalBody))

	resp, err := SafeDo(req, true, engine)
	if err != nil {
		utils.TacticalLog(fmt.Sprintf("[red]Step %d [%s] failed: %v[-]", i+1, step.Name, err))
		return
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if step.ExtractPath != "" {
		captured := gjson.Get(string(bodyBytes), step.ExtractPath)
		if captured.Exists() {
			FlowContext[step.ExtractPath] = captured.String()
			utils.TacticalLog(fmt.Sprintf("[green]Mapped State: {{%s}} = %s[-]", step.ExtractPath, captured.String()))
		}
	}

	utils.TacticalLog(fmt.Sprintf("[green]Step %d [%s]: %d[-]", i+1, step.Name, resp.StatusCode))
}

func RunRace(index int, threads int) {
	if index < 0 || index >= len(ActiveFlow) {
		utils.TacticalLog("[red]Invalid step index.[-]")
		return
	}

	step := ActiveFlow[index]
	var wg sync.WaitGroup
	startGate := make(chan struct{})

	utils.TacticalLog(fmt.Sprintf("[yellow]PHASE 7.3: Priming %d concurrent threads against [%s][-]", threads, step.Name))

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func(threadID int) {
			defer wg.Done()

			finalURL := step.URL
			finalBody := step.Body
			for k, v := range FlowContext {
				finalURL = strings.ReplaceAll(finalURL, "{{"+k+"}}", v)
				finalBody = strings.ReplaceAll(finalBody, "{{"+k+"}}", v)
			}

			req, _ := http.NewRequest(step.Method, finalURL, bytes.NewBufferString(finalBody))

			<-startGate

			resp, err := SafeDo(req, false, "RACE_ENGINE")

			if err == nil {
				if resp.StatusCode < 400 {
					utils.TacticalLog(fmt.Sprintf("[green]Thread %d | COLLISION SUCCESS: %d[-]", threadID, resp.StatusCode))
				}
			}
		}(i)
	}

	utils.TacticalLog("[blue]All threads ready. Releasing synchronizer...[-]")
	close(startGate)
	wg.Wait()
	utils.TacticalLog("[green]Race Condition probe sequence complete.[-]")
}
