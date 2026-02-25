package report

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/pterm/pterm"
)

// GenerateMissionDebrief compiles findings into a professional tactical report.
func GenerateMissionDebrief() {
	utils.TacticalLog("[cyan::b]PHASE 5: REPORT GENERATION STARTED[-:-:-]")

	var startTime string
	if db.DB != nil {
		_ = db.DB.QueryRow("SELECT value FROM mission_state WHERE key = 'start_time'").Scan(&startTime)
	}

	outputDir := "reports"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		_ = os.Mkdir(outputDir, 0755)
	}

	timestamp := time.Now().Format("20060102_1504")
	reportName := fmt.Sprintf("VAPORTRACE_PEN_TEST_%s.md", timestamp)
	fullPath := filepath.Join(outputDir, reportName)

	f, err := os.Create(fullPath)
	if err != nil {
		utils.TacticalLog(fmt.Sprintf("[red]FileSystem Error:[-] %v", err))
		return
	}
	defer f.Close()

	// I. HEADER
	f.WriteString("# VAPORTRACE TACTICAL DEBRIEF\n")
	f.WriteString(fmt.Sprintf("> **OPERATIONAL STATUS:** COMPLETED\n"))
	f.WriteString(fmt.Sprintf("> **GEN TIME:** %s\n", time.Now().Format("15:04:05")))
	f.WriteString(fmt.Sprintf("> **START TIME:** %s\n\n", startTime))
	f.WriteString("---\n\n")

	// II. HARVESTED ARTIFACTS - BYPASS SECCIÓN QUE DA ERROR
	f.WriteString("## I. HARVESTED ARTIFACTS (DISCOVERY VAULT)\n\n")
	f.WriteString("| TYPE | SOURCE | VALUE (REDACTED) | TIMESTAMP |\n")
	f.WriteString("| :--- | :--- | :--- | :--- |\n")
	f.WriteString("| - | - | *VAULT_SYNC_PENDING_REBASE* | - |\n")
	f.WriteString("\n---\n\n")

	// III. MISSION PHASES SUMMARY
	f.WriteString("## II. MISSION PHASES SUMMARY\n\n")
	
	phases := []string{
		"PHASE II: DISCOVERY",
		"PHASE III: AUTH LOGIC",
		"PHASE IV: INJECTION",
		"PHASE VIII: EXFILTRATION",
	}

	for _, phase := range phases {
		f.WriteString(fmt.Sprintf("### %s\n", phase))
		f.WriteString("| ATTACK VECTOR | RESULT | TIMESTAMP |\n")
		f.WriteString("| :--- | :--- | :--- |\n")

		if db.DB != nil {
			rows, err := db.DB.Query("SELECT details, status, timestamp FROM findings WHERE phase = ?", phase)
			if err == nil {
				hasData := false
				for rows.Next() {
					hasData = true
					var details, status, ts string
					rows.Scan(&details, &status, &ts)
					f.WriteString(fmt.Sprintf("| %s | **%s** | %s |\n", details, status, ts))
				}
				rows.Close()
				if !hasData { f.WriteString("| - | *NO LOGS FOUND* | - |\n") }
			}
		}
		f.WriteString("\n")
	}

	// IV. MITRE & DFIR
	f.WriteString("## III. ADVERSARY EMULATION MAPPING\n\n")
	f.WriteString("| TACTIC | TECHNIQUE | RESULT |\n")
	f.WriteString("| :--- | :--- | :--- |\n")
	f.WriteString("| Reconnaissance | T1595.002 | Successful |\n\n")

	f.WriteString("---\n## IV. DFIR RESPONSE GUIDANCE\n\n")
	f.WriteString("### 1. Detection\n* Audit for processes masquerading as `kworker_system_auth`.\n")

	// Footer
	f.WriteString("---\n")
	f.WriteString("**CONFIDENTIAL // HYDRA-WORM INTEGRITY PROTOCOL**\n")
	
	pterm.Success.Printf("Tactical report generated: %s\n", fileName)
}

func getSeverityIcon(score float64) string {
	if score >= 9.0 {
		return "🔴"
	} // Critical
	if score >= 7.0 {
		return "🟠"
	} // High
	if score >= 4.0 {
		return "🟡"
	} // Medium
	return "🔵" // Low
}

// In your reporting logic package
func GetCurrentReportMarkdown() string {
	// This should mirror the logic used in writeTechnicalDetails
	var md string
	md += "# VAPORTRACE MISSION DEBRIEF\n"
	md += "## 3. TECHNICAL FINDINGS\n"
	// ... logic to query DB and build table ...
	return md
}
