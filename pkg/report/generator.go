package report

import (
	"fmt"
	"os"
	"time"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/pterm/pterm"
)

// GenerateMissionDebrief compiles findings into a professional tactical report
func GenerateMissionDebrief() {
	pterm.DefaultHeader.WithFullWidth(false).Println("Compiling Strategic Mission Report")

	var startTime string
	db.DB.QueryRow("SELECT value FROM mission_state WHERE key = 'start_time'").Scan(&startTime)

	reportDate := time.Now().Format("2006-01-02")
	fileName := fmt.Sprintf("VAPOR_DEBRIEF_%s.md", reportDate)
	
	f, err := os.Create(fileName)
	if err != nil {
		pterm.Error.Printf("FileSystem Error: %v\n", err)
		return
	}
	defer f.Close()

	// Header - Red Team Format (Fixed Syntax)
	f.WriteString("# VAPORTRACE TACTICAL DEBRIEF\n")
	f.WriteString(fmt.Sprintf("> **OPERATIONAL STATUS:** COMPLETED\n"))
	f.WriteString(fmt.Sprintf("> **DATABASE ID:** 1 | **GEN TIME:** %s\n", time.Now().Format("15:04:05")))
	f.WriteString(fmt.Sprintf("> **START TIME:** %s\n\n", startTime))
	f.WriteString("---\n\n")

	// Phased Findings Table
	f.WriteString("## MISSION PHASES SUMMARY\n\n")
	
	phases := []string{
		"PHASE II: DISCOVERY",
		"PHASE III: AUTH LOGIC",
		"PHASE IV: INJECTION",
	}

	for _, phase := range phases {
		f.WriteString(fmt.Sprintf("### %s\n", phase))
		f.WriteString("| ATTACK VECTOR | TARGET ENDPOINT | RESULT | TIMESTAMP |\n")
		f.WriteString("| :--- | :--- | :--- | :--- |\n")

		rows, _ := db.DB.Query("SELECT details, target, status, timestamp FROM findings WHERE phase = ?", phase)
		
		hasData := false
		for rows.Next() {
			hasData = true
			var details, target, status, ts string
			rows.Scan(&details, &target, &status, &ts)
			
			// Format as a Markdown table row
			f.WriteString(fmt.Sprintf("| %s | `%s` | **%s** | %s |\n", details, target, status, ts))
		}
		rows.Close()

		if !hasData {
			f.WriteString("| - | - | *NO LOGS FOUND* | - |\n")
		}
		f.WriteString("\n")
	}

	// Footer
	f.WriteString("---\n")
	f.WriteString("**CONFIDENTIAL // INTERNAL RED TEAM USE ONLY**\n")
}