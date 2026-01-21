package logic

import (
	"regexp"
	"strings"
	"github.com/pterm/pterm"
)

type PipelineTarget struct {
	Path    string
	Engines []string
}

var PipelineQueue []PipelineTarget

func AnalyzeDiscovery() {
	GlobalDiscovery.mu.Lock()
	defer GlobalDiscovery.mu.Unlock()

	pterm.DefaultSection.Println("Phase 9.5: Tactical Pipeline Analysis")
	pterm.Info.Printfln("Analyzing %d endpoints for specialized testing...", len(GlobalDiscovery.Endpoints))

	// Reset queue to avoid duplicates on re-run
	PipelineQueue = []PipelineTarget{}

	idRegex := regexp.MustCompile(`\{.*\b(id|uuid|user|account)\b.*\}|[0-9]{3,}`)
	
	// Prepare table data for UI
	tableData := pterm.TableData{{"ENDPOINT", "ENGINES"}}

	for _, path := range GlobalDiscovery.Endpoints {
		target := PipelineTarget{Path: path, Engines: []string{"BFLA"}} 

		if idRegex.MatchString(path) {
			target.Engines = append(target.Engines, "BOLA")
		}

		lp := strings.ToLower(path)
		if strings.Contains(lp, "update") || strings.Contains(lp, "create") || strings.Contains(lp, "set") || strings.Contains(lp, "pet") {
			target.Engines = append(target.Engines, "BOPLA")
		}

		PipelineQueue = append(PipelineQueue, target)
		
		// Add to UI table
		tableData = append(tableData, []string{
			path, 
			strings.Join(target.Engines, ", "),
		})
	}

	// Render the findings
	pterm.DefaultTable.WithHasHeader().WithData(tableData).WithBoxed().Render()
	pterm.Success.Printfln("Analysis complete. %d targets categorized.", len(PipelineQueue))
}