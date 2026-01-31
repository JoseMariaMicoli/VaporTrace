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

	// Ensure DB is accessible
	if db.DB == nil {
		utils.TacticalLog("[red]Database Error:[-] Connection not initialized. Run 'init_db' first.")
		return
	}

	outputDir := "reports"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		_ = os.Mkdir(outputDir, 0755)
	}

	timestamp := time.Now().Format("20060102_1504")
	reportName := fmt.Sprintf("VAPORTRACE_TACTICAL_AUDIT_%s.md", timestamp)
	fullPath := filepath.Join(outputDir, reportName)

	var startTime string
	_ = db.DB.QueryRow("SELECT value FROM mission_state WHERE key = 'start_time'").Scan(&startTime)
	if startTime == "" {
		startTime = time.Now().Format("2006-01-02 15:04:05")
	}

	f, err := os.Create(fullPath)
	if err != nil {
		utils.TacticalLog(fmt.Sprintf("[red]FileSystem Error:[-] %v", err))
		return
	}
	defer f.Close()

	utils.TacticalLog(fmt.Sprintf("[blue]⠋[-] Writing to [white]%s/[-][yellow]%s[-]...", outputDir, reportName))

	// --- SECTION 1: HEADER ---
	f.WriteString("# VAPORTRACE TACTICAL AUDIT REPORT\n")
	f.WriteString("## CONFIDENTIAL - FOR INTERNAL USE ONLY\n\n")
	f.WriteString("| METADATA | VALUE |\n")
	f.WriteString("| :--- | :--- |\n")
	f.WriteString("| **AUDIT STATUS** | COMPLETED |\n")
	f.WriteString(fmt.Sprintf("| **MISSION START** | %s |\n", startTime))
	f.WriteString(fmt.Sprintf("| **GEN TIME (UTC)** | %s |\n", time.Now().Format("2006-01-02 15:04:05")))
	f.WriteString("| **CLASSIFICATION** | PROPRIETARY / ADVERSARY EMULATION |\n\n")
	f.WriteString("---\n\n")

	// --- SECTION 2: EXECUTIVE SUMMARY ---
	f.WriteString("## 1. EXECUTIVE SUMMARY\n")
	f.WriteString("This document provides a formal tactical debrief of offensive security operations. ")
	f.WriteString("Results are mapped to the **MITRE ATT&CK** and **OWASP API 2023** frameworks.\n\n")

	// --- SECTION 3: DETAILED FINDINGS ---
	f.WriteString("## 2. DETAILED TACTICAL FINDINGS\n\n")

	rows, err := db.DB.Query("SELECT DISTINCT phase FROM findings ORDER BY phase ASC")
	if err != nil {
		utils.TacticalLog("[red]Database Error:[-] Failed to fetch mission data.")
		return
	}
	defer rows.Close()

	phaseCount := 0
	for rows.Next() {
		var phase string
		rows.Scan(&phase)
		phaseCount++

		f.WriteString(fmt.Sprintf("### PHASE: %s\n", phase))
		f.WriteString("| TARGET | MITRE | OWASP | NIST | CVE | CVSS | STATUS | DETAILS |\n")
		f.WriteString("| :--- | :--- | :--- | :--- | :--- | :--- | :--- | :--- |\n")

		// Retrieve all columns including new CVE/CVSS
		fRows, err := db.DB.Query(`
			SELECT target, mitre_id, owasp_id, nist_tag, cve_id, cvss_score, status, details 
			FROM findings WHERE phase = ?`, phase)

		if err == nil {
			for fRows.Next() {
				var target, mitre, owasp, nist, cve, cvss, status, details string
				// Use pointer scanning to handle potential NULLs (though Schema sets defaults)
				fRows.Scan(&target, &mitre, &owasp, &nist, &cve, &cvss, &status, &details)

				// Formatting
				f.WriteString(fmt.Sprintf("| `%s` | %s | %s | %s | %s | %s | **%s** | %s |\n",
					target, mitre, owasp, nist, cve, cvss, status, details))
			}
			fRows.Close()
		}
		f.WriteString("\n")
	}

	if phaseCount == 0 {
		f.WriteString("> No findings recorded in this session.\n")
	}

	// --- SECTION 4: REMEDIATION TRACKER ---
	f.WriteString("## 3. REMEDIATION PRIORITY TRACKER\n")
	f.WriteString("The following findings are prioritized by **CVSS Score** (High to Low) to assist in triage.\n\n")
	f.WriteString("| SEV | CVSS | CVE ID | VULNERABILITY DETAIL | AFFECTED TARGET |\n")
	f.WriteString("| :--- | :--- | :--- | :--- | :--- |\n")

	// Query prioritizing Critical/Exploited findings by CVSS
	remRows, err := db.DB.Query(`
		SELECT status, cvss_score, cve_id, details, target 
		FROM findings 
		WHERE status IN ('CRITICAL', 'EXPLOITED', 'VULNERABLE') 
		ORDER BY CAST(cvss_score AS REAL) DESC
	`)

	if err == nil {
		count := 0
		for remRows.Next() {
			var status, cvss, cve, details, target string
			remRows.Scan(&status, &cvss, &cve, &details, &target)

			// Iconography
			icon := "⚪"
			if status == "EXPLOITED" {
				icon = "🟣"
			} // Purple Team / Exploited
			if status == "CRITICAL" {
				icon = "🔴"
			}
			if status == "VULNERABLE" {
				icon = "🟠"
			}

			f.WriteString(fmt.Sprintf("| %s | **%s** | %s | %s | `%s` |\n", icon, cvss, cve, details, target))
			count++
		}
		remRows.Close()
		if count == 0 {
			f.WriteString("| 🟢 | - | - | No Critical Vulnerabilities Detected | - |\n")
		}
	} else {
		f.WriteString(fmt.Sprintf("\n> Database Error: %v\n", err))
	}
	f.WriteString("\n")

	// --- SECTION 5: FOOTER ---
	f.WriteString("---\n")
	f.WriteString("### 4. FRAMEWORK ALIGNMENT\n\n")
	f.WriteString("| VAPORTRACE COMPONENT | MITRE TACTIC | OWASP TOP 10 | DEFENSIVE CONTEXT |\n")
	f.WriteString("| :--- | :--- | :--- | :--- |\n")
	f.WriteString("| I. INFIL | Reconnaissance | API9:2023 | Inventory Management |\n")
	f.WriteString("| II. EXPLOIT | Priv Escalation | API1 / API2 / API5 | Identity Assurance |\n")
	f.WriteString("| III. EXPAND | Discovery / Impact | API4 / API7 | Resource Hardening |\n")
	f.WriteString("| IV. OBFUSC | Defense Evasion | N/A | Stealth Analysis |\n\n")
	f.WriteString("---\n**UNAUTHORIZED DISCLOSURE OF THIS REPORT IS PROHIBITED**\n")

	utils.TacticalLog(fmt.Sprintf("[green]✔[-] Audit Report Generated: [yellow]%s[-]", fullPath))
}
