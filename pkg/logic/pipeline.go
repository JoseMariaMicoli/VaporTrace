package logic

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

// PipelineTarget maps a specific path to the engines authorized to attack it
type PipelineTarget struct {
	Path    string
	Engines []string
}

// PipelineQueue holds the analyzed targets ready for execution
var PipelineQueue []PipelineTarget

// AnalyzeDiscovery processes the GlobalDiscovery store and assigns attack vectors.
// It populates the metadata in GlobalDiscovery.Inventory so engines can self-filter.
func AnalyzeDiscovery() {
	GlobalDiscovery.mu.Lock()
	defer GlobalDiscovery.mu.Unlock()

	// Redirect feedback to LOG QUADRANT
	utils.TacticalLog("[cyan::b]PHASE 9.5: TACTICAL PIPELINE ANALYSIS[-:-:-]")

	inventoryCount := len(GlobalDiscovery.Inventory)
	if inventoryCount == 0 {
		utils.TacticalLog("[yellow]WARN:[-] Discovery store is empty. No endpoints to analyze.")
		return
	}

	utils.TacticalLog(fmt.Sprintf("[blue]⠋[-] Analyzing %d endpoints for specialized testing...", inventoryCount))

	// Reset queue to ensure a clean state if re-run
	PipelineQueue = []PipelineTarget{}

	// Heuristics: Catch Swagger placeholders like {petId} or numeric segments.
	// This ensures BOLA logic triggers for the correct RESTful patterns.
	idRegex := regexp.MustCompile(`\{.*\}|[0-9]{1,}`)

	for path, entry := range GlobalDiscovery.Inventory {
		// 1. Every endpoint is subject to BFLA (Method Matrix / Verb Tampering)
		engines := []string{"BFLA"}

		// 2. BOLA Detection: Look for ID-like patterns or variables in the path
		if idRegex.MatchString(path) {
			engines = append(engines, "BOLA")
		}

		// 3. BOPLA Detection: Look for keywords indicating resource mutation (Mass Assignment)
		lp := strings.ToLower(path)
		if strings.Contains(lp, "update") || strings.Contains(lp, "create") ||
			strings.Contains(lp, "set") || strings.Contains(lp, "profile") ||
			strings.Contains(lp, "pet") || strings.Contains(lp, "user") ||
			strings.Contains(lp, "order") {
			engines = append(engines, "BOPLA")
		}

		// Save tags back to the Inventory metadata for the specialized engines to read
		entry.Engines = engines

		// Build target for the PipelineQueue (used by the RunPipeline loop)
		target := PipelineTarget{
			Path:    path,
			Engines: engines,
		}
		PipelineQueue = append(PipelineQueue, target)
	}

	// Feedback to LOG QUADRANT
	utils.TacticalLog(fmt.Sprintf("[green]✔[-] Analysis complete. %d targets synchronized to pipeline.", len(PipelineQueue)))
}

// RunPipeline iterates through the queue and triggers the relevant industrialized engines
func RunPipeline(concurrency int) {
	// Always ensure analysis is synchronized before execution
	AnalyzeDiscovery()

	if len(PipelineQueue) == 0 {
		utils.TacticalLog("[red]ERROR:[-] Pipeline queue is empty. Run 'map' or 'swagger' first.")
		return
	}

	utils.TacticalLog("[cyan::b]EXECUTING INDUSTRIALIZED ATTACK PIPELINE[-:-:-]")

	// Trigger specialized mass execution engines.
	// These engines will internally call utils.RecordFinding, producing live logs.

	utils.TacticalLog("[blue]⠋[-] Launching Mass BOLA Engine (API1:2023)...")
	ExecuteMassBOLA(concurrency)

	utils.TacticalLog("[blue]⠋[-] Launching Mass BFLA Engine (API5:2023)...")
	ExecuteMassBFLA(concurrency)

	utils.TacticalLog("[blue]⠋[-] Launching Property Injection Engine (API3:2023)...")
	ExecuteMassBOPLA(concurrency)

	utils.TacticalLog("[green::b]PIPELINE EXECUTION FINISHED. Use 'report' to view findings.[-:-:-]")
}
