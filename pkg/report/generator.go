package report

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
)

// Stats holds the aggregate data for the Executive Summary
type MissionStats struct {
	TotalFindings int
	Critical      int
	High          int
	Medium        int
	Low           int
	AvgCVSS       float64
	UniqueTargets int
}

// GenerateMissionDebrief compiles findings into a professional C-Level + Technical report.
func GenerateMissionDebrief() {
	utils.TacticalLog("[cyan::b]PHASE 5: REPORT GENERATION STARTED[-:-:-]")

	if db.DB == nil {
		utils.TacticalLog("[red]Database Error:[-] Connection not initialized. Run 'init_db' first.")
		return
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

	// 1. GATHER STATISTICS
	stats := calculateStats()

	// 2. WRITE SECTIONS
	writeHeader(f)
	writeExecutiveSummary(f, stats)
	writeRemediationTracker(f)
	writeTechnicalDetails(f)
	writeMethodology(f)

	utils.TacticalLog(fmt.Sprintf("[green]✔[-] Audit Report Generated: [yellow]%s[-]", fullPath))
}

// --- HELPER FUNCTIONS ---

func calculateStats() MissionStats {
	var stats MissionStats

	// Get Totals and Average CVSS
	db.DB.QueryRow("SELECT COUNT(*), COALESCE(AVG(cvss_numeric), 0.0) FROM findings").Scan(&stats.TotalFindings, &stats.AvgCVSS)

	// Get Counts by CVSS Range
	db.DB.QueryRow("SELECT COUNT(*) FROM findings WHERE cvss_numeric >= 9.0").Scan(&stats.Critical)
	db.DB.QueryRow("SELECT COUNT(*) FROM findings WHERE cvss_numeric >= 7.0 AND cvss_numeric < 9.0").Scan(&stats.High)
	db.DB.QueryRow("SELECT COUNT(*) FROM findings WHERE cvss_numeric >= 4.0 AND cvss_numeric < 7.0").Scan(&stats.Medium)
	db.DB.QueryRow("SELECT COUNT(*) FROM findings WHERE cvss_numeric < 4.0").Scan(&stats.Low)

	// Get Unique Targets count
	db.DB.QueryRow("SELECT COUNT(DISTINCT target) FROM findings").Scan(&stats.UniqueTargets)

	return stats
}

func writeHeader(f *os.File) {
	startTime := time.Now().Format("2006-01-02 15:04:05")
	// Try to get actual start time from DB if available
	var dbStart string
	_ = db.DB.QueryRow("SELECT value FROM mission_state WHERE key = 'start_time'").Scan(&dbStart)
	if dbStart != "" {
		startTime = dbStart
	}

	f.WriteString("# VAPORTRACE PENETRATION TEST REPORT\n")
	f.WriteString("**CONFIDENTIAL - INTERNAL USE ONLY**\n\n")
	f.WriteString("| META | VALUE |\n")
	f.WriteString("| :--- | :--- |\n")
	f.WriteString(fmt.Sprintf("| **DATE** | %s |\n", time.Now().Format("2006-01-02")))
	f.WriteString(fmt.Sprintf("| **MISSION START** | %s |\n", startTime))
	f.WriteString("| **CLASSIFICATION** | PROPRIETARY / ADVERSARY EMULATION |\n")
	f.WriteString("| **ENGINE VERSION** | VaporTrace v3.1 (Tactical Suite) |\n\n")
	f.WriteString("---\n\n")
}

func writeExecutiveSummary(f *os.File, stats MissionStats) {
	f.WriteString("## 1. EXECUTIVE SUMMARY\n\n")
	f.WriteString("### 1.1 Risk Overview\n")
	f.WriteString("VaporTrace Tactical Suite performed an automated adversarial emulation against the target infrastructure. ")
	f.WriteString("This section provides a high-level overview of the security posture based on discovered vulnerabilities.\n\n")

	// Dynamic Risk Rating
	riskRating := "LOW"
	if stats.Critical > 0 {
		riskRating = "CRITICAL"
	} else if stats.High > 0 {
		riskRating = "HIGH"
	} else if stats.Medium > 0 {
		riskRating = "MODERATE"
	}

	f.WriteString(fmt.Sprintf("**OVERALL RISK RATING:** %s\n\n", riskRating))

	f.WriteString("| METRIC | VALUE |\n")
	f.WriteString("| :--- | :--- |\n")
	f.WriteString(fmt.Sprintf("| **Total Findings** | %d |\n", stats.TotalFindings))
	f.WriteString(fmt.Sprintf("| **Unique Targets** | %d |\n", stats.UniqueTargets))
	f.WriteString(fmt.Sprintf("| **Average CVSS** | %.1f / 10.0 |\n", stats.AvgCVSS))
	f.WriteString("\n")

	f.WriteString("### 1.2 Vulnerability Distribution\n")
	f.WriteString("Breakdown of findings by severity (CVSS v3.1):\n\n")

	// ASCII Chart Logic
	total := float64(stats.TotalFindings)
	if total == 0 {
		total = 1
	} // Prevent div by zero

	f.WriteString(fmt.Sprintf("- **CRITICAL (9.0+):** %d  (%s)\n", stats.Critical, progressBar(stats.Critical, int(total))))
	f.WriteString(fmt.Sprintf("- **HIGH (7.0-8.9):**     %d  (%s)\n", stats.High, progressBar(stats.High, int(total))))
	f.WriteString(fmt.Sprintf("- **MEDIUM (4.0-6.9):**   %d  (%s)\n", stats.Medium, progressBar(stats.Medium, int(total))))
	f.WriteString(fmt.Sprintf("- **LOW (0.0-3.9):**      %d  (%s)\n", stats.Low, progressBar(stats.Low, int(total))))
	f.WriteString("\n---\n\n")
}

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

		// Clean strings
		owaspShort := strings.Split(owasp, ":")[0] // Just take API1, API2 etc

		// FIX: Use owaspShort in the formatted string
		f.WriteString(fmt.Sprintf("| %s | %.1f | %s | %s | `%s` | %s |\n",
			icon, cvss, owaspShort, cve, target, action))
		count++
	}

	if count == 0 {
		f.WriteString("| 🟢 | - | No actionable vulnerabilities found above Low severity. | - | - | - |\n")
	}
	f.WriteString("\n---\n\n")
}

func writeTechnicalDetails(f *os.File) {
	f.WriteString("## 3. TECHNICAL FINDINGS (DEEP DIVE)\n")
	f.WriteString("Detailed evidence logs for engineering teams. ")
	f.WriteString("**Sorted Chronologically (Execution Order).**\n\n")

	// Query: Select ALL fields, Order by Timestamp ASC
	query := `
		SELECT 
			timestamp, command, target, details, status, 
			owasp_id, mitre_id, mitre_tactic, nist_control, cve_id, cvss_numeric
		FROM findings 
		ORDER BY timestamp ASC`

	rows, err := db.DB.Query(query)
	if err != nil {
		f.WriteString(fmt.Sprintf("> Error fetching details: %v\n", err))
		return
	}
	defer rows.Close()

	for rows.Next() {
		var ts, cmd, target, details, status, owasp, mitreId, mitreTac, nist, cve string
		var cvss float64

		rows.Scan(&ts, &cmd, &target, &details, &status, &owasp, &mitreId, &mitreTac, &nist, &cve, &cvss)

		// Create a "Card" for each finding
		f.WriteString(fmt.Sprintf("### [%s] %s on %s\n", status, owasp, target))
		f.WriteString(fmt.Sprintf("- **Timestamp:** %s\n", ts))
		f.WriteString(fmt.Sprintf("- **Vector/Command:** `%s`\n", cmd))
		f.WriteString(fmt.Sprintf("- **Target URL:** `%s`\n", target))
		f.WriteString(fmt.Sprintf("- **Details:** %s\n", details))

		f.WriteString("\n**Compliance Mapping:**\n")
		f.WriteString("| Framework | ID / Control | Description / Tactic |\n")
		f.WriteString("| :--- | :--- | :--- |\n")
		f.WriteString(fmt.Sprintf("| **MITRE ATT&CK** | `%s` | %s |\n", mitreId, mitreTac))
		f.WriteString(fmt.Sprintf("| **NIST CSF v2.0** | `%s` | Control Mapping |\n", nist))
		f.WriteString(fmt.Sprintf("| **CVE / CVSS** | `%s` | **%.1f** (Severity Score) |\n", cve, cvss))

		f.WriteString("\n---\n")
	}
}

func writeMethodology(f *os.File) {
	f.WriteString("## 4. METHODOLOGY & FRAMEWORK ALIGNMENT\n\n")
	f.WriteString("This assessment was conducted using the **VaporTrace Tactical Engine**, adhering to standard Adversary Emulation protocols.\n\n")

	f.WriteString("### 4.1 Framework Reference\n")
	f.WriteString("- **MITRE ATT&CK:** Used to classify adversary tactics and techniques (T-Codes).\n")
	f.WriteString("- **OWASP API Security Top 10 (2023):** Primary standard for API vulnerability classification.\n")
	f.WriteString("- **NIST CSF v2.0:** Used for mapping findings to defensive controls (Identify, Protect, Detect, Respond, Recover).\n")
	f.WriteString("- **CVSS v3.1:** Common Vulnerability Scoring System for severity quantification.\n\n")

	f.WriteString("**End of Report**\n")
}

// ASCII Progress Bar generator
func progressBar(count, total int) string {
	if total == 0 {
		return "░░░░░░░░░░"
	}
	percent := float64(count) / float64(total)
	barLen := 10
	filledLen := int(percent * float64(barLen))

	bar := ""
	for i := 0; i < barLen; i++ {
		if i < filledLen {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	return bar
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
