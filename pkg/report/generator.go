package report

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

// GenerateMissionDebrief compiles findings into a corporate tactical audit report.
func GenerateMissionDebrief() {
	utils.TacticalLog("[cyan::b]PHASE 5: REPORT GENERATION STARTED[-:-:-]")

	// 1. Setup Corporate Directory and Naming
	outputDir := "reports"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		_ = os.Mkdir(outputDir, 0755)
	}

	// Format: VAPORTRACE_TACTICAL_AUDIT_20260131_1030.md
	timestamp := time.Now().Format("20060102_1504")
	reportName := fmt.Sprintf("VAPORTRACE_TACTICAL_AUDIT_%s.md", timestamp)
	fullPath := filepath.Join(outputDir, reportName)

	var startTime string
	if db.DB != nil {
		_ = db.DB.QueryRow("SELECT value FROM mission_state WHERE key = 'start_time'").Scan(&startTime)
	}

	f, err := os.Create(fullPath)
	if err != nil {
		utils.TacticalLog(fmt.Sprintf("[red]FileSystem Error:[-] %v", err))
		return
	}
	defer f.Close()

	utils.TacticalLog(fmt.Sprintf("[blue]⠋[-] Writing to [white]%s/[-][yellow]%s[-]...", outputDir, reportName))

	// I. CORPORATE HEADER
	f.WriteString("# VAPORTRACE TACTICAL AUDIT REPORT\n")
	f.WriteString("## CONFIDENTIAL - FOR INTERNAL USE ONLY\n\n")
	f.WriteString("| METADATA | VALUE |\n")
	f.WriteString("| :--- | :--- |\n")
	f.WriteString(fmt.Sprintf("| **AUDIT STATUS** | COMPLETED |\n"))
	f.WriteString(fmt.Sprintf("| **MISSION START** | %s |\n", startTime))
	f.WriteString(fmt.Sprintf("| **GEN TIME (UTC)** | %s |\n", time.Now().Format("2006-01-02 15:04:05")))
	f.WriteString(fmt.Sprintf("| **CLASSIFICATION** | PROPRIETARY / ADVERSARY EMULATION |\n\n"))
	f.WriteString("---\n\n")

	// II. EXECUTIVE SUMMARY
	f.WriteString("## 1. EXECUTIVE SUMMARY\n")
	f.WriteString("This document provides a formal tactical debrief of offensive security operations. ")
	f.WriteString("Results are mapped to the **MITRE ATT&CK** and **OWASP API 2023** frameworks to facilitate risk-based remediation.\n\n")

	// III. TACTICAL FINDINGS (Grouped by Phase)
	f.WriteString("## 2. DETAILED TACTICAL FINDINGS\n\n")

	rows, err := db.DB.Query("SELECT DISTINCT phase FROM findings ORDER BY phase ASC")
	if err != nil {
		utils.TacticalLog("[red]Database Error:[-] Failed to fetch mission data.")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var phase string
		rows.Scan(&phase)

		f.WriteString(fmt.Sprintf("### PHASE: %s\n", phase))
		f.WriteString("| TARGET | MITRE ID | OWASP ID | NIST TAG | STATUS | OBSERVATION |\n")
		f.WriteString("| :--- | :--- | :--- | :--- | :--- | :--- |\n")

		fRows, err := db.DB.Query("SELECT target, mitre_id, owasp_id, nist_tag, status, details FROM findings WHERE phase = ?", phase)
		if err == nil {
			for fRows.Next() {
				var target, mitre, owasp, nist, status, details string
				fRows.Scan(&target, &mitre, &owasp, &nist, &status, &details)

				// Apply formatting for empty fields
				if owasp == "" {
					owasp = "N/A"
				}
				if nist == "" {
					nist = "N/A"
				}

				f.WriteString(fmt.Sprintf("| `%s` | %s | %s | %s | **%s** | %s |\n", target, mitre, owasp, nist, status, details))
			}
			fRows.Close()
		}
		f.WriteString("\n")
	}

	// IV. COMPLIANCE MAPPING
	f.WriteString("---\n## 3. FRAMEWORK ALIGNMENT\n\n")
	f.WriteString("| VAPORTRACE COMPONENT | MITRE TACTIC | OWASP TOP 10 | DEFENSIVE CONTEXT |\n")
	f.WriteString("| :--- | :--- | :--- | :--- |\n")
	f.WriteString("| I. INFIL | Reconnaissance | API9:2023 | Inventory Management |\n")
	f.WriteString("| II. EXPLOIT | Priv Escalation | API1 / API2 / API5 | Identity Assurance |\n")
	f.WriteString("| III. EXPAND | Discovery / Impact | API4 / API7 | Resource Hardening |\n")
	f.WriteString("| IV. OBFUSC | Defense Evasion | N/A | Stealth Analysis |\n\n")

	// V. REMEDIATION GUIDANCE
	f.WriteString("---\n## 4. REMEDIATION STRATEGY\n\n")
	f.WriteString("### Immediate Hardening\n")
	f.WriteString("1. **Inventory Control:** Remove all exposed Swagger documentation identified in Phase I.\n")
	f.WriteString("2. **Auth Validation:** Patch JWT algorithm processing to deny 'none' or 'weak' signing keys.\n")
	f.WriteString("3. **Rate Limiting:** Implement circuit-breaking on high-resource endpoints to mitigate DoS findings.\n\n")

	// Footer
	f.WriteString("---\n")
	f.WriteString("**UNAUTHORIZED DISCLOSURE OF THIS REPORT IS PROHIBITED**\n")

	utils.TacticalLog(fmt.Sprintf("[green]✔[-] Corporate Audit Package written to: [yellow]%s[-]", fullPath))
}
