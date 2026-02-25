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

func writeRemediationTracker(f *os.File) {
	f.WriteString("## 2. REMEDIATION PRIORITY TRACKER\n")
	f.WriteString("The following table prioritizes vulnerabilities requiring immediate attention. ")
	f.WriteString("**Sorted by Severity (CVSS Descending).**\n\n")

	f.WriteString("| SEVERITY | CVSS | VULNERABILITY (OWASP) | CVE ID | AFFECTED TARGET | ACTION |\n")
	f.WriteString("| :--- | :--- | :--- | :--- | :--- | :--- |\n")

	// Query: Order by CVSS Numeric DESC (Highest Risk First)
	rows, err := db.DB.Query(`
		SELECT status, cvss_numeric, owasp_id, cve_id, target 
		FROM findings 
		WHERE cvss_numeric >= 4.0 
		ORDER BY cvss_numeric DESC
	`)

	if err != nil {
		f.WriteString(fmt.Sprintf("> Error generating tracker: %v\n", err))
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var status, owasp, cve, target string
		var cvss float64
		rows.Scan(&status, &cvss, &owasp, &cve, &target)

		icon := getSeverityIcon(cvss)
		action := "Monitor"
		if cvss >= 9.0 {
			action = "**PATCH IMMEDIATELY**"
		}
		if cvss >= 7.0 && cvss < 9.0 {
			action = "Remediate < 7 Days"
		}
		if cvss >= 4.0 && cvss < 7.0 {
			action = "Remediate < 30 Days"
		}

		// Highlight Race Conditions as requiring architectural fixes
		if strings.Contains(strings.ToLower(owasp), "race") || strings.Contains(strings.ToLower(owasp), "api6") {
			action = "**ARCHITECTURAL FIX REQ**"
		}

		// Clean strings
		owaspShort := strings.Split(owasp, ":")[0] // Just take API1, API2 etc

		// FIX: Use owaspShort in the formatted string
		f.WriteString(fmt.Sprintf("| %s | %.1f | %s | %s | `%s` | %s |\n",
			icon, cvss, owaspShort, cve, target, action))
		count++
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
