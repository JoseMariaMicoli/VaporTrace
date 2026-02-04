package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
	"github.com/rivo/tview"
)

const (
	ReportsDir       = "reports"
	SessionDraftFile = "reports/current_session_draft.md"
)

// Global State for Dirty Checking (used by Status Bar/LoadFindings)
var IsReportDirty bool = false

// EnsureReportDir creates the output directory on startup
func EnsureReportDir() {
	if _, err := os.Stat(ReportsDir); os.IsNotExist(err) {
		os.Mkdir(ReportsDir, 0755)
	}
}

// GenerateBaseTemplate queries the DB (same data as generator.go) to create an editable string
// Task 3: Markdown generation for the TextArea
func GenerateBaseTemplate() string {
	var sb strings.Builder

	sb.WriteString("# VAPORTRACE TACTICAL DEBRIEF\n")
	sb.WriteString(fmt.Sprintf("**Date:** %s\n", time.Now().Format("2006-01-02 15:04:05")))
	sb.WriteString("**Status:** [DRAFT]\n\n")

	sb.WriteString("## 1. EXECUTIVE SUMMARY\n")
	sb.WriteString("<!-- EDITABLE: Write high-level business impact here -->\n")
	sb.WriteString("Security assessment conducted using VaporTrace. Analysis indicates several critical control failures regarding authentication and input validation.\n\n")

	sb.WriteString("## 2. AUTOMATED FINDINGS (DATABASE EXPORT)\n")
	sb.WriteString("| SEVERITY | OWASP | TARGET | DETAILS |\n")
	sb.WriteString("| :--- | :--- | :--- | :--- |\n")

	if db.DB != nil {
		// Mirroring logic from report/generator.go but for UI buffer
		rows, err := db.DB.Query("SELECT status, owasp_id, target, details, cvss_numeric FROM findings ORDER BY cvss_numeric DESC")
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var status, owasp, target, details string
				var cvss float64
				rows.Scan(&status, &owasp, &target, &details, &cvss)

				// ASCII Indicators for Table in Editor
				indicator := "🔵" // Low
				if cvss >= 9.0 {
					indicator = "🔴" // Critical
				} else if cvss >= 7.0 {
					indicator = "🟠" // High
				} else if cvss >= 4.0 {
					indicator = "🟡" // Medium
				}

				// Sanitize table columns (basic replacement to avoid breaking markdown tables)
				cleanDetails := strings.ReplaceAll(details, "|", "/")
				sb.WriteString(fmt.Sprintf("| %s %.1f | %s | %s | %s |\n", indicator, cvss, owasp, target, cleanDetails))
			}
		}
	} else {
		sb.WriteString("| - | - | DATABASE NOT CONNECTED | - |\n")
	}

	sb.WriteString("\n## 3. TECHNICAL CONCLUSION\n")
	sb.WriteString("<!-- EDITABLE: Add manual observations here -->\n")
	sb.WriteString("Remediation of Critical vulnerabilities is recommended within 24-48 hours.\n")

	return sb.String()
}

// SaveReportDisk writes the current buffer to a permanent timestamped file
func SaveReportDisk(content string) (string, error) {
	EnsureReportDir()
	filename := fmt.Sprintf("VaporTrace_Final_%s.md", time.Now().Format("20060102_150405"))
	fullPath := filepath.Join(ReportsDir, filename)

	err := os.WriteFile(fullPath, []byte(content), 0644)
	if err == nil {
		IsReportDirty = false // Reset dirty flag
	}
	return filename, err
}

// DeleteSession Logic - Purges Data and Views
// Requires references to the app and pages to display the confirmation modal
func DeleteSession(app *tview.Application, pages *tview.Pages, clearView func()) {
	// Show Confirmation Modal
	modal := tview.NewModal().
		SetText("[red::b]WARNING: PROCEED WITH SESSION PURGE?[-::-]\nThis will delete:\n1. In-Memory Report Draft\n2. Local Database Findings").
		AddButtons([]string{"CONFIRM", "CANCEL"}).
		SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "CONFIRM" {
				// 1. Wipe DB
				db.ResetDB()
				// 2. Try Remove Draft (if we were saving one)
				os.Remove(SessionDraftFile)
				// 3. Reset UI View via callback
				clearView()
				IsReportDirty = false
				utils.TacticalLog("[red]SESSION WIPED: All intel deleted.[-]")
			}
			// Close Modal (Return to previous page, handled by caller mostly or manual switch)
			pages.RemovePage("modal_purge")
			// Return focus using reportEditor because reportView is undefined
			if reportEditor != nil {
				app.SetFocus(reportEditor)
			}
		})

	pages.AddPage("modal_purge", modal, false, true)
}
