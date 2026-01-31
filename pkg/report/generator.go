package report

import (
	"fmt"
	"os"
	"time"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

// GenerateMissionDebrief compiles findings into a professional tactical report.
func GenerateMissionDebrief() {
	utils.TacticalLog("[cyan::b]PHASE 5: REPORT GENERATION STARTED[-:-:-]")

	var startTime string
	if db.DB != nil {
		_ = db.DB.QueryRow("SELECT value FROM mission_state WHERE key = 'start_time'").Scan(&startTime)
	}

	reportDate := time.Now().Format("2006-01-02")
	fileName := fmt.Sprintf("VAPOR_DEBRIEF_%s.md", reportDate)

	f, err := os.Create(fileName)
	if err != nil {
		utils.TacticalLog(fmt.Sprintf("[red]FileSystem Error:[-] %v", err))
		return
	}
	defer f.Close()

	utils.TacticalLog("[blue]⠋[-] Processing VaporTrace Tactical Mapping...")

	// I. HEADER & EXECUTIVE SUMMARY
	f.WriteString("# VAPORTRACE TACTICAL DEBRIEF\n")
	f.WriteString(fmt.Sprintf("> **OPERATIONAL STATUS:** COMPLETED\n"))
	f.WriteString(fmt.Sprintf("> **MISSION START:** %s\n", startTime))
	f.WriteString(fmt.Sprintf("> **REPORT GENERATED:** %s\n\n", time.Now().Format("2006-01-02 15:04:05")))
	f.WriteString("## EXECUTIVE SUMMARY\n")
	f.WriteString("This report details the offensive security testing results synchronized with the adversary lifecycle. ")
	f.WriteString("Findings are mapped to **MITRE ATT&CK**, **OWASP API Top 10**, and **NIST CSF** for compliance auditing.\n\n")
	f.WriteString("---\n\n")

	// II. TACTICAL FINDINGS BY VAPORTRACE BLOCK
	f.WriteString("## I. TACTICAL FINDINGS LOG\n\n")

	// We query the distinct VaporTrace phases present in the database
	rows, err := db.DB.Query("SELECT DISTINCT phase FROM findings ORDER BY phase ASC")
	if err != nil {
		utils.TacticalLog("[red]Database Error:[-] Failed to fetch phases.")
		return
	}
	defer rows.Close()

	var phases []string
	for rows.Next() {
		var p string
		rows.Scan(&p)
		phases = append(phases, p)
	}

	for _, phase := range phases {
		f.WriteString(fmt.Sprintf("### %s\n", phase))
		f.WriteString("| TARGET | MITRE | OWASP | NIST | STATUS | DETAILS |\n")
		f.WriteString("| :--- | :--- | :--- | :--- | :--- | :--- |\n")

		fRows, err := db.DB.Query("SELECT target, mitre_id, owasp_id, nist_tag, status, details FROM findings WHERE phase = ?", phase)
		if err == nil {
			for fRows.Next() {
				var target, mitre, owasp, nist, status, details string
				fRows.Scan(&target, &mitre, &owasp, &nist, &status, &details)
				// Clean empty fields
				if owasp == "" {
					owasp = "-"
				}
				if nist == "" {
					nist = "-"
				}

				f.WriteString(fmt.Sprintf("| `%s` | %s | %s | %s | **%s** | %s |\n", target, mitre, owasp, nist, status, details))
			}
			fRows.Close()
		}
		f.WriteString("\n")
	}

	// III. FRAMEWORK COMPLIANCE SUMMARY
	f.WriteString("---\n## II. FRAMEWORK MAPPING SUMMARY\n\n")

	// Quick Summary Table for Compliance Officers
	f.WriteString("### MITRE ATT&CK & OWASP Tactical Mapping\n")
	f.WriteString("| VAPORTRACE BLOCK | PRIMARY MITRE | PRIMARY OWASP | DEFENSIVE CONTEXT |\n")
	f.WriteString("| :--- | :--- | :--- | :--- |\n")
	f.WriteString("| I. INFIL | T1595 / T1592 | API9:2023 | Shadow API Discovery |\n")
	f.WriteString("| II. EXPLOIT | T1548 / T1606 | API1, API2, API5 | Identity & Access |\n")
	f.WriteString("| III. EXPAND | T1046 / T1499 | API4, API7 | Infrastructure Risk |\n")
	f.WriteString("| IV. OBFUSC | T1090 / T1562 | - | Stealth & Evasion |\n")
	f.WriteString("| V. COMPL | T1020 | - | Evidence Packaging |\n\n")

	// IV. DFIR / REMEDIATION GUIDANCE
	f.WriteString("---\n## III. DFIR RESPONSE & REMEDIATION\n\n")
	f.WriteString("### 1. Critical Hardening Actions\n")
	f.WriteString("- **API Inventory:** Decommission all endpoints identified in `I. INFIL` (Swagger leaks).\n")
	f.WriteString("- **Auth Logic:** Implement strict JWT signature validation to prevent `II. EXPLOIT` algorithm downgrades.\n")
	f.WriteString("- **Rate Limiting:** Apply aggressive 429 response codes to targets identified in `III. EXPAND` (DoS probes).\n\n")

	f.WriteString("### 2. Detection Logic\n")
	f.WriteString("* Monitor for anomalous User-Agent strings and high-frequency IP rotation (VaporTrace IV. OBFUSC).\n")
	f.WriteString("* Audit cloud IAM logs for calls to the Metadata Service (169.254.169.254).\n\n")

	// Footer
	f.WriteString("---\n")
	f.WriteString(fmt.Sprintf("**CONFIDENTIAL // VAPORTRACE GENERATED DEBRIEF // %s**\n", reportDate))

	utils.TacticalLog(fmt.Sprintf("[green]✔[-] Tactical report generated: [yellow]%s[-]", fileName))
}
