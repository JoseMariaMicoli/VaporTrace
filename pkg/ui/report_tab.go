package ui

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/JoseMariaMicoli/VaporTrace/pkg/db"
	"github.com/JoseMariaMicoli/VaporTrace/pkg/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// InitReportTab creates the component for the F7 page.
// It assigns the component to the global reportView defined in dashboard.go.
func InitReportTab() *tview.TextView {
	reportView = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetScrollable(true)

	reportView.SetTitle(" [white:red] MISSION REPORT (PREVIEW) [white] - [yellow]Ctrl+W: Save[-] | [red]Ctrl+X: Clear[-] ").
		SetBorder(true).
		SetBorderColor(tcell.ColorRed)

	// Input Capture for Editor Actions
	reportView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlW:
			SaveReport()
			return nil
		case tcell.KeyCtrlX:
			ClearReport()
			return nil
		}
		return event
	})

	return reportView
}

// LoadFindings fetches data from DB and renders Markdown with colors
func LoadFindings() {
	if db.DB == nil {
		if reportView != nil {
			reportView.SetText("[red]Database not initialized. Run 'init_db' first.[-]")
		}
		return
	}

	var sb strings.Builder
	sb.WriteString("[yellow]# VAPORTRACE MISSION DEBRIEF[-]\n")
	sb.WriteString(fmt.Sprintf("[gray]Generated: %s[-]\n\n", time.Now().Format("2006-01-02 15:04:05")))

	// Stats
	var count int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM findings").Scan(&count)
	if err != nil {
		count = 0
	}
	sb.WriteString(fmt.Sprintf("[blue]Total Findings:[-] [white]%d[-]\n", count))
	sb.WriteString("[white]--------------------------------------------------[-]\n\n")

	// Query Findings
	rows, err := db.DB.Query("SELECT status, owasp_id, target, details, cvss_numeric FROM findings ORDER BY cvss_numeric DESC")
	if err != nil {
		if reportView != nil {
			reportView.SetText(fmt.Sprintf("[red]Error fetching data: %v[-]", err))
		}
		return
	}
	defer rows.Close()

	for rows.Next() {
		var status, owasp, target, details string
		var cvss float64
		rows.Scan(&status, &owasp, &target, &details, &cvss)

		// Color Logic
		color := "[blue]"
		if cvss >= 9.0 {
			color = "[red]"
		} else if cvss >= 7.0 {
			color = "[orange]"
		} else if cvss >= 4.0 {
			color = "[yellow]"
		}

		// Markdown Formatting
		sb.WriteString(fmt.Sprintf("%s### [%s] %s (CVSS: %.1f)[-]\n", color, status, owasp, cvss))
		sb.WriteString(fmt.Sprintf("[green]Target:[-] %s\n", target))
		sb.WriteString(fmt.Sprintf("[white]Details:[-] %s\n", details))
		sb.WriteString("\n")
	}

	if reportView != nil {
		reportView.SetText(sb.String())
		reportView.ScrollToBeginning()
	}
}

// SaveReport exports the buffer to a file
func SaveReport() {
	if reportView == nil {
		return
	}
	content := reportView.GetText(true) // Get clean text without color tags
	filename := fmt.Sprintf("VaporTrace_Report_%s.md", time.Now().Format("20060102_150405"))

	err := os.WriteFile(filename, []byte(content), 0644)
	if err != nil {
		utils.TacticalLog(fmt.Sprintf("[red]REPORT ERROR: Could not save file: %v[-]", err))
	} else {
		utils.TacticalLog(fmt.Sprintf("[green]REPORT SAVED: %s[-]", filename))
	}
}

// ClearReport wipes the buffer
func ClearReport() {
	if reportView != nil {
		reportView.SetText("[gray]Report buffer cleared. Waiting for new data...[-]")
	}
	utils.TacticalLog("[yellow]REPORT: Buffer wiped.[-]")
}
