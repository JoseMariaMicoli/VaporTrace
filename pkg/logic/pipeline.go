package logic

import (
	"regexp"
	"strings"

	"github.com/pterm/pterm"
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

	pterm.DefaultSection.Println("Phase 9.5: Tactical Pipeline Analysis")
	
	inventoryCount := len(GlobalDiscovery.Inventory)
	if inventoryCount == 0 {
		pterm.Warning.Println("Discovery store is empty. No endpoints to analyze.")
		return
	}

	pterm.Info.Printfln("Analyzing %d endpoints for specialized testing...", inventoryCount)

	// Reset queue to ensure a clean state if re-run
	PipelineQueue = []PipelineTarget{}

	// Heuristics: Catch Swagger placeholders like {petId} or numeric segments.
	// This ensures BOLA logic triggers for the correct RESTful patterns.
	idRegex := regexp.MustCompile(`\{.*\}|[0-9]{1,}`)
	
	tableData := pterm.TableData{{"ENDPOINT", "ENGINES"}}

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
		
		// Add entry to the UI table
		tableData = append(tableData, []string{path, strings.Join(engines, ", ")})
	}

	// Render the tactical map to the console
	pterm.DefaultTable.WithHasHeader().WithData(tableData).WithBoxed().Render()
}

// RunPipeline iterates through the queue and triggers the relevant industrialized engines
func RunPipeline(concurrency int) {
	// Always ensure analysis is synchronized before execution
	AnalyzeDiscovery()

	if len(PipelineQueue) == 0 {
		pterm.Warning.Println("Pipeline queue is empty. Run 'map' or 'swagger' first.")
		return
	}

	pterm.DefaultSection.Println("Executing Industrialized Attack Pipeline")

	// Trigger specialized mass execution engines.
	// These functions (defined in bola.go, bfla.go, bopla.go) will now
	// pull targets directly from GlobalDiscovery.Inventory based on the tags.
	
	// Phase 9.7: Mass ID Enumeration (BOLA)
	ExecuteMassBOLA(concurrency)
	
	// Phase 9.9: Method Matrix Tampering (BFLA)
	ExecuteMassBFLA(concurrency)

	// Phase 9.8: Property Injection / Mass Assignment (BOPLA)
	ExecuteMassBOPLA(concurrency)

	pterm.Success.Println("Pipeline execution finished. Use 'report' to view findings.")
}